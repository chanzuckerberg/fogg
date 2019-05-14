package aws

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// STS is an STS client
type STS struct {
	Svc stsiface.STSAPI
}

// NewSTS returns an sts client
func NewSTS(c client.ConfigProvider, config *aws.Config) *STS {
	return &STS{Svc: sts.New(c, config)}
}

// GetSTSToken gets an sts token
func (s *STS) GetSTSToken(ctx context.Context, input *sts.GetSessionTokenInput) (*sts.Credentials, error) {
	output, err := s.Svc.GetSessionTokenWithContext(ctx, input)
	if err != nil {
		return nil, errors.Wrap(err, "Could not request sts tokens")
	}
	if output == nil {
		return nil, errors.New("Nil output from aws")
	}
	return output.Credentials, nil
}

const (
	// UserTokenProviderName is the name of this provider
	UserTokenProviderName = "UserTokenProvider"
)

// GetCallerIdentity gets the caller's identity
func (s *STS) GetCallerIdentity(ctx context.Context, input *sts.GetCallerIdentityInput) (*sts.GetCallerIdentityOutput, error) {
	output, err := s.Svc.GetCallerIdentityWithContext(ctx, input)
	if err != nil {
		return nil, errors.Wrap(err, "Could not get sts caller identity")
	}
	if output == nil {
		return nil, errors.New("Nil output from aws when calling sts get-caller-identity")
	}
	return output, nil
}

// UserTokenProviderCache caches mfa tokens
// Need this to json serialize/deserialize
type UserTokenProviderCache struct {
	Expiration      *time.Time `json:"expiration"`
	AccessKeyID     *string    `json:"access_key_id"`
	SecretAccessKey *string    `json:"secret_access_key"`
	SessionToken    *string    `json:"session_token"`
}

// UserTokenProvider is a token provider that gets sts tokens for a user
// Implementes the credentials.Provider interface
type UserTokenProvider struct {
	credentials.Expiry
	Client        *Client
	Duration      time.Duration
	cacheFile     string
	m             sync.RWMutex
	expireWindow  time.Duration
	isLogin       bool
	tokenProvider func() (string, error)
}

// NewUserTokenProvider returns a new user token provider
// Similar to doing an assume role operation but on a user instead
func NewUserTokenProvider(
	cacheFile string,
	client *Client,
	tokenProvider func() (string, error)) *UserTokenProvider {
	p := &UserTokenProvider{
		Client:        client,
		Duration:      stscreds.DefaultDuration,
		cacheFile:     cacheFile,
		expireWindow:  time.Minute,
		tokenProvider: tokenProvider,
	}
	return p
}

// try reading from file cache
func (p *UserTokenProvider) fromCache() (*sts.Credentials, error) {
	p.m.RLock()
	defer p.m.RUnlock()
	b, err := ioutil.ReadFile(p.cacheFile)
	if err != nil {
		// no cache - return nil credentials
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, errors.Wrapf(err, "Could not open mfa token cache %s", p.cacheFile)
	}
	var tokenCache UserTokenProviderCache
	err = json.Unmarshal(b, &tokenCache)
	if err != nil {
		log.Warnf("Cache corrupted at %s with error %s, deleting", p.cacheFile, err.Error())
		return nil, os.Remove(p.cacheFile)
	}
	// expired - return nil credentials
	if time.Now().After(tokenCache.Expiration.Add(-1 * p.expireWindow)) {
		return nil, nil
	}
	// else return cached
	creds := &sts.Credentials{
		AccessKeyId:     tokenCache.AccessKeyID,
		SecretAccessKey: tokenCache.SecretAccessKey,
		SessionToken:    tokenCache.SessionToken,
		Expiration:      tokenCache.Expiration,
	}
	return creds, nil
}

// toCache writes to cache
func (p *UserTokenProvider) toCache(creds *sts.Credentials) error {
	p.m.Lock()
	defer p.m.Unlock()
	cacheDir := path.Dir(p.cacheFile)
	err := os.MkdirAll(cacheDir, 0755) // #nosec
	if err != nil {
		return errors.Wrapf(err, "error creating cache dir %s", cacheDir)
	}

	tokenCache := &UserTokenProviderCache{
		AccessKeyID:     creds.AccessKeyId,
		SecretAccessKey: creds.SecretAccessKey,
		SessionToken:    creds.SessionToken,
		Expiration:      creds.Expiration,
	}
	b, err := json.Marshal(tokenCache)
	if err != nil {
		return errors.Wrap(err, "Could not marshal token to cache")
	}
	err = ioutil.WriteFile(p.cacheFile, b, 0644)
	return errors.Wrap(err, "Could not write token to cache")
}

// Retrieve generates a new set of temporary gredentials using STS.
func (p *UserTokenProvider) Retrieve() (credentials.Value, error) {
	creds := credentials.Value{}
	stsCreds, err := p.fromCache()
	if err != nil {
		return creds, err
	}

	if stsCreds == nil {
		// TODO: is there no better way than context.Background?
		user, err := p.Client.IAM.GetCurrentUser(context.Background())
		if err != nil {
			return creds, err
		}
		mfaSerial, err := p.Client.IAM.GetAnMFASerial(context.Background(), user.UserName)
		if err != nil {
			return creds, err
		}
		token, err := p.tokenProvider()
		if err != nil {
			return creds, errors.Wrap(err, "Could not read MFA token")
		}
		stsTokenInput := &sts.GetSessionTokenInput{}
		stsTokenInput.SetSerialNumber(mfaSerial).SetTokenCode(token)
		stsCreds, err = p.Client.STS.GetSTSToken(context.Background(), stsTokenInput)
		if err != nil {
			return creds, err
		}
	}
	// Check that we have all of these
	if stsCreds == nil ||
		stsCreds.AccessKeyId == nil ||
		stsCreds.Expiration == nil ||
		stsCreds.SecretAccessKey == nil ||
		stsCreds.SessionToken == nil {
		return creds, errors.New("Received malformed credentials from aws.Sts.GetSTSToken")
	}
	p.SetExpiration(*stsCreds.Expiration, p.expireWindow)
	creds.AccessKeyID = *stsCreds.AccessKeyId
	creds.SecretAccessKey = *stsCreds.SecretAccessKey
	creds.SessionToken = *stsCreds.SessionToken
	return creds, p.toCache(stsCreds)
}
