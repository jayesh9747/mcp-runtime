package cli

import (
	"bytes"
	"os"
	"testing"
)

func TestPrintSection(t *testing.T) {
	var buf bytes.Buffer
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	printSection("Test Section")

	w.Close()
	os.Stdout = old
	buf.ReadFrom(r)
	output := buf.String()

	if output == "" {
		t.Error("printSection should output something")
	}
}

func TestPrintStep(t *testing.T) {
	var buf bytes.Buffer
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	printStep("Test Step")

	w.Close()
	os.Stdout = old
	buf.ReadFrom(r)
	output := buf.String()

	if output == "" {
		t.Error("printStep should output something")
	}
}

func TestPrintInfo(t *testing.T) {
	var buf bytes.Buffer
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	printInfo("Test Info")

	w.Close()
	os.Stdout = old
	buf.ReadFrom(r)
	output := buf.String()

	if output == "" {
		t.Error("printInfo should output something")
	}
}

func TestPrintWarn(t *testing.T) {
	var buf bytes.Buffer
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	printWarn("Test Warn")

	w.Close()
	os.Stdout = old
	buf.ReadFrom(r)
	output := buf.String()

	if output == "" {
		t.Error("printWarn should output something")
	}
}

func TestPrintError(t *testing.T) {
	var buf bytes.Buffer
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	printError("Test Error")

	w.Close()
	os.Stdout = old
	buf.ReadFrom(r)
	output := buf.String()

	if output == "" {
		t.Error("printError should output something")
	}
}

func TestPrintSuccess(t *testing.T) {
	var buf bytes.Buffer
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	printSuccess("Test Success")

	w.Close()
	os.Stdout = old
	buf.ReadFrom(r)
	output := buf.String()

	if output == "" {
		t.Error("printSuccess should output something")
	}
}
