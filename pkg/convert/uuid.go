package convert

import (
	"encoding/hex"
	"strings"
)

func TrimUUID(uuid string) string {
	return strings.ReplaceAll(uuid, "-", "")
}

func UUIDToBytes(uuid string) []byte {
	uuid = TrimUUID(uuid)
	bin := make([]byte, 16)
	_, err := hex.Decode(bin, []byte(uuid))
	if err != nil {
		return nil
	}
	return bin
}

func BytesToUUID(bin []byte) string {
	return hex.EncodeToString(bin)
}
