package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func unzip(zipFile, destDir string) error {
	changeDIR("/data/automation/temp")
	reader, err := zip.OpenReader(zipFile)
	inspect(err)
	defer reader.Close()

	destDir = filepath.Clean(destDir)

	for _, file := range reader.File {
		if err := extractFile(file, destDir); err != nil {
			return err
		}
	}

	return nil
}

func extractFile(file *zip.File, destDir string) error {
	destPath := filepath.Join(destDir, file.Name)
	destPath = filepath.Clean(destPath)

	// Check for file traversal attack
	if !strings.HasPrefix(destPath, destDir) {
		return fmt.Errorf("invalid file path: %s", file.Name)
	}

	if file.FileInfo().IsDir() {
		if err := os.MkdirAll(destPath, file.Mode()); err != nil {
			return err
		}
	} else {
		if err := os.MkdirAll(filepath.Dir(destPath), os.ModePerm); err != nil {
			return err
		}

		destFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		inspect(err)
		defer destFile.Close()

		srcFile, err := file.Open()
		inspect(err)
		defer srcFile.Close()

		if _, err := io.Copy(destFile, srcFile); err != nil {
			return err
		}
	}

	return nil
}
