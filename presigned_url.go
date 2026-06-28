package main 

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"context"
	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/database"
	"strings"
	"time"
	"errors"
	"log"
)

func (cfg apiConfig) dbVideoToSignedVideo(video database.Video) (database.Video, error){

	if video.VideoURL == nil{
		return database.Video{}, errors.New("video URL does not exist")
	}
	bucketKey := strings.Split(*video.VideoURL, ",")
	log.Printf("this is the bucketKey variable: %v", bucketKey)
	var bucket string 
	var key string
	if len(bucketKey) == 2{
		bucket, key = bucketKey[0], bucketKey[1]
	} else{
		return database.Video{}, errors.New("video URL could not be split into bucket and key")
	}
	
	expireTime, _ := time.ParseDuration("1m")

	presignedURL, err := generatePresignedURL(cfg.s3Client, bucket, key, expireTime)
	if err!=nil{
		return database.Video{}, err
	}

	video.VideoURL = &presignedURL

	return video, nil
}
func generatePresignedURL(s3Client *s3.Client, bucket, key string, expireTime time.Duration) (string, error){

	presignClient := s3.NewPresignClient(s3Client)
	objectInput := &s3.GetObjectInput{
		Bucket: &bucket,
		Key: &key,
	}
	
	presignedReq, err := presignClient.PresignGetObject(context.Background(), objectInput, s3.WithPresignExpires(expireTime))
	if err!= nil{
		return "", err
	}
	return presignedReq.URL, nil

}