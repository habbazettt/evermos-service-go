package config

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// Variabel global untuk Cloudinary
var CLD *cloudinary.Cloudinary

// SetupCloudinary menginisialisasi Cloudinary
func SetupCloudinary() error {
	cld, err := cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))
	if err != nil {
		return fmt.Errorf("gagal menginisialisasi Cloudinary: %v", err)
	}
	CLD = cld
	fmt.Println("Cloudinary berhasil diinisialisasi")
	return nil
}

// UploadImage mengunggah gambar ke Cloudinary
func UploadImage(filePath string) (string, error) {
	uploadResult, err := CLD.Upload.Upload(context.Background(), filePath, uploader.UploadParams{})
	if err != nil {
		return "", fmt.Errorf("gagal mengunggah gambar: %v", err)
	}
	return uploadResult.SecureURL, nil
}
