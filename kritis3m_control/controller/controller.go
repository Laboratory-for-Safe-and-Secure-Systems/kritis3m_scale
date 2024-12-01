package controller

import (
	"github.com/Laboratory-for-Safe-and-Secure-Systems/kritis3m_scale/kritis3m_control/service"
	"github.com/gin-gonic/gin"
)

type NodeRegisterController interface {
	InitialAssignConfiguration(c *gin.Context)
	GetStatusReport(c *gin.Context)
}

type NodeHeartbeatController interface {
	RespondHeartbeatRequest(c *gin.Context)
}

type LogController interface {
	PushActiveConnections(c *gin.Context)
	PushConnectionRequests(c *gin.Context)
	PushInfo(c *gin.Context)
	PushWarn(c *gin.Context)
	PushErr(c *gin.Context)
	PushDebug(c *gin.Context)
}

// ------------------------   implementations ------------------------------

// ------------------------  Register

func NewNodeRegisterControllerImpl(svc service.NodeRegisterService) NodeRegisterControllerImpl {
	return NodeRegisterControllerImpl{svc: svc}
}

type NodeRegisterControllerImpl struct {
	svc service.NodeRegisterService
}

func (ctrl NodeRegisterControllerImpl) InitialAssignConfiguration(c *gin.Context) {
	ctrl.svc.InitialAssignConfiguration(c)
}

func (ctrl NodeRegisterControllerImpl) GetStatusReport(c *gin.Context) {
	ctrl.svc.GetStatusReport(c)
}

// ------------------------LogController

func NewLogControllerImpl(svc service.LogService) LogController { return LogControllerImpl{svc: svc} }

type LogControllerImpl struct {
	svc service.LogService
}

func (ctrl LogControllerImpl) PushActiveConnections(c *gin.Context) {
	ctrl.svc.PushActiveConnections(c)
}
func (ctrl LogControllerImpl) PushConnectionRequests(c *gin.Context) {
	ctrl.svc.PushConnectionRequests(c)
}
func (ctrl LogControllerImpl) PushInfo(c *gin.Context)  { ctrl.svc.PushInfo(c) }
func (ctrl LogControllerImpl) PushWarn(c *gin.Context)  { ctrl.svc.PushWarn(c) }
func (ctrl LogControllerImpl) PushErr(c *gin.Context)   { ctrl.svc.PushErr(c) }
func (ctrl LogControllerImpl) PushDebug(c *gin.Context) { ctrl.svc.PushErr(c) }
