// SPDX-License-Identifier: Apache-2.0
// Copyright 2026 alibaba/open-code-review Contributors

package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	customSkillOpenTag  = "<ocr_custom_skill>"
	customSkillCloseTag = "</ocr_custom_skill>"
	customSkillMaxChars = 12000
	customSkillMaxBytes = customSkillMaxChars * 4
)

func loadCustomSkillFile(path string) (string, error) {
	content, err := loadPromptFile(path)
	if err != nil {
		return "", fmt.Errorf("read custom skill file %q: %w", path, err)
	}
	if content == "" {
		return "", fmt.Errorf("custom skill file %q is empty after sanitisation", path)
	}
	if strings.Contains(content, customSkillOpenTag) || strings.Contains(content, customSkillCloseTag) {
		return "", fmt.Errorf("custom skill file %q must not contain the reserved delimiters %q or %q", path, customSkillOpenTag, customSkillCloseTag)
	}
	if n := len([]rune(content)); n > customSkillMaxChars {
		return "", fmt.Errorf("custom skill content is %d characters, exceeding the maximum of %d", n, customSkillMaxChars)
	}
	return customSkillOpenTag + "\n" + content + "\n" + customSkillCloseTag, nil
}

func loadPromptFile(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	if info.IsDir() {
		return "", fmt.Errorf("path is a directory, not a file")
	}
	if !info.Mode().IsRegular() {
		return "", fmt.Errorf("path is not a regular file")
	}
	if info.Size() > customSkillMaxBytes {
		return "", fmt.Errorf("file is %d bytes, exceeding the maximum of %d bytes", info.Size(), customSkillMaxBytes)
	}
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	raw, err := io.ReadAll(io.LimitReader(file, customSkillMaxBytes+1))
	if err != nil {
		return "", err
	}
	if len(raw) > customSkillMaxBytes {
		return "", fmt.Errorf("file exceeds the maximum of %d bytes", customSkillMaxBytes)
	}
	return sanitizeMarkdown(string(raw)), nil
}
