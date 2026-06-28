package main 

import (
	"os/exec"
)

func processVideoForFastStart(filePath string) (string, error) {
	newPath := filePath + ".processing"

	cmd := exec.Command("ffmpeg", "-i", filePath, "-c", "copy", "-movflags", "faststart", "-f", "mp4", newPath)

	err := cmd.Run()
	if err!=nil{
		return "", err
	}
	return newPath, nil
}