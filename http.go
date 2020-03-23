package mecachis

import (
	"context"
	"github.com/sonirico/mecachis/engines"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func (h *Hub) handleAdd(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	ns := ctx.Value("ns").(string)
	key := ctx.Value("key").(string)
	g, created := h.getOrCreateGroup(ns)
	if created {
		// Non-existent group. Check for configuration params
		// params := r.URL.Query()
		// TODO: Hardcoded values
		cap := uint64(2 << 10)
		g.Cap = cap
		g.Ct = engines.LRU
	}
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf(err.Error())
		http.Error(w, "error when reading request buffer", http.StatusInternalServerError)
		return
	}
	if err := g.Add(key, content); err != nil {
		log.Printf(err.Error())
		http.Error(w, "error when writing to response buffer", http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Hub) handleGet(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	ns := ctx.Value("ns").(string)
	g, ok := h.group(ns)
	if !ok {
		http.NotFound(w, r)
		return
	}
	key := ctx.Value("key").(string)
	value, ok := g.Get(key)
	if !ok {
		http.NotFound(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	if _, err := w.Write(value.Clone()); err != nil {
		log.Printf(err.Error())
		http.Error(w, "error when writing to response buffer", http.StatusInternalServerError)
	}
}

func (h *Hub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uri := r.URL.Path
	log.Printf("%s: %s\n", r.Method, uri)
	if !strings.HasPrefix(uri, basePath) {
		http.NotFound(w, r)
		return
	}
	uriParts := strings.SplitN(r.URL.Path[len(basePath):], "/", 2)
	log.Println(uriParts)
	if len(uriParts) < 2 {
		http.NotFound(w, r)
		return
	}
	ns := uriParts[0]
	key := uriParts[1]
	ctx := context.WithValue(r.Context(), "ns", ns)
	ctx = context.WithValue(ctx, "key", key)
	if r.Method == http.MethodGet {
		h.handleGet(ctx, w, r)
		return
	}
	if r.Method == http.MethodPost {
		h.handleAdd(ctx, w, r)
		return
	}
	http.NotFound(w, r)
}
