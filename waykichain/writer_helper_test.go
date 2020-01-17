package waykichain

import (
	"bytes"
	"fmt"
	"testing"
)

func TestWriteCompactSize(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	writer := NewWriterHelper(buf)

	writer.WriteCompactSize(256)

	fmt.Println("WriteCompactSize buf: ", buf.Bytes())
}

func TestWriteVarInt(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	writer := NewWriterHelper(buf)

	writer.WriteVarInt(256)

	fmt.Println("WriteVarInt buf: ", buf.Bytes())
}
