package types

import (
	"encoding/json"
	"fmt"
	"net"
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

// see linux/sys/socket.h PF_INET=2 &PF_INET6=10
type ProtoFamiliy uint8

const (
	AF_INET  ProtoFamiliy = 2
	AF_INET6 ProtoFamiliy = 10
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

type HeartbeatInstruction uint8

const (
	CallDistributionService HeartbeatInstruction = 0
	PostSystemStatus        HeartbeatInstruction = 1
	ChangeHeartbeatInterval HeartbeatInstruction = 2 //server and client tls endpoint
	ChangeLogLevel          HeartbeatInstruction = 3 //server and client tls endpoint
	UptoDate                HeartbeatInstruction = 4
)

type NodeState uint8

const (
	NotSeen  NodeState = 0
	Running  NodeState = 1
	Updating NodeState = 2 //server and client tls endpoint
)

func (a NodeState) String() string {
	switch a {
	case NotSeen:
		return "not seen"
	case Running:
		return "running"
	case Updating:
		return "updating"
	default:
		return "unknown state"
	}
}

func (a HeartbeatInstruction) String() string {
	switch a {
	case CallDistributionService:
		return "Call Distrib"
	case PostSystemStatus:
		return "post status"
	case ChangeHeartbeatInterval:
		return "change hb iv"
	case ChangeLogLevel:
		return "change log lvl"
	case UptoDate:
		return "no instruction"
	default:
		return "unknown instruction"
	}
}

type SelectedConfiguration struct {
	gorm.Model
	NodeID               uint
	Node                 DBNode `gorm:"foreignKey:NodeID"`
	ConfigID             uint
	Config               DBNodeConfig         `gorm:"foreignKey:ConfigID"`
	State                NodeState            `gorm:"default:0"` // Not seen
	HeartbeatInstruction HeartbeatInstruction `gorm:"default:0"` // cal distribution service
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

	ListeningIpPort Kritis3mAddr `gorm:"embedded;embeddedPrefix:listening_ip_" json:"server_ip_port"`
	ClientIpPort    Kritis3mAddr `gorm:"embedded;embeddedPrefix:client_ip_" json:"client_ip_port"`

	Ep1ID uint                 `json:"ep1_id,omitempty"`
	Ep1   *DBAslEndpointConfig `json:"-" gorm:"foreignKey:Ep1ID"`

	Ep2ID uint                 `json:"ep2_id,omitempty"`
	Ep2   *DBAslEndpointConfig `json:"-" gorm:"foreignKey:Ep2ID"`
}
type Kritis3mAddr struct {
	IP     net.IP       `json:"-" gorm:"type:varbinary(16)"` // To store up to 16 bytes (IPv6) // 0.0.0.0 for all ports
	IPStr  string       `json:"ip" gorm:"-" `
	Family ProtoFamiliy `json:"family"`
	Port   uint16       `json:"port"` // 0 for all ports
}

func (e Kritis3mAddr) MarshalJSON() ([]byte, error) {
	type Alias Kritis3mAddr
	aux := struct {
		IPStr string `json:"ip"`
		Alias
	}{
		IPStr: e.IP.String(),
		Alias: (Alias)(e),
	}
	return json.Marshal(aux)
}

// Custom JSON Unmarshaling
// json to struct
func (addr *Kritis3mAddr) UnmarshalJSON(data []byte) error {
	type RecursionBreaker *Kritis3mAddr
	var recBreaker RecursionBreaker
	recBreaker = (RecursionBreaker)(addr)

	if err := json.Unmarshal(data, recBreaker); err != nil {
		fmt.Print(err)
		return err
	}

	//ip and family are missing
	addr.IP = net.ParseIP(recBreaker.IPStr)
	if addr.IP == nil {
		return fmt.Errorf("can't parse ipstr to IP")
	}
	addr.IPStr = recBreaker.IPStr
	var family ProtoFamiliy

	if ip4 := recBreaker.IP.To4(); ip4 != nil {
		family = AF_INET
	} else if recBreaker.IP.To16() != nil {
		family = AF_INET6
	} else {
		return fmt.Errorf("invalid IP address: %v", recBreaker.IPStr)
	}
	addr.Family = family

	return nil
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
	ClientIpPort            Kritis3mAddr     `gorm:"embedded" json:"client_ip_port"`
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
	ServerAddr        Kritis3mAddr   `gorm:"embedded" json:"server_addr"`
	ServerUrl         string         `json:"server_url"`
	RevocationListUrl string         `json:"revocation_list_url"`
}
