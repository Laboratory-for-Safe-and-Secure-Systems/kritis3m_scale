package db

import (
	"github.com/philslol/kritis3m_scale/kritis3m_control/types"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

//-------------------- interfaces ------------------------

type TrustedClientsHandler interface {
	GetAllTC() ([]*types.DBTrustedClients, error)
	GetTCby_CfgID(cfg_id uint) ([]*types.DBTrustedClients, error)
	GetTCby_WhitelistID(whitelist_id uint) ([]*types.DBTrustedClients, error)
	GetTCby_ApplicationID(application_id uint) ([]*types.DBTrustedClients, error)
	GetTCby_ID(id uint) (*types.DBTrustedClients, error)

	AddTCto_Whitelist(whitelist_id uint, trusted_client *types.DBTrustedClients) error
	AddTCto_WhitelistValues(whitelist_id uint, client_ip_port string) error
	AddTCto_ApplicationID(application_id uint, trusted_client *types.DBTrustedClients) error              //client must be already present in a whitelist
	AddTCto_Application(application []*types.DBApplication, trusted_client *types.DBTrustedClients) error //client must be already present in a whitelist

	AddTCto_ApplicationsID(application_id []*uint, trusted_client *types.DBTrustedClients) error              //client must be already present in a whitelist
	AddTCbyIDto_ApplicationsID(application_id []*uint, trusted_client uint) error                             //client must be already present in a whitelist
	AddTCto_Applications(application_id []*types.DBApplication, trusted_client *types.DBTrustedClients) error //client must be already present in a whitelist
	AddTCbyIDto_Applications(application_id []*types.DBApplication, trusted_client uint) error                //client must be already present in a whitelist
}

func (db *KSDatabase) GetAllTC() ([]*types.DBTrustedClients, error) {
	gormDB := db.DB
	return Read(gormDB, func(tx *gorm.DB) ([]*types.DBTrustedClients, error) {
		return getAllTC(tx)
	})
}

func getAllTC(tx *gorm.DB) ([]*types.DBTrustedClients, error) {
	var clients []*types.DBTrustedClients
	if err := tx.Find(&clients).Error; err != nil {
		return nil, err
	}
	return clients, nil
}

func (db *KSDatabase) GetTCby_CfgID(cfg_id uint) ([]*types.DBTrustedClients, error) {
	gormDB := db.DB
	return Read(gormDB, func(tx *gorm.DB) ([]*types.DBTrustedClients, error) {
		return getTCby_CfgID(tx, cfg_id)
	})
}

func getTCby_CfgID(tx *gorm.DB, cfg_id uint) ([]*types.DBTrustedClients, error) {
	var clients []*types.DBTrustedClients
	if err := tx.Joins("JOIN db_whitelists ON db_trusted_clients.whitelist_id = db_whitelists.id").
		Where("db_whitelists.node_config_id = ?", cfg_id).
		Find(&clients).Error; err != nil {
		return nil, err
	}
	return clients, nil
}
func (db *KSDatabase) GetTCby_WhitelistID(whitelist_id uint) ([]*types.DBTrustedClients, error) {
	gormDB := db.DB
	return Read(gormDB, func(tx *gorm.DB) ([]*types.DBTrustedClients, error) {
		return getTCby_WhitelistID(tx, whitelist_id)
	})
}

func getTCby_WhitelistID(tx *gorm.DB, whitelist_id uint) ([]*types.DBTrustedClients, error) {
	var clients []*types.DBTrustedClients
	if err := tx.Where("whitelist_id = ?", whitelist_id).Find(&clients).Error; err != nil {
		return nil, err
	}
	return clients, nil
}
func (db *KSDatabase) GetTCby_ApplicationID(application_id uint) ([]*types.DBTrustedClients, error) {
	gormDB := db.DB
	return Read(gormDB, func(tx *gorm.DB) ([]*types.DBTrustedClients, error) {
		return getTCby_ApplicationID(tx, application_id)
	})
}

func getTCby_ApplicationID(tx *gorm.DB, application_id uint) ([]*types.DBTrustedClients, error) {
	var clients []*types.DBTrustedClients
	if err := tx.Joins("JOIN application_trusts_clients ON db_trusted_clients.id = application_trusts_clients.db_trusted_clients_id").
		Where("application_trusts_clients.db_application_id = ?", application_id).
		Find(&clients).Error; err != nil {
		return nil, err
	}
	return clients, nil
}
func (db *KSDatabase) GetTCby_ID(id uint) (*types.DBTrustedClients, error) {
	gormDB := db.DB
	return Read(gormDB, func(tx *gorm.DB) (*types.DBTrustedClients, error) {
		return getTCby_ID(tx, id)
	})
}

func getTCby_ID(tx *gorm.DB, id uint) (*types.DBTrustedClients, error) {
	var client types.DBTrustedClients
	if err := tx.Where("id = ?", id).First(&client).Error; err != nil {
		return nil, err
	}
	return &client, nil
}
func (db *KSDatabase) AddTCto_Whitelist(whitelist_id uint, trusted_client *types.DBTrustedClients) error {
	return db.Write(func(tx *gorm.DB) error {
		return addTCto_Whitelist(tx, whitelist_id, trusted_client)
	})
}

func addTCto_Whitelist(tx *gorm.DB, whitelist_id uint, trusted_client *types.DBTrustedClients) error {
	trusted_client.WhitelistID = whitelist_id
	return tx.Create(trusted_client).Error
}
func (db *KSDatabase) AddTCto_WhitelistValues(whitelist_id uint, client_ip_port string) error {
	return db.Write(func(tx *gorm.DB) error {
		return addTCto_WhitelistValues(tx, whitelist_id, client_ip_port)
	})
}

func addTCto_WhitelistValues(tx *gorm.DB, whitelist_id uint, client_ip_port string) error {
	trusted_client := &types.DBTrustedClients{
		WhitelistID:        whitelist_id,
		ClientEndpointAddr: client_ip_port,
	}
	return tx.Create(trusted_client).Error
}
func (db *KSDatabase) AddTCto_ApplicationID(application_id uint, trusted_client *types.DBTrustedClients) error {
	_, err := db.GetWhitelistby_ID(trusted_client.ID)
	if err != nil {
		log.Error().Msg("trusted client does not live in a whitelist. Please add the clieNewNodeHeartbeatServiceImplnt first to a whitelist before connecting to an application")
		return err
	}
	return db.Write(func(tx *gorm.DB) error {
		return addTCto_ApplicationID(tx, application_id, trusted_client)
	})
}

func addTCto_ApplicationID(tx *gorm.DB, application_id uint, trusted_client *types.DBTrustedClients) error {
	// Assuming the client must already be present in a whitelist, verify presence
	if err := tx.First(trusted_client, trusted_client.ID).Error; err != nil {
		return err
	}
	return tx.Model(&types.DBApplication{ID: application_id}).Association("TrustedClients").Append(trusted_client)
}
func (db *KSDatabase) AddTCto_ApplicationsID(application_ids []*uint, trusted_client *types.DBTrustedClients) error {
	_, err := db.GetWhitelistby_ID(trusted_client.ID)
	if err != nil {
		log.Error().Msg("trusted client does not live in a whitelist. Please add the client first to a whitelist before connecting to an application")
		return err
	}
	return db.Write(func(tx *gorm.DB) error {
		return addTCto_ApplicationsID(tx, application_ids, trusted_client)
	})
}

func addTCto_ApplicationsID(tx *gorm.DB, application_ids []*uint, trusted_client *types.DBTrustedClients) error {
	if err := tx.First(trusted_client, trusted_client.ID).Error; err != nil {
		return err
	}
	for _, app_id := range application_ids {
		if err := tx.Model(&types.DBApplication{ID: *app_id}).Association("TrustedClients").Append(trusted_client); err != nil {
			return err
		}
	}
	return nil
}

func (db *KSDatabase) AddTCbyIDto_ApplicationsID(application_ids []*uint, trusted_client_id uint) error {

	_, err := db.GetWhitelistby_ID(trusted_client_id)
	if err != nil {
		log.Error().Msg("trusted client does not live in a whitelist. Please add the client first to a whitelist before connecting to an application")
		return err
	}
	return db.Write(func(tx *gorm.DB) error {
		return addTCbyIDto_ApplicationsID(tx, application_ids, trusted_client_id)
	})
}

func addTCbyIDto_ApplicationsID(tx *gorm.DB, application_ids []*uint, trusted_client_id uint) error {
	var trusted_client types.DBTrustedClients
	if err := tx.First(&trusted_client, trusted_client_id).Error; err != nil {
		return err
	}
	for _, app_id := range application_ids {
		if err := tx.Model(&types.DBApplication{ID: *app_id}).Association("TrustedClients").Append(&trusted_client); err != nil {
			return err
		}
	}
	return nil
}

func (db *KSDatabase) AddTCto_Applications(application []*types.DBApplication, trusted_client *types.DBTrustedClients) error {
	_, err := db.GetWhitelistby_ID(trusted_client.ID)
	if err != nil {
		log.Error().Msg("trusted client does not live in a whitelist. Please add the client first to a whitelist before connecting to an application")
		return err
	}
	return db.Write(func(tx *gorm.DB) error {
		return addTCto_Applications(tx, application, trusted_client)
	})
}

func addTCto_Applications(tx *gorm.DB, application []*types.DBApplication, trusted_client *types.DBTrustedClients) error {
	if err := tx.First(trusted_client, trusted_client.ID).Error; err != nil {
		return err
	}
	for _, app := range application {
		if err := tx.Model(app).Association("TrustedClients").Append(trusted_client); err != nil {
			return err
		}
	}
	return nil
}
func (db *KSDatabase) AddTCbyIDto_Applications(application []*types.DBApplication, trusted_client_id uint) error {
	_, err := db.GetWhitelistby_ID(trusted_client_id)
	if err != nil {
		log.Error().Msg("trusted client does not live in a whitelist. Please add the client first to a whitelist before connecting to an application")
		return err
	}
	return db.Write(func(tx *gorm.DB) error {
		return addTCbyIDto_Applications(tx, application, trusted_client_id)
	})
}

func addTCbyIDto_Applications(tx *gorm.DB, applications []*types.DBApplication, trusted_client_id uint) error {
	var trusted_client types.DBTrustedClients
	if err := tx.First(&trusted_client, trusted_client_id).Error; err != nil {
		return err
	}

	for _, app := range applications {
		if err := tx.Model(app).Association("TrustedClients").Append(&trusted_client); err != nil {
			return err
		}
	}
	return nil
}
