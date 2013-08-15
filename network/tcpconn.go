package main

import (
	"flag"
	"fmt"
	"net"
	"os"
)

var hostname string
var port int
var verbose bool
func init() {
	flag.StringVar(&hostname, "hostname", "localhost", "hostname or ip to scan")
	flag.IntVar(&port, "port", 22, "port to try to connect to")
	flag.BoolVar(&verbose, "verbose", false,  "port to try to connect to")
	flag.Parse()
}

func main () {
	if verbose {
		fmt.Printf("trying to connect to %s:%d\n", hostname, port)
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", hostname, port))
	if err != nil {
		if verbose {
			fmt.Println("connection failed")
		}
		os.Exit(1)
	}
	if conn != nil {
		if verbose {
			fmt.Println("connection worked")
		}
		os.Exit(0)
	}
	
}

