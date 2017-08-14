package goarstream

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

const (
	host     = "192.168.1.1:"
	port     = "5555"
	connType = "tcp"
)

func main() {
	video, _, _, _ := ffmpeg()
	conn, err := net.Dial(connType, host+port)
	if err != nil {
		fmt.Println(err)
	}

	reader := bufio.NewReader(conn)

	vidc := make(chan frame)
	quit := make(chan bool)
	startWorker(video, vidc, quit)

	for {
		var f frame
		err := f.parse(reader)
		if err != nil {
			fmt.Println("Error parsing:", err.Error())
		} else {
			vidc <- f
		}
	}
}

func startWorker(video *bufio.Writer, videoFrames chan frame, quit chan bool) {
	go func() {
		for {
			select {
			case frame := <-videoFrames:
				start := time.Now()
				if i, err := video.Write(frame.Payload); err != nil {
					fmt.Println("Error writing to ffmpeg: ", err.Error(), i)
				} else {
					fmt.Printf("Wrote %d bytes to ffmpeg in %v \n", i, time.Since(start))
				}
				if err := video.Flush(); err != nil {
					log.Fatal(err)
				}
				//fmt.Printf("%+v\n", frame)
			case <-quit:
				fmt.Println("quit")
				os.Exit(0)
			}
		}
	}()
}
