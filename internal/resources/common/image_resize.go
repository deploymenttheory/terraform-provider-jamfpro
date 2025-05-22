package common

import (
	"fmt"
	"image"
	"image/gif"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/nfnt/resize"
)

// ResizeImage resizes an image to fit within maxWidth x maxHeight while maintaining aspect ratio
// returns the path to the resized image file and any error
func ResizeImage(img image.Image, format string, originalPath string, maxWidth, maxHeight uint) (string, error) {
	// Get original dimensions
	bounds := img.Bounds()
	origWidth := bounds.Dx()
	origHeight := bounds.Dy()

	// Calculate scaling to maintain aspect ratio while fitting within maxWidth x maxHeight
	ratio := float64(origWidth) / float64(origHeight)
	var newWidth, newHeight uint
	if ratio > 1 { // wider than tall
		newWidth = maxWidth
		newHeight = uint(float64(maxWidth) / ratio)
	} else { // taller than wide
		newHeight = maxHeight
		newWidth = uint(float64(maxHeight) * ratio)
	}

	// Resize the image while preserving aspect ratio
	resizedImg := resize.Resize(newWidth, newHeight, img, resize.Lanczos3)

	// Create a temporary file for the resized image
	ext := strings.ToLower(filepath.Ext(originalPath))
	tempDir := os.TempDir()
	tempFilename := filepath.Join(tempDir, "resized_"+filepath.Base(originalPath))

	tempFile, err := os.Create(tempFilename)
	if err != nil {
		//nolint:err113 // https://github.com/deploymenttheory/terraform-provider-jamfpro/issues/650
		return "", fmt.Errorf("failed to create temporary file for resized image: %v", err)
	}
	defer tempFile.Close()

	// Encode the image back to the temporary file
	if format == "png" || ext == ".png" {
		err = png.Encode(tempFile, resizedImg)
	} else if format == "gif" || ext == ".gif" {
		err = gif.Encode(tempFile, resizedImg, nil)
	} else {
		// Default to PNG if format is unexpected
		err = png.Encode(tempFile, resizedImg)
	}

	if err != nil {
		//nolint:err113 // https://github.com/deploymenttheory/terraform-provider-jamfpro/issues/650
		return "", fmt.Errorf("failed to encode processed image: %v", err)
	}

	return tempFilename, nil
}
