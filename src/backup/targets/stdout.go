package stdout

import (
	"fmt"
	"github.com/jcatana/reaper/config"
)

var GlobalCfg *config.Config

func init() {
	config.GlobalCfg.AddTarget("stdout")
}

func Backup(obj string) error {
	var err error
	err = nil
	fmt.Println(obj)
	return err
}
