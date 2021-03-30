package getproviders

import (
	"crypto"
	"crypto/sha256"
)

type packageAuthenticationResult int

const (
	verifiedChecksum packageAuthenticationResult = iota
	officialProvider
	partnerProvider
	communityProvider
)

type PackageAuthenticationResult struct {
	result packageAuthenticationResult
	KeyID  string
}

func (t *PackageAuthenticationResult) String() string {
	if t == nil {
		return "unauthenticated"
	}
	return []string{
		"verified checksum",
		"signed by HashiCorp",
		"signed by a HashiCorp partner",
		"self-signed",
	}[t.result]
}

func (t *PackageAuthenticationResult) SignedByHashiCorp() bool {
	if t == nil {
		return false
	}
	if t.result == officialProvider {
		return true
	}
	return false
}

func (t *PackageAuthenticationResult) SignedByAnyParty() bool {
	if t == nil {
		return false
	}
	if t.result == officialProvider || t.result == partnerProvider || t.result == communityProvider {
		return true
	}
	return false
}

type SigningKey struct {
	ASCIIArmor     string `json:"ascii_armor"`
	TrustSignature string `json:"trust_signature"`
}

type PackageAuthentication interface {
	AuthenticatePackage(localLocation PackageLocation) (*PackageAuthenticationResult, error)
}

type PackageAuthenticationHashes interface {
	PackageAuthentication
	AcceptableHashes() []Hash
}

type packageAuthenticationAll []PackageAuthentication

func PackageAuthenticationAll(check ...PackageAuthentication) PackageAuthentication {
	return packageAuthenticationAll(check)
}

func (checks packageAuthenticationAll) AuthenticatePackage(localLocation PackageLocation) (*PackageAuthenticationResult, error) {
	var authResult *PackageAuthenticationResult
	for _, check := range checks {
		var err error
		authResult, err = check.AuthenticatePackage(localLocation)
		if err != nil {
			return authResult, err
		}
	}
	return authResult, nil
}

func (checks packageAuthenticationAll) AcceptableHashes() []Hash {
	for i := len(checks) - 1; i >= 0; i-- {
		check, ok := checks[i].(PackageAuthenticationHashes)
		if !ok {
			continue
		}
		allHashes := check.AcceptableHashes()
		if len(allHashes) > 0 {
			return allHashes
		}
	}
	return nil
}

type packageHashAuthentication struct {
	Requirements []Hash
	AllHashes    []Hash
	Platform     Platform
}
