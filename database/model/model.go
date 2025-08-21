package model

import (
	"fmt"
	"x-ui/util/json_util"
	"x-ui/xray"
)

type Protocol string

const (
	VMess       Protocol = "vmess"
	VLESS       Protocol = "vless"
	Dokodemo    Protocol = "Dokodemo-door"
	Http        Protocol = "http"
	Socks       Protocol = "socks"
	Trojan      Protocol = "trojan"
	Shadowsocks Protocol = "shadowsocks"
)

type User struct {
	Id       int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Inbound struct {
	Id          int                  `json:"id" form:"id" gorm:"primaryKey;autoIncrement"`
	UserId      int                  `json:"-"`
	Up          int64                `json:"up" form:"up"`
	Down        int64                `json:"down" form:"down"`
	Total       int64                `json:"total" form:"total"`
	Remark      string               `json:"remark" form:"remark"`
	Enable      bool                 `json:"enable" form:"enable"`
	ExpiryTime  int64                `json:"expiryTime" form:"expiryTime"`
	ClientStats []xray.ClientTraffic `gorm:"foreignKey:InboundId;references:Id" json:"clientStats" form:"clientStats"`

	// config part
	Listen         string   `json:"listen" form:"listen"`
	Port           int      `json:"port" form:"port" gorm:"unique"`
	Protocol       Protocol `json:"protocol" form:"protocol"`
	Settings       string   `json:"settings" form:"settings"`
	StreamSettings string   `json:"streamSettings" form:"streamSettings"`
	Tag            string   `json:"tag" form:"tag" gorm:"unique"`
	Sniffing       string   `json:"sniffing" form:"sniffing"`

	// custom backend address fields
	BackendAddress  string `json:"backendAddress" form:"backendAddress" gorm:"column:backend_address"`
	BackendPort     int    `json:"backendPort" form:"backendPort" gorm:"column:backend_port"`
	BackendProtocol string `json:"backendProtocol" form:"backendProtocol" gorm:"column:backend_protocol"`
	EnableBackend   bool   `json:"enableBackend" form:"enableBackend" gorm:"column:enable_backend"`
}
type InboundClientIps struct {
	Id          int    `json:"id" gorm:"primaryKey;autoIncrement"`
	ClientEmail string `json:"clientEmail" form:"clientEmail" gorm:"unique"`
	Ips         string `json:"ips" form:"ips"`
}

func (i *Inbound) GenXrayInboundConfig() *xray.InboundConfig {
	listen := i.Listen
	if listen == "" {
		// 如果listen字段为空，设置默认监听地址
		listen = "\"0.0.0.0\""
	} else {
		listen = fmt.Sprintf("\"%v\"", listen)
	}

	config := &xray.InboundConfig{
		Listen:         json_util.RawMessage(listen),
		Port:           i.Port,
		Protocol:       string(i.Protocol),
		Settings:       json_util.RawMessage(i.Settings),
		StreamSettings: json_util.RawMessage(i.StreamSettings),
		Tag:            i.Tag,
		Sniffing:       json_util.RawMessage(i.Sniffing),
	}

	// Add backend address configuration if enabled (for routing purposes)
	if i.EnableBackend && i.BackendAddress != "" && i.BackendPort > 0 {
		config.BackendAddress = i.BackendAddress
		config.BackendPort = i.BackendPort
		config.EnableBackend = i.EnableBackend

		// 注意：不再修改入站协议的settings
		// 后端代理将通过专门的出站配置和路由规则实现
	}

	return config
}

type Setting struct {
	Id    int    `json:"id" form:"id" gorm:"primaryKey;autoIncrement"`
	Key   string `json:"key" form:"key"`
	Value string `json:"value" form:"value"`
}
type Client struct {
	ID         string `json:"id"`
	AlterIds   uint16 `json:"alterId"`
	Email      string `json:"email"`
	LimitIP    int    `json:"limitIp"`
	Security   string `json:"security"`
	TotalGB    int64  `json:"totalGB" form:"totalGB"`
	ExpiryTime int64  `json:"expiryTime" form:"expiryTime"`
}
