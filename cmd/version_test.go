package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func newTestVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version",
		Run: func(cmd *cobra.Command, args []string) {
			versionCmd.Run(cmd, args)
		},
	}
}

func TestVersion_DefaultOutput(t *testing.T) {
	buf := new(bytes.Buffer)
	cmd := newTestVersionCmd()
	cmd.SetOut(buf)
	cmd.Execute()

	out := buf.String()
	if !strings.HasPrefix(out, "mdpress ") {
		t.Errorf("expected output to start with 'mdpress ', got %q", out)
	}
	if !strings.Contains(out, "commit:") {
		t.Error("expected output to contain 'commit:'")
	}
	if !strings.Contains(out, "built:") {
		t.Error("expected output to contain 'built:'")
	}
}

func TestVersion_DevDefaults(t *testing.T) {
	buf := new(bytes.Buffer)
	cmd := newTestVersionCmd()
	cmd.SetOut(buf)
	cmd.Execute()

	out := buf.String()
	// Default ldflags values
	if !strings.Contains(out, "dev") {
		t.Error("expected default version to be 'dev'")
	}
	if !strings.Contains(out, "none") {
		t.Error("expected default commit to be 'none'")
	}
	if !strings.Contains(out, "unknown") {
		t.Error("expected default date to be 'unknown'")
	}
}

func TestVersion_OutputFormat(t *testing.T) {
	buf := new(bytes.Buffer)
	cmd := newTestVersionCmd()
	cmd.SetOut(buf)
	cmd.Execute()

	out := strings.TrimSpace(buf.String())
	// Should be a single line
	lines := strings.Split(out, "\n")
	if len(lines) != 1 {
		t.Errorf("expected 1 line of output, got %d", len(lines))
	}
}
