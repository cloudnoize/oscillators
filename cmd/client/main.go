package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
	"unsafe"

	"github.com/cloudnoize/conv"
	"github.com/cloudnoize/elport"
	locklessq "github.com/cloudnoize/locklessQ"
)

type streamImp struct {
	q32     *locklessq.Qfloat32
	q16     *locklessq.Qint16
	bitrate int
	frames  int
	test    bool
	count   int
	measure int
	start   time.Time
}

func (s *streamImp) CallBack(inputBuffer, outputBuffer unsafe.Pointer, frames uint64) {
	if s.bitrate == 16 {
		s.out16bit(outputBuffer, frames)
		return
	}
	s.out32bit(outputBuffer, frames)
}

func (s *streamImp) out32bit(outputBuffer unsafe.Pointer, frames uint64) {
	ob := (*[1024]float32)(outputBuffer)
	errNum := 0
	for i := 0; i < s.frames; i++ {
		val, ok := s.q32.Pop()
		if ok {
			if !s.test {
				(*ob)[i] = val
			} else {
				(*ob)[i] = 0
				if val != 1 {
					errNum++
					println("recieved ", val)
				}
			}

		} else {
			errNum++
		}
	}
	if errNum != 0 {
		if s.measure == 0 {
			s.start = time.Now()
		}
		s.measure++
	} else {
		if s.measure != 0 {
			ms := float64(time.Since(s.start) / time.Millisecond)
			fmt.Println("Time since last good packet ", ms, "ms frames elapsed ", s.measure)
			s.measure = 0
		}
	}
	println("frame ", s.count, " got ", errNum, " erred vals")

	s.count++
}

func (s *streamImp) out16bit(outputBuffer unsafe.Pointer, frames uint64) {
	ob := (*[1024]int16)(outputBuffer)
	for i := 0; i < s.frames; i++ {
		val, ok := s.q16.Pop()
		if ok {
			if !s.test {
				(*ob)[i] = int16(val)
			} else {
				(*ob)[i] = 0
				if val != 1 {
					println("recieved ", val)
				}
			}

		}
	}
}

func (s *streamImp) Write(b []byte) (n int, err error) {
	if s.bitrate == 16 {
		s.Write16int(b)
		return len(b), nil
	}
	s.Write32float(b)
	return len(b), nil
}

func (this *streamImp) Write32float(b []byte) {
	for i := 0; i < len(b)/4; i++ {
		f := conv.BytesToFloat32(b, i*4)
		this.q32.Insert(f)
	}
}

func (this *streamImp) Write16int(b []byte) {
	for i := 0; i < len(b)/2; i++ {
		s := conv.BytesToint16(b, i*2)
		this.q16.Insert(s)
	}
}

func main() {
	bitrate := int(32)
	if v := os.Getenv("BIT_RATE"); v != "" {
		brate, _ := strconv.Atoi(v)
		switch brate {
		case 16:
			bitrate = brate
		}
	}

	frames := 512
	if v := os.Getenv("FRAMES"); v != "" {
		frames, _ = strconv.Atoi(v)
	}
	println("frames - ", frames)
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

	test := false
	if v := os.Getenv("TEST"); v != "" {
		test = true
		println("test mode")
	}

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

	si := &streamImp{q32: locklessq.NewQfloat32(int32(44100)), q16: locklessq.NewQint16(int32(44100)), bitrate: bitrate, frames: frames, test: test}

	pa.CbStream = si

	pa.Initialize()
	sf := pa.Float32

	if bitrate == 16 {
		sf = pa.Int16
	}

	s, _ := pa.OpenDefaultStream(0, 1, sf, 44100, uint64(frames), nil)

	go func() {
		for {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Enter registration: ")
			text, _ := reader.ReadString('\n')
			println("sending ", text)
			conn.Write([]byte(text))
			fmt.Print("Enter signal: ")
			text, _ = reader.ReadString('\n')
			println("sending ", text)
			conn.Write([]byte(text))
			// time.Sleep(1 * time.Second)
			s.Start()
			return
		}
	}()

	defer func() {
		s.Stop()
		s.Close()
	}()

	go io.Copy(si, conn)
	fmt.Println("awaiting signal")
	<-done
	fmt.Println("exiting")

}
