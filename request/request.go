package request

import (
	"bytes"
	"fmt"
	"io"
)

type parserState string

const (
	Stateinit  parserState = "init"
	StateDone  parserState = "Done"
	StateError parserState = "error"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	State       parserState
}

var ErrorMalformedRequestLine = fmt.Errorf("malformed request-line")
var SEPERATOR = []byte("\r\n")
var ErrorUnsupportedHttpVertion = fmt.Errorf("unsupported http version")
var ErrorRequestInErrorState = fmt.Errorf("request in error state")

// func newRequest() *Request {
// 	return &Request{
// 		State: Stateinit,
// 	}
// }

func parseRequestLine(b []byte) (*RequestLine, int, error) {
	idx := bytes.Index(b, SEPERATOR)
	if idx == -1 {
		return nil, 0, nil
	}

	startLine := b[:idx]
	read := idx + len(SEPERATOR)

	parts := bytes.Split(startLine, []byte(" "))
	if len(parts) != 3 {
		return nil, 0, ErrorMalformedRequestLine
	}

	httpParts := bytes.Split(parts[2], []byte("/"))
	if len(httpParts) != 2 || string(httpParts[0]) != "HTTP" || string(httpParts[1]) != "1.1" {
		return nil, 0, ErrorMalformedRequestLine
	}

	rl := &RequestLine{
		Method:        string(parts[0]),
		RequestTarget: string(parts[1]),
		HttpVersion:   string(httpParts[1]),
	}

	return rl, read, nil

}

func (r *Request) parse(data []byte) (int, error) {
	read := 0
outer:
	for {
		switch r.State {
		case StateError:
			return 0, ErrorRequestInErrorState
		case Stateinit:
			rl, n, err := parseRequestLine(data[read:])
			if err != nil {
				r.State = StateError
				return 0, err
			}

			if n == 0 {
				break outer
			}

			r.RequestLine = *rl
			read += n
			r.State = StateDone
			break outer

		case StateDone:
			break outer
		}
	}
	return read, nil

}

func (r *Request) error() bool {
	return r.State == StateError
}

func (r *Request) done() bool {
	return r.State == StateDone
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := &Request{State: Stateinit}
	//NOTE: Buffer Could Get OVerrun
	buf := make([]byte, 1024)
	bufLen := 0
	for !request.done() && !request.error() {
		n, err := reader.Read(buf[bufLen:])
		if err != nil {
			return nil, err
		}

		bufLen += n
		readN, err := request.parse(buf[:bufLen+n])
		if err != nil {
			return nil, err
		}
		copy(buf, buf[readN:bufLen])
		bufLen -= readN
	}

	return request, nil

}
