package depmain_test

import (
	"os"
	"testing"

	"github.com/onemedical/depmain"
)

func TestEnv(t *testing.T) {
	os.Setenv("DEPMAIN_TEST_VALUE", "set")
	ext := depmain.New()
	if e := ext.Getenv("DEPMAIN_TEST_VALUE"); e != "set" {
		t.Errorf("Getenv: want %s got %s", "set", e)
	}
	if e := ext.Getenv("DEPMAIN_UNKNOWN_TEST_VALUE"); e != "" {
		t.Errorf("Getenv bad value: want %q got %s", "", e)
	}
	_, found := ext.LookupEnv("DEPMAIN_TEST_VALUE")
	if !found {
		t.Errorf("LookupEnv: should have found value, got false")
	}
	_, found = ext.LookupEnv("DEPMAIN_UNKNOWN_TEST_VALUE")
	if found {
		t.Errorf("LookupEnv: should not have found value, got true")
	}
}
