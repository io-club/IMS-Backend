package util

import (
	"encoding/base64"
	"encoding/binary"
	"strings"
)

func EncodeUserID(id uint64) string {
	buf := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(buf, id)
	urlBase := base64.StdEncoding.EncodeToString(buf)
	urlBase = strings.ReplaceAll(urlBase, "+", "-")
	urlBase = strings.ReplaceAll(urlBase, "/", "_")
	return urlBase
}
