package types

import (
	"gorm.io/gorm"
)

type DebugLog struct {
	gorm.Model
	id        int    `gorm:" primaryKey:key; automncrement"`
	Node_id   int    `gorm:"column:node_id" json:"node_id"  `
	Component string `gorm:"column:component" json:"component"`
	Msg       string `gorm:"column:msg" json:"msg" `
}
type InfoLog struct {
	gorm.Model
	id        int    `gorm:" primaryKey:key; autoIncrement"`
	Node_id   int    `gorm:"column:node_id" json:"node_id"  `
	Component string `gorm:"column:component" json:"component"`
	Msg       string `gorm:"column:msg" json:"msg"`
}

type ErrLog struct {
	gorm.Model
	id        int    `gorm:" primaryKey:key; autoIncrement"`
	Node_id   int    `gorm:"column:node_id" json:"node_id"`
	Component string `gorm:"column:component" json:"component"`
	Msg       string `gorm:"column:msg" json:"msg"`
}

type WarnLog struct {
	gorm.Model
	id        int    `gorm:" primaryKey:key; autoIncrement"`
	Node_id   int    `gorm:"column:node_id" json:"node_id"`
	Component string `gorm:"column:component" json:"component"`
	Msg       string `gorm:"column:msg" json:"msg"`
}

type Log_i interface {
	Has_level() bool
	Has_nodeid() bool
	Has_msg() bool

	Get_nodeid() int
	Get_component() *string
	Get_msg() *string
	Get_level() int
	Has_required_fields() bool
}

const (
	Debug int = iota
	Info
	Warn
	Err
)

func (d DebugLog) Has_level() bool {
	return true
}

func (d DebugLog) Has_nodeid() bool {
	return d.Node_id != 0
}

func (d DebugLog) Has_msg() bool {
	return d.Msg != ""
}

func (d DebugLog) Has_required_fields() bool {
	if d.Has_nodeid() && d.Has_msg() {
		return true
	} else {
		return false
	}
}

func (d DebugLog) Get_nodeid() int {
	return d.Node_id
}

func (d DebugLog) Get_component() *string {
	return &d.Component
}

func (d DebugLog) Get_msg() *string {
	return &d.Msg
}

func (d DebugLog) Get_level() int {
	return Debug
}

// Implementing Logable for Err_Log
// ****************************************ERR_LOG****************************************
// *																					 *
// *																					 *
// ****************************************ERR_LOG****************************************

func (e ErrLog) Has_level() bool {
	return true
}

func (e ErrLog) Has_nodeid() bool {
	return e.Node_id != 0
}

func (e ErrLog) Has_msg() bool {

	return e.Msg != ""
}

func (e ErrLog) Get_nodeid() int {
	return e.Node_id
}

func (e ErrLog) Get_component() *string {
	return &e.Component
}

func (e ErrLog) Get_msg() *string {
	return &e.Msg
}

func (e ErrLog) Get_level() int {
	return Err
}

func (e ErrLog) Has_required_fields() bool {
	if e.Has_nodeid() && e.Has_msg() {
		return true
	} else {
		return false
	}
}

// ****************************************Info_Log****************************************
// *																					 *
// *																					 *
// ****************************************Info_Log****************************************

// Implementing Logable for Info_Log

func (i InfoLog) Has_level() bool {
	return true
}

func (i InfoLog) Has_nodeid() bool {
	return i.Node_id != 0
}

func (i InfoLog) Has_msg() bool {
	return i.Msg != ""
}

func (i InfoLog) Get_nodeid() int {
	return i.Node_id
}

func (i InfoLog) Get_level() int {
	return Info
}

func (i InfoLog) Get_component() *string {
	return &i.Component
}

func (i InfoLog) Get_msg() *string {
	return &i.Msg
}

func (i InfoLog) Has_required_fields() bool {
	if i.Has_nodeid() && i.Has_msg() {
		return true
	} else {
		return false
	}
}

// ****************************************Warn_Log****************************************
// *																					 *
// *																					 *
// ****************************************Warn_Log****************************************

func (w WarnLog) Has_level() bool {
	return true
}

func (w WarnLog) Has_nodeid() bool {
	return w.Node_id != 0
}

func (w WarnLog) Has_msg() bool {
	return w.Msg != ""
}

func (w WarnLog) Get_nodeid() int {
	return w.Node_id
}

func (w WarnLog) Get_level() int {
	return Warn
}

func (w WarnLog) Get_component() *string {
	return &w.Component
}

func (w WarnLog) Get_msg() *string {
	return &w.Msg
}

func (w WarnLog) Has_required_fields() bool {
	if w.Has_nodeid() && w.Has_msg() {
		return true
	} else {
		return false
	}
}
