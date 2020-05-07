package archetype

import (
	"testing"
)

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func TestArchetype(t *testing.T) {
	assetDir, err := AssetDir("")
	if err != nil {
		t.Fatalf("expected no err, but got %v", err)
	}
	expected := "terraform-archetype-simple"
	if !contains(assetDir, expected) {
		t.Fatalf("expected assetDir to %s, but got %v", expected, assetDir)
	}
}
