package controller

import (
	"acortadorUrlService/components/config"
	"acortadorUrlService/components/logger"
	"acortadorUrlService/url-api/service"
	"encoding/json"
	"io"
	"net/http"

	"acortadorUrlService/components/metrics"

	"github.com/go-chi/chi/v5"
)

type UrlController struct {
	Shortener *service.UrlShortener
	Config    *config.AppConfig
}

func NewUrlController(s *service.UrlShortener, cfg *config.AppConfig) *UrlController {
	return &UrlController{Shortener: s, Config: cfg}
}

type CreateShortUrlRequest struct {
	URL string `json:"url"`
}

func (c *UrlController) MountIn(r chi.Router) {
	r.Post("/", c.CreateShortUrl)
	r.Delete("/", c.DeleteShortUrl)
	r.Get("/{hash}", c.ResolveShortUrl)
}

func (c *UrlController) CreateShortUrl(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		metrics.PutCountMetric(metrics.MetricPostShortUrlError, 1)
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		logger.LogError("Failed to read request body", err)
		return
	}
	defer r.Body.Close()

	var req CreateShortUrlRequest
	if err := json.Unmarshal(body, &req); err != nil {
		metrics.PutCountMetric(metrics.MetricPostShortUrlError, 1)
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		logger.LogError("Invalid JSON format", err)
		return
	}

	if req.URL == "" {
		metrics.PutCountMetric(metrics.MetricPostShortUrlMissingParam, 1)
		http.Error(w, "Missing 'url' in request body", http.StatusBadRequest)
		logger.LogError("Missing URL in request body", nil)
		return
	}

	shortened, err := c.Shortener.ShortenOrFetch(r.Context(), req.URL)
	if err != nil {
		metrics.PutCountMetric(metrics.MetricPostShortUrlError, 1)
		http.Error(w, "Failed to shorten URL", http.StatusInternalServerError)
		logger.LogError("ShortenOrFetch error", err)
		return
	}

	if !shortened.CreatedAt.IsZero() {
		metrics.PutCountMetric(metrics.MetricPostShortUrlCreatedNew, 1)
	} else {
		metrics.PutCountMetric(metrics.MetricPostShortUrlFoundExisting, 1)
	}

	metrics.PutCountMetric(metrics.MetricPostShortUrlSuccess, 1)
	response := map[string]string{
		"short_url": c.Config.BaseURL + "/" + shortened.Hash,
		"original":  shortened.Original,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (c *UrlController) DeleteShortUrl(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		metrics.PutCountMetric(metrics.MetricDeleteShortUrlError, 1)
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		logger.LogError("Failed to read request body", err)
		return
	}
	defer r.Body.Close()

	var req CreateShortUrlRequest
	if err := json.Unmarshal(body, &req); err != nil {
		metrics.PutCountMetric(metrics.MetricDeleteShortUrlError, 1)
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		logger.LogError("Invalid JSON format", err)
		return
	}

	if req.URL == "" {
		metrics.PutCountMetric(metrics.MetricDeleteShortUrlMissingParam, 1)
		http.Error(w, "Missing 'url' in request body", http.StatusBadRequest)
		logger.LogError("Missing URL in request body", nil)
		return
	}

	err = c.Shortener.Delete(r.Context(), req.URL)
	if err != nil {
		metrics.PutCountMetric(metrics.MetricDeleteShortUrlError, 1)
		http.Error(w, "Failed to delete short URL", http.StatusInternalServerError)
		logger.LogError("Delete error", err)
		return
	}

	metrics.PutCountMetric(metrics.MetricDeleteShortUrlSuccess, 1)
	w.WriteHeader(http.StatusNoContent)
}

func (c *UrlController) ResolveShortUrl(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")
	ctx := r.Context()

	originalUrl, err := c.Shortener.GetOriginalUrl(ctx, hash)
	if err != nil {
		metrics.PutCountMetric(metrics.MetricResolveShortUrlError, 1)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		logger.LogError("GetOriginalUrl error", err)
		return
	}

	if originalUrl == "" {
		metrics.PutCountMetric(metrics.MetricResolveShortUrlNotFound, 1)
		http.Error(w, "url not found", http.StatusNotFound)
		return
	}

	metrics.PutCountMetric(metrics.MetricResolveShortUrlSuccess, 1)
	// here is the magic to redirect to the long url recovered from the database
	http.Redirect(w, r, originalUrl, http.StatusFound)
}
