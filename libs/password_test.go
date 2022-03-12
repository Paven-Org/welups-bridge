package libs

import (
	"bridge/service-managers/logger"
	"testing"
)

func TestGeneratePassword(t *testing.T) {
	hs, err := HashPasswd("abcxyz")
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}
	logger.Get().Info().Msgf("Hashed password: %s", hs)
}
