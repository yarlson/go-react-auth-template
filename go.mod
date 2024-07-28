module goauth

go 1.22

require (
	github.com/go-chi/chi/v5 v5.1.0
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	github.com/lib/pq v1.10.9
	github.com/rs/cors v1.11.0
	golang.org/x/oauth2 v0.21.0
)

require cloud.google.com/go/compute/metadata v0.3.0 // indirect
