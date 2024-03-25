package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joaoribeirodasilva/sandstorm_web_admin/webadmin/services/admin_log"
	"github.com/joho/godotenv"
)

type Filesystem struct {
	Base   string `json:"base"`
	Server string `json:"server"`
}

type WebAdmin struct {
	Address          string `json:"address"`
	Port             int    `json:"port"`
	Password         string `json:"password"`
	SslUse           bool   `json:"sslUse"`
	SslVerify        bool   `json:"sslVerify"`
	SslCert          string `json:"sslCert"`
	SslKey           string `json:"sslKey"`
	AutomaticUpdates bool   `json:"automaticUpdates"`
	Dir              string `json:"dir"`
	ConfigDir        string `json:"configDir"`
	Env              string `json:"env"`
	Logs             string `json:"logs"`
}

type Steam struct {
	Installer        string            `json:"installer"`
	Dir              string            `json:"dir"`
	AutomaticUpdates bool              `json:"automaticUpdates"`
	DownloadUrls     map[string]string `json:"downloadUrls"`
}

type Sandstorm struct {
	Dir              string `json:"dir"`
	AutomaticUpdates bool   `json:"automaticUpdates"`
}

type Configuration struct {
	Directories Filesystem     `json:"directories"`
	WebAdmin    WebAdmin       `json:"webAdmin"`
	Steam       Steam          `json:"steam"`
	Sandstorm   Sandstorm      `json:"sandstorm"`
	log         *admin_log.Log `json:"-"`
}

const (
	MODULE = "config"

	FILESYSTEM_BASE   = ".."
	FILESYSTEM_SERVER = FILESYSTEM_BASE + "/server"

	ADMIN_ADDRESS           = "127.0.0.1"
	ADMIN_PORT              = "8080"
	ADMIN_PASSWORD          = ""
	ADMIN_SSL_USE           = "false"
	ADMIN_SSL_VERIFY        = "false"
	ADMIN_SSL_CERT          = ""
	ADMIN_SSL_KEY           = ""
	ADMIN_AUTOMATIC_UPDATES = "true"
	ADMIN_DIR               = "."
	ADMIN_ENV               = "./.env"
	ADMIN_LOGS              = ADMIN_DIR + "/logs"
	ADMIN_CONFIG_DIR        = ADMIN_DIR + "/config"

	STEAM_INSTALLER         = FILESYSTEM_SERVER + "/steam/installer"
	STEAM_DIR               = FILESYSTEM_SERVER + "/steam"
	STEAM_AUTOMATIC_UPDATES = false

	SANDSTORM_DIR               = FILESYSTEM_SERVER + "/sandstorm"
	SANDSTORM_AUTOMATIC_UPDATES = false
)

var envFile string = ADMIN_ENV

func New(log *admin_log.Log) *Configuration {

	var err error
	c := new(Configuration)
	c.log = log
	c.Directories.Base = FILESYSTEM_BASE
	c.Directories.Server = FILESYSTEM_SERVER
	c.WebAdmin.Address = ADMIN_ADDRESS
	c.WebAdmin.Port, err = strconv.Atoi(ADMIN_PORT)
	if err != nil {
		c.log.Write(err.Error(), MODULE, admin_log.LOG_CRITICAL)
	}
	c.WebAdmin.Password = ADMIN_PASSWORD
	c.WebAdmin.SslUse, err = strconv.ParseBool(ADMIN_SSL_USE)
	if err != nil {
		c.log.Write(err.Error(), MODULE, admin_log.LOG_CRITICAL)
	}
	c.WebAdmin.SslVerify, err = strconv.ParseBool(ADMIN_SSL_VERIFY)
	if err != nil {
		c.log.Write(err.Error(), MODULE, admin_log.LOG_CRITICAL)
	}
	c.WebAdmin.SslCert = ADMIN_SSL_CERT
	c.WebAdmin.SslKey = ADMIN_SSL_KEY
	c.WebAdmin.AutomaticUpdates, err = strconv.ParseBool(ADMIN_SSL_VERIFY)
	if err != nil {
		c.log.Write(err.Error(), MODULE, admin_log.LOG_CRITICAL)
	}
	c.WebAdmin.Dir = ADMIN_DIR
	c.WebAdmin.ConfigDir = ADMIN_CONFIG_DIR
	c.WebAdmin.Logs = ADMIN_LOGS

	c.Steam.DownloadUrls = make(map[string]string)
	c.Steam.Installer = STEAM_INSTALLER
	c.Steam.Dir = STEAM_DIR
	c.Steam.AutomaticUpdates = STEAM_AUTOMATIC_UPDATES

	c.Sandstorm.Dir = SANDSTORM_DIR
	c.Sandstorm.AutomaticUpdates = SANDSTORM_AUTOMATIC_UPDATES

	return c
}

func (c *Configuration) SetFile(path string) {
	envFile = path
}

func (c *Configuration) Read() error {

	var err error

	if err = godotenv.Load(envFile); err != nil {
		c.log.Write(fmt.Sprintf("failed to read env file '%s'. ERR: %s", envFile, err.Error()), MODULE, admin_log.LOG_WARNING)
		err = c.Write()
		if err != nil {
			return err
		}
		_ = godotenv.Load(envFile)
	}

	temp := os.Getenv("FILESYSTEM_BASE")
	if temp != "" {
		c.Directories.Base = temp
	}

	temp = os.Getenv("FILESYSTEM_SERVER")
	if temp != "" {
		c.Directories.Server = temp
	}

	temp = os.Getenv("ADMIN_ADDRESS")
	if temp != "" {
		c.WebAdmin.Address = temp
	}

	temp = os.Getenv("ADMIN_PORT")
	if temp != "" {
		c.WebAdmin.Port, err = strconv.Atoi(temp)
		if err != nil {
			err = fmt.Errorf("invalid ADMIN_PORT. ERR: %s", err.Error())
			c.log.Write(err.Error(), MODULE, admin_log.LOG_CRITICAL)
		}
	}

	temp = os.Getenv("ADMIN_PASSWORD")
	if temp != "" {
		c.WebAdmin.Password = temp
	}

	temp = os.Getenv("ADMIN_SSL_USE")
	if temp != "" {
		c.WebAdmin.SslUse, err = strconv.ParseBool(temp)
		if err != nil {
			err = fmt.Errorf("invalid ADMIN_SSL_USE. ERR: %s", err.Error())
			c.log.Write(err.Error(), MODULE, admin_log.LOG_CRITICAL)
		}
	}

	temp = os.Getenv("ADMIN_SSL_VERIFY")
	if temp != "" {
		c.WebAdmin.SslVerify, err = strconv.ParseBool(temp)
		if err != nil {
			err = fmt.Errorf("invalid ADMIN_SSL_VERIFY. ERR: %s", err.Error())
			c.log.Write(err.Error(), MODULE, admin_log.LOG_CRITICAL)
		}
	}

	c.WebAdmin.SslCert = os.Getenv("ADMIN_SSL_CERT")
	c.WebAdmin.SslKey = os.Getenv("ADMIN_SSL_KEY")

	temp = os.Getenv("ADMIN_AUTOMATIC_UPDATES")
	if temp != "" {
		c.WebAdmin.AutomaticUpdates, err = strconv.ParseBool(temp)
		if err != nil {
			err = fmt.Errorf("invalid ADMIN_AUTOMATIC_UPDATES: %s", err.Error())
			c.log.Write(err.Error(), MODULE, admin_log.LOG_CRITICAL)
		}
	}

	temp = os.Getenv("ADMIN_DIR")
	if temp != "" {
		c.WebAdmin.Dir = temp
	}

	temp = os.Getenv("ADMIN_LOGS")
	if temp != "" {
		c.WebAdmin.Logs = temp
	}

	c.Steam.Installer = os.Getenv("STEAM_INSTALLER")
	c.Steam.Dir = os.Getenv("STEAM_DIR")
	temp = os.Getenv("STEAM_AUTOMATIC_UPDATES")
	if temp != "" {
		c.Steam.AutomaticUpdates, err = strconv.ParseBool(temp)
		if err != nil {
			err = fmt.Errorf("invalid STEAM_AUTOMATIC_UPDATES: %s", err.Error())
			c.log.Write(err.Error(), MODULE, admin_log.LOG_CRITICAL)
		}
	}

	c.Steam.DownloadUrls["linux"] = os.Getenv("STEAM_CMD_LINUX")
	c.Steam.DownloadUrls["darwin"] = os.Getenv("STEAM_CMD_OSX")
	c.Steam.DownloadUrls["windows"] = os.Getenv("STEAM_CMD_WINDOWS")

	c.Sandstorm.Dir = os.Getenv("SANDSTORM_DIR")
	temp = os.Getenv("SANDSTORM_AUTOMATIC_UPDATES")
	if temp != "" {
		c.Sandstorm.AutomaticUpdates, err = strconv.ParseBool(temp)
		if err != nil {
			err = fmt.Errorf("invalid SANDSTORM_AUTOMATIC_UPDATES: %s", err.Error())
			c.log.Write(err.Error(), MODULE, admin_log.LOG_CRITICAL)
		}
	}

	return nil
}

func (c *Configuration) Write() error {

	f, err := os.OpenFile(ADMIN_ENV, os.O_RDWR|os.O_CREATE, 0660)
	if err != nil {
		return c.log.Write(err.Error(), MODULE, admin_log.LOG_CRITICAL)
	}
	defer f.Close()

	return nil
}
