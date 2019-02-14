package commons

import (
	"bytes"
	"encoding/binary"
	"regexp"
)

const (
	UINT16_MAX = ^uint16(0)
	UINT32_MAX = ^uint32(0)
	UINT_MAX   = ^uint(0)
	UINT64_MAX = ^uint64(0)
)

type WriterHelper struct {
	buf *bytes.Buffer
}

func NewWriterHelper(buf *bytes.Buffer) *WriterHelper {
	return &WriterHelper{buf}
}

func (writer WriterHelper) Write(p []byte) {
	writer.buf.Write(p)
}

func (writer WriterHelper) WriteByte(value byte) {
	writer.buf.WriteByte(value)
}

func (writer WriterHelper) WriteIntData(data interface{}) {
	binary.Write(writer.buf, binary.LittleEndian, data)
}

func (writer WriterHelper) WriteCompactSize(value uint64) {
	if value < 253 {
		writer.WriteByte(byte(value))
	} else if value <= uint64(UINT16_MAX) {
		writer.WriteByte(253) // size tag, byte
		writer.WriteIntData(uint16(value))
	} else if value <= uint64(UINT32_MAX) {
		writer.WriteByte(254) // size tag, byte
		writer.WriteIntData(uint32(value))
	} else {
		writer.WriteByte(255) // size tag, byte
		writer.WriteIntData(uint64(value))
	}
	return
}

func (writer WriterHelper) WriteBytes(value []byte) {
	len := uint64(len(value))
	writer.WriteCompactSize(len)
	if len > 0 {
		writer.Write(value)
	}
}

func (writer WriterHelper) WriteVarInt(value int64) {
	size := sizeofVarInt(value)
	tmp := make([]byte, ((size*8 + 6) / 7))
	len := 0
	n := value
	// h:=byte(0x00)
	for {
		h := byte(0)
		if len == 0 {
			h = 0x00
		} else {
			h = 0x80
		}

		tmp[len] = byte((byte(n) & 0x7f) | h)

		if n <= 0x7f {
			break
		}

		n = (n >> 7) - 1
		len += 1
	}

	for {
		writer.WriteByte(tmp[len])
		len -= 1
		if len < 0 {
			break
		}
	}
	return
}

func sizeofVarInt(value int64) int {
	ret := 0
	n := value

	for {
		ret += 1
		if n <= 0x7f {
			break
		}
		n = (n >> 7) - 1
	}

	return ret
}

/*
func parseRegId(value string) []int64 {
	regidStr := strings.Split(value, "-")
	regHeight, _ := strconv.ParseInt(regidStr[0], 10, 64)
	regIndex, _ := strconv.ParseInt(regidStr[1], 10, 64)
	return []int64{regHeight, regIndex}
}
*/
func isRegId(value string) bool {
	re := regexp.MustCompile(`^\s*(\d+)\-(\d+)\s*$`)
	return re.MatchString(value)
}

/*
func (writer WriterHelper) WriteRegId(value string) {
	regIdData := parseRegId(value)

}
*/

func (writer WriterHelper) WriteString(value string) {
	len := uint64(len(value))
	writer.WriteCompactSize(len)
	if len > 0 {
		writer.buf.WriteString(value)
	}
}

func (writer WriterHelper) WriteRegId(value RegId) {
	buf := bytes.NewBuffer([]byte{})
	idWriter := NewWriterHelper(buf)
	idWriter.WriteVarInt(int64(value.Height))
	idWriter.WriteVarInt(int64(value.Index))
	writer.WriteBytes(buf.Bytes())
}

func (writer WriterHelper) WriteContractScript(script []byte, description string) {
	buf := bytes.NewBuffer([]byte{})
	scriptWriter := NewWriterHelper(buf)
	scriptWriter.WriteBytes(script)
	scriptWriter.WriteString(description)
	writer.WriteBytes(buf.Bytes())
}
