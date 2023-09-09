package preview

import (
	"context"

	"github.com/bencleary/uploader"
)

var _ uploader.PreviewGeneratorService = (*ImagePreviewGenerator)(nil)

type ImagePreviewGenerator struct {
	Scaler uploader.ScalerService
}

func NewImagePreviewGenerator(scaler uploader.ScalerService) *ImagePreviewGenerator {
	return &ImagePreviewGenerator{
		scaler,
	}
}

func (i *ImagePreviewGenerator) Generate(ctx context.Context, attachment *uploader.Attachment, previewWidth int) error {
	return i.Scaler.Scale(ctx, attachment.PreviewLocalPath, previewWidth, attachment.MimeType)
}
