package uploader

import "context"

type ScalerService interface {
	Scale(ctx context.Context, filePath string, maxWidth int, mimeType string) error
	Supported(mimeType string) bool
}
