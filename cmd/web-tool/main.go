package main

import (
	"flag"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
)

func main() {
	var addr string
	flag.StringVar(&addr, "addr", "127.0.0.1:8088", "listen address")
	flag.Parse()

	e := echo.New()
	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})
	if err := e.Start(addr); err != nil {
		slog.Error("start server failed: error", err)
	}
}
