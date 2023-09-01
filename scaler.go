package uploader

type ScalerService interface {
	Scale(filePath string, maxWidth int) error
	Supported(mimeType string) bool
}
