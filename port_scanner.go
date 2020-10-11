package main

import (
	"fmt"
	"net"
	"os"
	"sort"
)

func worker(ip string, ports, results chan int) {
	for p := range ports {
		address := fmt.Sprintf("%s:%d", ip, p)
		fmt.Println(address)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			results <- 0
			continue
		}
		results <- p
		conn.Close()
	}
}

func distribute(ports chan int, r int) {
	for i := 0; i < r; i++ {
		ports <- i
	}
}

func collector(results chan int, openPorts []int, r int) []int {
	for i := 0; i < r; i++ {
		port := <-results
		if port != 0 {
			fmt.Println("Port : ", port)
			openPorts = append(openPorts, port)
		}
	}
	return openPorts
}

func main() {
	var ip string
	var ports = make(chan int, 100)
	var results = make(chan int)
	var openPorts []int
	if len(os.Args) == 1 {
		fmt.Println("Please provide an IP Address.")
		os.Exit(0)
	}
	ip = os.Args[1]
	fmt.Println("IP to scan : ", ip)

	for i := 0; i < 2048; i++ {
		go worker(ip, ports, results)
	}
	go distribute(ports, 1024)
	openPorts = collector(results, openPorts, 1024)
	close(ports)
	close(results)
	sort.Ints(openPorts)
	for _, port := range openPorts {
		fmt.Println("FOUND OPEN : ", port)
	}
}
