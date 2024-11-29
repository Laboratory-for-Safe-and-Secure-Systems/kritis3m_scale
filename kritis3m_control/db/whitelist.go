package db

import (
	"github.com/Laboratory-for-Safe-and-Secure-Systems/kritis3m_scale/kritis3m_control/types"
	"gorm.io/gorm"
)

//-------------------- interfaces ------------------------

type WhitelistHandler interface {
	GetAllWhitelists() ([]*types.DBWhitelist, error)
	GetWhitelistby_Cfgid(cfg_id uint) ([]*types.DBWhitelist, error)
	GetWhitelistby_ID(id uint) (types.DBWhitelist, error)

	AddEmptyWhitelistto_CfgID(cfg_id uint) error
	AddWhitelistto_Cfg(cfg *types.DBNodeConfig, whitelist *types.DBWhitelist) error
	AddWhitelistto_CfgID(cfg_id uint, whitelist *types.DBWhitelist) error

	UpdateWhitelistby_ID(whitelist_id uint, whitelist *types.DBWhitelist) error
}

func (db *KSDatabase) GetAllWhitelists() ([]*types.DBWhitelist, error) {
	gormDB := db.DB
	return Read(gormDB, func(tx *gorm.DB) ([]*types.DBWhitelist, error) {
		return getAllWhitelists(tx)
	})
}

func getAllWhitelists(tx *gorm.DB) ([]*types.DBWhitelist, error) {
	var whitelists []*types.DBWhitelist
	if err := tx.Find(&whitelists).Error; err != nil {
		return nil, err
	}
	return whitelists, nil
}
func (db *KSDatabase) GetWhitelistby_Cfgid(cfg_id uint) ([]*types.DBWhitelist, error) {
	gormDB := db.DB
	return Read(gormDB, func(tx *gorm.DB) ([]*types.DBWhitelist, error) {
		return getWhitelistby_Cfgid(tx, cfg_id)
	})
}

func getWhitelistby_Cfgid(tx *gorm.DB, cfg_id uint) ([]*types.DBWhitelist, error) {
	var whitelists []*types.DBWhitelist
	if err := tx.Where("node_config_id = ?", cfg_id).Find(&whitelists).Error; err != nil {
		return nil, err
	}
	return whitelists, nil
}
func (db *KSDatabase) GetWhitelistby_ID(id uint) (types.DBWhitelist, error) {
	gormDB := db.DB
	return Read(gormDB, func(tx *gorm.DB) (types.DBWhitelist, error) {
		return getWhitelistby_ID(tx, id)
	})
}

func getWhitelistby_ID(tx *gorm.DB, id uint) (types.DBWhitelist, error) {
	var whitelist types.DBWhitelist
	if err := tx.Where("id = ?", id).First(&whitelist).Error; err != nil {
		return types.DBWhitelist{}, err
	}
	return whitelist, nil
}
func (db *KSDatabase) AddEmptyWhitelistto_CfgID(cfg_id uint) error {
	return db.Write(func(tx *gorm.DB) error {
		return addEmptyWhitelistto_CfgID(tx, cfg_id)
	})
}

func addEmptyWhitelistto_CfgID(tx *gorm.DB, cfg_id uint) error {
	whitelist := &types.DBWhitelist{
		NodeConfigID: cfg_id,
	}
	return tx.Create(whitelist).Error
}
func (db *KSDatabase) AddWhitelistto_Cfg(cfg *types.DBNodeConfig, whitelist *types.DBWhitelist) error {
	return db.Write(func(tx *gorm.DB) error {
		return addWhitelistto_Cfg(tx, cfg, whitelist)
	})
}

func addWhitelistto_Cfg(tx *gorm.DB, cfg *types.DBNodeConfig, whitelist *types.DBWhitelist) error {
	whitelist.NodeConfigID = cfg.ID
	return tx.Create(whitelist).Error
}
func (db *KSDatabase) AddWhitelistto_CfgID(cfg_id uint, whitelist *types.DBWhitelist) error {
	return db.Write(func(tx *gorm.DB) error {
		return addWhitelistto_CfgID(tx, cfg_id, whitelist)
	})
}

func addWhitelistto_CfgID(tx *gorm.DB, cfg_id uint, whitelist *types.DBWhitelist) error {
	whitelist.NodeConfigID = cfg_id
	return tx.Create(whitelist).Error
}
func (db *KSDatabase) UpdateWhitelistby_ID(whitelist_id uint, whitelist *types.DBWhitelist) error {
	return db.Write(func(tx *gorm.DB) error {
		return updateWhitelistby_ID(tx, whitelist_id, whitelist)
	})
}

func updateWhitelistby_ID(tx *gorm.DB, whitelist_id uint, whitelist *types.DBWhitelist) error {
	var existingWhitelist types.DBWhitelist
	if err := tx.Where("id = ?", whitelist_id).First(&existingWhitelist).Error; err != nil {
		return err
	}
	return tx.Model(&existingWhitelist).Updates(whitelist).Error
}
