package utils

import (
	"strings"

	"github.com/gosimple/slug"
)

func GenerateSlug(nama string) string {
	return slug.Make(strings.ToLower(nama))
}
