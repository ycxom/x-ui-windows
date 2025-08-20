package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"x-ui/database/model"
	"x-ui/logger"
	"x-ui/util/json_util"
	"x-ui/xray"
	"go.uber.org/atomic"
)

var p *xray.Process
var lock sync.Mutex
var isNeedXrayRestart atomic.Bool
var result string

type XrayService struct {
	inboundService InboundService
	settingService SettingService
}

func (s *XrayService) IsXrayRunning() bool {
	return p != nil && p.IsRunning()
}

func (s *XrayService) GetXrayErr() error {
	if p == nil {
		return nil
	}
	return p.GetErr()
}

func (s *XrayService) GetXrayResult() string {
	if result != "" {
		return result
	}
	if s.IsXrayRunning() {
		return ""
	}
	if p == nil {
		return ""
	}
	result = p.GetResult()
	return result
}

func (s *XrayService) GetXrayVersion() string {
	if p == nil {
		return "Unknown"
	}
	return p.GetVersion()
}
func RemoveIndex(s []interface{}, index int) []interface{} {
	return append(s[:index], s[index+1:]...)
}

func (s *XrayService) GetXrayConfig() (*xray.Config, error) {
	templateConfig, err := s.settingService.GetXrayConfigTemplate()
	if err != nil {
		return nil, err
	}

	xrayConfig := &xray.Config{}
	err = json.Unmarshal([]byte(templateConfig), xrayConfig)
	if err != nil {
		return nil, err
	}

	s.inboundService.DisableInvalidClients()

	inbounds, err := s.inboundService.GetAllInbounds()
	if err != nil {
		return nil, err
	}
	
	// 收集需要后端代理的入站配置
	var backendProxyInbounds []*model.Inbound
	
	for _, inbound := range inbounds {
		if !inbound.Enable {
			continue
		}
		// get settings clients
		settings := map[string]interface{}{}
		json.Unmarshal([]byte(inbound.Settings), &settings)
		clients, ok :=  settings["clients"].([]interface{})
		if ok {
			// check users active or not

			clientStats := inbound.ClientStats
			for _, clientTraffic := range clientStats {
				
				for index, client := range clients {
					c := client.(map[string]interface{})
					if c["email"] == clientTraffic.Email {
						if ! clientTraffic.Enable {
							clients = RemoveIndex(clients,index)
							logger.Info("Remove Inbound User",c["email"] ,"due the expire or traffic limit")

						}

					}
				}
		

			}
			settings["clients"] = clients
			modifiedSettings, err := json.Marshal(settings)
			if err != nil {
				return nil, err
			}
		
			inbound.Settings = string(modifiedSettings)
		}
		inboundConfig := inbound.GenXrayInboundConfig()
		xrayConfig.InboundConfigs = append(xrayConfig.InboundConfigs, *inboundConfig)
		
		// 收集需要后端代理的入站
		if inbound.EnableBackend && inbound.BackendAddress != "" && inbound.BackendPort > 0 {
			backendProxyInbounds = append(backendProxyInbounds, inbound)
		}
	}
	
	// 添加后端代理出站配置和路由规则
	if len(backendProxyInbounds) > 0 {
		err = s.addBackendProxyConfig(xrayConfig, backendProxyInbounds)
		if err != nil {
			return nil, err
		}
	}
	
	return xrayConfig, nil
}

func (s *XrayService) GetXrayTraffic() ([]*xray.Traffic, []*xray.ClientTraffic, error) {
	if !s.IsXrayRunning() {
		return nil, nil, errors.New("xray is not running")
	}
	return p.GetTraffic(true)
}

func (s *XrayService) RestartXray(isForce bool) error {
	lock.Lock()
	defer lock.Unlock()
	logger.Debug("restart xray, force:", isForce)

	xrayConfig, err := s.GetXrayConfig()
	if err != nil {
		return err
	}

	if p != nil && p.IsRunning() {
		if !isForce && p.GetConfig().Equals(xrayConfig) {
			logger.Debug("not need to restart xray")
			return nil
		}
		p.Stop()
	}

	p = xray.NewProcess(xrayConfig)
	result = ""
	return p.Start()
}

func (s *XrayService) StopXray() error {
	lock.Lock()
	defer lock.Unlock()
	logger.Debug("stop xray")
	if s.IsXrayRunning() {
		return p.Stop()
	}
	return errors.New("xray is not running")
}

func (s *XrayService) SetToNeedRestart() {
	isNeedXrayRestart.Store(true)
}

func (s *XrayService) IsNeedRestartAndSetFalse() bool {
	return isNeedXrayRestart.CAS(true, false)
}

// addBackendProxyConfig 添加后端代理出站配置和路由规则
func (s *XrayService) addBackendProxyConfig(xrayConfig *xray.Config, backendProxyInbounds []*model.Inbound) error {
	// 解析现有的出站配置
	var outbounds []map[string]interface{}
	if len(xrayConfig.OutboundConfigs) > 0 {
		err := json.Unmarshal(xrayConfig.OutboundConfigs, &outbounds)
		if err != nil {
			return fmt.Errorf("解析出站配置失败: %v", err)
		}
	}
	
	// 解析现有的路由配置
	var routing map[string]interface{}
	if len(xrayConfig.RouterConfig) > 0 {
		err := json.Unmarshal(xrayConfig.RouterConfig, &routing)
		if err != nil {
			return fmt.Errorf("解析路由配置失败: %v", err)
		}
	}
	if routing == nil {
		routing = make(map[string]interface{})
	}
	
	// 获取现有路由规则
	rules, ok := routing["rules"].([]interface{})
	if !ok {
		rules = make([]interface{}, 0)
	}
	
	// 确保第一个出站配置有标签（避免成为默认出站）
	if len(outbounds) > 0 {
		firstOutbound := outbounds[0]
		if _, hasTag := firstOutbound["tag"]; !hasTag {
			firstOutbound["tag"] = "direct"
		}
	}
	
	// 为每个需要后端代理的入站添加出站配置和路由规则
	for _, inbound := range backendProxyInbounds {
		// 创建后端代理出站配置
		backendOutboundTag := fmt.Sprintf("backend-proxy-%s", inbound.Tag)
		
		// 根据后端协议类型创建出站配置
		backendProtocol := inbound.BackendProtocol
		if backendProtocol == "" {
			backendProtocol = "http" // 默认使用HTTP协议
		}
		
		var outboundConfig map[string]interface{}
		switch backendProtocol {
		case "http":
			outboundConfig = map[string]interface{}{
				"tag":      backendOutboundTag,
				"protocol": "http",
				"settings": map[string]interface{}{
					"servers": []map[string]interface{}{
						{
							"address": inbound.BackendAddress,
							"port":    inbound.BackendPort,
						},
					},
				},
			}
		case "socks":
			outboundConfig = map[string]interface{}{
				"tag":      backendOutboundTag,
				"protocol": "socks",
				"settings": map[string]interface{}{
					"servers": []map[string]interface{}{
						{
							"address": inbound.BackendAddress,
							"port":    inbound.BackendPort,
						},
					},
				},
			}
		case "dokodemo", "dokodemo-door":
			outboundConfig = map[string]interface{}{
				"tag":      backendOutboundTag,
				"protocol": "freedom",
				"settings": map[string]interface{}{
					"domainStrategy": "UseIP",
				},
				"streamSettings": map[string]interface{}{
					"sockopt": map[string]interface{}{
						"dialerProxy": fmt.Sprintf("%s:%d", inbound.BackendAddress, inbound.BackendPort),
					},
				},
			}
		default:
			// 默认使用HTTP协议
			outboundConfig = map[string]interface{}{
				"tag":      backendOutboundTag,
				"protocol": "http",
				"settings": map[string]interface{}{
					"servers": []map[string]interface{}{
						{
							"address": inbound.BackendAddress,
							"port":    inbound.BackendPort,
						},
					},
				},
			}
		}
		
		// 添加出站配置
		outbounds = append(outbounds, outboundConfig)
		
		// 创建路由规则：将特定入站的流量转发到后端代理出站
		routingRule := map[string]interface{}{
			"type":        "field",
			"inboundTag":  []string{inbound.Tag},
			"outboundTag": backendOutboundTag,
		}
		
		// 将新规则插入到现有规则的开头（优先级更高）
		rules = append([]interface{}{routingRule}, rules...)
	}
	
	// 不添加默认路由规则，让其他流量使用第一个出站（direct）
	
	// 更新配置
	updatedOutbounds, err := json.Marshal(outbounds)
	if err != nil {
		return fmt.Errorf("序列化出站配置失败: %v", err)
	}
	xrayConfig.OutboundConfigs = json_util.RawMessage(updatedOutbounds)
	
	routing["rules"] = rules
	updatedRouting, err := json.Marshal(routing)
	if err != nil {
		return fmt.Errorf("序列化路由配置失败: %v", err)
	}
	xrayConfig.RouterConfig = json_util.RawMessage(updatedRouting)
	
	return nil
}
