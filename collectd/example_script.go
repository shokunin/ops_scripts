package main

import "time"
import "fmt"
import "flag"
import "bufio"
import "os"
import "strconv"

var sleepTime int
var debug bool
var host string

func collectd(unixTs int, hostname string) {
	var hostlabel string
	if hostname == "localhost" {
		hostlabel, _ = os.Hostname()
	} else {
		hostlabel = hostname
	}

	f := bufio.NewWriter(os.Stdout)
	defer f.Flush()
	b := "PUTVAL " + hostlabel + "/" + "bar" + "/" + "gauge-name " + strconv.Itoa(unixTs) + ":" + "value\n"
	f.Write([]byte(b))
}

func init() {
	flag.BoolVar(&debug, "debug", false, "turn on debugging")
	flag.IntVar(&sleepTime, "sleep-time", 10, "Number of seconds between runs")
	flag.StringVar(&host, "host to connect to", "localhost", "Defaults to localhost")
	flag.Parse()
}

func main() {
	ticker := time.NewTicker(time.Second * time.Duration(sleepTime))
	// run once at the beginning
	collectd(int(time.Now().Unix()), host)
	go func() {
		//for t := range ticker.C {
		for t := range ticker.C {
			if debug {
				fmt.Println("Tick at", time.Now(), " DOH ", t)
			} else {
				collectd(int(time.Now().Unix()), host)
			}
		}
	}()

	// run for a year - as collectd will restart it
	time.Sleep(time.Second * 86400 * 365)
	ticker.Stop()
	fmt.Println("Ticker stopped")
}
