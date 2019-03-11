package transport

import (
	"net"

	"github.com/cloudnoize/oscillators/oscillators"
)

type Client struct {
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
		n, add, e := conn.ReadFrom(b[:])
		if e != nil {
			continue
		}
		go srv(&Client{Oscillator: oscillators.NewSinOsc(44100, 440), addr: add}, conn)
	}
}

func srv(cl *Client, pc net.PacketConn) {
	cont := true
	defer func() {
		cont = false
	}()
	//Writer
	go func() {
		var b [1024]byte
		for cont {
			s := cl.GetSample()
			//Convert
			//Fill buffer
			//Then write
			pc.WriteTo(b[:], cl.addr)
		}
	}()
	//Reader
	var b [1024]byte
	for {
		n, addr, e := pc.ReadFrom(b[:])
		if e != nil {
			return
		}
	}
}
