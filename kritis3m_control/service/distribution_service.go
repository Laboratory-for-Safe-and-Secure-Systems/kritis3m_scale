package service

import (
	"strconv"

	"github.com/Laboratory-for-Safe-and-Secure-Systems/kritis3m_scale/kritis3m_control/db"
	"github.com/Laboratory-for-Safe-and-Secure-Systems/kritis3m_scale/kritis3m_control/types"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	//types
)

type NodeRegisterService interface {
	InitialAssignConfiguration(c *gin.Context)
	GetStatusReport(c *gin.Context)
}

type NodeRegisterServiceImpl struct {
	db *db.KSDatabase
}

// StatusReportPayload represents the structure of the incoming status report
type StatusReportPayload struct {
	Status              int `json:"status"`
	RunningApplications int `json:"running_applications"`
}

func NewNodeRegisterServiceImpl(ks_db *db.KSDatabase) NodeRegisterServiceImpl {
	return NodeRegisterServiceImpl{db: ks_db}
}

func (svc NodeRegisterServiceImpl) GetStatusReport(c *gin.Context) {
	// Extract URL parameters
	serialNumber := c.Param("serialnumber")
	configID := c.Param("config_id")
	versionNumber := c.Param("version_number")

	// Validate URL parameters
	if serialNumber == "" || configID == "" || versionNumber == "" {
		c.Set("caller_error", true)
		c.Error(types.ErrInvalidParam).SetMeta("GetStatusReport: Imissing url params")
		return
	}

	// Parse version number to ensure it's a valid integer
	version, err := strconv.Atoi(versionNumber)
	if err != nil {
		c.Set("caller_error", true)
		c.Error(types.ErrInvalidParam).SetMeta("GetStatusReport: can't parse version number")
		return
	}
	log.Info().Msgf("version number of client is %d", version)

	// Parse version number to ensure it's a valid integer
	cfg_id, err := strconv.Atoi(configID)
	if err != nil {
		c.Set("caller_error", true)
		c.Error(types.ErrInvalidParam).SetMeta("GetStatusReport: can't parse config_id")
		return
	}

	// Parse the incoming JSON payload
	var payload StatusReportPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.Set("caller_error", true)
		c.Error(types.ErrInvalidParam).SetMeta("GetStatusReport: can't parse payload")
		return
	}
	if payload.Status == -1 {
		svc.db.ActiveConfigSetState_byCfgID(uint(cfg_id), types.ErrorState)
		log.Error().Msgf("error occured, at client with serial_number: %s, when applying new configuration(cfg_id: %d)", serialNumber, cfg_id)
		c.JSON(200, gin.H{
			"message": "Status report received successfully",
		})
		return
	}
	svc.db.ActiveConfigSetState_byCfgID(uint(cfg_id), types.Running)

	// Successful response
	c.JSON(200, gin.H{
		"message": "Status report received successfully",
	})
}

func (svc NodeRegisterServiceImpl) InitialAssignConfiguration(c *gin.Context) {
	serial_number, err := get_serialnumber(c)
	if err != nil {
		log.Err(err).Msg("cant get identity of request")

		c.Set("caller_error", true)
		c.Error(types.ErrInvalidParam).SetMeta("DistributionService: InitialAssignConfiguration")
		return
	}
	log.Info().Msgf("Initial Configuration Request from Node with Serialnumber %s", serial_number)
	//Get Config:

	node, err := svc.db.GetNodeby_SerialNumber(serial_number)
	if err != nil {
		log.Err(err).Msgf("No Node with Serialnumber %s in Database", serial_number)
		c.Set("internal_error", true)
		c.Error(types.ErrUnauthorized).SetMeta("DistributionService: InitialAssignConfiguration")
		return
	}

	//check if active config exists
	//TODO:check this function
	config, err := svc.db.GetConfigFor_DistributionService(node.ID)
	if err != nil {
		log.Err(err).Msgf("No Configuration available for Node with Serialnumber %s ", serial_number)
		c.Set("internal_error", true)
		c.Error(types.ErrInternalError).SetMeta("DistributionService: InitialAssignConfiguration")
		return
	}
	var endpoint_ids []uint
	//get relevant endpoints
	for _, app := range config.Application {
		endpoint_ids = addToListIfNotPresent(endpoint_ids, app.Ep1ID)
		endpoint_ids = addToListIfNotPresent(endpoint_ids, app.Ep2ID)
	}

	//returnvalue
	var distribution_response types.DistributionResponse
	for _, ep_id := range endpoint_ids {
		ep, err := svc.db.GetEPby_ID(ep_id)
		if err != nil {
			log.Err(err).Msgf("No EP Configuration found for config id  %d ", ep_id)
			c.Set("internal_error", true)
			c.Error(types.ErrInternalError).SetMeta("DistributionService: InitialAssignConfiguration")
			return
		}
		distribution_response.CryptoConfig = append(distribution_response.CryptoConfig, ep)
	}

	var identity_ids []uint
	for _, ep := range distribution_response.CryptoConfig {
		identity_ids = addToListIfNotPresent(identity_ids, ep.IdentityID)
	}
	distribution_response.Identities, err = svc.db.GetIdentities(identity_ids)
	if err != nil {
		log.Err(err).Msgf("Problem with identities ")
		c.Set("internal_error", true)
		c.Error(types.ErrInternalError).SetMeta("DistributionService: InitialAssignConfiguration")
		return
	}

	var lconfig []*types.DBNodeConfig
	lconfig = append(lconfig, config)
	node.Config = lconfig
	distribution_response.Node = *node

	//send out
	c.JSON(200, distribution_response)
	svc.db.ActiveConfigSetState_byCfgID(distribution_response.Node.Config[0].ID, types.NodeRequestedConfig)
}

func addToListIfNotPresent(list []uint, item uint) []uint {
	if item == 0 {
		return list
	}

	for _, v := range list {
		if v == item {
			return list // Item is already in the list, return unchanged list
		}
	}
	// Item is not in the list, append it
	return append(list, item)
}
