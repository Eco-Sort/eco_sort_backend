package image

import (
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/Eco-Sort/eco_sort_backend/domain"
	"github.com/go-resty/resty/v2"
)

type imageService struct {
	contextTimeout  time.Duration
	imageBucketRepo domain.ImageRepository
	Resty           *resty.Client
}

func NewImageService(t time.Duration, imageBucketRepo domain.ImageRepository, Resty *resty.Client) domain.ImageService {
	return &imageService{
		contextTimeout:  t,
		imageBucketRepo: imageBucketRepo,
		Resty:           Resty,
	}
}

func (i *imageService) UploadImage(file *os.File, fileName string, userId uint) (*string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), i.contextTimeout)
	defer cancel()

	timeNow := time.Now().Unix()
	nameWithoutExtension := fileName[:len(fileName)-len(filepath.Ext(fileName))]
	dir := fmt.Sprintf("%s_%d%s", nameWithoutExtension, timeNow, filepath.Ext(fileName))
	res, err := i.imageBucketRepo.UploadImage(ctx, file, userId, dir)
	if err != nil {
		return nil, err
	}
	file.Close()

	err = os.Remove(file.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to delete file: %w", err)
	}

	return &res, nil
}

func (i *imageService) ProcessImage(file *multipart.FileHeader) (domain.MLResponse, error) {
	fileData, err := file.Open()
	if err != nil {
		return domain.MLResponse{}, err
	}
	defer fileData.Close()

	response, err := i.Resty.
		SetBaseURL(domain.MLUrl).
		R().
		SetHeader("Content-Type", "application/json; charset=utf-8").
		SetFileReader("file", file.Filename, fileData).
		Post("/")
	if err != nil {
		return domain.MLResponse{}, err
	}
	var res domain.MLResponse
	err = json.Unmarshal(response.Body(), &res)
	if err != nil {
		return domain.MLResponse{}, err
	}
	return res, nil
}
