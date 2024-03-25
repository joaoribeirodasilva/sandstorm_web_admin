package utils

import (
	"io/fs"
	"os"
)

func FileExists(path string) bool {

	fi := info(path)
	if fi == nil {
		return false
	}

	return isFile(fi)
}

func DirectoryExists(path string) bool {

	fi := info(path)
	if fi == nil {
		return false
	}

	return isDirectory(fi)
}

func info(path string) fs.FileInfo {

	fi, err := os.Stat(path)
	if err != nil {
		return nil
	}

	return fi
}

func isDirectory(fi fs.FileInfo) bool {

	return fi.IsDir()
}

func isFile(fi fs.FileInfo) bool {

	return !fi.IsDir()
}
