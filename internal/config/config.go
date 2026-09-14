// Package config provides configuration loading for the parameters-kerberos service.
package config

import (
	coreconfig "github.com/jasonmiller-cc/parameters-core/pkg/config"
)

// KerberosConfig holds Kerberos KDC connection settings.
type KerberosConfig struct {
	Realm            string `yaml:"realm"             env:"KERBEROS_REALM"`
	KAdminServer     string `yaml:"kadmin_server"     env:"KERBEROS_KADMIN_SERVER"`
	KAdminPrincipal  string `yaml:"kadmin_principal"  env:"KERBEROS_KADMIN_PRINCIPAL"`
	KeytabPath       string `yaml:"keytab_path"       env:"KERBEROS_KEYTAB_PATH"`
	KDCHost          string `yaml:"kdc_host"          env:"KERBEROS_KDC_HOST"`
}

// Config is the top-level configuration for parameters-kerberos.
type Config struct {
	coreconfig.BaseConfig `yaml:",inline"`
	Kerberos              KerberosConfig `yaml:"kerberos"`
}

// Load reads config from path (falling back to standard locations) and returns
// a populated *Config.
func Load(path string) (*Config, error) {
	var cfg Config
	if err := coreconfig.Load(path, "PARAMS_KERBEROS", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
