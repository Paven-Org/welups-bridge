package libs

import (
	"bridge/logger"
	"testing"
)

func TestGeneratePassword(t *testing.T) {
	hs, err := HashPasswd("root")
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}
	logger.Get().Info().Msgf("Hashed password: %s", hs)
}
