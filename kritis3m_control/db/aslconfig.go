package db

import (
	"github.com/philslol/kritis3m_scale/kritis3m_control/types"
	"gorm.io/gorm"
)

type ASLEndpointHandler interface {
	GetAllEPs() ([]*types.DBAslEndpointConfig, error)
	GetEPby_ID(id uint) (*types.DBAslEndpointConfig, error)

	AddEP(ep *types.DBAslEndpointConfig) (*types.DBAslEndpointConfig, error)
	AddEPs(ep []*types.DBAslEndpointConfig) ([]*types.DBAslEndpointConfig, error)

	UpdateEpby_Ep(ep *types.DBAslEndpointConfig, update_ep *types.DBAslEndpointConfig) (*types.DBAslEndpointConfig, error)
	UpdateEpby_EpID(ep *types.DBAslEndpointConfig, update_ep_id uint) (*types.DBAslEndpointConfig, error)
}

func (db *KSDatabase) GetAllEPs() ([]*types.DBAslEndpointConfig, error) {
	gormDB := db.DB
	return Read(gormDB, getAllEPs)
}
func getAllEPs(tx *gorm.DB) ([]*types.DBAslEndpointConfig, error) {
	var eps []*types.DBAslEndpointConfig
	if err := tx.Find(&eps).Error; err != nil {
		return nil, err
	}
	return eps, nil
}
func (db *KSDatabase) GetEPby_ID(id uint) (*types.DBAslEndpointConfig, error) {
	gormDB := db.DB
	return Read(gormDB, func(tx *gorm.DB) (*types.DBAslEndpointConfig, error) {
		return getEPby_ID(tx, id)
	})
}
func getEPby_ID(tx *gorm.DB, id uint) (*types.DBAslEndpointConfig, error) {
	var ep types.DBAslEndpointConfig
	if err := tx.First(&ep, id).Error; err != nil {
		return nil, err
	}
	return &ep, nil
}
func (db *KSDatabase) AddEP(ep *types.DBAslEndpointConfig) (*types.DBAslEndpointConfig, error) {
	gormDB := db.DB
	return Write(gormDB, func(tx *gorm.DB) (*types.DBAslEndpointConfig, error) {
		return addEP(tx, ep)
	})
}
func addEP(tx *gorm.DB, ep *types.DBAslEndpointConfig) (*types.DBAslEndpointConfig, error) {
	if err := tx.Create(ep).Error; err != nil {
		return nil, err
	}
	return ep, nil
}
func (db *KSDatabase) AddEPs(eps []*types.DBAslEndpointConfig) ([]*types.DBAslEndpointConfig, error) {
	gormDB := db.DB
	return Write(gormDB, func(tx *gorm.DB) ([]*types.DBAslEndpointConfig, error) {
		return addEPs(tx, eps)
	})
}
func addEPs(tx *gorm.DB, eps []*types.DBAslEndpointConfig) ([]*types.DBAslEndpointConfig, error) {
	if err := tx.Create(&eps).Error; err != nil {
		return nil, err
	}
	return eps, nil
}

func (db *KSDatabase) AddIdentitys(identities []*types.DBIdentity) ([]*types.DBIdentity, error) {
	gormDB := db.DB
	return Write(gormDB, func(tx *gorm.DB) ([]*types.DBIdentity, error) {
		return addIdentities(tx, identities)
	})
}

func addIdentities(tx *gorm.DB, identities []*types.DBIdentity) ([]*types.DBIdentity, error) {
	if err := tx.Create(&identities).Error; err != nil {
		return nil, err
	}
	return identities, nil
}
func (db *KSDatabase) UpdateEpby_Ep(ep *types.DBAslEndpointConfig, update_ep *types.DBAslEndpointConfig) (*types.DBAslEndpointConfig, error) {
	gormDB := db.DB
	return Write(gormDB, func(tx *gorm.DB) (*types.DBAslEndpointConfig, error) {
		return updateEpby_Ep(tx, ep, update_ep)
	})
}
func updateEpby_Ep(tx *gorm.DB, ep *types.DBAslEndpointConfig, update_ep *types.DBAslEndpointConfig) (*types.DBAslEndpointConfig, error) {
	if err := tx.Model(ep).Updates(update_ep).Error; err != nil {
		return nil, err
	}
	return ep, nil
}
func (db *KSDatabase) UpdateEpby_EpID(ep *types.DBAslEndpointConfig, update_ep_id uint) (*types.DBAslEndpointConfig, error) {
	gormDB := db.DB
	return Write(gormDB, func(tx *gorm.DB) (*types.DBAslEndpointConfig, error) {
		return updateEpby_EpID(tx, ep, update_ep_id)
	})
}
func updateEpby_EpID(tx *gorm.DB, ep *types.DBAslEndpointConfig, update_ep_id uint) (*types.DBAslEndpointConfig, error) {
	var update_ep types.DBAslEndpointConfig
	if err := tx.First(&update_ep, update_ep_id).Error; err != nil {
		return nil, err
	}
	if err := tx.Model(&update_ep).Updates(ep).Error; err != nil {
		return nil, err
	}
	return &update_ep, nil
}
