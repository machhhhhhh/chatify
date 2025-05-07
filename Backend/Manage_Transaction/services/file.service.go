package services

import (
	"image"
	"mime/multipart"
	"os"

	"github.com/disintegration/imaging"
	"github.com/gofiber/fiber/v2"
	"github.com/nfnt/resize"
)

func CreateDirectory(folder string) error {
	return os.MkdirAll(folder, os.ModePerm)
}

var MAX_FILE_SIZE int64 = 1000000 * 5 // 5MB

func ResizeImage(ctx *fiber.Ctx, input_image *multipart.FileHeader, path string) error {
	if input_image.Size <= MAX_FILE_SIZE {
		// The file is received, so let's save it
		if err := ctx.SaveFile(input_image, path); err != nil {
			return err
		}
		return nil
	}

	// Open the uploaded image file
	uploaded_file, err := input_image.Open()
	if err != nil {
		return err
	}
	defer uploaded_file.Close()

	// Decode the image
	target_image, _, err := image.Decode(uploaded_file)
	if err != nil {
		return err
	}

	var quality imaging.EncodeOption = imaging.JPEGQuality(80)

	// Resize the image
	resize_image := resize.Resize(0, 0, target_image, resize.Lanczos3)

	// Save the resized image
	err = imaging.Save(resize_image, path, quality)
	if err != nil {
		return err
	}

	return nil
}
