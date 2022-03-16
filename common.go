package simple_go_redis

import (
	"net"
	"strconv"
)

type redisConn struct {
	address string
	conn    net.Conn
}

func New(address string) (*redisConn, error) {
	r := &redisConn{address: address}
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	r.conn = conn
	return r, nil
}

func (r *redisConn) Do(params ...string) (interface{}, error) {
	res := new(redisProtocol)
	res.numOfParams = len(params)
	for _, v := range params {
		res.params = append(res.params, paramLen{
			len: len(v),
			val: v,
		})
	}

	_, err := r.write(res)
	if err != nil {
		return nil, err
	}

	return parseResponse(r.conn)
}

func (r *redisConn) Select(db int) (interface{}, error) {
	return r.Do("select", strconv.Itoa(db))
}

func (r *redisConn) String() string {
	return r.address
}

func (r *redisConn) write(p *redisProtocol) (int, error) {
	return r.conn.Write(p.Bytes())
}

func (r *redisConn) Close() error {
	return r.conn.Close()
}
