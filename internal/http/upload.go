package http

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/context"
)

const (
	MAX_IMAGE_WIDTH = 2000
	PREVIEW_WIDTH   = 320
)

func (s *Server) upload(c echo.Context) error {
	// Encryption Key should be header
	key := c.FormValue("key")

	// Read file
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	attachment, err := s.storage.Hold(context.Background(), file)

	if err != nil {
		return err
	}

	err = attachment.CopyFileToPath(attachment.CreatePreviewLocalPath())

	if err != nil {
		return err
	}

	if s.scaler.Supported(attachment.MimeType) {
		err := s.scaler.Scale(context.Background(), attachment.LocalPath, MAX_IMAGE_WIDTH, attachment.MimeType)
		if err != nil {
			return err
		}
	}

	err = s.preview.Generate(context.Background(), attachment, PREVIEW_WIDTH)

	if err != nil {
		return err
	}

	err = s.storage.Upload(context.Background(), attachment, key)

	if err != nil {
		// writing files has failed
		return err
	}

	err = s.filer.Record(attachment)

	if err != nil {
		// writing files has failed
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"uid": attachment.UID.String()})
}

func (s *Server) download(c echo.Context) error {
	// Retrieve the file UID from the request parameter.
	uid := c.Param("uid")
	key := c.Request().Header.Get("key")

	attachment, err := s.filer.Fetch(uuid.MustParse(uid))

	if err != nil {
		return err
	}

	// // Load the attachment by UID.
	decrypted, err := s.storage.Download(c.Request().Context(), attachment, key)
	if err != nil {
		return err
	}

	contentType := attachment.MimeType
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Stream the file for download.
	return c.Stream(http.StatusOK, contentType, decrypted)
}
