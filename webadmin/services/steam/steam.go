package steam

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"time"

	"github.com/joaoribeirodasilva/sandstorm_web_admin/webadmin/services/admin_log"
	"github.com/joaoribeirodasilva/sandstorm_web_admin/webadmin/services/config"
	"github.com/joaoribeirodasilva/sandstorm_web_admin/webadmin/services/utils"
)

type Steam struct {
	Installer         string    `json:"installer"`
	Dir               string    `json:"dir"`
	AutomaticUpdates  bool      `json:"automaticUpdates"`
	LastUpdated       time.Time `json:"lastUpdated"`
	Downloading       bool      `json:"downloading"`
	Updating          bool      `json:"updating"`
	log               *admin_log.Log
	installerFilePath string
	DownloadUrls      map[string]string
}

const (
	MODULE = "steam"
)

func New(conf *config.Configuration, log *admin_log.Log) *Steam {

	s := new(Steam)

	s.Installer = conf.Steam.Installer
	s.Dir = conf.Steam.Dir
	s.AutomaticUpdates = conf.Steam.AutomaticUpdates
	s.DownloadUrls = conf.Steam.DownloadUrls
	s.log = log

	return s
}

func (s *Steam) HasInstaller() bool {

	var err error
	s.Installer, err = filepath.Abs(s.Installer)
	if err != nil {
		s.log.Write(fmt.Sprintf("failed to calculate absolute path from '%s' relative path. ERR: %s", s.Installer, err.Error()), MODULE, admin_log.LOG_ERROR)
		return false
	}

	if !utils.DirectoryExists(s.Installer) {
		s.log.Write(fmt.Sprintf("steam installer directory '%s' dosen't exists! Creating...", s.Installer), MODULE, admin_log.LOG_WARNING)
		if err := os.MkdirAll(s.Installer, 0660); err != nil {
			s.log.Write(fmt.Sprintf("failed to create steam installer directory '%s'. ERR: %s", s.Installer, err.Error()), MODULE, admin_log.LOG_ERROR)
			return false
		}
	}

	path := path.Base(s.DownloadUrls[runtime.GOOS])
	if path == "" {
		s.log.Write(fmt.Sprintf("failed to find installed file for '%s'", runtime.GOOS), MODULE, admin_log.LOG_ERROR)
		return false
	}

	s.installerFilePath, err = filepath.Abs(fmt.Sprintf("%s/%s", s.Installer, path))
	if err != nil {
		s.log.Write(fmt.Sprintf("failed to calculate absolute path from '%s/%s' relative path. ERR: %s", s.Installer, path, err.Error()), MODULE, admin_log.LOG_ERROR)
		return false
	}

	if !utils.FileExists(s.installerFilePath) {
		s.log.Write(fmt.Sprintf("steam installer directory '%s' dosen't exists! Downloading...", s.Installer), MODULE, admin_log.LOG_WARNING)
		if err := s.Download(); err != nil {
			s.log.Write(fmt.Sprintf("failed to get host ip addresses from network interfaces. ERR: %s", err), MODULE, admin_log.LOG_ERROR)
			return false
		}
	} else {
		s.log.Write(fmt.Sprintf("steamsmd installer is at '%s'", s.installerFilePath), MODULE, admin_log.LOG_INFO)
	}

	return true

}

func (s *Steam) IsInstalled() bool {

	var err error

	// ext := ""
	// if runtime.GOOS == "windows" {
	// 	ext = ".exe"
	// } else if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
	// 	ext = ".sh"
	// } else {
	// 	s.log.Write(fmt.Sprintf("operating system '%s' not supported.", runtime.GOOS), MODULE, admin_log.LOG_ERROR)
	// 	return false
	// }

	s.Dir, err = filepath.Abs(s.Dir)
	if err != nil {
		s.log.Write(fmt.Sprintf("failed to calculate absolute path from '%s' relative path. ERR: %s", s.Dir, err.Error()), MODULE, admin_log.LOG_ERROR)
		return false
	}

	path := fmt.Sprintf("%s/steamcmd.sh", s.Dir)
	isInstalled := utils.FileExists(path)
	if isInstalled {
		s.log.Write(fmt.Sprintf("steamcmd is installed at '%s'", path), MODULE, admin_log.LOG_INFO)
	}

	return isInstalled

}

func (s *Steam) Download() error {

	s.Downloading = true

	if !utils.DirectoryExists(s.Installer) {
		s.log.Write(fmt.Sprintf("steam installer directory '%s' dosen't exists! Creating...", s.Installer), MODULE, admin_log.LOG_WARNING)
		if err := os.MkdirAll(s.Installer, 0660); err != nil {
			s.Downloading = false
			return s.log.Write(fmt.Sprintf("failed to create steam installer directory '%s'. ERR: %s", s.Installer, err.Error()), MODULE, admin_log.LOG_ERROR)
		}
	}

	url := s.DownloadUrls[runtime.GOOS]
	if url == "" {
		s.Downloading = false
		return s.log.Write(fmt.Sprintf("invalid download path '%s'. ERR:", url), MODULE, admin_log.LOG_ERROR)
	}
	file := fmt.Sprintf("%s/%s", s.Installer, path.Base(url))
	if err := utils.Download(url, file, s.log); err != nil {
		s.Downloading = false
		return err
	}

	s.Downloading = false
	return nil
}

func (s *Steam) Install() error {

	s.Updating = true
	if err := utils.UntarGz(s.installerFilePath, s.Dir); err != nil {
		s.Updating = false
		return s.log.Write(fmt.Sprintf("failed to extract file '%s' into '%s'. ERR:", s.installerFilePath, s.Dir), MODULE, admin_log.LOG_ERROR)
	}

	s.Updating = false
	return nil
}
