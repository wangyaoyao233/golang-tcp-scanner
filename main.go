package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"sort"
	"time"
)

type Result struct {
	open bool
	port int
}

type Config struct {
	Addr    string
	MinPort int
	MaxPort int
}

var ConfigObj *Config

func initConfig() {
	ConfigObj = &Config{
		Addr:    "127.0.0.1",
		MinPort: 1,
		MaxPort: 1024,
	}

	ConfigObj.reloadConfig()
}

func (c *Config) reloadConfig() {
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, c)
	if err != nil {
		panic(err)
	}

	fmt.Printf("ip: %s\n", ConfigObj.Addr)
	fmt.Printf("scan from port: %d\nto port: %d\n", ConfigObj.MinPort, ConfigObj.MaxPort)
}

func worker(ports chan int, results chan Result) {
	for p := range ports {
		addr := fmt.Sprintf("%s:%d", ConfigObj.Addr, p)
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			result := Result{open: false, port: p}
			results <- result
			continue
		}
		conn.Close()
		result := Result{open: true, port: p}
		results <- result
	}
}

func main() {
	start := time.Now()

	initConfig()

	ports := make(chan int, 100)
	results := make(chan Result)
	var openports []int
	var closedports []int

	for i := 0; i < cap(ports); i++ {
		go worker(ports, results)
	}

	go func() {
		for i := ConfigObj.MinPort; i < ConfigObj.MaxPort; i++ {
			ports <- i
		}
	}()

	for i := ConfigObj.MinPort; i < ConfigObj.MaxPort; i++ {
		result := <-results
		if result.open {
			openports = append(openports, result.port)
		} else {
			closedports = append(closedports, result.port)
		}
	}

	close(ports)
	close(results)

	sort.Ints(openports)
	sort.Ints(closedports)

	for _, port := range openports {
		fmt.Printf("%d is opened!\n", port)
	}

	elapsed := time.Since(start) / 1e9
	fmt.Printf("\n\ntotal cost time %d seconds!", elapsed)
}

// func main() {
// 	start := time.Now()

// 	var wg sync.WaitGroup
// 	for i := 1; i < 65535; i++ {
// 		wg.Add(1)
// 		go func(j int) {
// 			defer wg.Done()
// 			addr := fmt.Sprintf("1.15.129.136:%d", j)
// 			conn, err := net.Dial("tcp", addr)
// 			if err != nil {
// 				fmt.Printf("%s is closed\n", addr)
// 				return
// 			}
// 			fmt.Printf("%s is opened\n", addr)
// 			conn.Close()
// 		}(i)
// 	}
// 	wg.Wait()

// 	elapsed := time.Since(start) / 1e9
// 	fmt.Printf("\n\ncost time %d seconds\n", elapsed)
// }

// func main() {
// 	for i := 21; i < 120; i++ {
// 		addr := fmt.Sprintf("1.15.129.136:%d", i)
// 		conn, err := net.Dial("tcp", addr)
// 		if err != nil {
// 			fmt.Printf("%s is closed\n", addr)
// 			continue
// 		}
// 		fmt.Printf("%s is opened\n", addr)
// 		conn.Close()
// 	}
// }
