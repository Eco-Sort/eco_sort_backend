package image

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/Eco-Sort/eco_sort_backend/domain"
)

type imageService struct {
	contextTimeout  time.Duration
	imageBucketRepo domain.ImageRepository
}

func NewImageService(t time.Duration, imageBucketRepo domain.ImageRepository) domain.ImageService {
	return &imageService{
		contextTimeout:  t,
		imageBucketRepo: imageBucketRepo,
	}
}

func (i *imageService) UploadImage(file *multipart.FileHeader, userId uint) (*string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), i.contextTimeout)
	defer cancel()
	fileContents, err := file.Open()
	if err != nil {
		return nil, err
	}
	timeNow := time.Now().Unix()
	fileName := file.Filename
	nameWithoutExtension := fileName[:len(fileName)-len(filepath.Ext(fileName))]
	dir := fmt.Sprintf("%s_%d%s", nameWithoutExtension, timeNow, filepath.Ext(fileName))
	res, err := i.imageBucketRepo.UploadImage(ctx, fileContents, userId, dir)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
