// Package service provides the Kerberos KDC business logic layer.
package service

import (
	"context"
	"fmt"
)

// Principal represents a Kerberos principal entry.
type Principal struct {
	Name               string `json:"name"`
	Policy             string `json:"policy"`
	MaxLife            string `json:"max_life"`
	MaxRenewLife       string `json:"max_renew_life"`
	LastPasswordChange string `json:"last_password_change"`
	ExpirationDate     string `json:"expiration_date"`
	Locked             bool   `json:"locked"`
}

// Policy represents a Kerberos policy.
type Policy struct {
	Name      string `json:"name"`
	MinLength int    `json:"min_length"`
	History   int    `json:"history"`
	MaxLife   int    `json:"max_life"`
}

// AddPrincipalRequest is the body for creating a new principal.
type AddPrincipalRequest struct {
	Principal    string `json:"principal"`
	Policy       string `json:"policy"`
	MaxLife      string `json:"maxlife"`
	MaxRenewLife string `json:"maxrenewlife"`
}

// ChangePWRequest is the body for changing a principal's password.
type ChangePWRequest struct {
	Password string `json:"password"`
}

// RealmInfo holds high-level information about the Kerberos realm.
type RealmInfo struct {
	Realm       string `json:"realm"`
	KDCHost     string `json:"kdc_host"`
	AdminServer string `json:"admin_server"`
}

// KerberosService performs Kerberos KDC operations via kadmin.
// All methods are currently stubs; replace with real kadmin invocations.
type KerberosService struct {
	realm       string
	kadminServer string
	principal   string
	keytabPath  string
	kdcHost     string
}

// New creates a KerberosService from the supplied connection parameters.
func New(realm, kadminServer, kadminPrincipal, keytabPath, kdcHost string) *KerberosService {
	return &KerberosService{
		realm:        realm,
		kadminServer: kadminServer,
		principal:    kadminPrincipal,
		keytabPath:   keytabPath,
		kdcHost:      kdcHost,
	}
}

// ListPrincipals returns all principal names in the realm.
func (s *KerberosService) ListPrincipals(_ context.Context) ([]string, error) {
	// TODO: exec kadmin -q "listprincs"
	return []string{}, nil
}

// GetPrincipal returns details about a single principal.
func (s *KerberosService) GetPrincipal(_ context.Context, name string) (*Principal, error) {
	// TODO: exec kadmin -q "getprinc <name>"
	if name == "" {
		return nil, fmt.Errorf("principal name is required")
	}
	return &Principal{Name: name}, nil
}

// AddPrincipal creates a new principal in the KDC.
func (s *KerberosService) AddPrincipal(_ context.Context, req AddPrincipalRequest) error {
	// TODO: exec kadmin -q "addprinc ..."
	if req.Principal == "" {
		return fmt.Errorf("principal name is required")
	}
	return nil
}

// DeletePrincipal removes a principal from the KDC.
func (s *KerberosService) DeletePrincipal(_ context.Context, name string) error {
	// TODO: exec kadmin -q "delprinc <name>"
	if name == "" {
		return fmt.Errorf("principal name is required")
	}
	return nil
}

// ChangePassword sets a new password for the named principal.
func (s *KerberosService) ChangePassword(_ context.Context, name, password string) error {
	// TODO: exec kadmin -q "cpw <name>"
	if name == "" || password == "" {
		return fmt.Errorf("principal name and password are required")
	}
	return nil
}

// GenerateKeytab produces a keytab file for the named principal.
// Returns the keytab bytes on success.
func (s *KerberosService) GenerateKeytab(_ context.Context, name string) ([]byte, error) {
	// TODO: exec kadmin -q "ktadd -k <tmpfile> <name>"
	if name == "" {
		return nil, fmt.Errorf("principal name is required")
	}
	return []byte{}, nil
}

// ListPolicies returns all policy names in the realm.
func (s *KerberosService) ListPolicies(_ context.Context) ([]string, error) {
	// TODO: exec kadmin -q "listpols"
	return []string{}, nil
}

// GetPolicy returns details about a single policy.
func (s *KerberosService) GetPolicy(_ context.Context, name string) (*Policy, error) {
	// TODO: exec kadmin -q "getpol <name>"
	if name == "" {
		return nil, fmt.Errorf("policy name is required")
	}
	return &Policy{Name: name}, nil
}

// CreatePolicy creates a new password policy.
func (s *KerberosService) CreatePolicy(_ context.Context, p Policy) error {
	// TODO: exec kadmin -q "addpol ..."
	if p.Name == "" {
		return fmt.Errorf("policy name is required")
	}
	return nil
}

// DeletePolicy removes a password policy.
func (s *KerberosService) DeletePolicy(_ context.Context, name string) error {
	// TODO: exec kadmin -q "delpol <name>"
	if name == "" {
		return fmt.Errorf("policy name is required")
	}
	return nil
}

// RealmInfo returns high-level information about the configured realm.
func (s *KerberosService) RealmInfo(_ context.Context) (*RealmInfo, error) {
	return &RealmInfo{
		Realm:       s.realm,
		KDCHost:     s.kdcHost,
		AdminServer: s.kadminServer,
	}, nil
}
