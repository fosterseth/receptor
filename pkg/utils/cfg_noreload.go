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

func MarkforNoReload(cfg interface{}) error {
	cfgStr, err := cfgToString(cfg)
	if err != nil {
		return err
	}
	cfgNotReloadable[cfgStr] = false

	return nil
}

func isPresent(cfg interface{}) bool {
	cfgStr, _ := cfgToString(cfg)
	_, ok := cfgNotReloadable[cfgStr]
	if ok {
		cfgNotReloadable[cfgStr] = true
	}

	return ok
}

func ErrorIfCfgChanged(cfg interface{}) error {
	fmt.Printf("%v", cfgNotReloadable)
	if !isPresent(cfg) {
		return fmt.Errorf("%v was modified or added. Must restart receptor for changes to take effect", reflect.TypeOf(cfg))
	}

	return nil
}

func reset() {
	for k := range cfgNotReloadable {
		cfgNotReloadable[k] = false
	}
}

func ErrorIfAbsent() error {
	defer reset()
	fmt.Printf("%v", cfgNotReloadable)
	for _, v := range cfgNotReloadable {
		if !v {
			return fmt.Errorf("non-reloadable items were removed from configuration file. Must restart receptor for changes to take effect")
		}
	}

	return nil
}

func EnableReload(cfg interface{}) error {
	cfgStr, err := cfgToString(cfg)
	if err != nil {
		return err
	}
	cfgNotReloadable[cfgStr] = true

	return nil
}

func init() {
	cfgNotReloadable = make(map[string]bool)
	cmdline.RegisterFuncForApp("ErrorIfCfgChanged", ErrorIfCfgChanged)
	cmdline.RegisterFuncForApp("MarkforNoReload", MarkforNoReload)
}
