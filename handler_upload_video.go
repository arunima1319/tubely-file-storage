package main

import (
	"fmt"
	"net/http"
	"io"
	"errors"
	"os"
	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/auth"
	"github.com/google/uuid"
	"mime"
	"crypto/rand"
	"encoding/base64"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"log"
)

func (cfg *apiConfig) handlerUploadVideo(w http.ResponseWriter, r *http.Request) {

	r.Body = http.MaxBytesReader(w, r.Body, 1 << 30)

	videoIDString := r.PathValue("videoID")
	videoID, err := uuid.Parse(videoIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	video, err := cfg.db.GetVideo(videoID)
	if err!=nil{
		respondWithError(w, http.StatusInternalServerError, "unable to get videos", err)
		return
	}
    fmt.Printf("jwt userID: %s, video ownerID: %s\n", userID, video.UserID)

	if video.UserID != userID{
		respondWithError(w, http.StatusUnauthorized, "authenticated user is not owner of video", errors.New("not authorized"))
		return
	}

	file, header, err := r.FormFile("video")
	if err!= nil{
		respondWithError(w, http.StatusBadRequest, "unable to parse form file", err)
		return
	}
	defer file.Close()

	mediaType, _, err := mime.ParseMediaType(header.Header.Get("Content-Type"))
	if err!=nil{
		respondWithError(w, http.StatusBadRequest, "invalid Content-Type", err)
		return
	}
	if mediaType != "video/mp4" {
		respondWithError(w, http.StatusBadRequest, "invalid media type", nil)
		return
	}

	tmpFile, err := os.CreateTemp("", "tubely-upload.mp4")
	if err!=nil{
		respondWithError(w, http.StatusInternalServerError, "could not create file on server", err)
		return
	}
	defer os.Remove("tubely-upload.mp4")
	defer tmpFile.Close()

	_, err = io.Copy(tmpFile, file)
	if err!=nil{
		respondWithError(w, http.StatusInternalServerError, "could not copy contents to file", err)
		return
	}

	tmpFile.Seek(0, io.SeekStart)

	aspectRatio, err := getVideoAspectRatio(tmpFile.Name())
	if err!=nil{
		respondWithError(w, http.StatusInternalServerError, "could not get aspect ratio", err)
	}
	log.Printf("aspectRatio: %s", aspectRatio)

	bucketString := "tubely-1123581321"
	b := make([]byte, 32)
	rand.Read(b)
	keyRand := base64.RawURLEncoding.EncodeToString(b)
	var key string
	if aspectRatio == "16:9"{
		key = "landscape/" + keyRand
	}else if aspectRatio == "9:16"{
		key = "portrait/" + keyRand
	}else{
		key = "other/" + keyRand
	}

	putObjectInput := &s3.PutObjectInput{
		Bucket: &bucketString,
		Key: &key,
		Body: tmpFile,
		ContentType: &mediaType,
	}

	_, err = cfg.s3Client.PutObject(r.Context(), putObjectInput)	
	if err!= nil{
		respondWithError(w, http.StatusInternalServerError, "could not put object in bucket", err)
		return
	}

	videoURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", cfg.s3Bucket, cfg.s3Region, key)
	video.VideoURL = &videoURL

	err = cfg.db.UpdateVideo(video)
	if err!= nil{
		respondWithError(w, http.StatusInternalServerError, "could not update video metadata", err)
		return
	}

	respondWithJSON(w, http.StatusOK, video)


}
