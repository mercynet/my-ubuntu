package sys

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"
)

func TestStreamLongLineDoesNotDeadlock(t *testing.T) {
	if os.Getenv("MEU_UBUNTU_TEST_LONG_LINE") == "1" {
		os.Stdout.WriteString(strings.Repeat("x", 2<<20) + "\n")
		os.Exit(0)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := (Real{Env: []string{"MEU_UBUNTU_TEST_LONG_LINE=1"}}).Stream(ctx, func(string) {}, os.Args[0], "-test.run=^TestStreamLongLineDoesNotDeadlock$")
	if ctx.Err() != nil {
		t.Fatal("long output blocked until process timeout")
	}
	if err == nil || !strings.Contains(err.Error(), "token too long") {
		t.Fatalf("expected visible log error, got %v", err)
	}
}
