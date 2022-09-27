package uuid

// V4 is a type representing a v4 UUID string
type V4 string

// GenV4Func generates a v4 UUID string
type GenV4Func func() V4
