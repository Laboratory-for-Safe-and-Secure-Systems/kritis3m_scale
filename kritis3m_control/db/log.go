package db

import (
	"fmt"

	"github.com/Laboratory-for-Safe-and-Secure-Systems/kritis3m_scale/kritis3m_control/types"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func (ks_db *KSDatabase) PushMsg(msg types.Log_i) (int, error) {
	gormDB := ks_db.DB
	return Write(gormDB, func(db *gorm.DB) (int, error) { return pushMsg(db, msg) })

}

func (ks_db *KSDatabase) PushInfo(info_msg types.InfoLog) (int, error) {
	gormDB := ks_db.DB
	return Write(gormDB, func(db *gorm.DB) (int, error) { return pushInfo(db, info_msg) })
}

func (ks_db *KSDatabase) PushDebug(info_msg types.DebugLog) (int, error) {
	gormDB := ks_db.DB
	return Write(gormDB, func(db *gorm.DB) (int, error) { return pushDebug(db, info_msg) })
}

func (ks_db *KSDatabase) PushWarn(info_msg types.WarnLog) (int, error) {
	gormDB := ks_db.DB
	return Write(gormDB, func(db *gorm.DB) (int, error) { return pushWarn(db, info_msg) })
}

func (ks_db *KSDatabase) PushErr(info_msg types.ErrLog) (int, error) {
	gormDB := ks_db.DB
	return Write(gormDB, func(db *gorm.DB) (int, error) { return pushErr(db, info_msg) })
}

func pushMsg(tx *gorm.DB, msg types.Log_i) (int, error) {
	var table string
	switch msg.Get_level() {
	case types.Err:
		table = "err_logs"
	case types.Warn:
		table = "warn_logs"
	case types.Info:
		table = "info_logs"
	case types.Debug:
		table = "debug_logs"
	}
	insert_into := fmt.Sprintf(`insert into %s (node_id,component,msg) values (%d, "%s", "%s")`, table, msg.Get_nodeid(), *msg.Get_component(), *msg.Get_msg())
	if err := tx.Exec(insert_into).Error; err != nil {
		log.Err(err).Msg("Error inserting log")
		return -1, err
	}
	return 1, nil
}

func pushDebug(tx *gorm.DB, info_msg types.DebugLog) (int, error) {
	if err := tx.Create(&info_msg).Error; err != nil {
		return -1, err
	}
	return 1, nil
}

func pushInfo(tx *gorm.DB, info_msg types.InfoLog) (int, error) {
	if err := tx.Create(&info_msg).Error; err != nil {
		return -1, err
	}
	return 1, nil
}

func pushErr(tx *gorm.DB, info_msg types.ErrLog) (int, error) {
	if err := tx.Create(&info_msg).Error; err != nil {
		return -1, err
	}
	return 1, nil
}

func pushWarn(tx *gorm.DB, info_msg types.WarnLog) (int, error) {
	if err := tx.Create(&info_msg).Error; err != nil {
		return -1, err
	}
	return 1, nil
}
