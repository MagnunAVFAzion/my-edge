package helpers

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func RenameFile(oldName, newName string) error {
	err := os.Rename(oldName, newName)
	if err != nil {
		return fmt.Errorf("error renaming file: %w", err)
	}

	return nil
}

func ReplaceStrInFile(filePath, oldPattern, newPattern string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	newContent := strings.ReplaceAll(string(content), oldPattern, newPattern)

	err = os.WriteFile(filePath, []byte(newContent), 0644)
	if err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}

	return nil
}

func CreateDirIfNotExists(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err := os.MkdirAll(dirPath, 0755)
		if err != nil {
			return fmt.Errorf("error creating directory: %w", err)
		}
	}

	fmt.Println("Dir created:", dirPath)

	return nil
}

func CreateFileIfNotExists(dirPath, filename, content string) error {
	err := CreateDirIfNotExists(dirPath)
	if err != nil {
		return err
	}

	filePath := dirPath + "/" + filename

	if _, err := os.Stat(filePath); err == nil {
		fmt.Println("File already exists:", filePath)
		return nil
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("error checking file existence: %w", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	fmt.Println("File created:", filePath)
	return nil
}

func ZipDir(srcDir, zipFilePath string) error {
	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		return fmt.Errorf("error creating zip file: %w", err)
	}
	defer zipFile.Close()

	zipWritter := zip.NewWriter(zipFile)
	defer zipWritter.Close()

	err = filepath.WalkDir(srcDir, func(filePath string, f fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if f.IsDir() {
			return nil
		}

		fileInfo, err := f.Info()
		if err != nil {
			return fmt.Errorf("error getting file info: %w", err)
		}

		header, err := zip.FileInfoHeader(fileInfo)
		if err != nil {
			return fmt.Errorf("error creating zip header: %w", err)
		}

		header.Name, err = filepath.Rel(srcDir, filePath)
		if err != nil {
			return fmt.Errorf("error getting zip header name: %w", err)
		}

		fileWriter, err := zipWritter.CreateHeader(header)
		if err != nil {
			return fmt.Errorf("error creating zip file writer: %w", err)
		}

		file, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("error opening file: %w", err)
		}
		defer file.Close()

		_, err = io.Copy(fileWriter, file)
		if err != nil {
			return fmt.Errorf("error copying file to zip: %w", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("error copying zip file: %w", err)
	}

	return nil
}

func RunShellCommand(command string) (string, error) {
	cmd := exec.Command("bash", "-c", command)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("error running command: %s, stderr: %s", err, stderr.String())
	}

	return out.String(), nil
}

func GetCurrentUnixTimestamp() string {
	now := time.Now()
	epochSeconds := now.Unix()
	timestampString := strconv.FormatInt(epochSeconds, 10)

	return timestampString
}
