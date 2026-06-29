package main

import (
	"os"
	"path/filepath"
	"fmt"
	"strings"
	
)


func (cfg apiConfig) ensureAssetsDir() error {
	if _, err := os.Stat(cfg.assetsRoot); os.IsNotExist(err) {
		return os.Mkdir(cfg.assetsRoot, 0755)
	}
	return nil
}

func getAssetPath(videoID string, mediaType string) string{
	extension := getMediaTypeExtension(mediaType)
	return fmt.Sprintf("%s%s", videoID, extension)
}

func (cfg apiConfig) getAssetDiskPath(assetPath string) string{
	return filepath.Join(cfg.assetsRoot, assetPath)
}

func (cfg apiConfig) getAssetURL(assetPath string) string{
	return fmt.Sprintf("http://localhost:%s/assets/%s", cfg.port, assetPath)
}

func getMediaTypeExtension(mediaType string) string{
	parts := strings.Split(mediaType, "/")
	if len(parts) != 2{
		return ".bin"
	}
	return fmt.Sprintf(".%s", parts[1])

}