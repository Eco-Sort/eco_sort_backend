package domain

import (
	"context"
	"mime/multipart"
	"os"
)

var (
	MLUrl   string
	MLImage string
)

type MLResponse struct {
	Filename  string     `json:"filename"`
	Filepath  string     `json:"filepath"`
	Result    []MLResult `json:"result"`
	TimeTaken string     `json:"time_taken"`
}

type MLResult struct {
	Bbox  []float64 `json:"bbox"`
	Class float64   `json:"class"`
	Score float64   `json:"score"`
}

type ImageRepository interface {
	UploadImage(ctx context.Context, file multipart.File, userId uint, fileName string) (string, error)
}

type ImageService interface {
	UploadImage(file *os.File, fileName string, userId uint) (*string, error)
	ProcessImage(file *multipart.FileHeader) (MLResponse, error)
}
