package main

import (
	"github.com/cloudnoize/oscillators/transport"
)

func main() {
	transport.ServeUdp(":9876")
}
