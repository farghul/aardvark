// zip.go
package main

import (
	"archive/zip"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func zipFiles(archive string, files ...string) (err error) {
	changeDIR("/data/automation/temp")
	// Create the ZIP file
	zipFile, err := os.Create(archive)
	inspect(err)
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for _, filePath := range files {
		err = processFile(zipWriter, filePath)
		inspect(err)
	}

	return
}

func processFile(zipWriter *zip.Writer, filePath string) error {
	var err error
	var fileInfo fs.FileInfo
	var header *zip.FileHeader
	var headerWriter io.Writer
	var file *os.File

	fileInfo, err = os.Stat(filePath)
	inspect(err)

	header, err = zip.FileInfoHeader(fileInfo)
	inspect(err)

	header.Method = zip.Deflate

	header.Name, err = filepath.Rel(filepath.Dir("."), filePath)
	inspect(err)

	if fileInfo.IsDir() {
		header.Name += "/"
	}

	headerWriter, err = zipWriter.CreateHeader(header)
	inspect(err)

	if fileInfo.IsDir() {
		return nil
	}

	file, err = os.Open(filePath)
	inspect(err)
	defer file.Close()

	_, err = io.Copy(headerWriter, file)
	return err
}
