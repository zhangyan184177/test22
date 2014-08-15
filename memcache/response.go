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

func (rsp *Response) WriteData(writer io.Writer) error {
	var err error
	var line int
	
	if rsp.result == Invaild {
		clierrs := []byte(CliErr + CRLF)
		for begin := 0; line < len(clierrs[begin:]); begin = line {
			line, err = writer.Write(clierrs)
			if err != nil {
				return err
			}
		}
	}
	
	switch rsp.cmd {
		case SET, ADD, REPLACE, DELETE, FLUSH_ALL:
			sets := []byte(rsp.result + CRLF)
			for begin := 0; line < len(sets[begin:]); begin = line {
				line, err = writer.Write(sets[begin:])
				if err != nil {
					return err
				}
			}
		case GET:
			s := ""
			if rsp.result != NotFound {
				s = fmt.Sprintf(" %s %s %d", rsp.key, string(rsp.value), len(rsp.value))
				s = Value + s + CRLF
			}
			gets := []byte(s)
			for begin := 0; line < len(gets[begin:]); begin = line {
				line, err = writer.Write(gets)
				if err != nil {
					return err
				}
			}
		default:
			errs := []byte(CliErr + CRLF)
			for begin := 0; line < len(errs[begin:]); begin = line {
				line, err = writer.Write(errs)
				if err != nil {
					return err
				}
			}
	}
	return err
}
