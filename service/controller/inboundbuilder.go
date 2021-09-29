//Package generate the InbounderConfig used by add inbound
package controller

import (
	"encoding/json"
	"fmt"

	"github.com/frainzy1477/Xray/api"
	"github.com/frainzy1477/Xray/common/legocmd"
	"github.com/xtls/xray-core/common/net"
	"github.com/xtls/xray-core/common/uuid"
	"github.com/xtls/xray-core/core"
	"github.com/xtls/xray-core/infra/conf"
)

//InboundBuilder build Inbound config for different protocol
func InboundBuilder(config *Config, nodeInfo *api.NodeInfo) (*core.InboundHandlerConfig, error) {
	inboundDetourConfig := &conf.InboundDetourConfig{}
	// Build Listen IP address
	if config.ListenIP != "" {
		ipAddress := net.ParseAddress(config.ListenIP)
		inboundDetourConfig.ListenOn = &conf.Address{ipAddress}
	}else{
		inboundDetourConfig.ListenOn = &conf.Address{net.ParseAddress("0.0.0.0")}
	}

	// Build Port
	portRange := &conf.PortRange{From: uint32(nodeInfo.Port), To: uint32(nodeInfo.Port)}
	inboundDetourConfig.PortRange = portRange
	// Build Tag
	inboundDetourConfig.Tag = fmt.Sprintf("%s|%d|%d", nodeInfo.NodeType, nodeInfo.Port, nodeInfo.NodeID)
	// SniffingConfig
	sniffingConfig := &conf.SniffingConfig{
		Enabled:      true,
		DestOverride: &conf.StringList{"http", "tls"},
	}
	inboundDetourConfig.SniffingConfig = sniffingConfig

	var (
		protocol      string
		streamSetting *conf.StreamConfig
		setting       json.RawMessage
	)

	var proxySetting interface{}
	// Build Protocol and Protocol setting
	if nodeInfo.NodeType == "V2ray" {
		if nodeInfo.EnableVless {
			protocol = "vless"
			// Enable fallback
			if config.EnableFallback {
				fallbackConfigs, err := buildVlessFallbacks(config.FallBackConfigs)
				if err == nil {
					proxySetting = &conf.VLessInboundConfig{
						Decryption: "none",
						Fallbacks:  fallbackConfigs,
					}
				} else {
					return nil, err
				}
			} else {
				proxySetting = &conf.VLessInboundConfig{
					Decryption: "none",
				}
			}
		} else {
			protocol = "vmess"
			proxySetting = &conf.VMessInboundConfig{}
		}
	} else if nodeInfo.NodeType == "Trojan" {
		protocol = "trojan"
		// Enable fallback
		if config.EnableFallback {
			fallbackConfigs, err := buildTrojanFallbacks(config.FallBackConfigs)
			if err == nil {
				proxySetting = &conf.TrojanServerConfig{
					Fallbacks: fallbackConfigs,
				}
			} else {
				return nil, err
			}
		} else {
			proxySetting = &conf.TrojanServerConfig{}
		}
	} else if nodeInfo.NodeType == "Shadowsocks" || nodeInfo.NodeType == "Shadowsocks-Plugin" {
		protocol = "shadowsocks"
		proxySetting = &conf.ShadowsocksServerConfig{}
		randomPasswd := uuid.New()
		defaultSSuser := &conf.ShadowsocksUserConfig{
			Cipher:   nodeInfo.CypherMethod,
			Password: randomPasswd.String(),
		}
		proxySetting, _ := proxySetting.(*conf.ShadowsocksServerConfig)
		proxySetting.Users = append(proxySetting.Users, defaultSSuser)
		proxySetting.NetworkList = &conf.NetworkList{"tcp", "udp"}
	} else if nodeInfo.NodeType == "dokodemo-door" {
		protocol = "dokodemo-door"
		proxySetting = struct {
			Host        string   `json:"address"`
			Redirect    bool     `json:"followRedirect"`
			NetworkList []string `json:"network"`
		}{
			Host:        "v1.mux.cool",
			Redirect:    false,
			NetworkList: []string{"tcp", "udp"},
		}
	}else {
		return nil, fmt.Errorf("Unsupported node type: %s, Only support: V2ray, Trojan, Shadowsocks, and Shadowsocks-Plugin", nodeInfo.NodeType)
	}

	setting, err := json.Marshal(proxySetting)
	if err != nil {
		return nil, fmt.Errorf("Marshal proxy %s config fialed: %s", nodeInfo.NodeType, err)
	}

	// Build streamSettings
	streamSetting = new(conf.StreamConfig)
	transportProtocol := conf.TransportProtocol(nodeInfo.TransportProtocol)
	networkType, err := transportProtocol.Build()
	if err != nil {
		return nil, fmt.Errorf("convert TransportProtocol failed: %s", err)
	}
	if networkType == "tcp" {
		headers := make(map[string]string)
		headers["type"] = nodeInfo.HeaderType
		var header json.RawMessage
		header, err  := json.Marshal(headers)
		if err != nil {
			return nil, fmt.Errorf("Marshal Header Type %s into config fialed: %s", header, err)
		}
		tcpSetting := &conf.TCPConfig{
			AcceptProxyProtocol: config.EnableProxyProtocol,
			HeaderConfig:        header,
		}
		streamSetting.TCPSettings = tcpSetting
	} else if networkType == "websocket" {
		headers := make(map[string]string)
		headers["Host"] = nodeInfo.Host
		wsSettings := &conf.WebSocketConfig{
			AcceptProxyProtocol: config.EnableProxyProtocol,
			Path:                nodeInfo.Path,
			Headers:             headers,
		}
		streamSetting.WSSettings = wsSettings
	} else if networkType == "http" {
		hosts := conf.StringList{nodeInfo.Host}
		httpSettings := &conf.HTTPConfig{
			Host: &hosts,
			Path: nodeInfo.Path,
		}
		streamSetting.HTTPSettings = httpSettings
	}else if networkType == "grpc" {
		grpcSettings := &conf.GRPCConfig{
			ServiceName: nodeInfo.ServiceName,
		}
		streamSetting.GRPCConfig = grpcSettings
	}

	streamSetting.Network = &transportProtocol
	// Build TLS and XTLS settings
	if nodeInfo.EnableTLS && config.CertConfig.CertMode != "none" {
		streamSetting.Security = nodeInfo.TLSType
		certFile, keyFile, err := getCertFile(config.CertConfig)
		if err != nil {
			return nil, err
		}
		if nodeInfo.TLSType == "tls" {
			tlsSettings := &conf.TLSConfig{}
			tlsSettings.Certs = append(tlsSettings.Certs, &conf.TLSCertConfig{CertFile: certFile, KeyFile: keyFile, OcspStapling: 3600})

			streamSetting.TLSSettings = tlsSettings
		} else if nodeInfo.TLSType == "xtls" {
			xtlsSettings := &conf.XTLSConfig{}
			xtlsSettings.Certs = append(xtlsSettings.Certs, &conf.XTLSCertConfig{CertFile: certFile, KeyFile: keyFile, OcspStapling: 3600})
			streamSetting.XTLSSettings = xtlsSettings
		}
	}

	inboundDetourConfig.Protocol = protocol
	inboundDetourConfig.StreamSetting = streamSetting
	inboundDetourConfig.Settings = &setting

	return inboundDetourConfig.Build()
}

func getCertFile(certConfig *CertConfig) (certFile string, keyFile string, err error) {
	if certConfig.CertMode == "file" {
		if certConfig.CertFile == "" || certConfig.KeyFile == "" {
			return "", "", fmt.Errorf("Cert file path or key file path not exist")
		}
		return certConfig.CertFile, certConfig.KeyFile, nil
	} else if certConfig.CertMode == "dns" {
		lego, err := legocmd.New()
		if err != nil {
			return "", "", err
		}
		certPath, keyPath, err := lego.DNSCert(certConfig.CertDomain, certConfig.Email, certConfig.Provider, certConfig.DNSEnv)
		if err != nil {
			return "", "", err
		}
		return certPath, keyPath, err
	} else if certConfig.CertMode == "http" {
		lego, err := legocmd.New()
		if err != nil {
			return "", "", err
		}
		certPath, keyPath, err := lego.HTTPCert(certConfig.CertDomain, certConfig.Email)
		if err != nil {
			return "", "", err
		}
		return certPath, keyPath, err
	}

	return "", "", fmt.Errorf("Unsupported certmode: %s", certConfig.CertMode)
}

func buildVlessFallbacks(fallbackConfigs []*FallBackConfig) ([]*conf.VLessInboundFallback, error) {
	if fallbackConfigs == nil {
		return nil, fmt.Errorf("You must provide FallBackConfigs")
	}
	vlessFallBacks := make([]*conf.VLessInboundFallback, len(fallbackConfigs))
	for i, c := range fallbackConfigs {

		if c.Dest == "" {
			return nil, fmt.Errorf("Dest is required for fallback fialed")
		}

		var dest json.RawMessage
		dest, err := json.Marshal(c.Dest)
		if err != nil {
			return nil, fmt.Errorf("Marshal dest %s config fialed: %s", dest, err)
		}
		vlessFallBacks[i] = &conf.VLessInboundFallback{
			Name: c.SNI,
			Path: c.Path,
			Dest: dest,
			Xver: c.ProxyProtocolVer,
		}
	}
	return vlessFallBacks, nil
}

func buildTrojanFallbacks(fallbackConfigs []*FallBackConfig) ([]*conf.TrojanInboundFallback, error) {
	if fallbackConfigs == nil {
		return nil, fmt.Errorf("You must provide FallBackConfigs")
	}	
	trojanFallBacks := make([]*conf.TrojanInboundFallback, len(fallbackConfigs))
	for i, c := range fallbackConfigs {

		if c.Dest == "" {
			return nil, fmt.Errorf("Dest is required for fallback fialed")
		}

		var dest json.RawMessage
		dest, err := json.Marshal(c.Dest)
		if err != nil {
			return nil, fmt.Errorf("Marshal dest %s config fialed: %s", dest, err)
		}
		trojanFallBacks[i] = &conf.TrojanInboundFallback{
			Name: c.SNI,
			Path: c.Path,
			Dest: dest,
			Xver: c.ProxyProtocolVer,
		}
	}
	return trojanFallBacks, nil
}