package getproviders

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
