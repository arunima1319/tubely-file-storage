package main 

import (
	"os/exec"
	"bytes"
	"encoding/json"
	"log"
)

func getVideoAspectRatio(filePath string) (string, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-print_format", "json", 	"-show_streams", filePath)
	var b bytes.Buffer
	cmd.Stdout = &b
	err := cmd.Run()
	if err!=nil{
		return "", err
	}

	type videoProbe struct{
		Streams []struct{
			Width int `json:"width,omitempty"`
			Height int `json:"height,omitempty"`
		}`json:"streams"`
	} 
		
	measures := videoProbe{}

	err = json.Unmarshal(b.Bytes(), &measures)
	if err!= nil{
		return "", err
	}

	width := float64(measures.Streams[0].Width)
	height := float64(measures.Streams[0].Height)
	ratio := width/height
	log.Printf("ratio: %v", ratio)
	if ratio < 1.79 && ratio > 1.75{
		return "16:9", nil
	}else if ratio > 0.55 && ratio < 0.57{
		return "9:16", nil
	}else{
		return "other", nil
	}


}