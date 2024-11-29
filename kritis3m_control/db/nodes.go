package db

import (
	"errors"

	"time"

	"github.com/Laboratory-for-Safe-and-Secure-Systems/kritis3m_scale/kritis3m_control/types"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

//-------------------- global Errors ------------------------

var (
	ErrNodeNotFound                  = errors.New("node not found")
	ErrNodeRouteIsNotAvailable       = errors.New("route is not available on node")
	ErrNodeNotFoundRegistrationCache = errors.New(
		"node not found in registration cache",
	)
	ErrCouldNotConvertNodeInterface = errors.New("failed to convert node interface")
	ErrDifferentRegisteredUser      = errors.New(
		"node was previously registered with a different user",
	)
)

//-------------------- interfaces ------------------------

type NodeHandler interface {
	GetAllNodes() ([]*types.DBNode, error)
	GetNodeby_SerialNumber(serial_number string) (*types.DBNode, error)
	GetNodeby_ID(id int) (*types.DBNode, error)

	UpdateLastSeenby_ID(serial_number string, time_stamp time.Time) error
	UpdateLastSeenby_SerialNumber(id uint, time_stamp time.Time) error

	AddNode(*types.DBNode) error
}

func (db *KSDatabase) ActivateConfig(node_id uint, cfg_id uint) (*types.SelectedConfiguration, error) {
	return Write(db.DB, func(tx *gorm.DB) (*types.SelectedConfiguration, error) {
		return activateConfig(tx, node_id, cfg_id)
	})
}

func activateConfig(tx *gorm.DB, node_id uint, cfg_id uint) (*types.SelectedConfiguration, error) {
	var node types.DBNode
	var selcfg types.SelectedConfiguration
	var cfg types.DBNodeConfig

	if err := tx.Model(&types.DBNodeConfig{}).Where("id = ?", cfg_id).First(&cfg).Error; err != nil {
		log.Err(err).Msgf("no config with config id %d not exists", cfg_id)
		return nil, err
	}
	if cfg.NodeID != node_id {
		return nil, errors.New("cfg does not correspont to node")
	}

	// Fetch the node using node_id
	if err := tx.Model(&types.DBNode{}).Where("id = ?", node_id).First(&node).Error; err != nil {
		return nil, err
	}

	// Prepare the selected configuration with the new config_id
	selcfg.NodeID = node.ID
	selcfg.ConfigID = cfg.ID

	// Save the selected configuration (this will update if it exists, or insert a new record if not)
	if err := tx.Save(&selcfg).Error; err != nil {
		return nil, err
	}

	return &selcfg, nil
}

func (db *KSDatabase) GetAllNodes() ([]*types.DBNode, error) {
	gormDB := db.DB
	return Read(gormDB, func(db *gorm.DB) ([]*types.DBNode, error) { return getAllNodes(db) })

}

func (db *KSDatabase) GetNodeby_SerialNumber(serial_number string) (*types.DBNode, error) {
	gormDB := db.DB
	return Read(gormDB, func(db *gorm.DB) (*types.DBNode, error) { return getNodeby_SerialNumber(db, serial_number) })
}

func getNodeby_SerialNumber(tx *gorm.DB, serial_number string) (*types.DBNode, error) {
	var node *types.DBNode
	if err := tx.Model(&types.DBNode{}).Where("serial_number= ?", serial_number).First(&node).Error; err != nil {
		return nil, err
	}
	return node, nil
}
func (db *KSDatabase) GetNodeby_ID(id uint) (*types.DBNode, error) {
	gormDB := db.DB
	return Read(gormDB, func(db *gorm.DB) (*types.DBNode, error) {
		return getNodeby_ID(db, id)
	})
}
func getNodeby_ID(tx *gorm.DB, id uint) (*types.DBNode, error) {
	var node types.DBNode
	if err := tx.Model(&types.DBNode{}).Where("id = ?", id).First(&node).Error; err != nil {
		return nil, err
	}
	return &node, nil
}

func (db *KSDatabase) UpdateLastSeenby_ID(id uint, time_stamp time.Time) error {
	return db.Write(func(tx *gorm.DB) error {
		return updateLastSeenby_ID(tx, id, time_stamp)
	})
}
func updateLastSeenby_ID(tx *gorm.DB, id uint, time_stamp time.Time) error {
	return tx.Model(&types.DBNode{}).Where("id = ?", id).Update("last_seen", time_stamp).Error
}
func (db *KSDatabase) UpdateLastSeenby_SerialNumber(serial_number string, time_stamp time.Time) error {
	return db.Write(func(tx *gorm.DB) error {
		return updateLastSeenby_SerialNumber(tx, serial_number, time_stamp)
	})
}
func updateLastSeenby_SerialNumber(tx *gorm.DB, serial_number string, time_stamp time.Time) error {
	return tx.Model(&types.DBNode{}).Where("serial_number = ?", serial_number).Update("last_seen", time_stamp).Error
}
func (db *KSDatabase) AddNode(node *types.DBNode) error {
	return db.Write(func(tx *gorm.DB) error {
		return addNode(tx, node)
	})
}
func addNode(tx *gorm.DB, node *types.DBNode) error {
	return tx.Create(node).Error
}

func getAllNodes(tx *gorm.DB) ([]*types.DBNode, error) {
	var nodes []*types.DBNode
	if err := tx.
		Find(&nodes).Error; err != nil {

		return nil, err
	}
	return nodes, nil
}
