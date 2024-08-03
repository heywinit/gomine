package types

import (
	"bytes"
	"reflect"

	"github.com/Tnze/go-mc/nbt"
)

func SerializeNbt(field reflect.Value, databuf *bytes.Buffer) error {
	encoder := nbt.NewEncoder(databuf)

	return encoder.Encode(field.Interface(), "")
}

func DeserializeNbt(field reflect.Value, databuf *bytes.Buffer) error {
	decoder := nbt.NewDecoder(databuf)
	_, err := decoder.Decode(field.Addr().Interface())
	return err
}
