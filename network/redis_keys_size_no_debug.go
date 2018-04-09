package main

// MAKE SURE TO SET export GOMAXPROCS=16  in your environment before running

import "fmt"
import "strings"
import "strconv"
import "golang.org/x/sync/syncmap"
import "flag"
import "time"
import "os"
import "unsafe"
import "github.com/fzzy/radix/redis"

func fetchAllKeys(hostame string, port int, database int) []string {
	c, err := redis.DialTimeout("tcp", fmt.Sprintf("%s:%d", hostname, port), time.Duration(300)*time.Second)
	errHndlr(err)
	keys := c.Cmd("SELECT", database)
	keys = c.Cmd("KEYS", "*")
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

func worker(id int, jobs <-chan string, results chan<- string, hostame string, port int, database int, sm *syncmap.Map) {
	c, err := redis.DialTimeout("tcp", fmt.Sprintf("%s:%d", hostname, port), time.Duration(3)*time.Second)
	errHndlr(err)
	for j := range jobs {
		rsize := uint64(0)
		errHndlr(err)
		c.Cmd("SELECT", database)
		z := c.Cmd("TYPE", j)
		w, _ := z.Bytes()
		if string(w) == "string" {
			k := c.Cmd("GET", j)
			rsize += uint64(unsafe.Sizeof(k))
			y, _ := k.Bytes()
			rsize += uint64(len(y))
			fmt.Println("total size: ", rsize, "customer: ", j)
		}
		if string(w) == "hash" {
			k := c.Cmd("HGETALL", j)
			s := strings.Split(j, ":")
			rsize += uint64(unsafe.Sizeof(k))
			if k.Type > 4 {
				for _, x := range k.Elems {
					y, _ := x.Bytes()
					rsize += uint64(len(y))
				}
			}
			val, hasVal := sm.Load(s[len(s)-1])
			if hasVal {
				i, _ := strconv.ParseUint(val.(string), 10, 64)
				sm.Store(s[len(s)-1], fmt.Sprintf("%d", rsize+i))
			} else {
				sm.Store(s[len(s)-1], fmt.Sprintf("%d", (rsize)))
			}
		}
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

	sm := new(syncmap.Map)

	keys := fetchAllKeys(hostname, port, database)
	// In order to use our pool of workers we need to send
	// them work and collect their results. We make 2
	// channels for this.
	jobs := make(chan string, len(keys))
	results := make(chan string, len(keys))

	for w := 0; w <= concurrent; w++ {
		go worker(w, jobs, results, hostname, port, database, sm)
	}

	for j := 0; j <= len(keys)-1; j++ {
		jobs <- keys[j]
	}
	close(jobs)

	// Finally we collect all the results of the work.
	for a := 0; a <= len(keys)-1; a++ {
		<-results
	}
	sm.Range(func(k, v interface{}) bool {
		fmt.Printf("%s,%s\n", v, k)
		return true
	})

	os.Exit(0)
}
