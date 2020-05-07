package version

import (
	"fmt"
	"strings"
	"testing"
)

func TestVersion(t *testing.T) {
	GitCommit = "a2c3e4"
	GitDescribe = "v1.2.3"
	expectedString := fmt.Sprintf("%s-%s (%s)", GitDescribe, VersionPrerelease, GitCommit)
	producedString := GetHumanVersion()
	if !strings.EqualFold(expectedString, producedString) {
		t.Fatalf("got unexpected user agent string: %s. Expected: %s.", producedString, expectedString)
	}
}
