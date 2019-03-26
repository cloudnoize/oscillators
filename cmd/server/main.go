package main

import (
	"net/http"

	"github.com/cloudnoize/oscillators/oscillators"
	"github.com/cloudnoize/oscillators/transport"
)

func main() {
	cc := make(oscillators.ClientContexts)
	go transport.ServeUdp(":9876", cc)
	f := oscillators.NewFreq()
	http.ListenAndServe(":9876", transport.NewHTTPHandler(cc, f))
}
