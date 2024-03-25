package insurgency

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/joaoribeirodasilva/sandstorm_web_admin/webadmin/services/admin_log"
	"github.com/joaoribeirodasilva/sandstorm_web_admin/webadmin/services/config"
)

type Insurgency struct {
	Dir              string `json:"dir"`
	AutomaticUpdates bool   `json:"automaticUpdates"`
	steamcmdPath     string
	log              *admin_log.Log
}

const (
	MODULE = "insurgency"
	GAMEID = 581320
)

func New(config *config.Configuration, log *admin_log.Log, steamcmdPath string) *Insurgency {

	i := new(Insurgency)

	i.log = log
	i.Dir = config.Sandstorm.Dir
	i.AutomaticUpdates = config.Sandstorm.AutomaticUpdates

	return i
}

func (i *Insurgency) IsInstalled() bool {

	var err error
	i.Dir, err = filepath.Abs(i.Dir)
	if err != nil {
		i.log.Write(fmt.Sprintf("failed to calculate absolute path from '%s' relative path. ERR: %s", i.Dir, err.Error()), MODULE, admin_log.LOG_ERROR)
		return false
	}
	return true
}

func (i *Insurgency) Install() bool {

	// steamcmd +force_install_dir ../cs1_ds +login anonymous +app_update 730 +quit
	cmd := exec.Command(i.steamcmdPath, "+force_install_dir", i.Dir, "+login", "anonymous", "+app_update", fmt.Sprintf("%d", 581320, "+quit"))

	return true
}
