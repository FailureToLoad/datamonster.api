package server

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/failuretoload/datamonster/web"
)

const (
	missingJWTErrorMessage = "Requires authentication"
	invalidJWTErrorMessage = "Bad credentials"
)

func MakeJwtMiddleware(audience, domain string) (*jwtmiddleware.JWTMiddleware, error) {
	issuerURL, err := url.Parse(domain)
	if err != nil {
		log.Fatalf("Failed to parse the issuer url: %v", err)
	}

	provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)

	jwtValidator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		domain,
		[]string{audience},
		validator.WithCustomClaims(func() validator.CustomClaims {
			return new(web.CustomClaims)
		}),
	)
	if err != nil {
		log.Fatalf("Failed to set up the jwt validator")
	}

	errorHandler := func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("Encountered error while validating JWT: %v", err)
		if errors.Is(err, jwtmiddleware.ErrJWTMissing) {
			errorMessage := web.ErrorMessage{Message: missingJWTErrorMessage}
			if err := web.WriteJSON(w, http.StatusUnauthorized, errorMessage); err != nil {
				log.Printf("Failed to write error message: %v", err)
			}
			return
		}
		if errors.Is(err, jwtmiddleware.ErrJWTInvalid) {
			errorMessage := web.ErrorMessage{Message: invalidJWTErrorMessage}
			if err := web.WriteJSON(w, http.StatusUnauthorized, errorMessage); err != nil {
				log.Printf("Failed to write error message: %v", err)
			}
			return
		}
		web.InternalServerError(w, err)
	}

	middleware := jwtmiddleware.New(
		jwtValidator.ValidateToken,
		jwtmiddleware.WithErrorHandler(errorHandler),
	)

	return middleware, nil
}

func ValidateJWT(middleware *jwtmiddleware.JWTMiddleware, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if authHeaderParts := strings.Fields(r.Header.Get("Authorization")); len(authHeaderParts) > 0 && strings.ToLower(authHeaderParts[0]) != "bearer" {
			errorMessage := web.ErrorMessage{Message: invalidJWTErrorMessage}
			if err := web.WriteJSON(w, http.StatusUnauthorized, errorMessage); err != nil {
				log.Printf("Failed to write error message: %v", err)
			}
			return
		}
		middleware.CheckJWT(next).ServeHTTP(w, r)
	}
}
