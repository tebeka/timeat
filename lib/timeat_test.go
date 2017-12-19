package timeat

import (
	"strings"
	"testing"
)

func init() {
	SetAPIKey("")
}

func TestTimeAt(t *testing.T) {
	addr := "sunnyvale"
	tis, err := TimeAt(addr)
	if err != nil {
		t.Fatalf("can't get time at %s - %s", addr, err)
	}
	found := false
	for _, ti := range tis {
		if strings.Contains(ti.Address, "Sunnyvale") {
			found = true
			break
		}
	}

	if !found {
		t.Fatalf("no 'Sunnyvale' in output: %v", tis)
	}
}
