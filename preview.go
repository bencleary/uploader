package uploader

import "context"

type PreviewGeneratorService interface {
	Generate(ctx context.Context, attachment *Attachment, previewWidth int) error
}

type PreviewService struct {
	handlers map[string]PreviewGeneratorService
}

func NewPreviewService() *PreviewService {
	return &PreviewService{
		handlers: make(map[string]PreviewGeneratorService),
	}
}

func (p *PreviewService) Register(name string, handler PreviewGeneratorService) {
	p.handlers[name] = handler
}

func (p *PreviewService) Generate(ctx context.Context, attachment *Attachment, previewWidth int) error {
	handler, ok := p.handlers[attachment.MimeType]
	if !ok {
		return nil
	}
	return handler.Generate(ctx, attachment, previewWidth)
}
