package confd

import (
	"reflect"
	"testing"

	"github.com/kelseyhightower/confd/log"
)

func TestInitConfigDefaultConfig(t *testing.T) {
	log.SetQuiet(true)
	want := Config{
		Backend:      "etcd",
		BackendNodes: []string{"http://127.0.0.1:4001"},
		ClientCaKeys: "",
		ClientCert:   "",
		ClientKey:    "",
		ConfDir:      "/etc/confd",
		Debug:        false,
		Interval:     600,
		Noop:         false,
		Prefix:       "/",
		Quiet:        false,
		SRVDomain:    "",
		Scheme:       "http",
		Verbose:      false,
	}
	if err := InitConfig(); err != nil {
		t.Errorf(err.Error())
	}
	if !reflect.DeepEqual(want, Cfg) {
		t.Errorf("initConfig() = %v, want %v", Cfg, want)
	}
}
