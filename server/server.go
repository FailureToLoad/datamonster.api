package server

import (
	"github.com/failuretoload/datamonster/helpers"
	"github.com/go-chi/cors"
	"log"
	"net/http"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/unrolled/secure"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	Mux        *chi.Mux
	middleware *jwtmiddleware.JWTMiddleware
}

func NewServer() Server {
	audience := helpers.SafeGetEnv("AUTH0_AUDIENCE")
	domain := helpers.SafeGetEnv("AUTH0_DOMAIN")

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(middleware.Timeout(10 * time.Second))

	secureMiddleware := secure.New(SecureOptions())
	router.Use(secureMiddleware.Handler)
	router.Use(HandleCacheControl)
	jwtAuthenticator, err := MakeJwtMiddleware(audience, domain)
	if err != nil {
		panic(err)
	}
	return Server{
		Mux:        router,
		middleware: jwtAuthenticator,
	}
}

func (s Server) Run() {
	log.Default().Println("Starting server on port 8080")
	err := http.ListenAndServe(":8080", finalHandler(ValidateJWT(s.middleware, s.Mux)))
	if err != nil {
		log.Default().Fatal(err)
	}
}

func finalHandler(next http.Handler) http.Handler {
	secOptionsHandler := secure.New(SecureOptions()).Handler
	corsHandler := CorsHandler()
	return corsHandler(
		secOptionsHandler(
			HandleCacheControl(next)))
}

func SecureOptions() secure.Options {
	return secure.Options{
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
}

func HandleCacheControl(next http.Handler) http.Handler {
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
