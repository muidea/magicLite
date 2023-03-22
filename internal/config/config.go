package config

import (
	log "github.com/cihub/seelog"
	fu "github.com/muidea/magicCommon/foundation/util"
)

var configItem *CfgItem

const cfgPath = "/var/app/config/cfg.json"

func init() {
	cfg := &CfgItem{}
	err := fu.LoadConfig(cfgPath, cfg)
	if err != nil {
		log.Errorf("load config failed, err:%s", err.Error())
		return
	}

	configItem = cfg
}

func GetAuthToken() []string {
	if configItem == nil {
		return []string{"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"}
	}

	return configItem.AuthToken
}

type CfgItem struct {
	AuthToken []string `json:"authToken"`
}
