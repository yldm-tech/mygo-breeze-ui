package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestListComponents(t *testing.T) {
	var out bytes.Buffer
	if err := run([]string{"list", "-kind", "component"}, &out); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "component\tButton") || strings.Contains(out.String(), "utility\t") {
		t.Fatalf("list output = %q", out.String())
	}
}

func TestAddWritesSnippet(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "button.go")
	var out bytes.Buffer
	if err := run([]string{"add", "Button", "-o", path}, &out); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "breeze.Button") || !strings.Contains(out.String(), "wrote") {
		t.Fatalf("snippet/output = %q / %q", data, out.String())
	}
}

func TestInitIsIdempotentUnlessForced(t *testing.T) {
	dir := t.TempDir()
	var out bytes.Buffer
	if err := run([]string{"init", dir}, &out); err != nil {
		t.Fatal(err)
	}
	if err := run([]string{"init", dir}, &out); err == nil {
		t.Fatal("second init should refuse to overwrite")
	}
	if err := run([]string{"init", dir, "--force"}, &out); err != nil {
		t.Fatal(err)
	}
}

func TestListPlatformsAndLooks(t *testing.T) {
	var out bytes.Buffer
	if err := run([]string{"list", "-kind", "platform"}, &out); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"platform\tmacOS\tmacOS 26+", "platform\tWindows\tWindows 11", "platform\tLinux\tGNOME 48"} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("list -kind platform = %q, missing %q", out.String(), want)
		}
	}
	out.Reset()
	if err := run([]string{"list"}, &out); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "component\tSwitch\tdrawn") || !strings.Contains(out.String(), "component\tTable\tthemed") {
		t.Fatalf("list = %q", out.String())
	}
}

func TestTokensCommandIsGone(t *testing.T) {
	if err := run([]string{"tokens"}, &bytes.Buffer{}); err == nil {
		t.Fatal("tokens is no longer a command")
	}
}
