package service

import (
	"github.com/Laboratory-for-Safe-and-Secure-Systems/kritis3m_scale/kritis3m_control/db"
	"github.com/Laboratory-for-Safe-and-Secure-Systems/kritis3m_scale/kritis3m_control/types"
	"github.com/gin-gonic/gin"
)

type LogService interface {
	PushActiveConnections(c *gin.Context)
	PushConnectionRequests(c *gin.Context)
	PushInfo(c *gin.Context)
	PushWarn(c *gin.Context)
	PushErr(c *gin.Context)
	PushDebug(c *gin.Context)
}

type LogServiceImpl struct {
	log_db *db.KSDatabase
}

func NewLogServiceImpl(ks_db *db.KSDatabase) LogServiceImpl {
	return LogServiceImpl{log_db: ks_db}
}

func (svc LogServiceImpl) PushActiveConnections(c *gin.Context) {
}
func (svc LogServiceImpl) PushConnectionRequests(c *gin.Context) {}

func (svc LogServiceImpl) PushInfo(c *gin.Context) {
	var info types.InfoLog
	if err := c.ShouldBindJSON(&info); err != nil {
		c.Set("caller_error", true)
		c.Error(types.ErrNoID).SetMeta("Info Service: PushInfoby_ID")
		return
	}
	svc.log_db.PushInfo(info)
}
func (svc LogServiceImpl) PushWarn(c *gin.Context) {
	var warn types.WarnLog
	if err := c.ShouldBindJSON(&warn); err != nil {
		c.Set("caller_error", true)
		c.Error(types.ErrNoID).SetMeta("Info Service: PushInfoby_ID")
		return
	}
	if warn.Has_required_fields() {
		svc.log_db.PushMsg(warn)
	} else {
		c.Set("caller_error", true)
		c.Error(types.ErrInvalidParam).SetMeta("Push warn log: missing required fields")
		return
	}
}

func (svc LogServiceImpl) PushErr(c *gin.Context) {
	var err types.ErrLog
	if err := c.ShouldBindJSON(&err); err != nil {
		c.Set("caller_error", true)
		c.Error(types.ErrNoID).SetMeta("Info Service: PushInfoby_ID")
		return
	}
	if err.Has_required_fields() {
		svc.log_db.PushMsg(err)
	}
}
func (svc LogServiceImpl) PushDebug(c *gin.Context) {
	var debug types.DebugLog
	if err := c.ShouldBindJSON(&debug); err != nil {
		c.Set("caller_error", true)
		c.Error(types.ErrNoID).SetMeta("Info Service: PushInfoby_ID")
		return
	}
	if debug.Has_required_fields() {
		svc.log_db.PushMsg(debug)
	}
}
