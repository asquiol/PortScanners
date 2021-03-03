package main

import (
	"context"
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

type Scanner struct {
	ip   string
	lock *semaphore.Weighted
}

func Abvalue() int64 {
	out, err := exec.Command("abvalue", "-n").Output()
	if err != nil {
		panic(err)
	}
	
	s := strings.TrimSpace(string(out))
	
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}
	
	return i
}

func PortScan(ip string, port int, timeout time.Duration) {
	target := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", target, timeout)

	if err != nil {
		if strings.Contains(err.Error(), "too many open files") {
			time.Sleep(timeout)
			PortScan(ip, port, timeout)
		} else {
			fmt.Println(port, "closed")
		}
		return
	}

	conn.Close()
	fmt.Println(port, "open")
}

func (ps *Scanner) Start(f, l int, timeout time.Duration) {
	wg := sync.WaitGroup{}
	defer wg.Wait()

	for port := f; port <= l; port++ {
		ps.lock.Acquire(context.TODO(), 1)
		wg.Add(1)
		go func(port int) {
			defer ps.lock.Release(1)
			defer wg.Done()
			PortScan(ps.ip, port, timeout)
		}(port)
	}
}

func main() {
	ps := &Scanner{
		ip:   "127.0.0.1",
		lock: semaphore.NewWeighted(abvalue()),
	}
	ps.Start(1, 65535, 500*time.Millisecond)
}