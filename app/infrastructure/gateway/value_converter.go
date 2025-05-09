package gateway

import (
	"io"
	"strings"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func ConvertFromEucJPToUtf8(eucjpStr string) string {
	reader := transform.NewReader(strings.NewReader(eucjpStr), japanese.EUCJP.NewDecoder())
	bytes, err := io.ReadAll(reader)
	if err != nil {
		panic(err)
	}

	str := string(bytes)
	str = strings.TrimSpace(str)
	str = strings.TrimRight(str, "\n")

	return str
}

func Trim(str string) string {
	str = strings.TrimSpace(str)
	str = strings.TrimRight(str, "\n")

	return str
}
