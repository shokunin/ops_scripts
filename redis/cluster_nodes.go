package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-redis/redis"
)

type clusterNode struct {
	id      string
	ip      string
	port    int
	cmdport int
	role    string
	slaves  []string
}

func listMasters(clusterNodes []clusterNode) []string {
	var masters []string
	for _, v := range clusterNodes {
		if v.role == "master" {
			masters = append(masters, v.ip+":"+strconv.Itoa(v.port))
		}
	}
	return masters
}

func parseNodes(nodes *redis.StringCmd) []clusterNode {
	var clusterNodes []clusterNode
	for _, line := range strings.Split(nodes.Val(), "\n") {
		ln := strings.Split(line, " ")
		if len(ln) > 1 {
			role := "slave"
			r := regexp.MustCompile(`(\S+):(\d+)@(\d+)`)
			res := r.FindStringSubmatch(ln[1])
			match, _ := regexp.MatchString("master", ln[2])
			if match {
				role = "master"
			}

			i, _ := strconv.Atoi(res[2])
			j, _ := strconv.Atoi(res[3])

			n := clusterNode{
				id:      ln[0],
				role:    role,
				ip:      res[1],
				port:    i,
				cmdport: j,
			}
			clusterNodes = append(clusterNodes, n)

			if role == "slave" {
				for i, v := range clusterNodes {
					if v.id == ln[3] {
						clusterNodes[i].slaves = append(clusterNodes[i].slaves, ln[0])
					}
				}
			}

		}
	}

	return clusterNodes
}

func getKeyspace(servers []string, password string) int {
	keys := 0
	for _, server := range servers {
		client := redis.NewClient(&redis.Options{
			Addr:     server,
			Password: password, // no password set
		})
		info := client.Info("keyspace")
		for _, line := range strings.Split(info.Val(), "\n") {
			r := regexp.MustCompile(`db\d+:keys=(\d+),`)
			res := r.FindStringSubmatch(line)
			if len(res) > 0 {
				j, _ := strconv.Atoi(res[1])
				keys += j
			}
		}
	}
	return keys
}

func getMemory(servers []string, password string) int {
	bytes := 0
	for _, server := range servers {
		client := redis.NewClient(&redis.Options{
			Addr:     server,
			Password: password, // no password set
		})
		info := client.Info("memory")
		for _, line := range strings.Split(info.Val(), "\n") {
			r := regexp.MustCompile(`used_memory:(\d+)`)
			res := r.FindStringSubmatch(line)
			if len(res) > 0 {
				j, _ := strconv.Atoi(res[1])
				bytes += j
			}
		}
	}
	return (bytes)
}

func getReplicationFactor(clusterNodes []clusterNode) int {
	var repFactor []int
	for _, v := range clusterNodes {
		if v.role == "master" {
			repFactor = append(repFactor, len(v.slaves))
		}
	}
	return (sliceMax(repFactor))
}

func sliceMax(s []int) int {
	m := 0
	for i, e := range s {
		if i == 0 || e > m {
			m = e
		}
	}
	return (m)
}

func main() {

	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{":30001"},
	})
	j := rdb.ClusterNodes()
	k := parseNodes(j)
	m := listMasters(k)
	fmt.Println("master_count", len(m))
	fmt.Println("replication_factor", getReplicationFactor(k))
	fmt.Println("total_key_count", getKeyspace(m, ""))
	fmt.Println("total_memory", getMemory(m, ""))
	os.Exit(0)
}
