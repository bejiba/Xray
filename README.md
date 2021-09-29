申请ssl证书文件

```
yum -y install socat curl

curl -sL https://get.acme.sh | bash

bash /root/.acme.sh/acme.sh --set-default-ca  --server  letsencrypt

bash /root/.acme.sh/acme.sh --issue -d uk.xxx.com --standalone --force

bash /root/.acme.sh/acme.sh --installcert -d uk.xxx.com --fullchainpath /etc/Xray/uk.xxx.com.crt --keypath /etc/Xray/uk.xxx.com.key

```

### 配置文件
config.yml

```
Log:
  Level: none 
  AccessPath: #/etc/Xray/access.Log
  ErrorPath: #/etc/Xray/error.log
DnsConfigPath: # /etc/Xray/dns.json Path to dns config, check https://xtls.github.io/config/base/dns/ for help
RouteConfigPath: # /etc/Xray/route.json # Path to route config, check https://xtls.github.io/config/base/route/ for help
OutboundConfigPath: # /etc/Xray/custom_outbound.json # Path to custom outbound config, check https://xtls.github.io/config/base/outbound/ for help
ConnetionConfig:
  Handshake: 4 # Handshake time limit, Second
  ConnIdle: 3 # Connection idle time limit, Second
  UplinkOnly: 1 # Time limit when the connection downstream is closed, Second
  DownlinkOnly: 1 # Time limit when the connection is closed after the uplink is closed, Second
  BufferSize: 64 # The internal cache size of each connection, kB 
Nodes:
  -
    ApiConfig:
      ApiHost: "https://xxx"
      ApiKey: "123456"
      NodeID: 2
      Timeout: 30 
      SpeedLimit: 0 
      DeviceLimit: 0 
      RuleListPath: #/etc/Xray/rulelist 
    ControllerConfig:
      ListenIP: 0.0.0.0 
      SendIP: 0.0.0.0 
      UpdatePeriodic: 60 
      EnableDNS: true
      DNSType: AsIs # AsIs, UseIP, UseIPv4, UseIPv6, DNS strategy
      EnableProxyProtocol: false # Only works for WebSocket and TCP
      EnableFallback: false # Only support for Trojan and Vless
      FallBackConfigs:  # Support multiple fallbacks
        -
          SNI: # TLS SNI(Server Name Indication), Empty for any
          Path: # HTTP PATH, Empty for any
          Dest: 80 # Required, Destination of fallback, check https://xtls.github.io/config/fallback/ for details.
          ProxyProtocolVer: 0 # Send PROXY protocol version, 0 for dsable
      CertConfig:
        CertMode: file # file, dns, http, none
        CertDomain: "uk.xxx.com" # Domain to cert	
        CertFile: /etc/Xray/uk.xxx.com.crt
        KeyFile: /etc/Xray/uk.xxx.com.key
        Provider: cloudflare #Get the full DNS cert provider support list here: https://go-acme.github.io/lego/dns/
        Email: 
        DNSEnv: # DNS ENV option used by DNS provider
          CLOUDFLARE_EMAIL: 
          CLOUDFLARE_API_KEY: 

```

### Docker Compose 文件

docker-compose.yml

```
version: '3'
services: 
  UK: 
    image: frainzy1477/xray:latest
    volumes:
      - ./config.yml:/etc/Xray/config.yml
      - ./dns.json:/etc/Xray/dns.json
      - ./uk.xxx.com.crt:/etc/Xray/uk.xxx.com.crt
      - ./uk.xxx.com.key:/etc/Xray/uk.xxx.com.key
    restart: always
    network_mode: host

```

# Credits
- 后端 [XrayR-project](https://github.com/XrayR-project/XrayR) 
