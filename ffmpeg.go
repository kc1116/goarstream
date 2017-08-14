package goarstream

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
)

//Trying to get png from stdout pipe

func ffmpeg() (bufStdin *bufio.Writer, stderr io.ReadCloser, stdout io.ReadCloser, err error) {
	//ffmpegCmd := exec.Command("ffmpeg", "-re", "-f", "h264", "-i", "pipe:0", "-f", "webm", "pipe:1")
	ffmpegCmd := exec.Command("ffmpeg", "-i", "-", "-f", "image2pipe", "png", "pipe:1")
	//ffmpegCmd := exec.Command("ffmpeg", "-i", "pipe:0", "-f", "image2pipe", "-codec:v", "png", "-strict", "-2", "$out > /dev/null 2>&1")

	//ffmpegCmd := exec.Command("ffplay", "-f", "h264", "-i", "pipe:0")
	//ffmpegCmd.Stdin = os.Stdin
	//ffmpegCmd.Stdout = os.Stdout

	stderr, err = ffmpegCmd.StderrPipe()
	if err != nil {
		return
	}

	stdin, err := ffmpegCmd.StdinPipe()
	if err != nil {
		return
	}

	bufStdin = bufio.NewWriter(stdin)

	stdout, err = ffmpegCmd.StdoutPipe()
	if err != nil {
		return
	}

	go func() {
		for {
			buf, err := ioutil.ReadAll(stderr)
			if err != nil {
				fmt.Println(err)
			}
			if len(buf) > 0 {
				fmt.Println(string(buf))
			}
		}
	}()

	go func() {
		for {
			count := 0
			buf, err := ioutil.ReadAll(stdout)
			if err != nil {
				fmt.Println(err)
			}
			if len(buf) > 0 {
				img, _, _ := image.Decode(bytes.NewReader(buf))

				//save the imgByte to file
				out, err := os.Create("./img-" + strconv.Itoa(count) + ".png")

				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				err = png.Encode(out, img)

				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}

		}
	}()

	err = ffmpegCmd.Start()
	if err != nil {
		fmt.Println(err)
	}

	return bufStdin, stderr, nil, nil
}
