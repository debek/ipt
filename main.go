package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

const (
	Red   = "\033[31m"
	Green = "\033[32m"
	Reset = "\033[0m"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: ipt <ip> <port> [<timeout>]")
		os.Exit(1)
	}

	ip := os.Args[1]
	port := os.Args[2]
	timeout := "5s"

	if len(os.Args) == 4 {
		timeout = os.Args[3] + "s"
	}

	timeoutDuration, err := time.ParseDuration(timeout)
	if err != nil {
		fmt.Printf(Red+"Error: Invalid timeout format: '%s'. Please provide time in format like '5s' for 5 seconds.\n"+Reset, timeout)
		os.Exit(1)
	}

	startTime := time.Now()
	fmt.Printf("[%s] START: Connection Test to %s:%s with a timeout of %s.\n", startTime.Format("2006-01-02 15:04:05"), ip, port, timeoutDuration)

	for {
		currentTime := time.Now()
		err := checkConnection(ip, port, timeoutDuration)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				fmt.Printf(Red+"[%s] FAILURE: Connection to %s:%s timed out (possible firewall rejection)\n"+Reset, currentTime.Format("2006-01-02 15:04:05"), ip, port)
			} else {
				fmt.Printf(Red+"[%s] FAILURE: Connection to %s:%s failed (service not running or blocked)\n"+Reset, currentTime.Format("2006-01-02 15:04:05"), ip, port)
			}
		} else {
			fmt.Printf(Green+"[%s] SUCCESS: Connected to %s:%s\n"+Reset, currentTime.Format("2006-01-02 15:04:05"), ip, port)
		}
		time.Sleep(1 * time.Second)
	}
}

func checkConnection(ip string, port string, timeout time.Duration) error {
	address := net.JoinHostPort(ip, port)
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return err
	}
	defer conn.Close()
	return nil
}
