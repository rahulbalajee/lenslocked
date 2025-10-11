package models

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Image struct {
	GalleryID int
	Path      string
	Filename  string
}

type ImageService struct {
	// ImagesDir is used to tell the GalleryService where to store and locate
	// images. If not set, the GalleryService will default to using the "images"
	// directory.
	Dir string

	// Image extensions that are allowed to be uploaded by the user
	// If this is not set will fetch a default value from this package
	Extensions []string

	// Image Content-Type that are allowed when user uploaded an image
	// If this is not set will fetch a default value from this package
	ContentTypes []string
}

func (is *ImageService) Image(galleryId int, filename string) (Image, error) {
	imagePath := filepath.Join(is.imagesDir(galleryId), filename)

	_, err := os.Stat(imagePath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return Image{}, ErrNotFound
		}
		return Image{}, fmt.Errorf("query image: %w", err)
	}

	return Image{
		Filename:  filename,
		GalleryID: galleryId,
		Path:      imagePath,
	}, nil
}

func (is *ImageService) Images(galleryID int) ([]Image, error) {
	globPattern := filepath.Join(is.imagesDir(galleryID), "*")

	allFiles, err := filepath.Glob(globPattern)
	if err != nil {
		return nil, fmt.Errorf("retrieving gallery images: %w", err)
	}

	extensions := is.defaultExtensions()
	if is.Extensions != nil {
		extensions = is.Extensions
	}

	var images []Image
	for _, file := range allFiles {
		if hasExtension(file, extensions) {
			images = append(images, Image{
				GalleryID: galleryID,
				Path:      file,
				Filename:  filepath.Base(file),
			})
		}
	}

	return images, nil
}

func (is *ImageService) CreateImage(galleryID int, filename string, contents io.ReadSeeker) error {
	contentType := is.defaultImageContentsType()
	if is.ContentTypes != nil {
		contentType = is.ContentTypes
	}

	err := checkContentType(contents, contentType)
	if err != nil {
		return fmt.Errorf("creating image %v: %w", filename, err)
	}

	extensions := is.defaultExtensions()
	if is.Extensions != nil {
		extensions = is.Extensions
	}
	err = checkExtension(filename, extensions)
	if err != nil {
		return fmt.Errorf("creating image %v: %w", filename, err)
	}

	galleryDir := is.imagesDir(galleryID)
	err = os.MkdirAll(galleryDir, 0755)
	if err != nil {
		return fmt.Errorf("creating gallery-%d images directory: %w", galleryID, err)
	}

	imagePath := filepath.Join(galleryDir, filename)
	dst, err := os.Create(imagePath)
	if err != nil {
		return fmt.Errorf("creating image file: %w", err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, contents)
	if err != nil {
		return fmt.Errorf("copying contents to image: %w", err)
	}

	return nil
}

func (is *ImageService) DeleteImage(galleryID int, filename string) error {
	image, err := is.Image(galleryID, filename)
	if err != nil {
		return fmt.Errorf("deleting image: %w", err)
	}

	err = os.Remove(image.Path)
	if err != nil {
		return fmt.Errorf("deleting image: %w", err)
	}

	return nil
}

// DeleteAllGalleryImages deletes all images for a given gallery by removing the gallery directory.
// It returns nil if the directory doesn't exist (idempotent).
func (is *ImageService) DeleteAllGalleryImages(galleryID int) error {
	imagesDir := is.imagesDir(galleryID)
	err := os.RemoveAll(imagesDir)
	if err != nil {
		return fmt.Errorf("deleting all gallery images: %w", err)
	}
	return nil
}

func (is *ImageService) defaultExtensions() []string {
	return []string{".png", ".jpg", ".jpeg", ".gif"}
}

func (is *ImageService) defaultImageContentsType() []string {
	return []string{"image/png", "image/jpg", "image/jpeg", "image/gif"}
}

func (is *ImageService) imagesDir(id int) string {
	imagesDir := is.Dir
	if imagesDir == "" {
		imagesDir = "images"
	}
	return filepath.Join(imagesDir, fmt.Sprintf("gallery-%d", id))
}

func hasExtension(file string, extensions []string) bool {
	fileExt := strings.ToLower(filepath.Ext(file))
	for _, ext := range extensions {
		if fileExt == strings.ToLower(ext) {
			return true
		}
	}
	return false
}
