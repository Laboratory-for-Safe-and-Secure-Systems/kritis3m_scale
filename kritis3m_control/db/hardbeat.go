package db

import (
	"time"

	"github.com/Laboratory-for-Safe-and-Secure-Systems/kritis3m_scale/kritis3m_control/types"
	"gorm.io/gorm"
)

func (ks_db *KSDatabase) GetTimeStampby_ID(id uint) (time.Time, error) {
	gormDB := ks_db.DB
	return Read(gormDB, func(db *gorm.DB) (time.Time, error) { return get_LatestNodeTimeStamp_ID(db, id) })
}

func (ks_db *KSDatabase) GetTimeStampby_sha(sha string) (time.Time, error) {
	gormDB := ks_db.DB
	return Read(gormDB, func(db *gorm.DB) (time.Time, error) { return get_LatestNodeTimeStamp_sha(db, sha) })
}

// local

func get_LatestNodeTimeStamp_sha(tx *gorm.DB, sha string) (time.Time, error) {
	var node types.DBNode

	err := tx.Select("updated_at").
		Where("NodeKey = ?", sha).
		Order("updated_at DESC").
		First(&node).Error
	if err != nil {
		return time.Time{}, err
	}
	return node.UpdatedAt, nil
}
func get_LatestNodeTimeStamp_ID(tx *gorm.DB, id uint) (time.Time, error) {
	var node types.DBNode
	// Query the database for the latest updated_at timestamp
	err := tx.Select("updated_at").
		Where("id = ?", id).
		Order("updated_at DESC").
		First(&node).Error
	if err != nil {
		return time.Time{}, err
	}
	return node.UpdatedAt, nil
}

//-------------------- implementations ------------------------
