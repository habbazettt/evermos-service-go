package services

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"path"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/habbazettt/evermos-service-go/config"
)

// UploadToCloudinary mengunggah gambar langsung ke Cloudinary dari memory
func UploadToCloudinary(file multipart.File) (string, error) {
	ctx := context.Background()

	uploadResult, err := config.CLD.Upload.Upload(ctx, file, uploader.UploadParams{})
	if err != nil {
		return "", err
	}

	return uploadResult.SecureURL, nil
}

// DeleteFromCloudinary menghapus gambar berdasarkan URL
func DeleteFromCloudinary(imageURL string) error {
	// Inisialisasi Cloudinary
	cld, err := cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))
	if err != nil {
		return fmt.Errorf("gagal menginisialisasi Cloudinary: %w", err)
	}

	// Ambil public ID
	publicID := extractPublicID(imageURL)
	if publicID == "" {
		return fmt.Errorf("gagal mendapatkan public ID dari URL")
	}

	// Hapus gambar dari Cloudinary
	_, err = cld.Upload.Destroy(context.Background(), uploader.DestroyParams{
		PublicID: publicID,
	})
	if err != nil {
		return fmt.Errorf("gagal menghapus gambar dari Cloudinary: %w", err)
	}

	return nil
}

// extractPublicID mengubah URL Cloudinary menjadi public ID
func extractPublicID(imageURL string) string {
	// Pisahkan URL berdasarkan "/upload/"
	parts := strings.Split(imageURL, "/upload/")
	if len(parts) < 2 {
		return ""
	}

	// Ambil bagian setelah "/upload/"
	publicIDWithExt := parts[1]

	// Hapus ekstensi file (misal .png, .jpg)
	publicID := strings.TrimSuffix(publicIDWithExt, path.Ext(publicIDWithExt))

	return publicID
}
