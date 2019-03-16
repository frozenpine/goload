package utils

import (
	"errors"
	"regexp"
)

// IdentityMap identity pattern map
type IdentityMap map[string]*regexp.Regexp

// CheckIdentity check & modify login map
func (idMap *IdentityMap) CheckIdentity(id string, login map[string]string) error {
	for name, pattern := range *idMap {
		if !pattern.MatchString(id) {
			continue
		}

		login[name] = id
		login["type"] = "account"
		login["verifyCode"] = ""

		return nil
	}

	return errors.New("invalid login identiy")
}

// NewIdentityMap generate identity pattern map
func NewIdentityMap() IdentityMap {
	mp := make(IdentityMap)

	mp["email"] = regexp.MustCompile(`[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)*`)
	mp["mobile"] = regexp.MustCompile(`(+[0-9]{2,3})?[0-9-]{6,13}`)

	return mp
}
