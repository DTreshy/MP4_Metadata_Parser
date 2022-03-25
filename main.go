package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	out, err := exec.Command("./ffprobe.exe", os.Args[1]).CombinedOutput()
	data := map[string]string{
		"type":        "Audio",
		"duration":    "0",
		"video_codec": "None",
		"audio_codec": "None",
		"width":       "0",
		"height":      "0",
	}
	if err != nil {
		log.Fatal(err)
	}
	sOut := strings.Split(string(out), "\n")
	for i := 0; i < len(sOut); i++ {
		if len(sOut[i]) > 7 && sOut[i][2:8] == "Stream" {
			splitted := strings.Split(sOut[i], " ")
			if splitted[4] == "Video:" {
				data["type"] = "Video"
				data["video_codec"] = splitted[5]
				regex, _ := regexp.Compile("^\\d+[x]\\d+,?$")
				for j := 6; j < len(splitted); j++ {
					if regex.MatchString(splitted[j]) {
						data["width"] = strings.Split(splitted[j], "x")[0]
						data["height"] = strings.Split(splitted[j], "x")[1][0:len(strings.Split(splitted[j], "x")[1])]
						data["height"] = strings.TrimRight(data["height"], ",")
					}
				}
			} else {
				data["audio_codec"] = splitted[5]
			}

		} else if len(sOut[i]) > 9 && sOut[i][2:10] == "Duration" {
			h_seconds, err := strconv.Atoi(sOut[i][12:14])
			m_seconds, err := strconv.Atoi(sOut[i][15:17])
			seconds, err := strconv.Atoi(sOut[i][18:20])

			if err != nil {
				log.Fatal(err)
			}
			data["duration"] = strconv.Itoa(h_seconds*60*60 + m_seconds*60 + seconds)
		}
	}
	fmt.Print("Type:         ", data["type"], "\n")
	fmt.Print("Duration:     ", data["duration"], "\n")
	fmt.Print("Video Codec:  ", data["video_codec"], "\n")
	fmt.Print("Audio Codec:  ", data["audio_codec"], "\n")
	fmt.Print("Width:        ", data["width"], "\n")
	fmt.Print("Height:       ", data["height"], "\n")
}
