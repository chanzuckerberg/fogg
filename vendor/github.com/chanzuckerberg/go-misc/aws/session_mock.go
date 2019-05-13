package aws

import (
	"net/http"
	"net/http/httptest"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

// NewMockSession returns a new mock aws session
// safe for concurrent use
// Returns the mock server as well - user is responsible for calling server.Close() to free resources
func NewMockSession() (*session.Session, *httptest.Server) {
	// server is the mock server that simply writes a 200 status back to the client
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	return session.Must(session.NewSession(&aws.Config{
		DisableSSL: aws.Bool(true),
		Endpoint:   aws.String(server.URL),
		Region:     aws.String("mock"),
	})), server

}
