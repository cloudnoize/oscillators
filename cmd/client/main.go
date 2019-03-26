package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
	"unsafe"

	"github.com/cloudnoize/conv"
	"github.com/cloudnoize/elport"
	locklessq "github.com/cloudnoize/locklessQ"
)

type streamImp struct {
	q *locklessq.Q
}

//mono float 32
func (s *streamImp) Cbb(inputBuffer, outputBuffer unsafe.Pointer, frames uint64) {
	ob := (*[512]float32)(outputBuffer)
	for i := uint64(0); i < frames; i++ {
		val, ok := s.q.Pop()
		if ok {
			(*ob)[i] = val

		}
	}
}

func (this *streamImp) Write(b []byte) (n int, err error) {
	for i := 0; i < len(b)/4; i++ {
		f := conv.BytesToFloat32(b, i*4)
		this.q.Insert(f)
	}
	return len(b), nil
}

func main() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	addr := os.Getenv("ADDR")

	udpaddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		println(err.Error())
		return
	}
	conn, err := net.DialUDP("udp", nil, udpaddr)
	if err != nil {
		println(err.Error())
		return
	}
	go func() {
		for {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Enter signal: ")
			text, _ := reader.ReadString('\n')
			println("sending ", text)
			conn.Write([]byte(text))
		}
	}()

	si := &streamImp{q: locklessq.New(int32(44100))}

	pa.Cba[0] = si

	pa.Initialize()
	s, _ := pa.OpenDefaultStream(0, 1, pa.Float32, 44100, 512, nil)
	s.Start()
	defer func() {
		s.Stop()
		s.Close()
	}()

	go io.Copy(si, conn)
	fmt.Println("awaiting signal")
	<-done
	fmt.Println("exiting")

}
