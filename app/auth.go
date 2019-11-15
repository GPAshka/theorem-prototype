package app

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"fmt"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"log"
	"net/http"
	"strings"
	"theorem-prototype/config"
)

var authClient *auth.Client

func init() {
	projectId := config.GetFirebaseProjectId()

	opt := option.WithCredentials(&google.Credentials{ProjectID: projectId})
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error while getting Firebase app: %v\n", err)
		return
	}

	authClient, err = app.Auth(context.Background())
	if err != nil {
		log.Fatalf("error while getting Firebase Auth client: %v\n", err)
		return
	}
}

func JwtAuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//List of endpoints that doesn't require auth
		notAuth := []string{"/api/v1/hc", "/api/v1/device", "/api/v1/device/sensors", "/api/v1/device/sensors/bulk"}
		requestPath := r.URL.Path //current request path

		//check if request does not need authentication, serve the request if it doesn't need it
		for _, value := range notAuth {
			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		tokenHeader := r.Header.Get("Authorization") //Grab the token from the header
		if tokenHeader == "" {                       //Token is missing, returns with error code 403 Unauthorized
			writeError(w, "Missing auth token")
			return
		}

		splited := strings.Split(tokenHeader, " ") //The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requirement
		if len(splited) != 2 {
			writeError(w, "Invalid/Malformed auth token")
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Add("Content-Type", "application/json")
			return
		}

		tokenBody := splited[1] //Grab the token part, what we are truly interested in
		token, err := authClient.VerifyIDToken(context.Background(), tokenBody)
		if err != nil {
			writeError(w, fmt.Sprintf("authentication token is not valid: %v\n", err))
			return
		}

		ctx := context.WithValue(r.Context(), "user", token.UID)
		r = r.WithContext(ctx)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func writeError(w http.ResponseWriter, errorMessage string) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Header().Add("Content-Type", "application/json")
	if _, err := w.Write([]byte(errorMessage)); err != nil {
		log.Fatal("error while writing JwtAuthentication message", err)
	}
}
