package main

// MAKE SURE TO SET export GOMAXPROCS=16  in your environment before running
// example config file
// { "redis1_0_6379": {"hostname":"redis1","port":6379,"database":0}, "redis1_1_63792": {"hostname":"redis1","port":6379,"database":1} }

import "fmt"
import "flag"
import "io/ioutil"
import "encoding/json"
import "os"
//import "reflect"

// import "github.com/fzzy/radix/redis"

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

func main() {

	boxen := new(RedisServers)
	cfg, err := ioutil.ReadFile(config)
	check(err)
	cfgerr := boxen.FromJson(string(cfg))
	check(cfgerr)
	for k, v := range boxen.Pool {
		fmt.Println(k)
		fmt.Println(v.Hostname)
		fmt.Println(v.Port)
		fmt.Println(v.Database)
	}
}
