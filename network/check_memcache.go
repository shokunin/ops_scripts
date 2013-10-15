package main

import (
	"flag"
	"github.com/bradfitz/gomemcache/memcache"
	"fmt"
	"os"
	"io"
	"time"
)

var hostname string
var port int
var timeout int

func init() {
	flag.StringVar(&hostname, "hostname", "localhost", "hostname or ip of memcache server")
	flag.IntVar(&port, "port", 11211, "port to try to connect to")
	flag.IntVar(&timeout, "timeout", 3, "number of seconds to wait for connection")
	flag.Parse()
}

func main() {

	c1 := make(chan bool)

	//##################################################################################
	go func() {
		for {
			mc := memcache.New(fmt.Sprintf("%s:%d", hostname, port))
			mc.Set(&memcache.Item{Key: "dtm-not-real", Value: []byte("setByNagios")})
			conn, err := mc.Get("dtm-not-real")
			if err != nil {
				io.WriteString(os.Stdout, (fmt.Sprintf("CRITICAL: %s\n", err)))
				os.Stdout.Sync()
				c1 <- false
			}
			if conn != nil {
				c1 <- true
			}

			mc.Delete("dtm-not-real")

		}
	}()
	//##################################################################################
	go func() {
		time.Sleep(time.Duration(timeout) * time.Second)
		c1 <- false
	}()
	//##################################################################################
	go func () {
		for {
			select {
				case msg1 := <- c1:
					if msg1 {
						io.WriteString(os.Stdout, "OK: able to create and fetch a record\n")
						os.Stdout.Sync()
						os.Exit(0)
					} else {
						os.Exit(2)
					}
				}
		}
	}()
	var input bool
	fmt.Scanln(&input)
}
