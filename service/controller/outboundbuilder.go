package controller

import (
	"encoding/json"
	"fmt"

	"github.com/frainzy1477/Xray/api"
	"github.com/xtls/xray-core/common/net"
	"github.com/xtls/xray-core/core"
	"github.com/xtls/xray-core/infra/conf"
)

//OutboundBuilder build freedom outbund config for addoutbound
func OutboundBuilder(config *Config, nodeInfo *api.NodeInfo) (*core.OutboundHandlerConfig, error) {
	outboundDetourConfig := &conf.OutboundDetourConfig{}
	outboundDetourConfig.Protocol = "freedom"
	outboundDetourConfig.Tag = fmt.Sprintf("%s|%d|%d", nodeInfo.NodeType, nodeInfo.Port, nodeInfo.NodeID)

	// Build Send IP address
	if config.SendIP != "" {
		ipAddress := net.ParseAddress(config.SendIP)
		outboundDetourConfig.SendThrough = &conf.Address{ipAddress}
	}

	// Freedom Protocol setting
	var domainStrategy string = "Asis"
	if config.EnableDNS {
		if config.DNSType != "" {
			domainStrategy = config.DNSType
		} else {
			domainStrategy = "UseIP"
		}
	}
	proxySetting := &conf.FreedomConfig{
		DomainStrategy: domainStrategy,
	}
	
	if nodeInfo.NodeType == "dokodemo-door" {
		proxySetting.Redirect = fmt.Sprintf("0.0.0.0:%d", nodeInfo.Port-1)
	}
	
	var setting json.RawMessage
	setting, err := json.Marshal(proxySetting)
	if err != nil {
		return nil, fmt.Errorf("Marshal proxy %s config fialed: %s", nodeInfo.NodeType, err)
	}

	outboundDetourConfig.Settings = &setting
	return outboundDetourConfig.Build()
}
