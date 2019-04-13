package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
)

const testFile = "video.mp4"
const testTime = 30
const videoSize = 0.43
const videoURL = "https://www.youtube.com/watch?v=zOWJqNPeifU" // Worst quality 0.43MB

func main() {
	for {
		tempDir, err := ioutil.TempDir("", "gfe-internet-check-tool")
		if err != nil {
			log.Fatal(errors.New(fmt.Sprintf("failed to get working directory %s", err)))
		}
		run(tempDir)
		os.RemoveAll(tempDir)
	}
}

func needForSpeed(size, timeSeconds float64) int {
	return (int)(size * 1024 / timeSeconds)
}

func run(dir string) {
	cmd := exec.Command("youtube-dl",
		"-r",
		strconv.Itoa(needForSpeed(videoSize, testTime)) + "K",
		"--socket-timeout",
		"3",
		"-o",
		dir + "/" + testFile,
		"-f",
		"worst[ext=mp4]",
		videoURL,
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(errors.New(fmt.Sprintf("failed to get pipe %s", err)))
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(errors.New(fmt.Sprintf("failed to get pipe %s", err)))
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	outputBuffer := make([]byte, 100)
	for {
		n, err := stdout.Read(outputBuffer)
		if err != nil {
			if err == io.EOF {
				fmt.Println(string(outputBuffer[:n]))
				break
			}
		}
		fmt.Printf("%s", string(outputBuffer[:n]))
	}

	if stdErr, err := ioutil.ReadAll(stderr); err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s\n", stdErr)
	}

	if err := cmd.Wait(); err != nil {
		fmt.Printf("Scientific browsing not good, download failed, retrying\n")
	} else {
		fmt.Printf("Scientific browsing still good, true internet")
	}
}
