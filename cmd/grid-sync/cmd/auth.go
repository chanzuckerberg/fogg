package cmd

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/terraconstructs/grid/pkg/sdk"
)

type sessionConfig struct {
	ServerURL    string
	ClientID     string
	ClientSecret string
}

type tokenProvider struct {
	issuer       string
	clientID     string
	clientSecret string

	mu    sync.Mutex
	creds *sdk.Credentials
}

func (p *tokenProvider) token(ctx context.Context) (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.creds != nil && time.Until(p.creds.ExpiresAt) > time.Minute {
		return p.creds.AccessToken, nil
	}

	creds, err := sdk.LoginWithServiceAccount(ctx, p.issuer, p.clientID, p.clientSecret)
	if err != nil {
		return "", fmt.Errorf("service account login failed: %w", err)
	}

	p.creds = creds
	return creds.AccessToken, nil
}

type authTransport struct {
	base     http.RoundTripper
	provider *tokenProvider
}

func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	base := t.base
	if base == nil {
		base = http.DefaultTransport
	}

	if t.provider == nil {
		return base.RoundTrip(req)
	}

	token, err := t.provider.token(req.Context())
	if err != nil {
		return nil, err
	}

	clone := req.Clone(req.Context())
	clone.Header = clone.Header.Clone()
	clone.Header.Set("Authorization", "Bearer "+token)

	return base.RoundTrip(clone)
}

func newGridClient(ctx context.Context, cfg sessionConfig) (*sdk.Client, error) {
	if cfg.ServerURL == "" {
		return nil, fmt.Errorf("server url is required (flag --server or GRID_API_URL)")
	}

	authCfg, err := sdk.DiscoverAuthConfig(ctx, cfg.ServerURL)
	if err != nil {
		if isAuthDisabledError(err) {
			return sdk.NewClient(cfg.ServerURL), nil
		}
		return nil, fmt.Errorf("failed to discover auth config: %w", err)
	}

	if cfg.ClientID == "" || cfg.ClientSecret == "" {
		return nil, fmt.Errorf("client credentials are required (flags --client-id/--client-secret or GRID_CLIENT_ID/GRID_CLIENT_SECRET)")
	}

	issuer := authCfg.Issuer
	if issuer == "" {
		issuer = cfg.ServerURL
	}

	httpClient := &http.Client{Transport: &authTransport{
		base:     http.DefaultTransport,
		provider: &tokenProvider{issuer: issuer, clientID: cfg.ClientID, clientSecret: cfg.ClientSecret},
	}}

	return sdk.NewClient(cfg.ServerURL, sdk.WithHTTPClient(httpClient)), nil
}

func isAuthDisabledError(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "server returned 404") || strings.Contains(msg, "server returned 503")
}
