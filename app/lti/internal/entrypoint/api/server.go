package api

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
)

var (
	sessionStore = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET_KEY")))
	oauthConfig  *oauth2.Config
)




func validateLTI11Request(r *http.Request) error {
	// Use an OAuth1 library (e.g., github.com/stretchr/gomniauth) to validate the OAuth1 signature
	// You'll need to configure the library with the correct Consumer Key and Shared Secret
	// provided by the LMS administrator

	// For example:
	// err := gomniauth.ValidateSignature(r, consumerKey, sharedSecret)
	// if err != nil {
	// 	 return err
	// }

	// Replace this with actual validation using an OAuth1 library
	log.Fatal("Not yet implemented")
	return nil
}

func validateLTI13Request(r *http.Request) error {
	// Extract the JWT token from the request
	tokenString := r.FormValue("id_token")
	if tokenString == "" {
		return fmt.Errorf("Missing JWT token in LTI 1.3 launch request")
	}

	// Use a JWT library (e.g., github.com/golang-jwt/jwt) to parse and validate the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check that the token uses the expected signing method
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// Retrieve the public key for the issuer (LMS) to validate the token
		// You can either hardcode the public key or retrieve it dynamically (e.g., via JWKS)
		// For example:
		issuer, ok := token.Claims.(jwt.MapClaims)["iss"].(string)
		if !ok || issuer == "" {
			return nil, fmt.Errorf("Missing or invalid 'iss' claim")
		}

		jwksURL := issuer + "/.well-known/jwks.json"

		keyID := token.Header["kid"].(string)
		publicKey, err := getPublicKeyForIssuer(jwksURL, keyID)
		if err != nil {
			return nil, err
		}

		// Replace this with actual public key retrieval
		return publicKey, nil
	})

	if err != nil || !token.Valid {
		return fmt.Errorf("Invalid JWT token in LTI 1.3 launch request")
	}

	// Check required LTI 1.3 claims
	claims := token.Claims.(jwt.MapClaims)
	requiredClaims := []string{"iss", "aud", "exp", "iat", "nonce", "azp", "message_type", "version", "resource_link"}

	for _, claim := range requiredClaims {
		if _, ok := claims[claim]; !ok {
			return fmt.Errorf("Missing required LTI 1.3 claim: %s", claim)
		}
	}

	return nil
}

func getPublicKeyForIssuer(jwksURL, keyID string) (*rsa.PublicKey, error) {
	resp, err := http.Get(jwksURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed to fetch JWKS: %s", resp.Status)
	}

	var jwks struct {
		Keys []struct {
			Kty string `json:"kty"`
			Alg string `json:"alg"`
			Use string `json:"use"`
			Kid string `json:"kid"`
			N   string `json:"n"`
			E   string `json:"e"`
		} `json:"keys"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, err
	}

	for _, key := range jwks.Keys {
		if key.Kid == keyID {
			if key.Kty != "RSA" {
				return nil, fmt.Errorf("Invalid key type: %s", key.Kty)
			}

			n, err := base64.RawURLEncoding.DecodeString(key.N)
			if err != nil {
				return nil, err
			}

			e, err := base64.RawURLEncoding.DecodeString(key.E)
			if err != nil {
				return nil, err
			}

			rsaPublicKey := &rsa.PublicKey{
				N: new(big.Int).SetBytes(n),
				E: int(new(big.Int).SetBytes(e).Int64()),
			}

			return rsaPublicKey, nil
		}
	}

	return nil, fmt.Errorf("Public key not found for kid: %s", keyID)
}
