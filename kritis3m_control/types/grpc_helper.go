package types

import (
	"fmt"

	v1 "github.com/Laboratory-for-Safe-and-Secure-Systems/kritis3m_scale/gen/go/v1"
)

// UpdateEndpointConfig to Proto conversion
func (db_type UpdateEndpointConfig) Proto() *v1.EndpointConfig {
	return &v1.EndpointConfig{
		Id:                   uint32(db_type.ID),
		Name:                 db_type.UpdateEndpointConfig.Name,
		MutualAuthentication: db_type.UpdateEndpointConfig.MutualAuthentication,
		NoEncryption:         db_type.UpdateEndpointConfig.NoEncryption,
		Kex:                  int32(db_type.UpdateEndpointConfig.ASLKeyExchangeMethod), // Fixed field name
		Cipher:               db_type.UpdateEndpointConfig.Cipher,
	}
}

// UpdateNode to Proto conversion

func (db_type UpdateNode) Proto() *v1.Node {
	var hw_ret []*v1.HardwareConfig
	var proxy_ret []*v1.Proxy

	// Fixed range syntax
	for _, hw_conf := range db_type.UpdateHWConfigs {
		hw_ret = append(hw_ret, hw_conf.Proto())
	}

	// Fixed range syntax and variable names
	for _, proxy := range db_type.UpdateProxies {
		proxy_ret = append(proxy_ret, proxy.Proto())
	}

	return &v1.Node{
		Id:               uint32(db_type.ID),
		SerialNumber:     db_type.UpdateNode.SerialNumber,
		NodeNetworkIndex: uint32(db_type.UpdateNode.NodeNetworkIndex),
		Locality:         db_type.UpdateNode.Locality,
		HwConfig:         hw_ret,
		Proxy:            proxy_ret,
		UserState:        "tobeimplemented",
		ConfigState:      "tobeimplemented",
	}
}

// UpdateHardwareConfig to Proto conversion
func (db_type *UpdateHardwareConfig) Proto() *v1.HardwareConfig {
	return &v1.HardwareConfig{
		Id:     uint32(db_type.ID),
		NodeId: uint32(db_type.UpdateNodeID),
		Device: db_type.UpdateHardwareConfig.Device,
		Cidr:   db_type.UpdateHardwareConfig.IpCidr,
	}
}

func (db_type ApplicationType) Proto() v1.ProxyType {
	switch db_type {
	case ForwardProxy:
		return v1.ProxyType_FORWARD_PROXY
	case ReverseProxy:

		return v1.ProxyType_REVERSE_PROXY
	case TLS_TLSProxy:
		return v1.ProxyType_TLS_TLS_PROXY
	}
	fmt.Println("Error!!!!!!!1 unknown appl type")
	return 3

}

func (db_type *UpdateProxy) Proto() *v1.Proxy {
	return &v1.Proxy{
		Id:                 uint32(db_type.ID),
		NodeId:             uint32(db_type.UpdateNodeID),
		State:              db_type.UpdateProxy.State,
		Type:               db_type.UpdateProxy.Type.Proto(),
		ServerEndpointAddr: db_type.UpdateProxy.ServerEndpointAddr,
		ClientEndpointAddr: db_type.UpdateProxy.ClientEndpointAddr,
		GroupId:            uint32(db_type.GroupID),
	}
}

// UpdateIdentity to Proto conversion
func (db_type UpdateIdentity) Proto() *v1.Identity {
	return &v1.Identity{
		Id:                uint32(db_type.ID),
		ServerUrl:         db_type.UpdateIdentity.ServerUrl,
		RevocationListUrl: db_type.UpdateIdentity.RevocationListUrl,
	}
}

// UpdateGroup to Proto conversion
func (db_type UpdateGroup) Proto() *v1.Group {
	if db_type.UpdateGroup.LegacyEpConfigID == 0 {
		ret := v1.Group{
			Id:                     uint32(db_type.ID),
			Loglevel:               uint32(db_type.UpdateGroup.Loglevel),
			EndpointConfigId:       uint32(db_type.UpdateGroup.EpConfigID),
			LegacyEndpointConfigId: nil,
		}
		return &ret

	} else {
		legacyid := uint32(db_type.UpdateGroup.LegacyEpConfigID)

		ret := v1.Group{
			Id:                     uint32(db_type.ID),
			Loglevel:               uint32(db_type.UpdateGroup.Loglevel),
			EndpointConfigId:       uint32(db_type.UpdateGroup.EpConfigID),
			LegacyEndpointConfigId: &legacyid,
		}
		return &ret
	}
}
