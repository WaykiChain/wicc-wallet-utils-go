// Package bip32 provides ...
package commons

import (
	"errors"
	"strconv"
	"strings"
)

var ErrKeyPathFormat = errors.New("Wallet Path Error")
var ErrParentKey = errors.New("Key must master")

var mapKey = make(map[string]*ExtendedKey)

// DerivePath return key by path : m/0'/1/2' etc...
func (key *ExtendedKey) DerivePath(pathStr string) (*ExtendedKey, error) {
	//fmt.Println("###########", len(mapKey))
	if key.childNum > 0 {
		return nil, ErrParentKey
	}
	keyTmp := mapKey[pathStr]
	path := strings.Split(pathStr, "/")
	err := vaildPath(path)
	if err != nil {
		return nil, err
	}
	tmpPath := []string{}
	var tmpPathStr string
	var tmpParentKey *ExtendedKey
	for _, childNumStr := range path {
		tmpPath = append(tmpPath, childNumStr)
		tmpPathStr = strings.Join(tmpPath, "/")
		keyTmp = mapKey[tmpPathStr]
		if tmpPathStr == "m" {
			keyTmp = key
		} else {
			isHardenedChild := false
			if strings.HasSuffix(childNumStr, "'") {
				childNumStr = strings.Replace(childNumStr, "'", "", -1)
				isHardenedChild = true
			}
			childNum, _ := strconv.Atoi(childNumStr)
			var err error
			if isHardenedChild {
				keyTmp, err = tmpParentKey.HardenedChild(uint32(childNum))
			} else {
				keyTmp, err = tmpParentKey.Child(uint32(childNum))
			}
			if err != nil {
				return nil, err
			}
		}
		mapKey[tmpPathStr] = keyTmp
		tmpParentKey = keyTmp
	}
	return keyTmp, nil
}

func vaildPath(path []string) error {
	if path[0] != "m" {
		return ErrKeyPathFormat
	}
	for i := 1; i < len(path); i++ {
		childNumStr := path[i]
		if strings.HasSuffix(childNumStr, "'") {
			childNumStr = strings.Replace(childNumStr, "'", "", -1)
		}
		childNum, err := strconv.Atoi(childNumStr)
		if err != nil {
			return ErrKeyPathFormat
		}
		if uint32(childNum) >= HardenedKeyStart || childNum < 0 {
			return ErrKeyPathFormat
		}
	}
	return nil
}
