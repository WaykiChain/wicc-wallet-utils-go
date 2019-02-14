package commons_test

import (
	"bytes"
	"fmt"
	"testing"
	"wicc-wallet-utils-go/commons"
)

func TestWriteCompactSize(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	writer := commons.NewWriterHelper(buf)

	writer.WriteCompactSize(256)

	fmt.Println("WriteCompactSize buf: ", buf.Bytes())
}

func TestWriteVarInt(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	writer := commons.NewWriterHelper(buf)

	writer.WriteVarInt(256)

	fmt.Println("WriteVarInt buf: ", buf.Bytes())
}
