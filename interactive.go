package simple_go_redis

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
)

// Interactive 交互式运行
func (r *redisConn) Interactive() {
	fmt.Print("\nThis is a simple redis client.\n\n")

	defer r.Close()

	tip := func() {
		fmt.Printf("%s> ", r)
	}

	for {
		tip()
		f := bufio.NewReader(os.Stdin)
		n, _, err := f.ReadLine()
		s := bytes.TrimSpace(n)
		if len(s) == 0 {
			continue
		}

		if err != nil {
			fmt.Printf("error: %v\n", err)
			continue
		}

		c := parseInput(n)

		_, err = r.Write([]byte(c.String()))
		if err != nil {
			fmt.Println("request error:", err)
			continue
		}

		res, err := parseResponse(r)
		if err != nil {
			fmt.Println("response error:", err)
			continue
		}

		formatShowResponse(c, res)
	}
}

// parseInput 解析用户输入
func parseInput(input []byte) *redisProtocol {
	ns := bytes.Split(input, []byte{' '})
	res := new(redisProtocol)
	for _, v := range ns {
		if len(v) == 0 {
			continue
		}
		res.numOfParams++
		res.params = append(res.params, paramLen{
			len: len(v),
			val: string(v),
		})
	}

	return res
}

// formatShowResponse 结构化输出
func formatShowResponse(c *redisProtocol, i interface{}) {
	switch i.(type) {
	case string:
	case []byte:
		// 内部命令不添加引号了
		if c.params[0].val == "info" {
			fmt.Printf("%s\n", i)
		} else {
			fmt.Printf("\"%s\"\n", i)
		}
	case int:
		fmt.Printf("(integer) %d\n", i)
	case []interface{}:
		for k, v := range i.([]interface{}) {
			fmt.Printf("%d) \"%s\"\n", k+1, v)
		}
	}
}
