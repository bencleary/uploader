package uploader

type IssueDetected struct{}

type ScannerService interface {
	Scan(filePath string) IssueDetected
}
