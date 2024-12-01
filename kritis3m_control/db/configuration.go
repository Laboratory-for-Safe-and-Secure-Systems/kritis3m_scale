package db

import (
	"errors"
	"fmt"

	"github.com/Laboratory-for-Safe-and-Secure-Systems/kritis3m_scale/kritis3m_control/types"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

//-------------------- interfaces ------------------------

type ConfigurationHandler interface {
	GetAllConfigs() ([]*types.DBNodeConfig, error)
	GetAllConfigsOfNodeby_ID(id uint) ([]*types.DBNodeConfig, error)
	GetAllConfigsOfNodeby_SerialNumber(serial_number string) ([]*types.DBNodeConfig, error)
	GetActiveConfigOfNodeby_ID(id uint) (*types.DBNodeConfig, error)
	GetConfigby_ID(cfg_id uint) (*types.DBNodeConfig, error)

	GetHwConfigby_ID(cfg_id uint) (*types.HardwareConfig, error)
	GetHwConfigby_ConfigID(cfg_id uint) ([]*types.HardwareConfig, error)

	AddConfigto_NodeSerialNumber(noide_serialnumber string, cfg *types.DBNodeConfig) error
	AddConfigto_NodeID(noide_id uint, cfg *types.DBNodeConfig) error
	AddConfigto_Node(node *types.DBNode, cfg *types.DBNodeConfig) error
	AddHwConfigto_Config(config_id uint, hw_cfg *types.HardwareConfig) error

	UpdateConfig_byID(hw_id uint, cfg *types.DBNodeConfig) error
	UpdateHwConfig(config_id uint, hw_cfg *types.HardwareConfig) error

	ActiveConfigSetState_byCfgID(cfg_id uint, new_state types.NodeState) error
	ActivateConfig_byCfgID(cfg_id uint, node_id uint, new_state types.NodeState) error
}

func (db *KSDatabase) ActiveConfigSetState_byCfgID(cfg_id uint, new_state types.NodeState) error {
	return db.Write(func(tx *gorm.DB) error {
		return activeconfigsetstate_bycfgid(tx, cfg_id, new_state)
	})

}

func (db *KSDatabase) ActivateConfig_byCfgID(cfg_id uint, serial_number string) error {
	node, err := db.GetNodeby_SerialNumber(serial_number)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error().Msgf("Node with serial number %s, not known", serial_number)
		}
		return err
	}
	return db.Write(func(tx *gorm.DB) error {
		return activateconfig_bycfgid(tx, cfg_id, node.ID)
	})
}

// activateconfig_bycfgid manages the activation of a configuration for a specific node
// It either updates an existing entry or creates a new one if no entry exists
func activateconfig_bycfgid(tx *gorm.DB, cfg_id uint, node_id uint) error {
	// First, check if an entry for this node already exists
	var existingConfig types.SelectedConfiguration
	result := tx.Where("node_id = ?", node_id).First(&existingConfig)

	if result.Error != nil {
		// No existing configuration found - create a new entry
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {

			newConfig := types.SelectedConfiguration{
				NodeID:    node_id,
				ConfigID:  cfg_id,
				NodeState: types.NotSeen,
			}
			if err := tx.Create(&newConfig).Error; err != nil {
				return err
			}
			return nil
		}
		// Other database error
		return result.Error
	}

	// If existing configuration has a different ConfigID, update it
	if existingConfig.ConfigID != cfg_id {
		log.Info().Msgf("Node with node id %d, hat cfg with cfg id %d. New cfg id is: %d", node_id, existingConfig.ConfigID, cfg_id)
		existingConfig.ConfigID = cfg_id
		existingConfig.NodeState = types.NotSeen
		if err := tx.Save(&existingConfig).Error; err != nil {
			return err
		}
	}
	return nil
}

// activeconfigsetstate_bycfgid updates the state of a specific configuration
// It throws an error if the configuration is not found
func activeconfigsetstate_bycfgid(tx *gorm.DB, config_id uint, new_state types.NodeState) error {
	// Find the configuration first
	var selectedConfig types.SelectedConfiguration
	result := tx.Where("config_id = ?", config_id).First(&selectedConfig)

	// Check if configuration exists
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf("configuration with ID %d not found in selected configurations", config_id)
		}
		return fmt.Errorf("database error while finding configuration: %w", result.Error)
	}

	// Update the state
	selectedConfig.NodeState = new_state
	if err := tx.Save(&selectedConfig).Error; err != nil {
		return fmt.Errorf("failed to update configuration state: %w", err)
	}

	return nil
}

func (db *KSDatabase) AddHwConfigto_Config(config_id uint, hw_cfg *types.HardwareConfig) error {
	return db.Write(func(tx *gorm.DB) error {
		return addhwconfigto_config(tx, config_id, hw_cfg)
	})
}
func addhwconfigto_config(tx *gorm.DB, config_id uint, hw_cfg *types.HardwareConfig) error {
	hw_cfg.ConfigID = config_id
	return tx.Create(&hw_cfg).Error
}

func (db *KSDatabase) UpdateHwConfig(hw_cfg_id uint, hw_cfg *types.HardwareConfig) error {
	return db.Write(func(tx *gorm.DB) error {
		return updatehwconfig(tx, hw_cfg_id, hw_cfg)
	})
}
func updatehwconfig(tx *gorm.DB, hw_cfg_id uint, hw_cfg *types.HardwareConfig) error {
	hw_cfg.ID = hw_cfg_id
	return tx.Save(&hw_cfg).Error

}

func (db *KSDatabase) GetHwConfigby_ID(cfg_id uint) (*types.HardwareConfig, error) {
	gormDB := db.DB
	return Read(gormDB, func(tx *gorm.DB) (*types.HardwareConfig, error) {
		return gethwconfigby_id(tx, cfg_id)
	})
}
func gethwconfigby_id(rx *gorm.DB, cfg_id uint) (*types.HardwareConfig, error) {
	var hw_config *types.HardwareConfig
	if err := rx.Where("id = ", cfg_id).Find(hw_config).Error; err != nil {
		return nil, err
	}
	return hw_config, nil
}

func (db *KSDatabase) GetHwConfigby_ConfigID(cfg_id uint) ([]*types.HardwareConfig, error) {
	gormDB := db.DB
	return Read(gormDB, func(tx *gorm.DB) ([]*types.HardwareConfig, error) {
		return gethwconfigby_configid(gormDB, cfg_id)
	})
}
func gethwconfigby_configid(rx *gorm.DB, cfg_id uint) ([]*types.HardwareConfig, error) {
	var hw_configs []*types.HardwareConfig
	if err := rx.Where("config_id = ?", cfg_id).Find(hw_configs).Error; err != nil {
		return nil, err
	}
	return hw_configs, nil
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

func (db *KSDatabase) AddConfigto_NodeValues(node *types.DBNode, version uint) error {
	return db.Write(func(tx *gorm.DB) error {
		return addConfigto_NodeValues(tx, node, version)
	})
}

func addConfigto_NodeValues(tx *gorm.DB, node *types.DBNode, version uint) error {
	config := &types.DBNodeConfig{
		NodeID:  node.ID,
		Version: version,
	}
	return tx.Create(config).Error
}
