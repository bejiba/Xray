package panel

import (
	"github.com/frainzy1477/Xray/api"
	"github.com/frainzy1477/Xray/service/controller"
)

type Config struct {
	LogConfig          *LogConfig       `mapstructure:"Log"`
	DnsConfigPath      string           `mapstructure:"DnsConfigPath"`
	OutboundConfigPath string           `mapstructure:"OutboundConfigPath"`
	RouteConfigPath    string           `mapstructure:"RouteConfigPath"`
	ConnetionConfig    *ConnetionConfig `mapstructure:"ConnetionConfig"`
	NodesConfig        []*NodesConfig   `mapstructure:"Nodes"`
}

type NodesConfig struct {
	ApiConfig        *api.Config        `mapstructure:"ApiConfig"`
	ControllerConfig *controller.Config `mapstructure:"ControllerConfig"`
}

type LogConfig struct {
	Level      string `mapstructure:"Level"`
	AccessPath string `mapstructure:"AccessPath"`
	ErrorPath  string `mapstructure:"ErrorPath"`
}

type ConnetionConfig struct {
	Handshake    uint32 `mapstructure:"handshake"`
	ConnIdle     uint32 `mapstructure:"connIdle"`
	UplinkOnly   uint32 `mapstructure:"uplinkOnly"`
	DownlinkOnly uint32 `mapstructure:"downlinkOnly"`
	BufferSize   int32  `mapstructure:"bufferSize"`
}