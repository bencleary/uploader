package http

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/bencleary/uploader"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const (
	MAX_IMAGE_WIDTH = 2000
	PREVIEW_WIDTH   = 320
)

func (s *Server) upload(c echo.Context) error {
	// Encryption Key should be header
	key := c.Request().Header.Get("key")

	if key == "" {
		return uploader.Errorf(uploader.INVALID, "")
	}

	// Read file
	file, err := c.FormFile("file")
	if err != nil {
		return echo.NewHTTPError(http.StatusNoContent, "Reading file upload failed")
	}

	attachment, err := s.storage.Hold(c.Request().Context(), file)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Storing uploaded file failed")
	}

	err = attachment.CopyFileToPath(attachment.CreatePreviewLocalPath())

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Storing uploaded file failed")
	}

	if s.scaler.Supported(attachment.MimeType) {
		err := s.scaler.Scale(c.Request().Context(), attachment.LocalPath, MAX_IMAGE_WIDTH, attachment.MimeType)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Storing uploaded file failed")
		}
	} else {
		return echo.NewHTTPError(http.StatusBadRequest, "File type is not supported")
	}

	err = s.preview.Generate(c.Request().Context(), attachment, PREVIEW_WIDTH)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Processing uploaded file has failed, please try again.")
	}

	err = s.storage.Upload(c.Request().Context(), attachment, key)

	if err != nil {
		// writing files has failed
		return echo.NewHTTPError(http.StatusInternalServerError, "Processing uploaded file has failed, please try again.")
	}

	err = s.filer.Record(attachment)

	if err != nil {
		// writing files has failed
		return echo.NewHTTPError(http.StatusInternalServerError, "Processing uploaded file has failed, please try again.")
	}

	tempURL := url.URL{
		Scheme: "http",
		Host:   c.Request().Host,
		Path:   "/file/" + attachment.UID.String(),
	}
	tempPreviewURL := tempURL.String() + "?preview=true"

	upload, err := uploader.NewUpload(attachment, tempPreviewURL, tempURL.String())

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, upload)
}

func (s *Server) download(c echo.Context) error {
	// Retrieve the file UID from the request parameter.
	uid := c.Param("uid")
	key := c.Request().Header.Get("key")
	preview := c.QueryParam("preview")

	var previewValue bool
	previewValue, err := strconv.ParseBool(preview)

	if err != nil {
		previewValue = false
	}

	// TODO: Implement validation for UUID as Must Parse panics...
	attachment, err := s.filer.Fetch(uuid.MustParse(uid))

	if err != nil {
		return err
	}

	// // Load the attachment by UID.
	decrypted, err := s.storage.Download(c.Request().Context(), attachment, previewValue, key)
	if err != nil {
		// return uploader.Errorf(uploader.INVALID, "Decryption failed")
		return echo.NewHTTPError(echo.ErrBadRequest.Code, "Decryption failed")
	}

	contentType := attachment.MimeType
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Stream the file for download.
	return c.Stream(http.StatusOK, contentType, decrypted)
}
