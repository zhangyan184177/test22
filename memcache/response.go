package memcache

import (
	"fmt"
	"io"
)

type Response struct {
	cmd string
	keys string
	value string
	length int
	result string
}

func (rsp *Response) WriteData(writer io.Writer) error {
	var err error
//	log.Println("write data rsp: ", rsp)
	switch rsp.cmd {
		case SET, ADD, REPLACE, DELETE, FLUSH_ALL:	
			s := rsp.result + CRLF
			_, err := writer.Write([]byte(s))
			if err != nil {
				return err	
			}
		case GET:
			gets := ""	
			if rsp.result != NotFound {
				gets = fmt.Sprintf(" %s %s %d", rsp.keys, rsp.value, len(rsp.value))	
				gets = Value + gets + CRLF	
			}
			_, err := writer.Write([]byte(gets))
			if err != nil {
				return err	
			}
		case INVAILD:
			errs := CommonErr + CRLF
			_, err := writer.Write([]byte(errs))
			if err != nil {
				return err	
			}
			
	}
	return err	
}
