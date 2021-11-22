package controlsvc

import (
	"testing"
)

func TestReload(t *testing.T) {
	type yamltest struct {
		filename    string
		expectError bool
	}

	scenarios := []yamltest{
		{filename: "reload_test_yml/init.yml", expectError: false},
		{filename: "reload_test_yml/add_cfg.yml", expectError: true},
		{filename: "reload_test_yml/drop_cfg.yml", expectError: true},
		{filename: "reload_test_yml/modify_cfg.yml", expectError: true},
		{filename: "reload_test_yml/syntax_error.yml", expectError: true},
		{filename: "reload_test_yml/successful_reload.yml", expectError: false},
		{filename: "reload_test_yml/change_log.yml", expectError: false},
		{filename: "reload_test_yml/remove_log.yml", expectError: false},
	}
	err := parseConfig("reload_test_yml/init.yml", cfgPrevious)
	if err != nil {
		t.Fatal("could not parse a good-syntax yaml")
	}
	if len(cfgPrevious) != 6 {
		t.Fatal("incorrect cfgPrevious length")
	}

	for i := range scenarios {
		err = parseConfig("reload_test_yml/init.yml", cfgPrevious)
		if err != nil {
			t.Fatalf("could not parse a good-syntax file")
		}
		err = parseConfig(scenarios[i].filename, cfgNext)
		if err != nil && scenarios[i].expectError == false {
			t.Fatal("Could not parse the modified file")
		}
		err = checkReload()
		if err != nil && scenarios[i].expectError == false {
			t.Fatal("Expected error did not occur")
		}
		if err == nil && scenarios[i].expectError == true {
			t.Fatal("Error did not occur , where it was expected")
		}
		resetAfterReload()
	}
}
