package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/joaoribeirodasilva/sandstorm_web_admin/webadmin/services/admin_log"
	"github.com/walle/targz"
)

func TarGz(source string, dest string) error {
	return targz.Compress(source, dest)
}

func UntarGz(source string, dest string) error {
	return targz.Extract(source, dest)
}

func Unzip(source string, dest string, log *admin_log.Log) error {

	if log != nil {
		log.Write(fmt.Sprintf("check if source file '%s' exists", source), MODULE, admin_log.LOG_DEBUG)
	}

	if !FileExists(source) {
		return fmt.Errorf("source file '%s' doesn't exist", source)
	}

	if log != nil {
		log.Write(fmt.Sprintf("check if destination directory '%s' exists", dest), MODULE, admin_log.LOG_DEBUG)
	}

	if !DirectoryExists(dest) {
		return fmt.Errorf("destination directory '%s' doesn't exist", dest)
	}

	if log != nil {
		log.Write(fmt.Sprintf("opening zip file '%s'", source), MODULE, admin_log.LOG_DEBUG)
	}

	archive, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer archive.Close()

	for _, f := range archive.File {

		filePath := filepath.Join(dest, f.Name)

		if log != nil {
			log.Write(fmt.Sprintf("unzipping file '%s'", f.Name), MODULE, admin_log.LOG_DEBUG)
		}

		if !strings.HasPrefix(filePath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid path '%s'", filePath)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return fmt.Errorf("error creating path '%s'. ERR: %s", filePath, err.Error())
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return fmt.Errorf("error writing file '%s'. ERR: %s", filePath, err.Error())
		}

		fileInArchive, err := f.Open()
		if err != nil {
			return fmt.Errorf("failed to open file '%s' from within zip file. ERR: %s", filePath, err.Error())
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			return fmt.Errorf("failed to copy file '%s' from within zip file. ERR: %s", filePath, err.Error())
		}

		dstFile.Close()
		fileInArchive.Close()

		if log != nil {
			log.Write(fmt.Sprintf("file '%s' unzipped into '%s'", f.Name, filePath), MODULE, admin_log.LOG_DEBUG)
		}
	}

	return nil
}
