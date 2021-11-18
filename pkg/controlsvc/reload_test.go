package controlsvc

import (
	"testing"
)

func TestReload(t *testing.T) {
	type yamltest struct {
		filename    string
		modifyError bool
		absentError bool
	}

	scenarios := []yamltest{
		{filename: "reload_test_yml/init.yml", modifyError: false, absentError: false},
		{filename: "reload_test_yml/add_cfg.yml", modifyError: true, absentError: false},
		{filename: "reload_test_yml/drop_cfg.yml", modifyError: false, absentError: true},
		{filename: "reload_test_yml/modify_cfg.yml", modifyError: true, absentError: true},
		{filename: "reload_test_yml/syntax_error.yml", modifyError: true, absentError: true},
		{filename: "reload_test_yml/successful_reload.yml", modifyError: false, absentError: false},
	}
	err := parseConfig("reload_test_yml/init.yml", cfgPrevious)
	if err != nil {
		t.Fatal("could not parse a good-syntax yaml")
	}
	if len(cfgPrevious) != 6 {
		t.Fatal("incorrect cfgPrevious length")
	}

	for _, s := range scenarios {
		t.Logf("%s", s.filename)
		parseConfig(s.filename, cfgNext)
		err = checkReload()
		// t.Logf("%v\n", cfgNext)
		if s.modifyError || s.absentError {
			if err == nil {
				t.Fatal("error expected")
			}
		} else {
			if err != nil {
				t.Fatal("did not expect error")
			}
		}
		resetAfterReload()
	}
}
