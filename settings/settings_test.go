package settings

import (
	"os"
	"strings"
	"testing"
)

func TestSettings(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current testing directory: %v", err)
	}
	defer os.Chdir(wd) // Restore WD

	type TestObj struct {
		Hello string
	}

	for name, test := range map[string]struct {
		Prelude   func(t *testing.T)
		Compare   func(t *testing.T, to TestObj)
		Directory string
		ErrorLike string
	}{
		"basic file": {
			Prelude: func(t *testing.T) {},
			Compare: func(t *testing.T, to TestObj) {
				if to.Hello != "World" {
					t.Fatalf("Expected \"World\" got \"%v\"", to.Hello)
				}
			},
			Directory: "testing/basic",
		},
		"bad file": {
			Prelude:   func(t *testing.T) {},
			Compare:   func(t *testing.T, to TestObj) {},
			Directory: "testing/badfile",
			ErrorLike: "unexpected character \"!\"",
		},
		"no file": {
			Prelude: func(t *testing.T) {
				if err := os.Setenv("HELLO", "No File"); err != nil {
					t.Fatalf("Failed to set env for test: %v", err)
				}
			},
			Compare: func(t *testing.T, to TestObj) {
				if to.Hello != "No File" {
					t.Fatalf("Expected \"No File\" got \"%v\"", to.Hello)
				}
			},
			Directory: "testing/nofile",
		},
	} {
		test := test
		t.Run(name, func(t *testing.T) {
			if err := os.Chdir(wd); err != nil {
				t.Fatalf("Failed to reset WD: %v", err)
			}

			if test.Directory != "" {
				if err := os.Chdir(test.Directory); err != nil {
					t.Fatalf("Failed to set Test Directory: %v", err)
				}
			}

			test.Prelude(t)
			conf, err := Process[TestObj]()
			if (test.ErrorLike != "") != (err != nil) {
				t.Fatalf("Expect Error %v got %v", (test.ErrorLike != ""), err)
			}
			if test.ErrorLike != "" {
				if !strings.Contains(err.Error(), test.ErrorLike) {
					t.Fatalf("Expected error to contain \"%v\" but it didn't. Error: %v", test.ErrorLike, err)
				}
				return
			}
			test.Compare(t, conf)
		})
	}
}
