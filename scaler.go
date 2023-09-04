package uploader

import "context"

type ScalerService interface {
	Scale(ctx context.Context, filePath string, maxWidth int) error
	Supported(mimeType string) bool
}
