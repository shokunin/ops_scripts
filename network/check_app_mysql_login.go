package main

import (
	"flag"
	"github.com/kylelemons/go-gypsy/yaml"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"os"
	"io"
	"time"
)

var hostname string
var configfile string
var port int
var timeout int

func init() {
	flag.StringVar(&configfile, "configfile", "config.yaml", "(Simple) YAML file to read")
	flag.IntVar(&timeout, "timeout", 3, "number of seconds to wait for connection")
	flag.Parse()
}

func main() {


	config, err := yaml.ReadFile(configfile)

	c1 := make(chan bool)

	//##################################################################################
	go func() {
		for {
			db, err := sql.Open("mysql",
				"user:password@tcp(127.0.0.1:3306)/hello")
			if err != nil {
				io.WriteString(os.Stdout, (fmt.Sprintf("CRITICAL: %s\n", err)))
				os.Stdout.Sync()
				c1 <- false
			}

			defer db.Close

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
