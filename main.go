package main

import (
	"monitor/server"
	"runtime"
	"fmt"
)

func main() {
	server.Run()
	fmt.Println(runtime.NumGoroutine())
}
