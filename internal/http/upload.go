package http

import (
	"fmt"
	"net/http"

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

	fmt.Println(key)

	// Read file
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	attachment, err := s.storage.Hold(context.Background(), file)

	if err != nil {
		return err
	}

	// Virus scanning

	// err := s.filer.Scan(filePath)
	if err != nil {
		// virus validation failed...
		return err
	}

	if s.scaler.Supported(attachment.MimeType) {
		err := s.scaler.Scale(context.Background(), attachment, MAX_IMAGE_WIDTH)
		if err != nil {
			return err
		}
	}

	// TODO - fix overwriting file, create new preview
	err = s.preview.Generate(context.Background(), attachment, PREVIEW_WIDTH)

	if err != nil {
		return err
	}

	err = s.storage.Save(context.Background(), attachment, key)

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

// func (s *Server) preview(c echo.Context) error {
// 	// fetch attachment and return preview url
// 	c.Response().Header().Set("Content-Type", "")
// 	return nil
// }
