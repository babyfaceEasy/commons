package commonerror

import (
	"encoding/json"
	"fmt"

	pkgErr "github.com/pkg/errors"
)

const (
	invalidReqBodyErrCode     = ErrorCode("invalid request body")
	invalidReqBodyErrMesssage = ErrorMessage("One or more inputs are invalid, please enter valid information")
	serverErrCode             = ErrorCode("Internal server error")
	serverErrMessage          = ErrorMessage("Something unplanned for has gone wrong")
	unauthErrMsg              = ErrorMessage("Authentication details were not provided")
	authenticationErrMsg      = ErrorMessage("Authentication credentials are invalid")
)

// ErrorParams represents the functionality of error parameters
type ErrorParams map[string]interface{}

// ErrorCode is a machine friendly description of an error
type ErrorCode string

// ErrorMessage is a human friendly description of an error
type ErrorMessage string

// Error is a representation of an api Error
type Error struct {
	Code    ErrorCode    `json:"code,omitempty"`
	Message ErrorMessage `json:"message"`
	Params  ErrorParams  `json:"params,omitempty"`
}

func (e Error) Error() string {
	msg, err := json.MarshalIndent(&e, "", "\t")
	if err != nil {
		return fmt.Sprint("Unmarshed error message: ", e.Code)
	}
	return string(msg)
}

// NewErrorParams initializes an error params value
func NewErrorParams(k, v string) ErrorParams {
	return map[string]interface{}{k: v}
}

// ToBadRequest converts error params to a Bad Request Error
func (p ErrorParams) ToBadRequest() Error {
	return BadRequestError(p)
}

// ToUnauthorizedRequest converts error params to an unauthorized request
func (p ErrorParams) ToUnauthorizedRequest() Error {
	return UnauthorizedError(p)
}

// ToUnAuthenticatedRequest converts error params to an authentication error request
func (p ErrorParams) ToUnAuthenticatedRequest() Error {
	return UnAuthenticatedError(p)
}

// ToServerError converts error params to an internal server error
func (p ErrorParams) ToServerError() Error {
	return ServerError(p)
}

// BadRequestError creates a new bad request error
func BadRequestError(p ErrorParams) Error {
	return Error{
		Code:    invalidReqBodyErrCode,
		Message: invalidReqBodyErrMesssage,
		Params:  p,
	}
}

// UnauthorizedError creates a unauthorized error
func UnauthorizedError(p ErrorParams) Error {
	return Error{Message: unauthErrMsg}
}

// UnAuthenticatedError creates an authentication error message
func UnAuthenticatedError(p ErrorParams) Error {
	return Error{Message: authenticationErrMsg}
}

// ServerError creates a new server down error
func ServerError(p ErrorParams) Error {
	return Error{
		Code:    serverErrCode,
		Message: serverErrMessage,
		Params:  p,
	}
}

// IsBadRequestError checks if an error is a bad request error
func IsBadRequestError(err error) bool {
	cause := pkgErr.Cause(err)
	customErr, ok := cause.(Error)
	return ok && customErr.Code == invalidReqBodyErrCode
}

// IsUnathourizedError checks if an error is an unauthorized error
func IsUnathourizedError(err error) bool {
	cause := pkgErr.Cause(err)
	customErr, ok := cause.(Error)
	return ok && customErr.Message == unauthErrMsg
}

// IsUnAuthenticatedError checks if an error is an authtentication error
func IsUnAuthenticatedError(err error) bool {
	cause := pkgErr.Cause(err)
	customErr, ok := cause.(Error)
	return ok && customErr.Message == authenticationErrMsg
}

// IsServerError checks if an error is an internal server error
func IsServerError(err error) bool {
	cause := pkgErr.Cause(err)
	customErr, ok := cause.(Error)
	return ok && customErr.Code == serverErrCode
}
