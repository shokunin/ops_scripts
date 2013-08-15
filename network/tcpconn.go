package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

var hostname string
var port int
var timeout int
var verbose bool

func init() {
	flag.StringVar(&hostname, "hostname", "localhost", "hostname or ip to scan")
	flag.IntVar(&port, "port", 22, "port to try to connect to")
	flag.IntVar(&timeout, "timeout", 22, "number of seconds to wait for connection")
	flag.BoolVar(&verbose, "verbose", false,  "port to try to connect to")
	flag.Parse()
}

func main() {
	if verbose {
		fmt.Printf("trying to connect to %s:%d\n", hostname, port)
	}
	c1 := make(chan bool)

	//##################################################################################
	go func() {
		for {
			conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", hostname, port))
			if err != nil {
				c1 <- false
			}
			if conn != nil {
				c1 <- true
			}

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
						if verbose {
							fmt.Println("connection sucessful")
						}
						os.Exit(0)
					} else {
						if verbose {
							fmt.Println("connection failed")
						}
						os.Exit(1)
					}
				}
		}
	}()
	var input bool
	fmt.Scanln(&input)
}
