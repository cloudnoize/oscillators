package transport

import (
	"net"

	"github.com/cloudnoize/conv"

	"github.com/cloudnoize/oscillators/oscillators"
)

type ClientCon struct {
	oscillators.Oscillator
	addr net.Addr
}

func ServeUdp(addr string) {
	conn, e := net.ListenPacket("udp", addr)
	if e != nil {
		println(e.Error())
		return
	}

	for {
		var b [1024]byte
		_, add, e := conn.ReadFrom(b[:])
		println("Init...")
		if e != nil {
			continue
		}
		cc := &ClientCon{Oscillator: oscillators.NewSinOsc(44100, 440), addr: add}
		go srv(cc, conn)
		println("Ready, send me signal to start")
		conn.ReadFrom(b[:])
	}
}

func srv(cl *ClientCon, pc net.PacketConn) {
	cont := true
	defer func() {
		cont = false
		cl.Close()
	}()
	//Writer
	go func() {
		var b [1024]byte
		for cont {
			for i := 0; i < len(b)/4; i++ {
				s := cl.GetSample()
				conv.Float32ToBytes(s, b[:], i*4)
			}
			pc.WriteTo(b[:], cl.addr)
		}
	}()
	//Reader
	var b [1024]byte
	for {
		_, _, e := pc.ReadFrom(b[:])
		if e != nil {
			return
		}
	}
}
