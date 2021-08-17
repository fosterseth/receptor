package utils

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/ghjm/cmdline"
)

var cfgNotReloadable map[string]bool

func cfgToString(cfg interface{}) (string, error) {
	cfgBytes, err := json.Marshal(cfg)
	if err != nil {
		return "", err
	}
	t := reflect.TypeOf(cfg)
	cfgStr := t.Name() + string(cfgBytes)

	return cfgStr, nil
}

func DisableReload(cfg interface{}) error {
	cfgStr, err := cfgToString(cfg)
	if err != nil {
		return err
	}
	cfgNotReloadable[cfgStr] = false

	return nil
}

func ErrorCfgChangedOrNew(cfg interface{}) error {
	cfgStr, _ := cfgToString(cfg)
	if _, ok := cfgNotReloadable[cfgStr]; ok {
		cfgNotReloadable[cfgStr] = true

		return nil
	}

	return fmt.Errorf("%v was modified or added. Must restart receptor for changes to take effect", reflect.TypeOf(cfg))
}

func reset() {
	for k := range cfgNotReloadable {
		cfgNotReloadable[k] = false
	}
}

func ErrorCfgAbsent() error {
	defer reset()
	for _, v := range cfgNotReloadable {
		if !v {
			return fmt.Errorf("non-reloadable items were removed from configuration file. Must restart receptor for changes to take effect")
		}
	}

	return nil
}

func init() {
	cfgNotReloadable = make(map[string]bool)
	cmdline.RegisterFuncForApp("ErrorCfgChangedOrNew", ErrorCfgChangedOrNew)
	cmdline.RegisterFuncForApp("DisableReload", DisableReload)
}
