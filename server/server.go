package server

import (
	"context"
	"github.com/failuretoload/datamonster/web"
	"log"
	"net/http"
	"time"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/failuretoload/datamonster/helpers"
	"github.com/go-chi/cors"

	"github.com/unrolled/secure"

	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	Mux *chi.Mux
}

func NewServer() Server {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(middleware.Timeout(10 * time.Second))
	router.Use(SecureOptions())
	router.Use(CacheControl)
	router.Use(CorsHandler())
	router.Use(clerkhttp.WithHeaderAuthorization())
	router.Use(UserIdExtractor)

	return Server{
		Mux: router,
	}
}

func (s Server) Run() {
	clerk.SetKey(helpers.SafeGetEnv("CLERK_SECRET_KEY"))
	log.Default().Println("Starting server on port 8080")
	err := http.ListenAndServe(":8080", s.Mux)
	if err != nil {
		log.Default().Fatal(err)
	}
}

func SecureOptions() func(http.Handler) http.Handler {

	options := secure.Options{
		STSSeconds:            31536000,
		STSIncludeSubdomains:  true,
		STSPreload:            true,
		FrameDeny:             true,
		ForceSTSHeader:        true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		CustomBrowserXssValue: "0",
		ContentSecurityPolicy: "default-src 'self', frame-ancestors 'none'",
	}
	return secure.New(options).Handler
}

func CacheControl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		headers := rw.Header()
		headers.Set("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate")
		headers.Set("Pragma", "no-cache")
		headers.Set("Expires", "0")
		next.ServeHTTP(rw, req)
	})
}

func CorsHandler() func(http.Handler) http.Handler {
	client := helpers.SafeGetEnv("WEB_CLIENT")
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{client},
		AllowedMethods:   []string{"HEAD", "GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           3599, // Maximum value not ignored by any of major browsers
	})
	return c.Handler
}

func UserIdExtractor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		claims, ok := clerk.SessionClaimsFromContext(req.Context())
		if !ok {
			web.Unauthorized(rw, nil)
			return
		}
		userID := claims.RegisteredClaims.Subject
		if userID == "" {
			reason := "user id not found"
			web.Unauthorized(rw, &reason)
			return
		}
		ctx := context.WithValue(req.Context(), web.UserIdKey, userID)
		next.ServeHTTP(rw, req.WithContext(ctx))
	})
}
