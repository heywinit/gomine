package types

import "errors"

var (
	ErrNotSlice           = errors.New("gomine: struct field is not a slice")
	ErrMissingLen         = errors.New("gomine: missing len struct tag where absolutely necessary")
	ErrIncorrectFieldType = errors.New(
		"gomine: the target field type does not correspond to the one specified in the type tag",
	)
)
