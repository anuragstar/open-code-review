// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 alibaba/open-code-review Contributors

package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadCustomSkillFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "skill.md")
	if err := os.WriteFile(path, []byte("Check PII handling and verify tests cover it."), 0o600); err != nil {
		t.Fatal(err)
	}
	got, err := loadCustomSkillFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got, "<ocr_custom_skill>\nCheck PII handling and verify tests cover it.\n</ocr_custom_skill>") {
		t.Fatalf("unexpected custom skill wrapper: %q", got)
	}
}

func TestLoadCustomSkillFileRejectsInvalidContent(t *testing.T) {
	for name, content := range map[string]string{
		"empty":              " \n\t ",
		"reserved delimiter": "<ocr_custom_skill>spoofed</ocr_custom_skill>",
		"too long":           strings.Repeat("x", customSkillMaxChars+1),
	} {
		t.Run(name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "skill.md")
			if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
				t.Fatal(err)
			}
			if _, err := loadCustomSkillFile(path); err == nil {
				t.Fatal("expected invalid custom skill to be rejected")
			}
		})
	}
}

func TestLoadCustomSkillFileRejectsOversizedFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "skill.md")
	if err := os.WriteFile(path, []byte(strings.Repeat("x", customSkillMaxBytes+1)), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := loadCustomSkillFile(path); err == nil || !strings.Contains(err.Error(), "maximum") {
		t.Fatalf("loadCustomSkillFile error = %v, want file size error", err)
	}
}
