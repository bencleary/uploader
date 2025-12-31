package http

import (
	"net/url"

	"github.com/bencleary/uploader"
	middlewareValidator "github.com/bencleary/uploader/internal/middleware"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	http    *echo.Echo
	filer   uploader.FilerService
	storage uploader.StorageService
	scaler  uploader.ScalerService
	preview *uploader.PreviewService
}

func NewServer(filer uploader.FilerService, storage uploader.StorageService, scaler uploader.ScalerService, preview *uploader.PreviewService) *Server {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOriginFunc: func(origin string) (bool, error) {
			// Non-browser clients (curl, server-to-server) typically send no Origin header.
			if origin == "" {
				return true, nil
			}

			parsed, err := url.Parse(origin)
			if err != nil {
				return false, nil
			}

			switch parsed.Hostname() {
			case "localhost", "127.0.0.1", "::1":
				return true, nil
			default:
				return false, nil
			}
		},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Content-Type", "key"},
		ExposeHeaders: []string{
			"Content-Type",
			"Content-Disposition",
		},
	}))

	e.Use(middlewareValidator.ValidateEncryptionKey)

	server := &Server{
		http:    e,
		filer:   filer,
		storage: storage,
		scaler:  scaler,
		preview: preview,
	}

	server.http.POST("/file/upload", server.upload)
	server.http.GET("/file/:uid", server.download)

	return server
}

func (s *Server) Start() {
	s.http.Logger.Fatal(s.http.Start(":1323"))
}
