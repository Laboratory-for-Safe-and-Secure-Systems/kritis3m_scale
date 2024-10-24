package db

import (
	"time"

	"github.com/philslol/kritis3m_scale/kritis3m_control/types"
	"gorm.io/gorm"
)

//-------------------- interfaces ------------------------

type ConfigurationHandler interface {
	GetAllConfigs() ([]*types.DBNodeConfig, error)
	GetAllConfigsOfNodeby_ID(id uint) ([]*types.DBNodeConfig, error)
	GetAllConfigsOfNodeby_SerialNumber(serial_number string) ([]*types.DBNodeConfig, error)
	GetActiveConfigOfNodeby_ID(id uint) (*types.DBNodeConfig, error)
	GetConfigby_ID(cfg_id uint) (*types.DBNodeConfig, error)
	AddConfigto_NodeSerialNumber(noide_serialnumber string, cfg *types.DBNodeConfig) error
	AddConfigto_NodeID(noide_id uint, cfg *types.DBNodeConfig) error
	AddConfigto_Node(node *types.DBNode, cfg *types.DBNodeConfig) error
	AddConfigto_NodeValues(node *types.DBNode, hb_inteval time.Duration, version uint) error
	UpdateConfig_byID(id uint, cfg *types.DBNodeConfig) error
}

func (db *KSDatabase) UpdateConfig_byID(id uint, cfg *types.DBNodeConfig) error {
	return db.Write(func(tx *gorm.DB) error {
		return updateConfig_byID(tx, id, cfg)
	})
}

func updateConfig_byID(tx *gorm.DB, id uint, cfg *types.DBNodeConfig) error {
	// Find the existing config by ID
	var existingConfig types.DBNodeConfig
	if err := tx.Where("id = ?", id).First(&existingConfig).Error; err != nil {
		return err
	}
	// Update the existing config with the new values from cfg
	return tx.Model(&existingConfig).Updates(cfg).Error
}

func (db *KSDatabase) GetAllConfigs() ([]*types.DBNodeConfig, error) {
	gormDB := db.DB
	return Read(gormDB, func(tx *gorm.DB) ([]*types.DBNodeConfig, error) {
		return getAllConfigs(tx)
	})
}

func getAllConfigs(tx *gorm.DB) ([]*types.DBNodeConfig, error) {
	var configs []*types.DBNodeConfig
	if err := tx.Find(&configs).Error; err != nil {
		return nil, err
	}
	return configs, nil
}
func (db *KSDatabase) GetAllConfigsOfNodeby_ID(id uint) ([]*types.DBNodeConfig, error) {
	gormDB := db.DB
	return Read(gormDB, func(tx *gorm.DB) ([]*types.DBNodeConfig, error) {
		return getAllConfigsOfNodeby_ID(tx, id)
	})
}

func getAllConfigsOfNodeby_ID(tx *gorm.DB, id uint) ([]*types.DBNodeConfig, error) {
	var configs []*types.DBNodeConfig
	if err := tx.Where("node_id = ?", id).Find(&configs).Error; err != nil {
		return nil, err
	}
	return configs, nil
}

func (db *KSDatabase) GetAllConfigsOfNodeby_SerialNumber(serial_number string) ([]*types.DBNodeConfig, error) {
	gormDB := db.DB
	return Read(gormDB, func(tx *gorm.DB) ([]*types.DBNodeConfig, error) {
		return getAllConfigsOfNodeby_SerialNumber(tx, serial_number)
	})
}

func getAllConfigsOfNodeby_SerialNumber(tx *gorm.DB, serial_number string) ([]*types.DBNodeConfig, error) {
	var node types.DBNode
	if err := tx.Where("serial_number = ?", serial_number).First(&node).Error; err != nil {
		return nil, err
	}

	var configs []*types.DBNodeConfig
	if err := tx.Where("node_id = ?", node.ID).Find(&configs).Error; err != nil {
		return nil, err
	}
	return configs, nil
}

func (db *KSDatabase) GetActiveConfigs() ([]*types.SelectedConfiguration, error) {
	gormDB := db.DB
	return Read(gormDB, func(tx *gorm.DB) ([]*types.SelectedConfiguration, error) {
		return getActiveConfigs(tx)
	})
}

func getActiveConfigs(tx *gorm.DB) ([]*types.SelectedConfiguration, error) {
	var configs []*types.SelectedConfiguration

	// Query to get active configurations
	if err := tx.Find(&configs).Error; err != nil {
		return nil, err
	}

	return configs, nil
}

func (db *KSDatabase) GetActiveConfigOfNodeby_SerialNumber(serialnumber string) (*types.DBNodeConfig, error) {
	gormDB := db.DB
	return Read(gormDB, func(tx *gorm.DB) (*types.DBNodeConfig, error) {
		return getActiveConfigOfNodeby_SerialNumber(tx, serialnumber)
	})
}

func getActiveConfigOfNodeby_SerialNumber(tx *gorm.DB, serialnumber string) (*types.DBNodeConfig, error) {
	//
	var cfg *types.SelectedConfiguration
	if err := tx.Where("serialnumber=?", serialnumber).First(&cfg).Error; err != nil {
		return nil, err
	}
	return &cfg.Config, nil
}

func (db *KSDatabase) GetActiveConfigOfNodeby_ID(id uint) (*types.DBNodeConfig, error) {
	gormDB := db.DB
	return Read(gormDB, func(tx *gorm.DB) (*types.DBNodeConfig, error) {
		return getActiveConfigOfNodeby_ID(tx, id)
	})
}

func getActiveConfigOfNodeby_ID(tx *gorm.DB, id uint) (*types.DBNodeConfig, error) {
	//
	var cfg *types.SelectedConfiguration
	if err := tx.Where("noide_id=?", id).First(&cfg).Error; err != nil {
		return nil, err
	}
	return &cfg.Config, nil
}

func (db *KSDatabase) GetConfigby_ID(cfg_id uint) (*types.DBNodeConfig, error) {
	gormDB := db.DB
	return Read(gormDB, func(tx *gorm.DB) (*types.DBNodeConfig, error) {
		return getConfigby_ID(tx, cfg_id)
	})
}

func getConfigby_ID(tx *gorm.DB, cfg_id uint) (*types.DBNodeConfig, error) {
	var config types.DBNodeConfig
	if err := tx.Where("id = ?", cfg_id).First(&config).Error; err != nil {
		return nil, err
	}
	return &config, nil
}

func (db *KSDatabase) AddConfigto_NodeSerialNumber(node_serialnumber string, cfg *types.DBNodeConfig) error {
	return db.Write(func(tx *gorm.DB) error {
		return addConfigto_NodeSerialNumber(tx, node_serialnumber, cfg)
	})
}

func addConfigto_NodeSerialNumber(tx *gorm.DB, serial_number string, cfg *types.DBNodeConfig) error {
	var node types.DBNode
	if err := tx.Where("serial_number = ?", serial_number).First(&node).Error; err != nil {
		return err
	}
	cfg.NodeID = node.ID
	return tx.Create(cfg).Error
}
func (db *KSDatabase) AddConfigto_NodeID(node_id uint, cfg *types.DBNodeConfig) error {
	return db.Write(func(tx *gorm.DB) error {
		return addConfigto_NodeID(tx, node_id, cfg)
	})
}

func addConfigto_NodeID(tx *gorm.DB, node_id uint, cfg *types.DBNodeConfig) error {
	cfg.NodeID = node_id
	return tx.Create(cfg).Error
}
func (db *KSDatabase) AddConfigto_Node(node *types.DBNode, cfg *types.DBNodeConfig) error {
	return db.Write(func(tx *gorm.DB) error {
		return addConfigto_Node(tx, node, cfg)
	})
}

func addConfigto_Node(tx *gorm.DB, node *types.DBNode, cfg *types.DBNodeConfig) error {
	cfg.NodeID = node.ID
	return tx.Create(cfg).Error
}

func (db *KSDatabase) AddConfigto_NodeValues(node *types.DBNode, hb_interval time.Duration, version uint) error {
	return db.Write(func(tx *gorm.DB) error {
		return addConfigto_NodeValues(tx, node, hb_interval, version)
	})
}

func addConfigto_NodeValues(tx *gorm.DB, node *types.DBNode, hb_interval time.Duration, version uint) error {
	config := &types.DBNodeConfig{
		NodeID:           node.ID,
		HardbeatInterval: hb_interval,
		Version:          version,
	}
	return tx.Create(config).Error
}
