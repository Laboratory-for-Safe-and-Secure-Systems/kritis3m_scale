package types

import (
	"time"

	"github.com/Laboratory-for-Safe-and-Secure-Systems/go-wolfssl/asl"
	"gorm.io/gorm"
)

type (
	DBNodes []DBNode
)

type ApplicationType uint8

const (
	ForwardProxy   ApplicationType = 0
	ReverseProxy   ApplicationType = 1
	TLS_TLSProxy   ApplicationType = 2 //server and client tls endpoint
	EchoServer     ApplicationType = 3
	L2Bridge       ApplicationType = 4
	NetworkTester  ApplicationType = 5 //server and client tls endpoint
	TcpStdinBridge ApplicationType = 6 //server and client tls endpoint

)

func (a ApplicationType) String() string {
	switch a {
	case ForwardProxy:
		return "Forward Proxy"
	case ReverseProxy:
		return "Reverse Proxy"
	case TLS_TLSProxy:
		return "TLS-TLS Proxy"
	case EchoServer:
		return "Echo Server"
	case L2Bridge:
		return "L2 Bridge"
	case NetworkTester:
		return "Network Tester"
	case TcpStdinBridge:
		return "TCP Stdin Bridge"
	default:
		return "Unknown"
	}
}

type ImportStructure struct {
	Node         []*DBNode              `json:"nodes"`
	CryptoConfig []*DBAslEndpointConfig `json:"crypto"`
	Identites    []*DBIdentity          `json:"pki_identities"`
}

type DistributionResponse struct {
	Node         DBNode                 `json:"node"`
	CryptoConfig []*DBAslEndpointConfig `json:"crypto_config"`
	Identities   []*DBIdentity          `json:"identities"`
}

type SelectedConfiguration struct {
	gorm.Model
	NodeID   uint
	Node     DBNode `gorm:"foreignKey:NodeID"`
	ConfigID uint
	Config   DBNodeConfig `gorm:"foreignKey:ConfigID"`
} // Node represents a node within a network

type DBNode struct {
	CreatedAt        time.Time       `json:"-"`
	UpdatedAt        time.Time       `json:"-"`
	DeletedAt        gorm.DeletedAt  `gorm:"index" json:"-"`
	ID               uint            `gorm:"primarykey" json:"id,omitempty"`
	SerialNumber     string          `gorm:"uniqueindex" json:"serial_number"`
	NodeNetworkIndex uint            `json:"network_index"`
	Locality         string          `json:"locality,omitempty"`
	LastSeen         time.Time       `json:"-"`
	Config           []*DBNodeConfig `gorm:"foreignKey:NodeID" json:"configs"`
}

// Node represents a node within a network
type DBNodeConfig struct {
	CreatedAt        time.Time        `json:"-"`
	DeletedAt        gorm.DeletedAt   `json:"-" gorm:"index"`
	ID               uint             `gorm:"primarykey" json:"id"`
	NodeID           uint             `json:"-"`
	UpdatedAt        time.Time        `json:"updated_at,omitempty"`
	Version          uint             `gorm:"default:0" json:"version,omitempty"`
	HardbeatInterval time.Duration    `json:"hb_interval"`
	Whitelist        DBWhitelist      `gorm:"foreignKey:NodeConfigID" json:"whitelist"`
	Application      []*DBApplication `gorm:"foreignKey:NodeConfigID" json:"applications"`
}

type DBApplication struct {
	CreatedAt      time.Time           `json:"-"`
	UpdatedAt      time.Time           `json:"-"`
	DeletedAt      gorm.DeletedAt      `json:"-" gorm:"index"`
	ID             uint                `gorm:"primarykey" json:"id"`
	NodeConfigID   uint                `json:"config_id,omitempty"`
	State          bool                `json:",omitempty"`
	TrustedClients []*DBTrustedClients `gorm:"many2many:application_trusts_clients" json:"-"`
	Type           ApplicationType     `json:"type"`

	ListeningIpPort string `json:"server_ip_port"`
	ClientIpPort    string `json:"client_ip_port"`

	Ep1ID uint                 `json:"ep1_id,omitempty"`
	Ep1   *DBAslEndpointConfig `json:"-" gorm:"foreignKey:Ep1ID"`

	Ep2ID uint                 `json:"ep2_id,omitempty"`
	Ep2   *DBAslEndpointConfig `json:"-" gorm:"foreignKey:Ep2ID"`
}

type DBWhitelist struct {
	CreatedAt      time.Time           `json:"-"`
	UpdatedAt      time.Time           `json:"-"`
	DeletedAt      gorm.DeletedAt      `json:"-" gorm:"index"`
	ID             uint                `gorm:"primarykey" json:"id"`
	NodeConfigID   uint                `json:"config_id,omitempty"`
	TrustedClients []*DBTrustedClients `gorm:"foreignKey:WhitelistID" json:"trusted_clients"`
}
type ApplicationTrustsClients struct {
	DBApplicationID    uint `gorm:"primaryKey"`
	DBTrustedClientsID uint `gorm:"primaryKey"`
}

type DBTrustedClients struct {
	CreatedAt               time.Time        `json:"-"`
	UpdatedAt               time.Time        `json:"-"`
	DeletedAt               gorm.DeletedAt   `json:"-" gorm:"index"`
	ID                      uint             `gorm:"primarykey:id" json:"id"`
	WhitelistID             uint             `json:"-"`
	ClientIpPort            string           `json:"client_ip_port"`
	ApplicationIDs          []uint           `gorm:"-" json:"application_ids" `
	ApplicationTrustsClient []*DBApplication `gorm:"many2many:application_trusts_clients;" json:"-"`
}

// ProxyApplication defines settings for a proxy application

// StandardApplication defines settings for a standard application
type DBAslEndpointConfig struct {
	CreatedAt            time.Time                `json:"-"`
	UpdatedAt            time.Time                `json:"-"`
	DeletedAt            gorm.DeletedAt           `json:"-" gorm:"index"`
	ID                   uint                     `gorm:"primarykey" json:"id"`
	Name                 string                   `json:"name"`
	MutualAuthentication bool                     `json:"mutual_auth"`
	NoEncryption         bool                     `json:"no_encrypt"`
	ASLKeyExchangeMethod asl.ASLKeyExchangeMethod `json:"kex"`
	UseSecureElement     bool                     `json:"use_secure_elem"`
	HybridSignatureMode  asl.HybridSignatureMode  `json:"signature_mode"`
	Keylog               bool                     `json:"keylog"`

	IdentityID uint        `json:"identity_id"`
	Identity   *DBIdentity `json:"-" gorm:"foreignKey:IdentityID"`
}

type DBIdentity struct {
	CreatedAt         time.Time      `json:"-"`
	UpdatedAt         time.Time      `json:"-"`
	DeletedAt         gorm.DeletedAt `json:"-" gorm:"index"`
	ID                uint           `gorm:"primarykey" json:"id"`
	Identity          uint           `json:"identity"`
	ServerAddr        string         `json:"server_addr"`
	ServerUrl         string         `json:"server_url"`
	RevocationListUrl string         `json:"revocation_list_url"`
}
