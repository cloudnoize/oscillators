package transport

import (
	"net"
	"strconv"

	"github.com/cloudnoize/conv"

	"github.com/cloudnoize/oscillators/oscillators"
)

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
		cc.Oscillator = oscillators.NewSinOsc(44100, 440)
		println("Ready, send me signal to start")
		conn.ReadFrom(b[:])
		go srv(cc, conn)
	}
}

func srv(cl *oscillators.ClientContext, pc net.PacketConn) {
	cont := true
	defer func() {
		cont = false
		cl.Close()
	}()
	//Writer
	var b [1024]byte
	for cont {
		for i := 0; i < len(b)/4; i++ {
			s := cl.GetSample()
			conv.Float32ToBytes(s, b[:], i*4)
		}
		_, e := pc.WriteTo(b[:], cl.Addr)
		if e != nil {
			println("rec ", e.Error())
			return
		}
	}
}
