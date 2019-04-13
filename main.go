package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
)

const testFile = "video.mp4"
const testTime = 5
const videoSize = 0.28
const videoURL = "https://www.youtube.com/watch?v=NQ9RtLrapzc"

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

func needForSpeed(sizeMB, timeSeconds float64) int {
	return (int)(sizeMB * 1024 / timeSeconds)
}

func run(dir string) {
	cmd := exec.Command("youtube-dl",
		"-r",
		strconv.Itoa(needForSpeed(videoSize, testTime))+"K",
		"--socket-timeout",
		"3",
		"-o",
		dir+"/"+testFile,
		"-f",
		"worst[ext=mp4]",
		videoURL,
	)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(errors.New(fmt.Sprintf("failed to get pipe %s", err)))
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	if stdErr, err := ioutil.ReadAll(stderr); err != nil {
		log.Fatal(err)
	} else if len(stdErr) != 0 {
		fmt.Printf("%s\n", stdErr)
	}

	if err := cmd.Wait(); err != nil {
		log.Printf("Scientific browsing not good, download failed, retrying")
	} else {
		log.Printf("Scientific browsing still good, true internet")
	}
}
