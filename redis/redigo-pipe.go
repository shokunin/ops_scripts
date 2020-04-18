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
	c, err = redis.Dial("tcp", ":6379")
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	defer c.Close()
	c.Send("GET", "beer")
	c.Send("GET", "wine")
	c.Send("GET", "bourbon")
	c.Send("GET", "gin")
	c.Flush()
	for i := 0; i < 4; i++ {
		v, err := c.Receive()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(v)
	}
}
