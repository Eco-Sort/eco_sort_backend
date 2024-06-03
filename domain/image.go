package domain

import (
	"context"
	"mime/multipart"
)

type ImageRepository interface {
	UploadImage(ctx context.Context, file multipart.File, userId uint, fileName string) (string, error)
}

type ImageService interface {
	UploadImage(file *multipart.FileHeader, userId uint) (*string, error)
}
