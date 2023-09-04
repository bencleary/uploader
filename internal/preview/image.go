package preview

import (
	"context"

	"github.com/bencleary/uploader"
)

const (
	MAX_PREVIEW_WIDTH = 320
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

func (i *ImagePreviewGenerator) Generate(ctx context.Context, attachment *uploader.Attachment) error {
	return nil
	// format, err := resize.GetImageFormat(attachment.MimeType)
	// if err != nil {
	// 	return err
	// }
	// resizer := resize.NewResizerService()
	// previewPath := strings.Join([]string{attachment.VaultPath, "preview"}, ".")
	// return resizer.Resize(attachment.VaultPath, previewPath, MAX_PREVIEW_WIDTH, format)
}
