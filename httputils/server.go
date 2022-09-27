package httputils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/Uchencho/commons/commonerror"
	"github.com/Uchencho/commons/uuid"
	"github.com/gabriel-vasile/mimetype"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

const (
	successMessage  = "success"
	failureMessage  = "error"
	noReqBody       = "EOF"
	maxUploadMemory = 20000000
)

// GenericResponse is a representation of a server response
type GenericResponse struct {
	Status string      `json:"status,omitempty"`
	Data   interface{} `json:"data,omitempty"`
	Error  error       `json:"error,omitempty"`
}

// FileDetails is a representation of the details of a file
type FileDetails struct {
	UploadKey   string `json:"uploadKey"`
	FileName    string `json:"filename"`
	Data        []byte `json:"data"`
	ContentType string `json:"contentType"`
}

// FileWithBodyResult is a representation of the result of extracting both text and file uploads from a request
type FileWithBodyResult struct {
	Files []FileDetails
	Body  map[string][]string
}

func setStandardHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS,POST,GET")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
}

// NewSuccessResponse creates a new success response
func NewSuccessResponse(d interface{}) GenericResponse {
	return GenericResponse{
		Status: successMessage,
		Data:   d,
	}
}

// NewErrorResponse creates a new error response
func NewErrorResponse(err error) GenericResponse {
	return GenericResponse{
		Status: failureMessage,
		Error:  err,
	}
}

// ServeInternalError serves an internal server error
func serveInternalError(err error, w http.ResponseWriter) {

	log.Println("Error is: ", err)

	errDTO := commonerror.Error{
		Code:    commonerror.ErrorCode("Internal server error"),
		Message: commonerror.ErrorMessage("Something unplanned for has gone wrong"),
	}

	res := NewErrorResponse(errDTO)
	bb, err := json.Marshal(res)
	if err != nil {
		serveInternalError(err, w)
	}

	setStandardHeaders(w)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(bb)
}

// ServeBadRequestError serves a bad request error
func serveBadRequestError(err error, w http.ResponseWriter) {

	log.Println("Error is: ", err)

	errDTO, ok := err.(commonerror.Error)
	if !ok {
		errDTO = commonerror.Error{
			Message: commonerror.ErrorMessage(err.Error()),
		}
	}
	res := NewErrorResponse(errDTO)
	bb, err := json.Marshal(res)
	if err != nil {
		serveBadRequestError(err, w)
	}

	setStandardHeaders(w)
	w.WriteHeader(http.StatusBadRequest)
	w.Write(bb)
}

// served when user did not provide authorization
func serveUnauthorizedResponse(err error, w http.ResponseWriter) {
	log.Println("Error is: ", err)

	res := NewErrorResponse(err)
	bb, err := json.Marshal(res)
	if err != nil {
		serveInternalError(err, w)
	}

	setStandardHeaders(w)
	w.WriteHeader(http.StatusUnauthorized)
	w.Write(bb)
}

// served when user passed in authentication but they are invalid
func serveAuthenticationErrResponse(err error, w http.ResponseWriter) {
	log.Println("Error is : ", err)

	res := NewErrorResponse(err)
	bb, err := json.Marshal(res)
	if err != nil {
		serveInternalError(err, w)
	}

	setStandardHeaders(w)
	w.WriteHeader(http.StatusForbidden)
	w.Write(bb)
}

// JSONToDTO decodes a request body into a Data Transfer Object
func JSONToDTO(DTO interface{}, w http.ResponseWriter, r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(&DTO)
	if err != nil && err.Error() == noReqBody {
		return commonerror.NewErrorParams("error", "No request body was passed").ToBadRequest()
	}
	return err
}

// FileUploadToBytes extracts from a request an uploaded file
func FileUploadToBytes(r *http.Request, filename string) ([]byte, error) {

	file, _, err := r.FormFile(filename)
	if err != nil {
		return []byte{}, errors.Wrapf(err, "Unable to get file: %s", filename)
	}
	defer file.Close()

	return readFile(file)
}

// ExtractMultipleFileUploads extracts multiple uploaded files from the request body
func ExtractMultipleFileUploads(r *http.Request, uploadKeys []string) ([]FileDetails, error) {

	if err := r.ParseMultipartForm(maxUploadMemory); err != nil {
		return []FileDetails{}, err
	}

	fd := []FileDetails{}

	for _, uploadKey := range uploadKeys {

		for _, fileHeader := range r.MultipartForm.File[uploadKey] {
			file, err := fileHeader.Open()
			if err != nil {
				return []FileDetails{}, errors.Wrapf(err, "Unable to open file %s", fileHeader.Filename)
			}
			byteSlice, err := readFile(file)
			if err != nil {
				return []FileDetails{}, errors.Wrap(err, "Unable to read file")
			}
			var buf bytes.Buffer
			tee := io.TeeReader(bytes.NewReader(byteSlice), &buf)
			contentType := GetFileContentType(tee)

			r := FileDetails{
				UploadKey:   uploadKey,
				FileName:    fileHeader.Filename,
				Data:        byteSlice,
				ContentType: contentType,
			}
			fd = append(fd, r)
		}
	}
	return fd, nil
}

// ExtractBodyAndFileUploads extracts both text and file uploads in a request body
func ExtractBodyAndFileUploads(r *http.Request, fileUploadKeys []string, textKeys ...string) (FileWithBodyResult, error) {

	if err := r.ParseMultipartForm(maxUploadMemory); err != nil {
		return FileWithBodyResult{}, err
	}

	textResult := map[string][]string{}
	for _, textKey := range textKeys {
		v := r.MultipartForm.Value[textKey]
		textResult[textKey] = v
	}

	fd := []FileDetails{}

	for _, fileUploadKey := range fileUploadKeys {
		for _, fileHeader := range r.MultipartForm.File[fileUploadKey] {
			file, err := fileHeader.Open()
			if err != nil {
				return FileWithBodyResult{}, errors.Wrapf(err, "Unable to open file %s", fileHeader.Filename)
			}
			byteSlice, err := readFile(file)
			if err != nil {
				return FileWithBodyResult{}, errors.Wrap(err, "Unable to read file")
			}
			var buf bytes.Buffer
			tee := io.TeeReader(bytes.NewReader(byteSlice), &buf)
			contentType := GetFileContentType(tee)

			r := FileDetails{
				UploadKey:   fileUploadKey,
				FileName:    fileHeader.Filename,
				Data:        byteSlice,
				ContentType: contentType,
			}
			fd = append(fd, r)
		}
	}
	return FileWithBodyResult{Files: fd, Body: textResult}, nil
}

// ServeError is a generic error function that serves a custom error depending on params
func ServeError(err error, w http.ResponseWriter) {
	switch {
	case commonerror.IsBadRequestError(err):
		serveBadRequestError(err, w)
		return
	case commonerror.IsUnathourizedError(err):
		serveUnauthorizedResponse(err, w)
		return
	case commonerror.IsUnAuthenticatedError(err):
		serveAuthenticationErrResponse(err, w)
		return
	default:
		serveInternalError(err, w)
		return
	}
}

//ServeJSON returns a JSON response for an http request
func ServeJSON(res interface{}, w http.ResponseWriter, statusCode int) {
	bb, err := json.Marshal(res)
	if err != nil {
		serveInternalError(err, w)
		return
	}
	setStandardHeaders(w)
	w.WriteHeader(statusCode)
	w.Write(bb)
}

// ServeGeneralJSON serves the generic response as a json format
func ServeGeneralJSON(res interface{}, w http.ResponseWriter, statusCode int) {
	r := NewSuccessResponse(res)
	ServeJSON(r, w, statusCode)
}

// ServeNoContent serves no content with standard headers
func ServeNoContent(w http.ResponseWriter) {
	setStandardHeaders(w)
	w.WriteHeader(http.StatusNoContent)
}

func readFile(file multipart.File) ([]byte, error) {
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return []byte{}, errors.Wrap(err, "Unable to read file content")
	}
	return fileBytes, nil
}

// GetFileContentType returns the file content
func GetFileContentType(r io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	mimeType, _ := mimetype.DetectReader(buf)
	return mimeType.String()
}

// ServeFile returns a File Download Capability for an http request
func ServeFile(w http.ResponseWriter, fileName string, data []byte) {

	var buf bytes.Buffer
	tee := io.TeeReader(bytes.NewReader(data), &buf)
	contentType := GetFileContentType(tee)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%v"`, fileName))
	w.Write(data)
}

// RetrieveUUIDResource retrieves a resource of type uuid from an incoming request
func RetrieveUUIDResource(r *http.Request, idKey string) (uuid.V4, error) {
	id := httprouter.ParamsFromContext(r.Context()).ByName(idKey)
	if id == "" {
		return "", commonerror.NewErrorParams(idKey, fmt.Sprintf("%s not found in request", idKey)).ToBadRequest()
	}

	uID, err := uuid.GenFromString(id)
	if err != nil {
		return "", commonerror.NewErrorParams(idKey, "Invalid ID. Expected type UUID").ToBadRequest()
	}
	return uID, nil
}
