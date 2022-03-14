package simple_go_redis

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"
)

const (
	delimiter   = "\r\n"
	noKeyReturn = -1
)

type paramLen struct {
	len int
	val string
}

type redisProtocol struct {
	numOfParams int
	params      []paramLen
}

func (rdc *redisProtocol) String() string {
	s := strings.Builder{}
	s.WriteString("*")
	s.WriteString(strconv.Itoa(rdc.numOfParams))
	s.WriteString(delimiter)
	for _, v := range rdc.params {
		s.WriteString("$")
		s.WriteString(strconv.Itoa(v.len))
		s.WriteString(delimiter)
		s.WriteString(v.val)
		s.WriteString(delimiter)
	}

	return s.String()
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
		case '+':
			return nil, nil
		case '-':
			return nil, errors.New(string(r[1:]))
		case ':':
			l, err := strconv.Atoi(string(r[1:]))
			if err != nil {
				return nil, err
			}

			return l, nil
		case '$':
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
		case '*':
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
