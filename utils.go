package goefibootmgr

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

func bootnumToHexString(bootnum uint16) string {
	return fmt.Sprintf("%04X", bootnum)
}

func hexStringToBootNum(s string) uint16 {
	b, _ := hex.DecodeString(s)
	return binary.BigEndian.Uint16(b)
}
