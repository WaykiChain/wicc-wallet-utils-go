package common

import (
	"github.com/blocktree/openwallet/common"
	"strings"
)

// []byte 倒序
func Reverse(src []byte) []byte{
	len := len(src)
	dst := make([]byte,len)
	for _ ,a := range src{
		dst[len-1] = a
		len--
	}
	return dst
}

func RemoveOxFromHex(value string) string {
	result := value
	if strings.Index(value, "0x") != -1 {
		result = common.Substr(value, 2, len(value))
	}
	return result
}