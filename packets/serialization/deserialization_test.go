package serialization_test

import (
	"bytes"
	"errors"
	"reflect"
	"testing"

	"github.com/heywinit/gomine/packets/serialization"
	"github.com/heywinit/gomine/packets/serialization/tagutils"
)

func TestArrDeserialization(t *testing.T) {
	type TestChildStruct struct {
		VarInt1 int32  `mc:"varint"`
		Text    string `mc:"string"`
	}

	type TestStruct struct {
		Nested []TestChildStruct `mc:"array" len:"2"`
	}

	testStruct := new(TestStruct)
	testBuffer := new(bytes.Buffer)

	testBuffer.Write([]byte{0x01, 0x02, 0x59, 0x68, 0x02, 0x00})

	err := serialization.DeserializeFields(reflect.ValueOf(testStruct).Elem(), testBuffer)
	if err != nil {
		t.Error(err)
	}
}

func TestDeserialization(t *testing.T) {
	type NbtStruct struct {
		String1 string `nbt:"stringone"`
		String2 string `nbt:"stringtwo"`
	}

	type TestStruct struct {
		VarInt  int32       `mc:"varint"`
		VarLong int64       `mc:"varlong"`
		String  string      `mc:"string"`
		Inherit uint32      `mc:"inherit"`
		Ignore  interface{} `mc:"ignore" len:"6"`
		Bytes   []byte      `mc:"bytes" len:"3"`
		Nbt     NbtStruct   `mc:"nbt"`
	}

	testBuffer := new(bytes.Buffer)
	testStruct := new(TestStruct)

	testBuffer.Write([]byte{
		0x80, 0x01, // varint
		0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x01, // varlong
		0x04, 0x59, 0x59, 0x59, 0x59, // string
		0xff, 0x00, 0xff, 0x00, // inherit
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // ignored
		0x01, 0x02, 0x03, // byte
	})
	//
	//err, _ := nbt.Marshal(testBuffer, NbtStruct{String1: "Hello", String2: "World"})
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//for i := 0; i < 3; i++ {
	//	err, _ := nbt.Marshal(testBuffer, NbtStruct{String1: "ArrayTest", String2: "ArrayTest2"})
	//	if err != nil {
	//		t.Fatal(err)
	//	}
	//}
	//
	//err = serialization.DeserializeFields(reflect.ValueOf(testStruct).Elem(), testBuffer)
	//if err != nil {
	//	t.Fatal(err)
	//}

	// VarInt
	if testStruct.VarInt != 128 {
		t.Fatal("deserialized varint mismatch")
	}

	if testStruct.VarLong != -9223372036854775808 {
		t.Fatal("deserialized varlong mismatch")
	}

	// String
	if testStruct.String != "YYYY" {
		t.Fatal("deserialized string mismatch")
	}

	// Inherit uint32
	if testStruct.Inherit != 4278255360 {
		t.Fatal("deserialized inherit (uint32) mismatch")
	}

	// Bytes
	if !bytes.Equal(testStruct.Bytes, []byte{0x01, 0x02, 0x03}) {
		t.Fatal("deserialized bytes mismatch")
	}

	if testStruct.Nbt.String1 != "Hello" || testStruct.Nbt.String2 != "World" {
		t.Fatal("deserialized nbt mismatch")
	}
}

func TestInvalidLength(t *testing.T) {
	type StructWithLength struct {
		Bytes []byte `mc:"bytes" len:"invalid"`
	}

	testBuffer := new(bytes.Buffer)
	testStruct := new(StructWithLength)

	err := serialization.DeserializeFields(reflect.ValueOf(testStruct).Elem(), testBuffer)
	if !errors.Is(err, tagutils.ErrInvalidLen) {
		t.Fatal(err)
	}
}
