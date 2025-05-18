package handlers

import (
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"

	"go.temporal.io/sdk/client"

	"github.com/go-microfrontend/images-provider/internal/domain"
)

const (
	ImageEndpoint = "GET /{bucket}/{object}"
)

type ImageHandler struct {
	client  client.Client
	options *client.StartWorkflowOptions
}

func NewImage(client client.Client, options *client.StartWorkflowOptions) *ImageHandler {
	return &ImageHandler{client: client, options: options}
}

func (h *ImageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	params := domain.GetFileParams{
		BucketName: r.PathValue("bucket"),
		ObjectName: r.PathValue("object"),
	}

	we, _ := h.client.ExecuteWorkflow(r.Context(), *h.options, "GetImage", params)

	var imageURL string
	we.Get(r.Context(), &imageURL)
	slog.Info("Generated URL", "url", imageURL)

	target, _ := url.Parse(imageURL)
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Director = func(req *http.Request) {
		req.URL = target
		req.Host = target.Host

		req.Header.Del("X-Forwarded-For")
		req.Header.Del("X-Forwarded-Host")
		req.Header.Del("Forwarded")

		req.Header.Set("Host", target.Host)
		req.Header.Set("User-Agent", "images-provider-proxy")
	}

	proxy.ServeHTTP(w, r)
}
