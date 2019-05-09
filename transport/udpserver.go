package transport

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/cloudnoize/conv"

	"github.com/cloudnoize/oscillators/oscillators"
)

var (
	bitrate int
	buf     []byte
	osc     string
	latency float32
)

func init() {
	bitrate = 32
	if v := os.Getenv("BIT_RATE"); v != "" {
		brate, _ := strconv.Atoi(v)
		switch brate {
		case 16:
			bitrate = brate
		}
	}
	bytesPerSample := bitrate / 8

	bufsizeSamples := 1024
	if v := os.Getenv("BUF_SIZE_SAMPLES"); v != "" {
		bsize, _ := strconv.Atoi(v)
		bufsizeSamples = bsize
	}

	bufsizeBytes := bufsizeSamples * bytesPerSample

	//22.675 per sample
	latency = 22.675 * float32(bufsizeSamples)
	dur := time.Duration(latency) * time.Microsecond
	println("Sleeping ", dur.String(), " between calls")

	buf = make([]byte, bufsizeBytes, bufsizeBytes)
	fmt.Printf("Bitrate was set to %d buffer size is %d samples \n", bitrate, bufsizeSamples)
}

func ServeUdp(addr string, contexts oscillators.ClientContexts) {
	conn, e := net.ListenPacket("udp", addr)
	if e != nil {
		println(e.Error())
		return
	}

	for {
		var b [1024]byte
		n, add, e := conn.ReadFrom(b[:])
		println("Init...")
		if e != nil {
			continue
		}
		icl, e := strconv.Atoi(string(b[:n-1]))
		println("got ", icl)
		if e != nil {
			println(e.Error())
			continue
		}
		cc, ok := contexts[uint(icl)]
		if cc == nil || !ok {
			println("got nil for cl ", icl)
			return
		}
		cc.Addr = add
		cc.Oscillator = oscillators.GetOsc(cc.Osc)
		println("Ready, send me signal to start")
		conn.ReadFrom(b[:])
		go srv(cc, conn)
	}
}

func srv(cl *oscillators.ClientContext, pc net.PacketConn) {
	f := fill32bit
	println("Serving...")
	if bitrate == 16 {
		println("serve 16 bits")
		f = fill16bit
	}
	cont := true
	defer func() {
		cont = false
		cl.Close()
	}()
	dur := time.Duration(latency) * time.Microsecond
	println("will sleep for ", dur.String(), " between calls")
	//Writer
	for cont {
		//Simulates sending sample rate num samples per sec, by filling buffer then sleeping.
		f(cl, buf)
		time.Sleep(dur)

		_, e := pc.WriteTo(buf, cl.Addr)
		if e != nil {
			println("rec ", e.Error())
			return
		}
	}
}

func fill32bit(cl *oscillators.ClientContext, buf []byte) {
	for i := 0; i < len(buf)/4; i++ {
		s := cl.GetSample()
		conv.Float32ToBytes(s, buf, i*4)
	}
}

func fill16bit(cl *oscillators.ClientContext, buf []byte) {
	for i := 0; i < len(buf)/2; i++ {
		s := cl.GetSample()
		conv.Float32To16intBytes(s, buf, i*2)
	}
}
