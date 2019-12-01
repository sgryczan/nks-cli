package extensions

import (
	"fmt"

	nks "github.com/NetApp/nks-sdk-go/nks"
)

// GetUserProfileKeysetID - extended
// This method in the SDK Client returns the first matching keyset for the specified provider
// regardless of what organization it belongs to. This can result in errors such as
// "Provider credential does not belong to organization: $ORG"
func GetUserProfileKeysetID(up *nks.UserProfile, org int, prov string) (int, error) {
	if up == nil {
		return 0, fmt.Errorf("userprofile object is nil")
	}
	for _, ks := range up.Keysets {
		if (prov == "user_ssh" && ks.Category == "user_ssh" && ks.IsDefault && ks.Org == org) || (ks.Category == "provider" && ks.Entity == prov && ks.Org == org) {
			return ks.ID, nil
		}
	}
	return 0, fmt.Errorf("no %s keyset found in userprofile", prov)
}
