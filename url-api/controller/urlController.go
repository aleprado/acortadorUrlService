package controller

import (
	"acortadorUrlService/components/config"
	"acortadorUrlService/components/logger"
	"acortadorUrlService/url-api/service"
	"encoding/json"
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

func (c *UrlController) MountIn(r chi.Router) {
	r.Route("/shorten", func(r chi.Router) {
		r.Get("/", c.GetShortUrl)
		r.Delete("/", c.DeleteShortUrl)
		r.Get("/{hash}", c.ResolveShortUrl)
	})
}

func (c *UrlController) GetShortUrl(w http.ResponseWriter, r *http.Request) {
	originalURL := r.URL.Query().Get("url")
	if originalURL == "" {
		metrics.PutCountMetric(metrics.MetricGetShortUrlMissingParam, 1)
		http.Error(w, "Missing 'url' query parameter", http.StatusBadRequest)
		logger.LogError("Missing URL parameter", nil)
		return
	}

	shortened, err := c.Shortener.ShortenOrFetch(r.Context(), originalURL)
	if err != nil {
		metrics.PutCountMetric(metrics.MetricGetShortUrlError, 1)
		http.Error(w, "Failed to shorten URL", http.StatusInternalServerError)
		logger.LogError("ShortenOrFetch error", err)
		return
	}

	if !shortened.CreatedAt.IsZero() {
		metrics.PutCountMetric(metrics.MetricGetShortUrlCreatedNew, 1)
	} else {
		metrics.PutCountMetric(metrics.MetricGetShortUrlFoundExisting, 1)
	}

	metrics.PutCountMetric(metrics.MetricGetShortUrlSuccess, 1)
	response := map[string]string{
		"short_url": c.Config.BaseURL + "/" + shortened.Hash,
		"original":  shortened.Original,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (c *UrlController) DeleteShortUrl(w http.ResponseWriter, r *http.Request) {
	originalURL := r.URL.Query().Get("url")
	if originalURL == "" {
		metrics.PutCountMetric(metrics.MetricDeleteShortUrlMissingParam, 1)
		http.Error(w, "Missing 'url' query parameter", http.StatusBadRequest)
		logger.LogError("Missing URL parameter", nil)
		return
	}

	err := c.Shortener.Delete(r.Context(), originalURL)
	if err != nil {
		metrics.PutCountMetric(metrics.MetricDeleteShortUrlError, 1)
		http.Error(w, "Failed to delete short URL", http.StatusInternalServerError)
		logger.LogError("Delete error", err)
		return
	}

	metrics.PutCountMetric(metrics.MetricDeleteShortUrlSuccess, 1)
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
