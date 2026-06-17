package main

import (
	"fmt"
	"net/http"
	"io"
	"errors"

	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUploadThumbnail(w http.ResponseWriter, r *http.Request) {
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


	fmt.Println("uploading thumbnail for video", videoID, "by user", userID)

	// TODO: implement the upload here
	const maxMemory = 10 << 20
	r.ParseMultipartForm(maxMemory)

	file, header, err := r.FormFile("thumbnail")
	if err!= nil{
		respondWithError(w, http.StatusBadRequest, "unable to parse form file", err)
		return
	}
	defer file.Close()

	mediaType := header.Header.Get("Content-Type")
	if mediaType == ""{
		respondWithError(w, http.StatusBadRequest, "missing content type for thumbnail", nil)
		return
	}

	dat, err := io.ReadAll(file)
	if err!=nil{
		respondWithError(w, http.StatusInternalServerError, "unable to read data from file", err)
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


	thumbnailNew := thumbnail{
		data: dat,
		mediaType: mediaType,
	}

	videoThumbnails[video.ID] = thumbnailNew

	

	url := fmt.Sprintf("http://localhost:8091/api/thumbnails/%s", video.ID.String())
	video.ThumbnailURL = &url


	err = cfg.db.UpdateVideo(video)
	if err!= nil{
		delete(videoThumbnails, videoID)
		respondWithError(w, http.StatusInternalServerError, "could not update video metadata", err)
		return
	}

	respondWithJSON(w, http.StatusOK, video)

}
