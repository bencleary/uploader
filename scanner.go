package uploader

import "context"

type IssueDetected struct{}

type ScannerService interface {
	Scan(ctx context.Context, filePath string) IssueDetected
}
