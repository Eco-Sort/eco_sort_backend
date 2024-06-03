package gcstorage

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"

	"cloud.google.com/go/storage"
	"github.com/Eco-Sort/eco_sort_backend/domain"
)

var _imageBucket = "c241-ps362-imageupload"
var imageUrl = "https://storage.googleapis.com/c241-ps362-imageupload/%s"

type gcStorageRepository struct {
	gcStorage *storage.Client
}

func NewGcStorageRepository(gcStorage *storage.Client) domain.ImageRepository {
	return &gcStorageRepository{
		gcStorage: gcStorage,
	}
}

func (g *gcStorageRepository) UploadImage(ctx context.Context, file multipart.File, userId uint, fileName string) (string, error) {
	bucket := g.gcStorage.Bucket(_imageBucket)
	objectName := fmt.Sprintf("user_%d/%s", userId, fileName)
	object := bucket.Object(objectName)
	writer := object.NewWriter(ctx)

	if _, err := io.Copy(writer, file); err != nil {
		return "", err
	}

	if err := writer.Close(); err != nil {
		return "", err
	}
	return fmt.Sprintf(imageUrl, objectName), nil
}
