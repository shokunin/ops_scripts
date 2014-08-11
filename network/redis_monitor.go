package main

// MAKE SURE TO SET export GOMAXPROCS=16  in your environment before running
// example config file
// { "redis1_0_6379": {"hostname":"redis1","port":6379,"database":0}, "redis1_1_63792": {"hostname":"redis1","port":6379,"database":1} }

import "fmt"
import "flag"
import "time"
import "os"
import "io/ioutil"
import "encoding/json"
import "github.com/fzzy/radix/redis"
//import "reflect"

func check(err error) {
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}

var config string

func init() {
	flag.StringVar(&config, "config", "/tmp/config.json", "JSON formatted config file")
	flag.Parse()
}

// ********************   Configuration information
type Redisbox struct {
	Hostname string `json:"hostname"`
	Port     int32  `json:"port"`
	Database int32  `json:"database"`
}

type RedisServers struct {
	Pool map[string]Redisbox
}

func (rc *RedisServers) FromJson(jsonStr string) error {
	var data = &rc.Pool
	b := []byte(jsonStr)
	return json.Unmarshal(b, data)
}

// ********************   Worker Pool information
func worker(id int, jobs <-chan Redisbox, results chan<- string) {
	for j := range jobs {
	  c, err := redis.DialTimeout("tcp", fmt.Sprintf("%s:%d", j.Hostname, j.Port), time.Duration(3)*time.Second)
		check(err)
		c.Cmd("PING")
		results <- "OK"
	  c.Close()
	}
}


func main() {

	boxen := new(RedisServers)
	cfg, err := ioutil.ReadFile(config)
	check(err)
	cfgerr := boxen.FromJson(string(cfg))
	check(cfgerr)
	jobs := make(chan Redisbox, len(boxen.Pool))
	results := make(chan string, len(boxen.Pool))

    for w := 0; w <= 8; w++ {
        go worker(w, jobs, results)
    }

	for k, v := range boxen.Pool {
		fmt.Println(k)
		jobs <- v
	}
    // Finally we collect all the results of the work.
    for a := 0; a <= len(boxen.Pool)-1; a++ {
        <-results
    }
	os.Exit(0)
}
