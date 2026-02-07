package handler

import (
	"net/http"

	"github.com/davidjchavez/url-shortener/internal/model"
	"github.com/davidjchavez/url-shortener/internal/service"
	"github.com/labstack/echo/v5"
)

type URLHandler struct {
	service *service.URLService
}

func NewURLHandler(service *service.URLService) *URLHandler {
	return &URLHandler{
		service: service,
	}
}

func (h *URLHandler) CreateShortURL(c *echo.Context) error {
	var req model.CreateURLRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	res, err := h.service.CreateShortURL(req.URL)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, res)
}

func (h *URLHandler) Redirect(c *echo.Context) error {
	code := c.Param("code")
	if code == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "code is required",
		})
	}

	originalURL, err := h.service.GetOriginalURL(code)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"message": "original url not found",
		})
	}

	return c.Redirect(http.StatusMovedPermanently, originalURL)
}

func (h *URLHandler) GetStats(c *echo.Context) error {
	code := c.Param("code")
	if code == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "code is required",
		})
	}

	res, err := h.service.GetStats(code)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"message": "code not found",
		})
	}

	return c.JSON(http.StatusOK, res)
}
