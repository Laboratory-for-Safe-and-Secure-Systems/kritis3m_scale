package db

import (
	"github.com/philslol/kritis3m_scale/kritis3m_control/types"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func (db *KSDatabase) AddEp1to_Application(ep_id uint, application *types.DBApplication) (*types.DBApplication, error) {
	gormDB := db.DB
	return Write(gormDB, func(gormDB *gorm.DB) (*types.DBApplication, error) {
		return addEp1to_application(gormDB, ep_id, application)
	})
}
func addEp1to_application(tx *gorm.DB, ep_id uint, application *types.DBApplication) (*types.DBApplication, error) {
	error := tx.Model(&application).Where("id = ?", application.ID).Update("ep1_id", ep_id).Error
	return application, error
}

func (db *KSDatabase) AddEp1to_ApplicationID(ep_id uint, applicationID uint) (*types.DBApplication, error) {
	gormDB := db.DB
	application, err := db.GetApplicationby_ID(applicationID)
	if err != nil {
		log.Error().Msgf("no application found for application id %d", applicationID)
		return nil, err

	}
	return Write(gormDB, func(gormDB *gorm.DB) (*types.DBApplication, error) {
		return addEp1to_application(gormDB, ep_id, application)
	})
}

func (db *KSDatabase) AddEp2to_Application(ep_id uint, application *types.DBApplication) (*types.DBApplication, error) {
	gormDB := db.DB
	return Write(gormDB, func(gormDB *gorm.DB) (*types.DBApplication, error) {
		return addEp2to_application(gormDB, ep_id, application)
	})
}
func addEp2to_application(tx *gorm.DB, ep_id uint, application *types.DBApplication) (*types.DBApplication, error) {
	error := tx.Model(&application).Where("id = ?", application.ID).Update("ep2_id", ep_id).Error
	return application, error
}

func (db *KSDatabase) AddEp2to_ApplicationID(ep_id uint, applicationID uint) (*types.DBApplication, error) {
	gormDB := db.DB
	application, err := db.GetApplicationby_ID(applicationID)
	if err != nil {
		log.Error().Msgf("no application found for application id %d", applicationID)
		return nil, err

	}
	return Write(gormDB, func(gormDB *gorm.DB) (*types.DBApplication, error) {
		return addEp2to_application(gormDB, ep_id, application)
	})
}
