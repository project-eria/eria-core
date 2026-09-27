package eria

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestControllerID(t *testing.T) {
	host, _ := os.Hostname()
	host, _, _ = strings.Cut(host, ".")
	bin := strings.ReplaceAll(filepath.Base(os.Args[0]), ".", "-")

	// Unique per binary and host, prefix kept for the HA integration
	if got, want := controllerID(""), "eria:app:controler:"+bin+"@"+host; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
	if got := controllerID("test:1"); got != "eria:app:controler:"+bin+"@"+host+":test:1" {
		t.Fatalf("instance not appended: %q", got)
	}
	if strings.Contains(controllerID(""), ".") {
		t.Fatal("dots in the id would be escaped in the mDNS instance name")
	}
}
