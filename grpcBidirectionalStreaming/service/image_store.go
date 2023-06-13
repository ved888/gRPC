package service

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"os"
	"sync"
)

// ImageStore is an interface to store laptop image
type ImageStore interface {
	// Save saves a new laptop image to the store
	Save(laptopId string, imageType string, imageData bytes.Buffer) (string, error)
}

// DiskImageStore stores images on disk and its info on memory
type DiskImageStore struct {
	mutex       sync.RWMutex
	imageFolder string
	images      map[string]*ImageInfo
}

// ImageInfo contains information of the laptop image
type ImageInfo struct {
	LaptopId string
	Type     string
	Path     string
}

// NewDiskImageStore returns a new DiskImagStore
func NewDiskImageStore(imageFolder string) *DiskImageStore {
	return &DiskImageStore{
		imageFolder: imageFolder,
		images:      make(map[string]*ImageInfo),
	}
}

// Save saves a new laptop image to the store
func (store *DiskImageStore) Save(
	laptopId string,
	imageType string,
	imageData bytes.Buffer,
) (string, error) {
	imageId, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("cannot generate image id:%v", err)
	}

	imagePath := fmt.Sprintf("%s/%s%s", store.imageFolder, imageId, imageType)

	file, err := os.Create(imagePath)
	if err != nil {
		return "", fmt.Errorf("cannot create image to file:%v", err)
	}

	_, err = imageData.WriteTo(file)
	if err != nil {
		return "", fmt.Errorf("cannpt write image to file:%v", err)
	}

	store.mutex.Lock()
	defer store.mutex.Unlock()

	store.images[imageId.String()] = &ImageInfo{
		LaptopId: laptopId,
		Type:     imageType,
		Path:     imagePath,
	}
	return imageId.String(), nil
}
