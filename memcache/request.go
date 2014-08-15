package memcache

import (
	"strconv"
	"strings"
	"errors"
	"bufio"
	"io"
)

type Request struct {
	cmd string
	key string
	delay int64
	value []byte
	result string
	interval int64
	clientchan chan bool
}

func (req *Request) ReadData(read io.Reader) error {
	req.clientchan = make(chan bool)
	reader := bufio.NewReaderSize(read, NormalReaderSize)
	line, err := reader.ReadBytes('\n')
	if err != nil {
		return err
	}
	if len(line) < 2 || line[len(line)-1] != '\n' || line[len(line)-2] != '\r' {
		req.result = Invaild
		return errors.New("command format is invaild")
	}
	
	params := strings.Fields(string(line))
	if len(params) < 1 {
		req.result = Invaild
		return errors.New("lack of instruction")
	}
	req.cmd = params[0]
	switch req.cmd {
		case SET, ADD, REPLACE:
			if len(params) != StoreLen {
				req.result = Invaild
				return errors.New("commond format of store is wrong")
			}
			req.key = params[1]
			req.value = []byte(params[2])
			req.interval, _ = strconv.ParseInt(params[3], 10, 64)
		case GET:
			if len(params) != GetLen {
				req.result = Invaild
				return errors.New("commond format of get is wrong")
			}
			req.key = params[1]
		case DELETE:
			if len(params) != DelNowLen && len(params) != DelDelayLen {
				req.result = Invaild
				return errors.New("commond format of delete is wrong")
			}
			req.key = params[1]
			if len(params) == DelDelayLen {
				req.delay, _ = strconv.ParseInt(params[2], 10, 64)
			}
		case FLUSH_ALL:
			if len(params) != FlushAllLen {
				req.result = Invaild
				return errors.New("commond format of flush_all is wrong")
			}
		default:
			return errors.New("cmd is invaild")
	}
	return err
}
