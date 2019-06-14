package hash

import (
	"fmt"
	"testing"
)

func TestHash256(t *testing.T) {
	fmt.Printf("%x\n", Hash256([]byte{0, 0}))
}
