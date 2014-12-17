package main

import (
	"strings"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"
)

var hostname string
var port int
var downtime int
var flapjack string
var hn string

type Entities struct {
	Entities []struct {
		Id    string `json:"id"`
		Name  string `json:"name"`
		Links struct {
			Contacts []string `json:"contacts"`
		} `json:"links"`
	} `json:"entities"`
}

func init() {

	// We need to make sure that we get the longest name returned
	sn, err := os.Hostname()
	if err != nil {
		fmt.Println(err)
	}

	cmd := exec.Command("/bin/hostname", "-f")
	out, err := cmd.Output()

	if err != nil {
		fmt.Println("using hosntame of " + sn + "since /bin/hostname -f failed")
	}
	ln := strings.Replace(string(out),"\n","",-1)

	//print(string(out))

	if len(ln) > len(sn) {
		hn = ln
	} else {
		hn = sn
	}

	flag.StringVar(&flapjack, "flapjack", "flapjack", "Home of the flapjack server")
	flag.StringVar(&hostname, "hostname", hn, "hostname to subscribe to")
	flag.IntVar(&downtime, "downtime", 7200, "seconds to set downtime for")
	flag.IntVar(&port, "port", 3081, "port to try to connect to")
	flag.Parse()
}

func get_host_id(hn string, fj string, prt int) string {
	url := "http://" + fj + ":" + strconv.Itoa(port) + "/entities"
	res, err := http.Get(url)
	var hostid string
	//	filtered := []Entities{}

	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		panic(err.Error())
	}

	var data Entities

	json.Unmarshal(body, &data)

	for _, e := range data.Entities {
		if e.Name == hn {
			hostid = e.Id
		}
	}

	return hostid

}

func set_downtime(hostid string, fj string, prt int, dt int) {
	url := "http://" + fj + ":" + strconv.Itoa(port) + "/scheduled_maintenances/entities/" + hostid
	//fmt.Println(url)
	lt := time.Now()
	t := lt.UTC()
	start_time := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d-00:00",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())

	jsonStr := []byte("{ \"scheduled_maintenances\": [ { \"start_time\" : \"" + start_time + "\", \"duration\" : " + strconv.Itoa(dt) + ", \"summary\" : \"STFU\" } ] }")
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/vnd.api+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.Status == "204 No Content" {
		fmt.Println("Downtime set")
	} else {
		fmt.Println("Error Occured:")
		fmt.Println("response Status:", resp.Status)
	}

}

func main() {
	host_id := get_host_id(hostname, flapjack, port)
	set_downtime(host_id, flapjack, port, downtime)
}
