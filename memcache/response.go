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
}

func (rsp *Response) Write(writer io.Writer) error {
	if rsp.result == Invaild {
		clierrs := []byte(CliErr + CRLF)
		err := doWrite(writer, clierrs)
		if err != nil {
			return err
		}
	}

	switch rsp.cmd {
		case SET, ADD, REPLACE, DELETE, FLUSH_ALL:
			sets := []byte(rsp.result + CRLF)
			err := doWrite(writer, sets)
			if err != nil {
				return err
			}
		case GET:
			s := ""
			if rsp.result != NotFound {
				s = fmt.Sprintf(" %s %s %d", rsp.key, string(rsp.value), len(rsp.value))
				s = Value + s + CRLF
			}
			gets := []byte(s)
			err := doWrite(writer, gets)
			if err != nil {
				return err
			}
		default:
			errs := []byte(CliErr + CRLF)
			err := doWrite(writer, errs)
			if err != nil {
				return err
			}
	}
	return nil
}

func doWrite(writer io.Writer, s []byte) error {
	var err error
	var line int
	for begin := 0; line < len(s[begin:]); begin = line {
		line, err = writer.Write(s[begin:])
		if err != nil {
			return err
		}
	}
	return err
}
