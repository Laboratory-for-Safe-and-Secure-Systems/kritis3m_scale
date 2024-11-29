package db

import (
	"github.com/Laboratory-for-Safe-and-Secure-Systems/kritis3m_scale/kritis3m_control/types"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func (db *KSDatabase) GetConfigFor_DistributionService(node_id uint) (*types.DBNodeConfig, error) {
	gormDB := db.DB
	return Read(gormDB, func(db *gorm.DB) (*types.DBNodeConfig, error) {
		return getconfigfor_distributionservice(db, node_id)
	})
}

func getconfigfor_distributionservice(tx *gorm.DB, node_id uint) (*types.DBNodeConfig, error) {
	var nodeConfig *types.DBNodeConfig

	if err := tx.Preload("Whitelist.TrustedClients").Preload("HardwareConfig").Preload("Application").Where(&types.DBNodeConfig{NodeID: node_id}).First(&nodeConfig).Error; err != nil {
		log.Err(err).Msgf("error getting configuration")
		return nil, err
	}

	// Fetch only Application IDs for each trusted client
	for i, trustedClient := range nodeConfig.Whitelist.TrustedClients {
		var appIDs []uint
		if err := tx.Model(&types.ApplicationTrustsClients{}).
			Where("db_trusted_clients_id = ?", trustedClient.ID).
			Pluck("db_application_id", &appIDs).Error; err != nil {
			log.Err(err).Msgf("error getting application IDs for trusted client %d", trustedClient.ID)
			return nil, err
		}
		nodeConfig.Whitelist.TrustedClients[i].ApplicationIDs = appIDs
	}

	return nodeConfig, nil
}

func (db *KSDatabase) GetIdentities(identity_ids []uint) ([]*types.DBIdentity, error) {
	gormDB := db.DB
	return Read(gormDB, func(db *gorm.DB) ([]*types.DBIdentity, error) {
		return getIdentities(db, identity_ids)
	})
}

func getIdentities(tx *gorm.DB, identity_ids []uint) ([]*types.DBIdentity, error) {
	var identities []*types.DBIdentity
	err := tx.Where(identity_ids).Find(&identities).Error
	if err != nil {
		log.Err(err).Msgf("error in getting identities")
	}

	return identities, nil
}
