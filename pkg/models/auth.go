package models

import (
	"context"
	"time"
)

type AuthInfo struct {
	Username          string    `json:"username,omitempty"`
	AccessKeyId       string    `json:"accessKeyId,omitempty"`
	SecretAccessKey   string    `json:"secretAccessKey,omitempty"`
	SessionToken      string    `json:"sessionToken,omitempty"`
	SessionExpiryTime time.Time `json:"sessionExpiryTime,omitempty"`
	AuthMechanism     string    `json:"authMechanism,omitempty"`
}

type Authenticator interface {
	Authenticate(ctx context.Context) (AuthInfo, error)
}
