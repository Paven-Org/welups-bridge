package libs

import (
	"bridge/service-managers/logger"
	"testing"
	"fmt"
)

func TestGeneratePassword(t *testing.T) {
	hs, err := HashPasswd("abcxyz")
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}
	logger.Get().Info().Msgf("Hashed password: %s", hs)
}

func TestStrongPasswd(t *testing.T) {
	type testCase struct {
		plain string
		expected bool
	}
	suite := []testCase {
		testCase{
			plain: "aA1_",
			expected: false,
		},
		testCase{
			plain: "aaabbbcc",
			expected: false,
		},
		testCase{
			plain: "aaab1bbc",
			expected: false,
		},
		testCase{
			plain: "Aaab1bbc",
			expected: false,
		},
		testCase{
			plain: "Aaa*1bbc",
			expected: true,
		},
	}
	for _, tcase := range suite {
		fmt.Printf("test case: %+v\n", tcase)
		res := StrongPasswd(tcase.plain)
		if res != tcase.expected {
			t.Fatal("FAILED")
		}
		fmt.Println("PASSED with result: ", res)
	}
}
