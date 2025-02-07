package types

import (
	"time"

	asl "github.com/Laboratory-for-Safe-and-Secure-Systems/go-asl"
	"gorm.io/gorm"
)

type ApplicationType uint8

const (
	ForwardProxy ApplicationType = 0
	ReverseProxy ApplicationType = 1
	TLS_TLSProxy ApplicationType = 2 //server and client tls endpoint
)

type IdentityType uint8

const (
	Production IdentityType = 0
	Management IdentityType = 1
)

type NodeState int8

const (
	ErrorState          NodeState = 0
	NotSeen             NodeState = 1
	NodeRequestedConfig NodeState = 2
	Running             NodeState = 3
)

type HardwareConfig struct {
	ID     uint `gorm:"primarykey" json:"-"`
	NodeID uint `json:"-"`

	Device string `json:"device"`
	IpCidr string `json:"cidr"`

	UpdateHardwareConfig   *HardwareConfig `gorm:"foreignKey:UpdateHardwareConfigID"`
	UpdateHardwareConfigID uint
} // Node represents a node within a network

type UpdateHardwareConfig struct {
	ID uint `gorm:"primarykey" json:"-"`

	UpdateHardwareConfig HardwareConfig `gorm:"embedded"`
	UpdateNodeID         uint           `json:"-"`

	ActiveHardwareConfig   *HardwareConfig `gorm:"foreignKey:ActiveHardwareConfigID"`
	ActiveHardwareConfigID uint
} // Node represents a node within a network

type Node struct {
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	ID        uint           `gorm:"primarykey" json:"id,omitempty"`

	SerialNumber     string `gorm:"uniqueindex" json:"serial_number"`
	NodeNetworkIndex uint   `json:"network_index"`
	Locality         string `json:"locality,omitempty"`

	LastSeen  time.Time         `json:"-"`
	HWConfigs []*HardwareConfig `gorm:"foreignKey:NodeID"`
	Proxies   []*Proxy          `gorm:"foreignKey:NodeID" json:"applications"`

	UpdateID   uint
	UpdateNode *UpdateNode `gorm:"foreignKey:UpdateID"`
}

type UpdateNode struct {
	ID uint

	UpdateNode Node `gorm:"embedded"`

	UpdateHWConfigs []*UpdateHardwareConfig `gorm:"foreignKey:UpdateNodeID"`
	UpdateProxies   []*UpdateProxy          `gorm:"foreignKey:UpdateNodeID" json:"applications"`

	Userstate bool
	Nodestate bool

	ActiveID   uint
	ActiveNode *Node `gorm:"foreignKey:ActiveID"`
}

type Proxy struct {
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	LogLevel  uint           `json:"log_level,omitempty" gorm:"default:3"`

	ID     uint `gorm:"primarykey" json:"id"`
	NodeID uint

	State bool            `json:",omitempty"`
	Type  ApplicationType `json:"type"`

	ServerEndpointAddr string `json:"server_endpoint_addr"`
	ClientEndpointAddr string `json:"client_endpoint_addr"`

	Group   *Group `gorm:"foreignKey:GroupID"`
	GroupID uint

	UpdateProxy   *UpdateProxy `gorm:"foreignKey:UpdateProxyID"`
	UpdateProxyID uint
}

type UpdateProxy struct {
	UpdateProxy  Proxy `gorm:"embedded"`
	ID           uint
	UpdateNodeID uint

	Group   *UpdateGroup `gorm:"foreignKey:GroupID"`
	GroupID uint

	ActiveProxy   *Proxy `gorm:"foreignKey:ActiveProxyID"`
	ActiveProxyID uint
}

type Group struct {
	gorm.Model
	ID uint

	EpConfigID uint
	EpConfig   *EndpointConfig `gorm:"foreignKey:EpConfigID"`

	LegacyEpConfigID uint
	LegacyEpConfig   *EndpointConfig `gorm:"foreignKey:LegacyEpConfigID"`
	Loglevel         uint

	UpdateGroup   *UpdateGroup `gorm:"foreignKey:UpdateGroupID"`
	UpdateGroupID uint
}

type EndpointConfig struct {
	CreatedAt            time.Time                `json:"-"`
	UpdatedAt            time.Time                `json:"-"`
	DeletedAt            gorm.DeletedAt           `json:"-" gorm:"index"`
	ID                   uint                     `gorm:"primarykey" json:"id"`
	Name                 string                   `json:"name"`
	MutualAuthentication bool                     `json:"mutual_auth"`
	NoEncryption         bool                     `json:"no_encrypt"`
	ASLKeyExchangeMethod asl.ASLKeyExchangeMethod `json:"kex"`
	Cipher               string
}

type Identity struct {
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	ID        uint           `gorm:"primarykey" json:"id"`
	Type      IdentityType   `json:"identity"`

	ServerEndpointAddr string ` json:"server_endpoint_addr"`
	ServerUrl          string `json:"server_url"`
	RevocationListUrl  string `json:"revocation_list_url"`

	UpdateIdentity   *UpdateIdentity `gorm:"foreignKey:UpdateIdentityID"`
	UpdateIdentityID uint
}

type UpdateEndpointConfig struct {
	ID                   uint
	UpdateEndpointConfig EndpointConfig `gorm:"embedded"`

	ActiveEndpointConfig   *EndpointConfig `gorm:"foreignKey:ActiveEndpointConfigID"`
	ActiveEndpointConfigID uint
}

type UpdateIdentity struct {
	ID             uint     `gorm:"primarykey" json:"id"`
	UpdateIdentity Identity `gorm:"embedded" `

	ActiveIdentity   *Identity `gorm:"foreignKey:ActiveIdentityID"`
	ActiveIdentityID uint
}

type UpdateGroup struct {
	ID            uint
	UpdateGroup   Group  `gorm:"embedded"`
	ActiveGroup   *Group `gorm:"foreignKey:ActiveGroupID"`
	ActiveGroupID uint
}
