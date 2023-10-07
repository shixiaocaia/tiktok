package utils

import (
	"github.com/shixiaocaia/tiktok/cmd/videosvr/config"
	"github.com/shixiaocaia/tiktok/cmd/videosvr/log"
	"os/exec"
	"path/filepath"
	"strings"
)

func GetImageFile(videoPath string) (string, error) {
	temp := strings.Split(videoPath, "/")
	videoName := temp[len(temp)-1]
	b := []byte(videoName)
	videoName = string(b[:len(b)-3]) + "jpg"
	picPath := config.GetGlobalConfig().MinioConfig.PicPath
	picName := filepath.Join(picPath, videoName)
	cmd := exec.Command("ffmpeg", "-i", videoPath, "-ss", "1", "-f", "image2", "-t", "0.01", "-y", picName)
	err := cmd.Run()
	if err != nil {
		log.Errorf("When run ffmpeg, cmd.Run() failed with %s\n", err)
		return "", err
	}
	return picName, nil
}
