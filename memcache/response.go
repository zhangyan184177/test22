package memcache

import (
	"fmt"
	"io"
)

type Response struct {
	cmd string
	key string
	value []byte
	result string
	params []byte
}

func (rsp *Response) DealProtocol() {
	var s []byte
	if rsp.result == Invaild {
		s = []byte(CliErr + CRLF)
	}

	switch rsp.cmd {
		case SET, ADD, REPLACE, DELETE, FLUSH_ALL:
			s = []byte(rsp.result + CRLF)
		case GET:
			gets := ""
			if rsp.result != NotFound {
				gets = fmt.Sprintf(" %s %s %d", rsp.key, string(rsp.value), len(rsp.value))
				gets = Value + gets + CRLF
			}
			s = []byte(gets)
		default:
			s = []byte(CliErr + CRLF)
	}
	rsp.params = s
}

func (rsp *Response)Write(writer io.Writer) error {
	var err error
	var line int
	for begin := 0; line < len(rsp.params[begin:]); begin = line {
		line, err = writer.Write(rsp.params[begin:])
		if err != nil {
			return err
		}
	}
	return err
}
