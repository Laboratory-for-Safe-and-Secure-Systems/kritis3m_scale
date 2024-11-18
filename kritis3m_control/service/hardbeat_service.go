package service

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/philslol/kritis3m_scale/kritis3m_control/db"
	"github.com/philslol/kritis3m_scale/kritis3m_control/types"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

/***
reasons to request a new configuration:
1. updated_at != updated_at
2. configuration_id != configuration_id
*/

type NodeHardbeatService interface {
	RespondHardbeatRequest(c *gin.Context)
}
type NodeHardbeatServiceImpl struct {
	db *db.KSDatabase
}

func NewNodeHardbeatServiceImpl(ks_db *db.KSDatabase) NodeHardbeatServiceImpl {
	return NodeHardbeatServiceImpl{db: ks_db}
}

func (svc NodeHardbeatServiceImpl) RespondHardbeatRequest(c *gin.Context) {
	identity, err := get_identity(c)
	if err != nil {
		log.Err(err).Msg("Identity is not complete. Request is missing either serialnumber or cfg id or version/updatedate")
		c.Set("caller_error", true)
		c.Error(types.ErrUnauthorized).SetMeta("DistributionService: RespondHardbeatRequest")
		return
	}
	log.Info().Msgf("Node with Serialnumber %s calls hardbeat service", identity.Serialnumber)

	// node is known but not configured
	// therefore Sleepmode is ordered
	is_up2date, err := svc.isconfiguration_up2date(identity)
	if err != nil {
		//No configuration found for Node
		if err == gorm.ErrRecordNotFound {
			log.Err(err).Msgf("No active configuration available for Node %s, respond shutdown", identity.Serialnumber)
			var response = HardbeatResponse{
				HBInstruction: HB_NOCONFIGAVAILABLE,
			}
			c.JSON(200, response)
			return
		} else {
			log.Err(err).Msg("Internal Error in RespondHardbeatRequest")
			c.Set("internal_error", true)
			c.Error(types.ErrInternalError).SetMeta("Hardbeat Service: Internal Error occured")
			return
		}
	}
	//when should a new version of the hardbeat service be called?
	//--1. newer version available
	//--2. completely new config available
	svc.db.UpdateLastSeenby_SerialNumber(identity.Serialnumber, time.Now())

	if is_up2date {
		var response = HardbeatResponse{
			HBInstruction: HB_NOTHING,
		}
		c.JSON(200, response)
		return
	} else {
		log.Info().Msgf("Node with Serialnumber %s is ordered to request a new configuration", identity.Serialnumber)
		var response = HardbeatResponse{
			HBInstruction: HB_REQUESTPOLICIES,
		}
		c.JSON(200, response)
		return
	}
}

func (svc NodeHardbeatServiceImpl) isconfiguration_up2date(identity Identity) (bool, error) {
	latest_config, err := svc.db.GetActiveConfigOfNodeby_SerialNumber(identity.Serialnumber)
	if err != nil {
		return false, err
	}
	if latest_config.ID != identity.Config_id ||
		latest_config.UpdatedAt != identity.Updated_at {
		return false, nil
	}
	return true, nil
}