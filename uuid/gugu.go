package uuid

import (
	"fmt"

	"github.com/google/uuid"
)

// GenV4 returns v4 UUID
func GenV4() V4 {
	return V4(uuid.New().String())
}

// IsValidUUID checks if a given string is a valid v4 UUID
func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

// GenFromString returns a uuid  from string
func GenFromString(uuid string) (V4, error) {
	ok := IsValidUUID(uuid)
	if !ok {
		return "", fmt.Errorf("Unable to generate UUID V4 from invalid value=%s", uuid)
	}
	return V4(uuid), nil
}
