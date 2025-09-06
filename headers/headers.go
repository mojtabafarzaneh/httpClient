package headers

import (
	"bytes"
	"fmt"
	"strings"
)

var rn = []byte("\r\n")

type Headers struct {
	headers map[string]string
}

func NewHeaders() *Headers {
	return &Headers{
		headers: map[string]string{},
	}
}

func ParseHeader(fieldLine []byte) (string, string, error) {
	parts := bytes.SplitN(fieldLine, []byte(":"), 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("malformed header field line")
	}
	name := parts[0]
	value := bytes.TrimSpace(parts[1])

	if bytes.HasSuffix(name, []byte(" ")) {
		return "", "", fmt.Errorf("malformed header field name")

	}

	return string(name), string(value), nil
}

func (h *Headers) Get(name string) string {
	return h.headers[strings.ToLower(name)]
}

func (h *Headers) Set(name, value string) {
	h.headers[strings.ToLower(name)] = value
}

func (h *Headers) Parse(data []byte) (int, bool, error) {
	read := 0
	done := false

	for {
		idx := bytes.Index(data[read:], rn)

		if idx == -1 {
			break
		}

		//Empty header
		if idx == 0 {
			done = true
			read += len(rn)
			break
		}
		name, value, err := ParseHeader(data[read : read+idx])
		if err != nil {
			return 0, false, err
		}

		read += idx + len(rn)
		h.Set(name, value)

	}

	return read, done, nil
}
