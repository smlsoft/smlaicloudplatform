package firebase_test

import (
	"smlaicloudplatform/internal/firebase"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerifyToken(t *testing.T) {

	adapter := firebase.NewFirebaseAdapter()
	user, err := adapter.ValidateToken("eyJhbGciOiJSUzI1NiIsImtpZCI6IjYyM2YzNmM4MTZlZTNkZWQ2YzU0NTkyZTM4ZGFlZjcyZjE1YTBmMTMiLCJ0eXAiOiJKV1QifQ.eyJuYW1lIjoic21sc29mdCBkZXYiLCJwaWN0dXJlIjoiaHR0cHM6Ly9saDMuZ29vZ2xldXNlcmNvbnRlbnQuY29tL2EvQUdObXl4WjJWRzVMM250Z0lCWkNodDNkbDRhRGFmU1N0WXI4VmhablRLWEs9czk2LWMiLCJpc3MiOiJodHRwczovL3NlY3VyZXRva2VuLmdvb2dsZS5jb20vZGVkZXBvcyIsImF1ZCI6ImRlZGVwb3MiLCJhdXRoX3RpbWUiOjE2OTEwMzU5MTMsInVzZXJfaWQiOiJHbXNuS1cxc0I3ZmNVeE1Rb045c2luWmNZemcxIiwic3ViIjoiR21zbktXMXNCN2ZjVXhNUW9OOXNpblpjWXpnMSIsImlhdCI6MTY5MTAzNTkxMywiZXhwIjoxNjkxMDM5NTEzLCJlbWFpbCI6InNtbHNvZnRkZXZAZ21haWwuY29tIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsImZpcmViYXNlIjp7ImlkZW50aXRpZXMiOnsiZ29vZ2xlLmNvbSI6WyIxMDAwNTAyODM3Mjc1NjUyNjE0MTgiXSwiZW1haWwiOlsic21sc29mdGRldkBnbWFpbC5jb20iXX0sInNpZ25faW5fcHJvdmlkZXIiOiJnb29nbGUuY29tIn19.ZDQbYsyksrLkUNcsEruloHqzSROAaOAa3HC9F_zuR91RuzuM_BhkcqMpFZnJ4YZySePKicIEHcDcbotcsimcKlqmaq_rpk1xhpB9RIIYwyI_G1F9veN9rFHi_gCgvDZJaBvf3DIeT7XXuQB5JE0v2DxLhZPqFfS2aXNZ8vqZ36YyjsMj0OzooEqkYsFvh102yyYdLvWoV2Pj9Cf0TL95uPy5BZK319_-fm2MLyFhU7mOZT7474P9Hw0AaoRRtzaWvwoed4MjlPVMl5O3oxytu4y12VQfEc-6FW1W92UE37XNWyLNbtBm_HlO6eNiqNR8ublgsi-9gjPO1Ad-efDBwA")
	assert.Nil(t, err, "error not nil")
	assert.NotNil(t, user, "user is nil")
}
