package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
)

var (
	c     redis.Conn
	err   error
	reply interface{}
)

func init() {
	c, err = redis.Dial("tcp", "20.36.28.197:10003")
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	defer c.Close()
	for i := 0; i < 100; i++ {
		c.Send("GET", fmt.Sprintf("PIPELINE-%d", i))
	}
	c.Flush()
	for i := 0; i < 100; i++ {
		v, err := c.Receive()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(v)
	}
}
