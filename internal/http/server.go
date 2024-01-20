package http

import (
	"github.com/bencleary/uploader"
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
