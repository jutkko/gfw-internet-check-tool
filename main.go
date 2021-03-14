package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/tidwall/gjson"
)

const debug = false
const tempFile = "video.mp4"
const testTime = 5
const videoSize = 0.094
const youtubeAPI = "http://youtube.com/get_video_info?video_id="
const videoID = "NQ9RtLrapzc"
const timeout = 5

func main() {
	client := &http.Client{
		Timeout: timeout * time.Second,
	}

	for {
		run(client)

		<-time.After(1 * time.Second)
	}
}

func run(client *http.Client) {
	dir, err := ioutil.TempDir("", "gfe-internet-check-tool")
	if err != nil {
		log.Fatal(errors.New(fmt.Sprintf("failed to get working directory %s", err)))
	}

	defer os.RemoveAll(dir)

	info, err := getVideoInfo(client)
	if err != nil {
		log.Printf("scientific browsing not good, download failed get video info %s, retrying", err.Error())
		return
	}

	url, err := getDownloadURLFromVideoInfo(info)
	if err != nil {
		log.Printf("scientific browsing not good, download failed get video url %s, retrying", err.Error())
	}

	cmd := exec.Command("wget",
		"--limit-rate",
		strconv.Itoa(needForSpeed(videoSize, testTime))+"K",
		"--tries",
		"0",
		"--timeout",
		strconv.Itoa(timeout),
		"-O",
		dir+"/"+tempFile,
		url,
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
	} else if len(stdErr) != 0 && debug {
		log.Printf("%s", stdErr)
	}

	if err := cmd.Wait(); err != nil {
		log.Printf("scientific browsing not good, download failed wget %s, retrying", err.Error())
	} else {
		log.Printf("scientific browsing still good, true internet")
	}
}

func needForSpeed(sizeMB, timeSeconds float64) int {
	return (int)(sizeMB * 1024 / timeSeconds)
}

func getVideoInfo(client *http.Client) (string, error) {
	resp, err := client.Get(youtubeAPI + videoID)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func getDownloadURLFromVideoInfo(videoInfo string) (string, error) {
	answer, err := url.ParseQuery(videoInfo)
	if err != nil {
		return "", err
	}

	status, ok := answer["status"]
	if !ok {
		err = fmt.Errorf("no response status found in the server's answer")
		return "", err
	}
	if status[0] == "fail" {
		reason, ok := answer["reason"]
		if ok {
			err = fmt.Errorf("'fail' response status found in the server's answer, reason: '%s'", reason[0])
		} else {
			err = errors.New(fmt.Sprint("'fail' response status found in the server's answer, no reason given"))
		}

		return "", err
	}
	if status[0] != "ok" {
		err = fmt.Errorf("non-success response status found in the server's answer (status: '%s')", status)
		return "", err
	}

	// read the streams map
	playerResponse, ok := answer["player_response"]
	if !ok {
		err = errors.New(fmt.Sprint("no stream map found in the server's answer"))
		return "", err
	}

	streamingURL := gjson.Get(playerResponse[0], "streamingData.formats.#(url).url")

	url, err := url.Parse(streamingURL.String())
	if err != nil {
		err = fmt.Errorf("no url found in the stream")
		return "", err
	}

	return url.String(), nil
}
