package http

import (
	"fmt"
	"net/http"

	"github.com/bencleary/uploader"
	"github.com/labstack/echo/v4"
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

	attachment := uploader.NewAttachment(file, 1)

	if err != nil {
		return err
	}

	// vaultPath, err := s.storage.Hold(attachment)
	// filePath := strings.Join([]string{vaultPath, attachment.FileName}, "/")
	// attachment.VaultPath = filePath

	// // Virus scanning

	// // err := s.filer.Scan(filePath)
	// if err != nil {
	// 	// virus validation failed...
	// 	return err
	// }

	// // if strings.Contains(attachment.MimeType, IMAGE_MIMETYPE) {
	// // 	r := resize.NewResizerService()

	// // 	format, err := resize.GetImageFormat(attachment.MimeType)
	// // 	if err != nil {
	// // 		return err
	// // 	}
	// // 	r.Resize(filePath, filePath, MAX_IMAGE_WIDTH, format)
	// // }

	// err = s.preview.Generate(attachment)

	// if err != nil {
	// 	return err
	// }

	// err = s.storage.Save(filePath, attachment, key)
	// // err = s.filer.Store(filePath, attachment, key)

	// if err != nil {
	// 	// writing files has failed
	// 	return err
	// }

	return c.JSON(http.StatusOK, map[string]string{"uid": attachment.UID.String()})
}

// func (s *Server) preview(c echo.Context) error {
// 	// fetch attachment and return preview url
// 	c.Response().Header().Set("Content-Type", "")
// 	return nil
// }
