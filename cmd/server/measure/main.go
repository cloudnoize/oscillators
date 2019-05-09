package main

import (
	"os"
	"strconv"
	"time"
)

func main() {
	Gap := 1
	if v := os.Getenv("GAP"); v != "" {
		Gap, _ = strconv.Atoi(v)
	}
	dur := time.Duration(Gap) * time.Microsecond
	println("tikcer dur is ", dur.String())
	//t := time.NewTicker(dur)
	count := uint(0)
	st := time.Now()
	for true {
		count++
		if time.Since(st).Seconds() >= 1 {
			println("elapsed ", time.Since(st).String(), " count ", count)
			count = 0
			st = time.Now()
		}
	}

}
