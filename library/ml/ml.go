package ml

import (
	"errors"
	"log"
	"os"

	"github.com/Eco-Sort/eco_sort_backend/domain"
)

func GetMLUrl() {
	url := os.Getenv("ML_URL")
	if url == "" {
		e := errors.New("undefined ML_URL")
		log.Fatal(e)
	}
	domain.MLUrl = url
}

func GetMLImage() {
	image := os.Getenv("ML_IMAGE")
	if image == "" {
		e := errors.New("undefined ML_IMAGE")
		log.Fatal(e)
	}
	domain.MLUrl = image
}
