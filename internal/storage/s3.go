package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/bencleary/uploader"
)

var _ uploader.StorageService = (*S3Storage)(nil)

type S3Options struct {
	Endpoint       string
	Bucket         string
	Region         string
	Prefix         string
	ForcePathStyle bool
	AccessKeyID    string
	SecretKey      string
}

type S3Storage struct {
	staging    *LocalStorage
	client     *s3.Client
	options    *S3Options
	encryption uploader.EncryptionService
}

func NewS3Storage(options *S3Options, encryption uploader.EncryptionService) *S3Storage {

	if options == nil {
		return nil
	}

	if encryption == nil {
		return nil
	}

	staging := NewLocalStorage(options.Prefix, options.Prefix, encryption)
	if staging == nil {
		return nil
	}

	// Build config options
	configOpts := []func(*config.LoadOptions) error{
		config.WithRegion(options.Region),
	}

	// Add static credentials if provided
	if options.AccessKeyID != "" && options.SecretKey != "" {
		configOpts = append(configOpts, config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(options.AccessKeyID, options.SecretKey, ""),
		))
	}

	cfg, err := config.LoadDefaultConfig(context.Background(), configOpts...)
	if err != nil {
		return nil
	}

	// Configure custom endpoint and path style if provided
	clientOptions := func(o *s3.Options) {
		if options.ForcePathStyle {
			o.UsePathStyle = true
		}
		if options.Endpoint != "" {
			o.BaseEndpoint = aws.String(options.Endpoint)
		}
	}

	client := s3.NewFromConfig(cfg, clientOptions)

	return &S3Storage{
		options:    options,
		encryption: encryption,
		staging:    staging,
		client:     client,
	}
}

func (s *S3Storage) Initialise(ctx context.Context) error {
	if err := s.staging.Initialise(ctx); err != nil {
		return err
	}

	return nil
}

func (s *S3Storage) Hold(ctx context.Context, attachment *multipart.FileHeader) (*uploader.Attachment, error) {
	return s.staging.Hold(ctx, attachment)
}

func (s *S3Storage) Upload(ctx context.Context, attachment *uploader.Attachment, key string) error {
	if attachment == nil {
		return uploader.Errorf(uploader.INVALID, "attachment is required")
	}

	// Upload main file
	if err := s.uploadEncryptedFile(ctx, attachment.LocalPath, attachment.UID.String(), false, key); err != nil {
		return err
	}

	// Upload preview file if it exists
	if attachment.PreviewLocalPath != "" {
		if err := s.uploadEncryptedFile(ctx, attachment.PreviewLocalPath, attachment.UID.String(), true, key); err != nil {
			return err
		}
	}

	return nil
}

// uploadEncryptedFile encrypts a local file and uploads it to S3
func (s *S3Storage) uploadEncryptedFile(ctx context.Context, filePath, uid string, isPreview bool, encryptionKey string) error {
	// Open the source file
	source, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer source.Close()

	// Encrypt the file
	encrypted, err := s.encryption.EncryptStream(ctx, source, encryptionKey)
	if err != nil {
		return err
	}
	defer encrypted.Close()

	// Read encrypted data into a buffer to make it seekable for S3 retries
	encryptedData, err := io.ReadAll(encrypted)
	if err != nil {
		return err
	}

	// Construct the S3 object key
	objectKey := s.objectKey(uid, isPreview)

	// Upload to S3
	_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.options.Bucket),
		Key:    aws.String(objectKey),
		Body:   bytes.NewReader(encryptedData),
	})
	if err != nil {
		return err
	}

	return nil
}

// objectKey constructs the S3 object key for a file
func (s *S3Storage) objectKey(uid string, isPreview bool) string {
	var key string
	if s.options.Prefix != "" {
		key = strings.TrimSuffix(s.options.Prefix, "/") + "/"
	}

	if isPreview {
		key = fmt.Sprintf("%s%s.preview.enc", key, uid)
	} else {
		key = fmt.Sprintf("%s%s.enc", key, uid)
	}

	return key
}

func (s *S3Storage) Download(ctx context.Context, attachment *uploader.Attachment, preview bool, key string) (io.ReadCloser, error) {
	if attachment == nil {
		return nil, uploader.Errorf(uploader.INVALID, "attachment is required")
	}

	// Construct the S3 object key
	objectKey := s.objectKey(attachment.UID.String(), preview)

	// Download from S3
	result, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.options.Bucket),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return nil, err
	}

	// Decrypt the stream
	decrypted, err := s.encryption.DecryptStream(ctx, result.Body, key)
	if err != nil {
		_ = result.Body.Close()
		return nil, err
	}

	return newChainedReadCloser(decrypted, decrypted, result.Body), nil
}

func (s *S3Storage) Delete(ctx context.Context, attachmentUID string) error {
	// Delete main file
	mainKey := s.objectKey(attachmentUID, false)
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.options.Bucket),
		Key:    aws.String(mainKey),
	})
	if err != nil {
		return err
	}

	// Delete preview file (ignore error if it doesn't exist)
	previewKey := s.objectKey(attachmentUID, true)
	_, _ = s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.options.Bucket),
		Key:    aws.String(previewKey),
	})

	return nil
}
