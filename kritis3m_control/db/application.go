package db

import (
	"github.com/Laboratory-for-Safe-and-Secure-Systems/kritis3m_scale/kritis3m_control/types"
	"gorm.io/gorm"
)

//-------------------- interfaces ------------------------

func (db *KSDatabase) AddTrustedClientsto_Application(applications *types.DBApplication, trusted_clients []*types.DBTrustedClients) (*types.DBApplication, error) {
	gormDB := db.DB
	return Write(gormDB, func(tx *gorm.DB) (*types.DBApplication, error) {
		return addtrustedclientsto_applications(tx, applications, trusted_clients)
	})

}

func addtrustedclientsto_applications(tx *gorm.DB, applications *types.DBApplication, trusted_clients []*types.DBTrustedClients) (*types.DBApplication, error) {
	err := tx.Model(&applications).Association("TrustedClients").Append(trusted_clients)
	if err != nil {
		return nil, err
	}
	return applications, err
}

func (db *KSDatabase) GetAllApplications() ([]*types.DBApplication, error) {
	gormDB := db.DB
	return Read(gormDB, func(tx *gorm.DB) ([]*types.DBApplication, error) {
		return getAllApplications(tx)
	})
}
func getAllApplications(tx *gorm.DB) ([]*types.DBApplication, error) {
	var applications []*types.DBApplication
	if err := tx.Find(&applications).Error; err != nil {
		return nil, err
	}
	return applications, nil
}
func (db *KSDatabase) GetApplicationsByCfgID(cfg_id uint) ([]*types.DBApplication, error) {
	gormDB := db.DB
	return Read(gormDB, func(tx *gorm.DB) ([]*types.DBApplication, error) {
		return getApplicationsByCfgID(tx, cfg_id)
	})
}
func getApplicationsByCfgID(tx *gorm.DB, cfg_id uint) ([]*types.DBApplication, error) {
	var applications []*types.DBApplication
	if err := tx.Where("node_config_id = ?", cfg_id).Find(&applications).Error; err != nil {
		return nil, err
	}
	return applications, nil
}

func (db *KSDatabase) GetApplicationsby_Cfg(cfg *types.DBNodeConfig) ([]*types.DBApplication, error) {
	gormDB := db.DB
	return Read(gormDB, func(tx *gorm.DB) ([]*types.DBApplication, error) {
		return getApplicationsby_Cfg(tx, cfg)
	})
}

func getApplicationsby_Cfg(tx *gorm.DB, cfg *types.DBNodeConfig) ([]*types.DBApplication, error) {
	var applications []*types.DBApplication
	if err := tx.Where("node_config_id = ?", cfg.ID).Find(&applications).Error; err != nil {
		return nil, err
	}
	return applications, nil
}

func (db *KSDatabase) GetApplicationby_ID(id uint) (*types.DBApplication, error) {
	gormDB := db.DB
	return Read(gormDB, func(tx *gorm.DB) (*types.DBApplication, error) {
		return getApplicationby_ID(tx, id)
	})
}
func getApplicationby_ID(tx *gorm.DB, id uint) (*types.DBApplication, error) {
	var application types.DBApplication
	if err := tx.First(&application, id).Error; err != nil {
		return nil, err
	}
	return &application, nil
}

func (db *KSDatabase) AddApplicationsto_Cfg(applications []*types.DBApplication, trusted_clients []*types.DBTrustedClients, cfg *types.DBNodeConfig) ([]*types.DBApplication, error) {
	gormDB := db.DB
	return Write(gormDB, func(tx *gorm.DB) ([]*types.DBApplication, error) {
		return addApplicationsto_Cfg(tx, applications, trusted_clients, cfg)
	})
}
func addApplicationsto_Cfg(tx *gorm.DB, applications []*types.DBApplication, trusted_client []*types.DBTrustedClients, cfg *types.DBNodeConfig) ([]*types.DBApplication, error) {
	for _, app := range applications {
		app.NodeConfigID = cfg.ID
		if err := tx.Model(app).Association("TrustedClients").Append(trusted_client); err != nil {
			return nil, err
		}
	}
	return applications, nil
}

func (db *KSDatabase) AddApplicationto_Cfg(application *types.DBApplication, trusted_clients []*types.DBTrustedClients, cfg *types.DBNodeConfig) (*types.DBApplication, error) {
	gormDB := db.DB
	return Write(gormDB, func(tx *gorm.DB) (*types.DBApplication, error) {
		return addApplicationto_Cfg(tx, application, trusted_clients, cfg)
	})
}
func addApplicationto_Cfg(tx *gorm.DB, application *types.DBApplication, trusted_clients []*types.DBTrustedClients, cfg *types.DBNodeConfig) (*types.DBApplication, error) {
	application.NodeConfigID = cfg.ID
	if err := tx.Model(application).Association("TrustedClients").Append(trusted_clients); err != nil {
		return nil, err
	}
	return application, nil
}
func (db *KSDatabase) UpdateApplicationby_ApplicationID(application_id uint, update_values *types.DBApplication) (*types.DBApplication, error) {
	gormDB := db.DB
	return Write(gormDB, func(tx *gorm.DB) (*types.DBApplication, error) {
		return updateApplicationby_ApplicationID(tx, application_id, update_values)
	})
}
func updateApplicationby_ApplicationID(tx *gorm.DB, application_id uint, update_values *types.DBApplication) (*types.DBApplication, error) {
	var application types.DBApplication
	if err := tx.First(&application, application_id).Error; err != nil {
		return nil, err
	}
	if err := tx.Model(&application).Updates(update_values).Error; err != nil {
		return nil, err
	}
	return &application, nil
}
func (db *KSDatabase) UpdateApplicationby_Application(appl *types.DBApplication, update_values *types.DBApplication) (*types.DBApplication, error) {
	gormDB := db.DB
	return Write(gormDB, func(tx *gorm.DB) (*types.DBApplication, error) {
		return updateApplicationby_Application(tx, appl, update_values)
	})
}
func updateApplicationby_Application(tx *gorm.DB, appl *types.DBApplication, update_values *types.DBApplication) (*types.DBApplication, error) {
	if err := tx.Model(appl).Updates(update_values).Error; err != nil {
		return nil, err
	}
	return appl, nil
}
