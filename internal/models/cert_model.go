package models

import (
	"crypto"

	"github.com/go-acme/lego/v4/registration"
)

type CertUser struct {
	Email        string
	Registration *registration.Resource
	PrivateKey   crypto.PrivateKey
}

func (u *CertUser) GetEmail() string {
	return u.Email
}

func (u *CertUser) GetRegistration() *registration.Resource {
	return u.Registration
}

func (u *CertUser) GetPrivateKey() crypto.PrivateKey {
	return u.PrivateKey
}
