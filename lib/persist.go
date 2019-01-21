package configclient

import (
	lib "github.com/cmiceli/configserver/lib"
	"io/ioutil"
)

func WriteConfig(path string, cfg lib.Config) error {
	return ioutil.WriteFile(path, []byte(cfg.Config), 0755)
}
