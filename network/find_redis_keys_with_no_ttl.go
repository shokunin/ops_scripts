package main

// MAKE SURE TO SET export GOMAXPROCS=16  in your environment before running
// This iterates through the redis keys and deletes with an lru_idle longer than x seconds

import "fmt"
import "flag"
import "time"
import "os"

import "github.com/fzzy/radix/redis"

func fetchAllKeys(hostame string, port int, database int) []string {
	c, err := redis.DialTimeout("tcp", fmt.Sprintf("%s:%d", hostname, port), time.Duration(300)*time.Second)
	errHndlr(err)
	keys := c.Cmd("SELECT", database)
	keys = c.Cmd("KEYS", "*")
	//fmt.Println("type:", reflect.TypeOf(keys))
	j := keys.Elems
	redis_keys := make([]string, len(j), len(j))
	for i := 0; i < len(j); i++ {
		redis_keys[i] = fmt.Sprintf("%s", j[i])
	}
	return redis_keys
}

func errHndlr(err error) {
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}

func worker(id int, jobs <-chan string, results chan<- string, hostame string, port int, database int) {
	c, err := redis.DialTimeout("tcp", fmt.Sprintf("%s:%d", hostname, port), time.Duration(3)*time.Second)
	errHndlr(err)
	for j := range jobs {
		errHndlr(err)
		k := c.Cmd("SELECT", database)
		k = c.Cmd("TTL", j)
		fmt.Println(j, ",", k)
		results <- "OK"
	}
	c.Close()
}

var hostname string
var port int
var concurrent int
var database int

func init() {
	flag.StringVar(&hostname, "hostname", "localhost", "hostname or ip to scan")
	flag.IntVar(&port, "port", 6379, "port to try to connect to")
	flag.IntVar(&concurrent, "concurrent", 10, "number of workers to run")
	flag.IntVar(&database, "database", 0, "Redis database to use.  DB 0 is the default")
	flag.Parse()
}

func main() {

	keys := fetchAllKeys(hostname, port, database)
	// In order to use our pool of workers we need to send
	// them work and collect their results. We make 2
	// channels for this.
	jobs := make(chan string, len(keys))
	results := make(chan string, len(keys))

	for w := 0; w <= concurrent; w++ {
		go worker(w, jobs, results, hostname, port, database)
	}

	for j := 0; j <= len(keys)-1; j++ {
		jobs <- keys[j]
	}
	close(jobs)

	// Finally we collect all the results of the work.
	for a := 0; a <= len(keys)-1; a++ {
		<-results
	}
	os.Exit(0)
}
