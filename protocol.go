package simple_go_redis

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"strconv"
)

const (
	CRLF        = "\r\n"
	noKeyReturn = -1

	flagStart        = '*'
	flagParamLen     = '$'
	flagSimpleReturn = '+'
	flagErrorReturn  = '-'
	flagIntReturn    = ':'
)

type paramLen struct {
	len int
	val string
}

type redisProtocol struct {
	numOfParams int
	params      []paramLen
}

func (rdc *redisProtocol) Bytes() []byte {
	s := bytes.Buffer{}
	s.WriteByte(flagStart)
	s.WriteString(strconv.Itoa(rdc.numOfParams))
	s.WriteString(CRLF)
	for _, v := range rdc.params {
		s.WriteByte(flagParamLen)
		s.WriteString(strconv.Itoa(v.len))
		s.WriteString(CRLF)
		s.WriteString(v.val)
		s.WriteString(CRLF)
	}

	return s.Bytes()
}

// parseResponse 解析输出
func parseResponse(reader io.Reader) (interface{}, error) {
	s := bufio.NewReader(reader)
	for {
		r, _, err := s.ReadLine()
		if err != nil {
			return nil, err
		}

		switch r[0] {
		case flagSimpleReturn:
			return nil, nil
		case flagErrorReturn:
			return nil, errors.New(string(r[1:]))
		case flagIntReturn:
			l, err := strconv.Atoi(string(r[1:]))
			if err != nil {
				return nil, err
			}

			return l, nil
		case flagParamLen:
			l, err := strconv.Atoi(string(r[1:]))
			if err != nil {
				return nil, err
			}

			if l == noKeyReturn {
				return nil, nil
			}

			r, _, err = s.ReadLine()
			if err != nil {
				return nil, err
			}

			return r[:l], nil
		case flagStart:
			l, err := strconv.Atoi(string(r[1:]))
			if err != nil {
				return nil, err
			}

			if l == noKeyReturn {
				return nil, nil
			}

			rs := make([]interface{}, l)

			for i := 0; i < l; i++ {
				r, _, err = s.ReadLine()
				pl, err := strconv.Atoi(string(r[1:]))
				if err != nil {
					return nil, err
				}

				if pl == noKeyReturn {
					rs[i] = nil
					continue
				}

				r, _, err = s.ReadLine()
				if err != nil {
					return nil, err
				}

				rs[i] = r[:pl]
			}

			return rs, nil
		}
	}
}
