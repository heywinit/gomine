package types

import (
	"bytes"
	"errors"
)

var (
	ErrIgnoreLenUnknown = errors.New("mcproto: ignore type specified but no len tag present")
)

func SerializeIgnore(ignoreLen int, databuf *bytes.Buffer) error {
	if ignoreLen < 0 {
		return ErrIgnoreLenUnknown
	}

	ignoreBuf := make([]byte, ignoreLen)

	_, err := databuf.Write(ignoreBuf)
	return err
}

func DeserializeIgnore(ignoreLen int, databuf *bytes.Buffer) error {
	if ignoreLen < 0 {
		return ErrIgnoreLenUnknown
	}

	ignoreBuf := make([]byte, ignoreLen)

	_, err := databuf.Read(ignoreBuf)
	return err
}
