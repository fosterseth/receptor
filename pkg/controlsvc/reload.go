package controlsvc

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"
	"sync"

	"github.com/ansible/receptor/pkg/logger"
	"github.com/ansible/receptor/pkg/netceptor"
	"gopkg.in/yaml.v2"
)

type (
	reloadCommandType struct{}
	reloadCommand     struct{}
)

var configPath = ""

var mu sync.Mutex

var reloadParseAndRun = func(toRun []string) error {
	return fmt.Errorf("no configuration file was provided, reload function not set")
}

var (
	cfgPrevious = make(map[string]struct{})
	cfgNext     = make(map[string]struct{})
)

type actionCallables struct {
	callWhenModifiedorAdded string
	callWhenAbsent          string
}

var reloadableActions = map[string]actionCallables{
	"tcp-peer":     {callWhenModifiedorAdded: "ReloadBackend", callWhenAbsent: ""},
	"tcp-listener": {callWhenModifiedorAdded: "ReloadBackend", callWhenAbsent: ""},
	"ws-peer":      {callWhenModifiedorAdded: "ReloadBackend", callWhenAbsent: ""},
	"ws-listener":  {callWhenModifiedorAdded: "ReloadBackend", callWhenAbsent: ""},
	"udp-peer":     {callWhenModifiedorAdded: "ReloadBackend", callWhenAbsent: ""},
	"udp-listener": {callWhenModifiedorAdded: "ReloadBackend", callWhenAbsent: ""},
	"local-only":   {callWhenModifiedorAdded: "ReloadBackend", callWhenAbsent: ""},
	"log-level":    {callWhenModifiedorAdded: "ReloadLogger", callWhenAbsent: "InitLogger"},
}

func getActionKeyword(cfg string) string {
	// extracts top-level key from the full configuration item
	cfgSplit := strings.Split(cfg, ":")
	var action string
	if len(cfgSplit) == 0 {
		action = cfg
	} else {
		action = cfgSplit[0]
	}

	return action
}

var toRun = make(map[string]struct{})

func parseConfig(filename string, cfgMap map[string]struct{}) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	m := make([]interface{}, 0)
	err = yaml.Unmarshal(data, &m)
	if err != nil {
		return err
	}
	for i := range m {
		cfgBytes, err := yaml.Marshal(&m[i])
		if err != nil {
			return err
		}
		cfg := string(cfgBytes)
		cfgMap[cfg] = struct{}{}
	}

	return nil
}

func checkReload() error {
	// Determine which items from the old config have been modified or added
	for cfg := range cfgNext {
		action := getActionKeyword(cfg)
		_, isReloadable := reloadableActions[action]
		_, inPrevious := cfgPrevious[cfg]
		if !isReloadable && !inPrevious {
			return fmt.Errorf("a non-reloadable config action '%s' was modified or added. Must restart receptor for these changes to take effect", action)
		}
		if isReloadable && !inPrevious {
			callableStr := reloadableActions[action].callWhenModifiedorAdded
			toRun[callableStr] = struct{}{}
		}
	}

	// Determine which items from the old config are absent
	for cfg := range cfgPrevious {
		action := getActionKeyword(cfg)
		_, isReloadable := reloadableActions[action]
		_, inNext := cfgNext[cfg]
		if !isReloadable && !inNext {
			return fmt.Errorf("a non-reloadable config action '%s' was removed. Must restart receptor for changes to take effect", action)
		}
		if isReloadable && !inNext {
			callableStr := reloadableActions[action].callWhenAbsent
			toRun[callableStr] = struct{}{}
		}
	}

	return nil
}

func resetAfterReload() {
	cfgNext = make(map[string]struct{})
	toRun = make(map[string]struct{})
}

// InitReload initializes objects required before reload commands are issued.
func InitReload(cPath string, fParseAndRun func([]string) error) error {
	configPath = cPath
	reloadParseAndRun = fParseAndRun

	return parseConfig(configPath, cfgPrevious)
}

func (t *reloadCommandType) InitFromString(params string) (ControlCommand, error) {
	c := &reloadCommand{}

	return c, nil
}

func (t *reloadCommandType) InitFromJSON(config map[string]interface{}) (ControlCommand, error) {
	c := &reloadCommand{}

	return c, nil
}

func handleError(err error, errorcode int) (map[string]interface{}, error) {
	cfr := make(map[string]interface{})
	cfr["Success"] = false
	cfr["Error"] = fmt.Sprintf("%s ERRORCODE %d", err.Error(), errorcode)
	logger.Warning("Reload not successful: %s", err.Error())

	return cfr, nil
}

func (c *reloadCommand) ControlFunc(ctx context.Context, nc *netceptor.Netceptor, cfo ControlFuncOperations) (map[string]interface{}, error) {
	// grab a mutex, so that only one goroutine can call reload at a time
	mu.Lock()
	defer mu.Unlock()

	logger.Debug("Reloading")
	defer resetAfterReload()

	cfr := make(map[string]interface{})
	cfr["Success"] = true

	// do a quick check to catch any yaml errors before canceling backends
	err := reloadParseAndRun([]string{"PreReload"})
	if err != nil {
		return handleError(err, 4)
	}

	err = parseConfig(configPath, cfgNext)
	if err != nil {
		return handleError(err, 4)
	}

	// check if non-reloadable items have been added or modified
	err = checkReload()
	if err != nil {
		return handleError(err, 3)
	}

	if len(toRun) == 0 {
		logger.Debug("Nothing to reload")

		return cfr, nil
	}

	if _, ok := toRun["ReloadBackend"]; ok {
		nc.CancelBackends()
	}

	// convert the map into a string, which is what the ParseAndRun expects
	toRunStr := []string{}
	for callableStr := range toRun {
		toRunStr = append(toRunStr, callableStr)
	}
	// reloadParseAndRun is a ParseAndRun closure, set in receptor.go/main()
	fmt.Printf("toRun %v\n", toRun)
	err = reloadParseAndRun(toRunStr)
	if err != nil {
		return handleError(err, 4)
	}

	// set old config to new config, only if successful
	cfgPrevious = make(map[string]struct{})
	for cfg := range cfgNext {
		cfgPrevious[cfg] = struct{}{}
	}

	return cfr, nil
}
