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
	params []string
	clientchan chan bool
}

func (req *Request) Read(read io.Reader) error {
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
	req.params = params
	return err
}

func (req *Request) DealProtocol() error {
	req.cmd = req.params[0]
	switch req.cmd {
		case SET, ADD, REPLACE:
			if len(req.params) != StoreLen {
				req.result = Invaild
				return errors.New("commond format of store is wrong")
			}
			req.key = req.params[1]
			req.value = []byte(req.params[2])
			req.interval, _ = strconv.ParseInt(req.params[3], 10, 64)
		case GET:
			if len(req.params) != GetLen {
				req.result = Invaild
				return errors.New("commond format of get is wrong")
			}
			req.key = req.params[1]
		case DELETE:
			if len(req.params) != DelNowLen && len(req.params) != DelDelayLen {
				req.result = Invaild
				return errors.New("commond format of delete is wrong")
			}
			req.key = req.params[1]
			if len(req.params) == DelDelayLen {
				req.delay, _ = strconv.ParseInt(req.params[2], 10, 64)
			}
		case FLUSH_ALL:
			if len(req.params) != FlushAllLen {
				req.result = Invaild
				return errors.New("commond format of flush_all is wrong")
			}
		default:
			return errors.New("cmd is invaild")
	}
	return nil
}

func (req *Request)OperateMap(datachan chan *Request) {
	datachan <- req
	<-req.clientchan
}
