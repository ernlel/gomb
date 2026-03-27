// Package main demonstrates several caching strategies for gomb-rendered HTML.
//
// Because gomb.Element is an immutable value type, rendered HTML strings can be
// cached freely. Three patterns are shown:
//
//  1. Static component cache – build once at startup with sync.Once.
//  2. Per-key fragment cache – cache rendered fragments keyed by a string with
//     sync.Map plus a TTL-based eviction helper.
//  3. http.Handler middleware – a simple caching middleware that stores the full
//     response body and replays it for subsequent requests.
//
// Run:
//
//	go run main.go
//	open http://localhost:8080
package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	. "github.com/ernlel/gomb"
)

// ── 1. Static component cache (sync.Once) ────────────────────────────────────
//
// Use when a component never changes after startup (e.g. a site-wide nav bar).
// The HTML is rendered exactly once; every handler call reuses the string.

var (
	navbarOnce sync.Once
	navbarHTML string
)

func navbarComponent() string {
	navbarOnce.Do(func() {
		navbarHTML = E("nav").
			A("class", "bg-blue-700 text-white px-6 py-3 flex gap-6").
			C(
				E("a").A("href", "/").A("class", "hover:underline").T("Home"),
				E("a").A("href", "/products").A("class", "hover:underline").T("Products"),
				E("a").A("href", "/contact").A("class", "hover:underline").T("Contact"),
			).ToString()
	})
	return navbarHTML
}

// ── 2. Per-key fragment cache with TTL ───────────────────────────────────────
//
// Use for fragments that are expensive to render and change infrequently
// (e.g. rendered markdown, user profile cards, product tiles).

type cacheEntry struct {
	html    string
	expires time.Time
}

var fragmentCache sync.Map

// cachedFragment returns the cached HTML for key if still valid, otherwise calls
// build(), caches the result for ttl, and returns it.
func cachedFragment(key string, ttl time.Duration, build func() Element) string {
	if v, ok := fragmentCache.Load(key); ok {
		e := v.(cacheEntry)
		if time.Now().Before(e.expires) {
			return e.html
		}
	}
	html := build().ToString()
	fragmentCache.Store(key, cacheEntry{html: html, expires: time.Now().Add(ttl)})
	return html
}

// productCard is an expensive-to-render component cached per product ID.
func productCard(id int, name, desc string, price float64) string {
	key := fmt.Sprintf("product:%d", id)
	return cachedFragment(key, 5*time.Minute, func() Element {
		return E("div").
			A("class", "border rounded p-4 shadow-sm").
			C(
				E("h3").A("class", "text-lg font-semibold").T(name),
				E("p").A("class", "text-gray-600 text-sm mt-1").T(desc),
				E("p").A("class", "text-blue-700 font-bold mt-2").T(fmt.Sprintf("$%.2f", price)),
				E("button").
					A("class", "mt-3 bg-blue-600 text-white px-4 py-1 rounded text-sm hover:bg-blue-700").
					A("onclick", fmt.Sprintf(`addToCart(%d)`, id)).
					T("Add to Cart"),
			)
	})
}

// ── 3. HTTP handler response cache ───────────────────────────────────────────
//
// Use for full pages or API responses that should be cached at the handler level.
// The middleware stores the response body once and writes it directly on subsequent
// requests, skipping all rendering work.

type responseCache struct {
	mu       sync.RWMutex
	entries  map[string]responseCacheEntry
	defaultTTL time.Duration
}

type responseCacheEntry struct {
	body        []byte
	contentType string
	status      int
	expires     time.Time
}

func newResponseCache(defaultTTL time.Duration) *responseCache {
	return &responseCache{
		entries:    make(map[string]responseCacheEntry),
		defaultTTL: defaultTTL,
	}
}

// Wrap returns an http.Handler that caches the response of next for defaultTTL.
func (c *responseCache) Wrap(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.RequestURI()

		c.mu.RLock()
		entry, ok := c.entries[key]
		c.mu.RUnlock()

		if ok && time.Now().Before(entry.expires) {
			w.Header().Set("Content-Type", entry.contentType)
			w.Header().Set("X-Cache", "HIT")
			w.WriteHeader(entry.status)
			w.Write(entry.body)
			return
		}

		rec := &responseRecorder{header: make(http.Header), status: http.StatusOK}
		next(rec, r)

		c.mu.Lock()
		c.entries[key] = responseCacheEntry{
			body:        rec.body,
			contentType: rec.header.Get("Content-Type"),
			status:      rec.status,
			expires:     time.Now().Add(c.defaultTTL),
		}
		c.mu.Unlock()

		w.Header().Set("Content-Type", rec.header.Get("Content-Type"))
		w.Header().Set("X-Cache", "MISS")
		w.WriteHeader(rec.status)
		w.Write(rec.body)
	}
}

// responseRecorder captures the response written by an http.Handler.
type responseRecorder struct {
	body   []byte
	header http.Header
	status int
}

func (r *responseRecorder) Header() http.Header { return r.header }
func (r *responseRecorder) WriteHeader(code int) {
	if r.status == 0 {
		r.status = code
	}
}
func (r *responseRecorder) Write(b []byte) (int, error) {
	if r.status == 0 {
		r.status = http.StatusOK
	}
	r.body = append(r.body, b...)
	return len(b), nil
}

// ── demo server ──────────────────────────────────────────────────────────────

var products = []struct {
	ID    int
	Name  string
	Desc  string
	Price float64
}{
	{1, "Gopher Plushie", "Soft stuffed Gopher toy.", 19.99},
	{2, "Go T-Shirt", "100% cotton, sizes S–XXL.", 24.99},
	{3, "Mechanical Keyboard", "Cherry MX Brown, TKL layout.", 129.99},
}

func homePage() Element {
	cards := Map(products, func(p struct {
		ID    int
		Name  string
		Desc  string
		Price float64
	}) Element {
		// productCard returns a cached HTML string; wrap it with Raw() so it is
		// inserted verbatim into the parent element.
		return Raw(productCard(p.ID, p.Name, p.Desc, p.Price))
	})

	return E("html").A("lang", "en").C(
		E("head").C(
			E("meta").A("charset", "UTF-8"),
			E("meta").A("name", "viewport").A("content", "width=device-width, initial-scale=1"),
			E("title").T("gomb caching demo"),
			E("link").
				A("rel", "stylesheet").
				A("href", "https://cdn.jsdelivr.net/npm/tailwindcss@3/dist/tailwind.min.css"),
		),
		E("body").A("class", "bg-gray-50 min-h-screen").C(
			// Static cached navbar
			Raw(navbarComponent()),
			E("main").A("class", "max-w-3xl mx-auto p-8").C(
				E("h1").A("class", "text-3xl font-bold mb-6").T("Products"),
				E("p").A("class", "text-sm text-gray-500 mb-4").
					T("Product cards are cached for 5 minutes. The full page response is cached for 10 seconds (see X-Cache header)."),
				E("div").A("class", "grid grid-cols-1 sm:grid-cols-3 gap-4").C(cards...),
			),
		),
	)
}

func main() {
	cache := newResponseCache(10 * time.Second)

	http.HandleFunc("/", cache.Wrap(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		homePage().Render(w)
	}))

	fmt.Println("Caching example running at http://localhost:8080")
	fmt.Println("Tip: check the X-Cache response header (HIT/MISS)")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("server error: %v\n", err)
	}
}
