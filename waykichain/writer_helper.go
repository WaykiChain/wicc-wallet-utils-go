package waykichain

import (
	"bytes"
	"encoding/binary"
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

func (writer WriterHelper) GetBuf() *bytes.Buffer {
	return writer.buf
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
		len++
	}

	for {
		writer.WriteByte(tmp[len])
		len--
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
		ret++
		if n <= 0x7f {
			break
		}
		n = (n >> 7) - 1
	}

	return ret
}

func (writer WriterHelper) WriteRegId(value RegId) {
	buf := bytes.NewBuffer([]byte{})
	idWriter := NewWriterHelper(buf)
	idWriter.WriteVarInt(int64(value.Height))
	idWriter.WriteVarInt(int64(value.Index))
	writer.WriteBytes(buf.Bytes())
}

func (writer WriterHelper) WritePubKeyId(value PubKeyId) {
	writer.WriteBytes(value)
}

func (writer WriterHelper) WriteReverse(value []byte) {
	length := uint64(len(value))
	buf := bytes.NewBuffer([]byte{})
	//writer.WriteVarInt(int64(len(value)))
	for i := 0; i < int(length); i++ {
		a := len(value)
		buf.WriteByte(value[int(a)-i-1]);
	}

	if length > 0 {
		writer.Write(buf.Bytes())
	}
}

func (write WriterHelper) WriteCdpAsset(assets []AssetModel) {
     len:= len(assets)
     write.WriteVarInt(int64(len))
	for i:=0; i<len;i++  {
		write.WriteString(assets[i].AssetSymbol)
		write.WriteVarInt(assets[i].AssetAmount)
	}
}

func (writer WriterHelper) WriteAddressId(value AddressId) {
	writer.WriteBytes(value)
}

func (writer WriterHelper) WriteUserId(value *UserIdWraper) {

	if value != nil {
		switch value.GetType() {
		case UID_REG:
			writer.WriteRegId(value.GetId().(RegId))
		case UID_PUB_KEY:
			writer.WritePubKeyId(value.GetId().(PubKeyId))
		case UID_ADDRESS:
			writer.WriteAddressId(value.GetId().(AddressId))
		}
	} else {
		writer.WriteVarInt(0) // write empty bytes
	}
}

func (writer WriterHelper) WriteString(value string) {
	len := uint64(len(value))
	writer.WriteCompactSize(len)
	if len > 0 {
		writer.buf.WriteString(value)
	}
}

func (writer WriterHelper) WriteContractScript(script []byte, description string) {
	buf := bytes.NewBuffer([]byte{})
	scriptWriter := NewWriterHelper(buf)
	scriptWriter.WriteBytes(script)
	scriptWriter.WriteString(description)
	writer.WriteBytes(buf.Bytes())
}
