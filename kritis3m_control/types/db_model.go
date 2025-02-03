package types

import (
	"time"

	asl "github.com/Laboratory-for-Safe-and-Secure-Systems/go-asl"
	v1 "github.com/Laboratory-for-Safe-and-Secure-Systems/kritis3m_scale/gen/go/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
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

type Identity uint8

const (
	MANAGEMENT_SERVICE uint = iota
	MANAGEMENT
	REMOTE
	PRODUCTION
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
func convertToUint32Slice(input []uint) []uint32 {
	output := make([]uint32, len(input))
	for i, v := range input {
		output[i] = uint32(v)
	}
	return output
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

type NodeState int8

const (
	ErrorState          NodeState = 0
	NotSeen             NodeState = 1
	NodeRequestedConfig NodeState = 2
	Running             NodeState = 3
)

func (a NodeState) String() string {
	switch a {
	case ErrorState:
		return "Error"
	case NotSeen:
		return "not seen"
	case Running:
		return "running"
	case NodeRequestedConfig:
		return "node requested configuration"
	default:
		return "unknown state"
	}
}

type SelectedConfiguration struct {
	gorm.Model
	NodeID    uint
	Node      DBNode `gorm:"foreignKey:NodeID"`
	ConfigID  uint
	Config    DBNodeConfig `gorm:"foreignKey:ConfigID"`
	NodeState NodeState    `gorm:"default:0"` // cal distribution service
} // Node represents a node within a network

type HardwareConfig struct {
	ID       uint         `gorm:"primarykey" json:"-"`
	ConfigID uint         `json:"-"`
	Config   DBNodeConfig `gorm:"foreignKey:ConfigID" json:"-"`

	Device string `json:"device"`
	IpCidr string `json:"cidr"`
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
	CreatedAt  time.Time      `json:"-"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
	ID         uint           `gorm:"primarykey" json:"id"`
	NodeID     uint           `json:"-"`
	LogLevel   uint           `gorm:"default:3" json:"log_level,omitempty"`
	UpdatedAt  time.Time      `json:"updated_at,omitempty"`
	ConfigName string         `json:"config_name"`
	Version    uint           `gorm:"default:0" json:"version,omitempty"`

	HardwareConfig []*HardwareConfig `gorm:"foreignKey:ConfigID" json:"hw_config"`

	Whitelist DBWhitelist `gorm:"foreignKey:NodeConfigID" json:"whitelist"`

	Application []*DBApplication `gorm:"foreignKey:NodeConfigID" json:"applications"`
}

type DBApplication struct {
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	ID        uint           `gorm:"primarykey" json:"id"`

	NodeConfigID   uint                `json:"config_id,omitempty"`
	State          bool                `json:",omitempty"`
	TrustedClients []*DBTrustedClients `gorm:"many2many:application_trusts_clients" json:"-"`
	Type           ApplicationType     `json:"type"`

	ServerEndpointAddr string `json:"server_endpoint_addr"`
	ClientEndpointAddr string `json:"client_endpoint_addr"`

	Ep1ID uint                 `json:"ep1_id,omitempty"`
	Ep1   *DBAslEndpointConfig `json:"-" gorm:"foreignKey:Ep1ID"`

	Ep2ID    uint                 `json:"ep2_id,omitempty"`
	Ep2      *DBAslEndpointConfig `json:"-" gorm:"foreignKey:Ep2ID"`
	LogLevel uint                 `json:"log_level,omitempty" gorm:"default:3"`
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
	ClientEndpointAddr      string           `json:"client_endpoint_addr"`
	ApplicationIDs          []uint           `gorm:"-" json:"application_ids" `
	ApplicationTrustsClient []*DBApplication `gorm:"many2many:application_trusts_clients;" json:"-"`
}

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
	// HybridSignatureMode  asl.HybridSignatureMode  `json:"signature_mode"`
	Keylog bool `json:"keylog"`

	IdentityID uint        `json:"identity_id"`
	Identity   *DBIdentity `json:"-" gorm:"foreignKey:IdentityID"`
}

type DBIdentity struct {
	CreatedAt          time.Time      `json:"-"`
	UpdatedAt          time.Time      `json:"-"`
	DeletedAt          gorm.DeletedAt `json:"-" gorm:"index"`
	ID                 uint           `gorm:"primarykey" json:"id"`
	Identity           Identity       `json:"identity"`
	ServerEndpointAddr string         ` json:"server_endpoint_addr"`
	ServerUrl          string         `json:"server_url"`
	RevocationListUrl  string         `json:"revocation_list_url"`
}

func (config DBAslEndpointConfig) Proto() *v1.DBAslEndpointConfig {
	return &v1.DBAslEndpointConfig{
		Id:                   uint32(config.ID),
		Name:                 config.Name,
		MutualAuthentication: config.MutualAuthentication,
		NoEncryption:         config.NoEncryption,
		AslKeyExchangeMethod: uint32(config.ASLKeyExchangeMethod),
		UseSecureElement:     config.UseSecureElement,
		Keylog:               config.Keylog,
		IdentityId:           uint32(config.IdentityID),
	}
}

func (identity DBIdentity) Proto() *v1.DBIdentity {
	return &v1.DBIdentity{
		Id:                 uint32(identity.ID),
		Identity:           v1.Identity(identity.Identity),
		ServerEndpointAddr: identity.ServerEndpointAddr,
		ServerUrl:          identity.ServerUrl,
		RevocationListUrl:  identity.RevocationListUrl,
	}
}
func (whitelist DBWhitelist) Proto() *v1.DBWhitelist {
	var trustedClients []*v1.DBTrustedClients
	for _, client := range whitelist.TrustedClients {
		trustedClients = append(trustedClients, client.Proto())
	}
	return &v1.DBWhitelist{
		Id:             uint32(whitelist.ID),
		TrustedClients: trustedClients,
		ConfigId:       uint32(whitelist.NodeConfigID),
	}
}
func (client DBTrustedClients) Proto() *v1.DBTrustedClients {
	return &v1.DBTrustedClients{
		Id:                 uint32(client.ID),
		ClientEndpointAddr: client.ClientEndpointAddr,
		ApplicationIds:     convertToUint32Slice(client.ApplicationIDs),
	}
}
func (sel HardwareConfig) Proto() *v1.HardwareConfig {
	hw_config := v1.HardwareConfig{
		Id:       uint32(sel.ID),
		ConfigId: uint32(sel.ConfigID),
		Device:   sel.Device,
		Cidr:     sel.IpCidr,
	}
	return &hw_config
}

func (config DBNodeConfig) Proto() *v1.DBNodeConfig {
	var hwConfigs []*v1.HardwareConfig
	for _, hw := range config.HardwareConfig {
		hwConfigs = append(hwConfigs, hw.Proto()) // Assuming HardwareConfig has a Proto() method
	}

	var applications []*v1.DBApplication
	for _, app := range config.Application {
		applications = append(applications, app.Proto()) // Assuming DBApplication has a Proto() method
	}

	return &v1.DBNodeConfig{
		Id:           uint32(config.ID),
		LogLevel:     uint32(config.LogLevel),
		ConfigName:   config.ConfigName,
		Version:      uint32(config.Version),
		HwConfig:     hwConfigs,
		Whitelist:    config.Whitelist.Proto(), // Assuming DBWhitelist has a Proto() method
		Applications: applications,
	}
}
func (app DBApplication) Proto() *v1.DBApplication {
	return &v1.DBApplication{
		Id:                 uint32(app.ID),
		ConfigId:           uint32(app.NodeConfigID),
		State:              app.State,
		Type:               v1.ApplicationType(app.Type),
		ServerEndpointAddr: app.ServerEndpointAddr,
		ClientEndpointAddr: app.ClientEndpointAddr,
		Ep1Id:              uint32(app.Ep1ID),
		Ep2Id:              uint32(app.Ep2ID),
		LogLevel:           uint32(app.LogLevel),
	}
}
func (node DBNode) Proto() v1.DBNode {
	var configs []*v1.DBNodeConfig
	for _, config := range node.Config {
		configs = append(configs, config.Proto())
	}
	return v1.DBNode{
		Id:               uint32(node.ID),
		SerialNumber:     node.SerialNumber,
		NodeNetworkIndex: uint32(node.NodeNetworkIndex),
		Locality:         node.Locality,
		LastSeen:         timestamppb.New(node.LastSeen),
		Configs:          configs,
	}
}
func (sel SelectedConfiguration) Proto() v1.SelectedConfiguration {
	node := sel.Node.Proto()
	return v1.SelectedConfiguration{
		NodeId:    uint64(sel.NodeID),
		Node:      &node,
		ConfigId:  uint64(sel.ConfigID),
		Config:    sel.Config.Proto(),
		NodeState: v1.NodeState(sel.NodeState),
	}
}
