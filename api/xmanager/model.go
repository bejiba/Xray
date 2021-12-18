package xmanager

import "encoding/json"

// NodeInfoResponse is the response of node
type NodeInfoResponse struct {
	SpeedLimit      float64 `json:"speedlimit"`
	Method		    string  `json:"method"`
	Port		    int     `json:"port"`
	Address         string  `json:"server"`
	Type            string  `json:"type"`
	Security		string	`json:"security"`
	Host            string  `json:"host"`
	Path            string  `json:"path"`
	Headertype      string  `json:"headertype"`
    Protocol        string  `json:"protocol"`	
	AllowInsecure   bool    `json:"allowinsecure"`
	Relay			int     `json:"relay"`
	RelayNodeID		int     `json:"relayid"`
	ListenIP        string  `json:"listenip"`
}

type RelayNodeInfoResponse struct {
	SpeedLimit      float64 `json:"speedlimit"`
	Method		    string  `json:"method"`
	Port		    int     `json:"port"`
	Address         string  `json:"server"`
	Type            string  `json:"type"`
	Security		string	`json:"security"`
	Host            string  `json:"host"`
	Path            string  `json:"path"`
	Headertype      string  `json:"headertype"`
    Protocol        string  `json:"protocol"`	
	AllowInsecure   bool    `json:"allowinsecure"`
	Relay			int     `json:"relay"`
	RelayNodeID		int     `json:"relayid"`
	ListenIP        string  `json:"listenip"`
}

// UserResponse is the response of user
type UserResponse struct {
	ID            int    `json:"id"`
	Email         string `json:"email"`
	Passwd        string `json:"passwd"`
	SpeedLimit    float64 `json:"speedlimit"`
	DeviceLimit   int    `json:"connector"`
	UUID          string `json:"uuid"`
	IPcount       int    `json:"ip_count"`
	LimitType     int    `json:"limit_type"`
	RelayUser     int    `json:"relay_user"` 
}

// Response is the common response
type Response struct {
	Ret  uint            `json:"ret"`
	Data json.RawMessage `json:"data"`
}

// PostData is the data structure of post data
type PostData struct {
	Data interface{} `json:"data"`
}

// SystemLoad is the data structure of systemload
type SystemLoad struct {
	Uptime string `json:"uptime"`
	Load   string `json:"load"`
}

// OnlineUser is the data structure of online user
type OnlineUser struct {
	UID int    `json:"user_id"`
	IP  string `json:"ip"`
}

// UserTraffic is the data structure of traffic
type UserTraffic struct {
	UID      int   `json:"user_id"`
	Upload   int64 `json:"u"`
	Download int64 `json:"d"`
}

type RuleItem struct {
	ID      int    `json:"id"`
	Content string `json:"regex"`
}

type IllegalItem struct {
	ID  int `json:"list_id"`
	UID int `json:"user_id"`
}
