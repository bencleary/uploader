package uploader

import "context"

type ScalerService interface {
	Scale(ctx context.Context, attachment *Attachment, maxWidth int) error
	Supported(mimeType string) bool
}
