package getproviders

import (
	"crypto"
	"crypto/sha256"
	"fmt"
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
	RequiredHashes []Hash
	AllHashes      []Hash
	Platform       Platform
}

func NewPackageHashAuthentication(platform Platform, validHashes []Hash) PackageAuthentication {
	requiredHashes := PreferredHashes(validHashes)
	return packageHashAuthentication{
		RequiredHashes: requiredHashes,
		AllHashes:      validHashes,
		Platform:       platform,
	}
}

func (a packageHashAuthentication) AuthenticatePackage(localLocation PackageLocation) (*PackageAuthenticationResult, error) {
	if len(a.RequiredHashes) == 0 {
		return nil, fmt.Errorf("this version of Terraform does not upport any of the checksum formats given for this provider")
	}
	matches, err := PackageMatchesAnyHash(localLocation, a.RequiredHashes)
	if err != nil {
		return nil, fmt.Errorf("failed to verify provider package checksums: %s", err)
	}
	if matches {
		return &PackageAuthenticationResult{result: verifiedChecksum}, nil
	}
	if len(a.RequiredHashes) == 1 {
		return nil, fmt.Errorf("provider package doesn't match the expected checksum %q", a.RequiredHashes[0].String())
	}

	return nil, fmt.Errorf("provider package doesn't match the any of the expected checksums")
}

func (a packageHashAuthentication) AcceptableHashes() []Hash {
	return a.AllHashes
}

type archiveHashAuthentication struct {
	Platform      Platform
	WantSHA256Sum [sha256.Size]byte
}

func NewArchiveChecksumAuthentication(platform Platform, wantSha256Sum [sha256.Size]byte) PackageAuthentication {
	return archiveHashAuthentication{Platform: platform, WantSHA256Sum: wantSha256Sum}
}

func (a archiveHashAuthentication) AuthenticatePackage(localLocation PackageLocation) (*PackageAuthenticationResult, error) {
	archiveLocation, ok := localLocation.(PackageLocalArchive)
	if !ok {
		return nil, fmt.Errorf("cannot check archive hash for non-archive location %s", localLocation)
	}

	gotHash, err := PackageHashLegacyZipSHA(archiveLocation)
	if err != nil {
		return nil, fmt.Errorf("failed to compute checksum for %s: %s", archiveLocation, err)
	}

	wantHash := HashLegacyZipSHAFromSHA(a.WantSHA256Sum)
	if gotHash != wantHash {
		return nil, fmt.Errorf("archive has incorrect checksum %s (expected %s)", gotHash, wantHash)
	}
	return &PackageAuthenticationResult{
		result: verifiedChecksum,
	}, nil
}

type matchingChecksumAuthentication struct {
	Document      []byte
	Filename      string
	WantSHA256Sum [sha256.Size]byte
}

func NewMatchingChecksumAuthentication(document []byte, filename string, wantSHA256Sum [sha256.Size]byte) PackageAuthentication {

}
