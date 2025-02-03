package kritis3m_control

import (
	"context"
	"fmt"

	"reflect"

	v1 "github.com/Laboratory-for-Safe-and-Secure-Systems/kritis3m_scale/gen/go/v1"
	database "github.com/Laboratory-for-Safe-and-Secure-Systems/kritis3m_scale/kritis3m_control/db"
	"github.com/Laboratory-for-Safe-and-Secure-Systems/kritis3m_scale/kritis3m_control/types"
	empty "github.com/golang/protobuf/ptypes/empty"
	"gorm.io/gorm"
)

// structToMap converts a struct to a map[string]interface{}, omitting zero-value fields.
func structToMap(s interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	val := reflect.ValueOf(s)
	typ := reflect.TypeOf(s)

	// Ensure that s is a pointer to a struct
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected a struct or pointer to struct, got %s", val.Kind())
	}
	count := 0
	// Iterate through struct fields
	for i := 0; i < val.NumField(); i++ {
		fieldValue := val.Field(i)
		fieldName := typ.Field(i).Name

		// Check if the field is non-zero
		if !fieldValue.IsZero() {
			count++
			result[fieldName] = fieldValue.Interface()
		}
	}

	if count == 0 {
		return nil, fmt.Errorf("error: all fields are zero values")
	}

	return result, nil
}

/******* GET *****************/
func (ks *Kritis3m_Scale) GetSelectedConfigurations(ctx context.Context, req *v1.GetSelectedConfigurationsRequest) (*v1.GetSelectedConfigurationsResponse, error) {
	var getSelConfig *types.SelectedConfiguration
	var includeAssociated bool = false

	if req.Queryoptions != nil {
		includeAssociated = req.Queryoptions.IncludeAssociatons
	}

	queryParams, err := structToMap(req.SelectedConfiguration)
	if err != nil {
		return nil, fmt.Errorf("failed to convert struct to map: %v", err)
	}

	if includeAssociated {
		getSelConfig, err = database.Read(ks.db.DB, func(rx *gorm.DB) (*types.SelectedConfiguration, error) {
			var selected types.SelectedConfiguration
			err := rx.Preload("Node").Preload("Config").Where(queryParams).First(&selected).Error
			if err != nil {
				return nil, err
			}
			return &selected, nil
		})
	} else {
		getSelConfig, err = database.Read(ks.db.DB, func(rx *gorm.DB) (*types.SelectedConfiguration, error) {
			var selected types.SelectedConfiguration
			err := rx.Where(queryParams).First(&selected).Error
			if err != nil {
				return nil, err
			}
			return &selected, nil
		})
	}

	if err != nil {
		return nil, fmt.Errorf("database query failed: %v", err)
	}
	returnvalue := getSelConfig.Proto()

	// Construct the response based on `getSelConfig`
	response := &v1.GetSelectedConfigurationsResponse{
		SelectedConfiguration: &returnvalue,
	}
	return response, nil
}

func (ks *Kritis3m_Scale) GetNode(ctx context.Context, req *v1.GetNodeRequest) (*v1.GetNodeResponse, error) {
	var includeAssociated bool = false

	if req.Queryoptions != nil {
		includeAssociated = req.Queryoptions.IncludeAssociatons
	}

	// Convert the request's Node struct to query parameters
	queryParams, err := structToMap(req.Node)
	if err != nil {
		return nil, fmt.Errorf("failed to convert struct to map: %v", err)
	}

	var getNode *types.DBNode
	if includeAssociated {
		// If including associated configurations
		getNode, err = database.Read(ks.db.DB, func(rx *gorm.DB) (*types.DBNode, error) {
			var node types.DBNode
			err := rx.Preload("Config").Where(queryParams).First(&node).Error
			if err != nil {
				return nil, err
			}
			return &node, nil
		})
	} else {
		// If not including associated configurations
		getNode, err = database.Read(ks.db.DB, func(rx *gorm.DB) (*types.DBNode, error) {
			var node types.DBNode
			err := rx.Where(queryParams).First(&node).Error
			if err != nil {
				return nil, err
			}
			return &node, nil
		})
	}

	if err != nil {
		return nil, fmt.Errorf("database query failed: %v", err)
	}

	// Convert to proto type
	returnValue := getNode.Proto()

	// Construct the response
	response := &v1.GetNodeResponse{
		Node: &returnValue,
	}

	return response, nil
}

func (ks *Kritis3m_Scale) GetConfig(ctx context.Context, req *v1.GetConfigRequest) (*v1.GetConfigResponse, error) {
	var includeAssociated bool = false
	var getAll bool = false

	if req.Queryoptions != nil {
		includeAssociated = req.Queryoptions.IncludeAssociatons
		getAll = req.Queryoptions.GetAll
	}

	queryParams, err := structToMap(req.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to convert struct to map: %v", err)
	}

	var getConfig []*types.DBNodeConfig
	if includeAssociated {
		getConfig, err = database.Read(ks.db.DB, func(rx *gorm.DB) ([]*types.DBNodeConfig, error) {
			var config []*types.DBNodeConfig
			query := rx.Preload("HardwareConfig").Preload("Whitelist").Preload("Application")
			if getAll {
				err := query.Find(&config).Error
				if err != nil {
					return nil, err
				}
			} else {
				err := query.Where(queryParams).First(&config).Error
				if err != nil {
					return nil, err
				}
			}
			return config, nil
		})
	} else {
		getConfig, err = database.Read(ks.db.DB, func(rx *gorm.DB) ([]*types.DBNodeConfig, error) {
			var config []*types.DBNodeConfig
			query := rx
			if getAll {
				err := query.Find(&config).Error
				if err != nil {
					return nil, err
				}
			} else {
				err := query.Where(queryParams).First(&config).Error
				if err != nil {
					return nil, err
				}
			}
			return config, nil
		})
	}

	if err != nil {
		return nil, fmt.Errorf("database query failed: %v", err)
	}
	var returnvalue []*v1.DBNodeConfig
	for _, node_config := range getConfig {
		v1_config := node_config.Proto()
		returnvalue = append(returnvalue, v1_config)
	}

	response := &v1.GetConfigResponse{
		Config: returnvalue,
	}

	return response, nil
}
func (ks *Kritis3m_Scale) GetApplication(ctx context.Context, req *v1.GetApplicationRequest) (*v1.GetApplicationResponse, error) {
	var includeAssociated bool = false
	var getAll bool = false

	if req.Queryoptions != nil {
		includeAssociated = req.Queryoptions.IncludeAssociatons
		getAll = req.Queryoptions.GetAll
	}

	queryParams, err := structToMap(req.Application)
	if err != nil {
		return nil, fmt.Errorf("failed to convert struct to map: %v", err)
	}

	var getApplication []*types.DBApplication
	if includeAssociated {
		getApplication, err = database.Read(ks.db.DB, func(rx *gorm.DB) ([]*types.DBApplication, error) {
			var application []*types.DBApplication
			query := rx.Preload("TrustedClients").Preload("Ep1").Preload("Ep2").Preload("Whitelist").Preload("HardwareConfig")
			if getAll {
				err := query.Find(&application).Error
				if err != nil {
					return nil, err
				}
			} else {
				err := query.Where(queryParams).First(&application).Error
				if err != nil {
					return nil, err
				}
			}
			return application, nil
		})
	} else {
		getApplication, err = database.Read(ks.db.DB, func(rx *gorm.DB) ([]*types.DBApplication, error) {
			var application []*types.DBApplication
			query := rx
			if getAll {
				err := query.Find(&application).Error
				if err != nil {
					return nil, err
				}
			} else {
				err := query.Where(queryParams).First(&application).Error
				if err != nil {
					return nil, err
				}
			}
			return application, nil
		})
	}

	if err != nil {
		return nil, fmt.Errorf("database query failed: %v", err)
	}
	var returnvalue []*v1.DBApplication
	for _, app := range getApplication {
		v1_app := app.Proto()
		returnvalue = append(returnvalue, v1_app)
	}

	response := &v1.GetApplicationResponse{
		Application: returnvalue,
	}

	return response, nil
}
func (ks *Kritis3m_Scale) GetWhitelist(ctx context.Context, req *v1.GetWhitelistRequest) (*v1.GetWhitelistResponse, error) {
	var includeAssociated bool = false

	if req.Queryoptions != nil {
		includeAssociated = req.Queryoptions.IncludeAssociatons
	}

	queryParams, err := structToMap(req.Whitelist)
	if err != nil {
		return nil, fmt.Errorf("failed to convert struct to map: %v", err)
	}

	var getWhitelist *types.DBWhitelist
	if includeAssociated {
		getWhitelist, err = database.Read(ks.db.DB, func(rx *gorm.DB) (*types.DBWhitelist, error) {
			var whitelist *types.DBWhitelist
			query := rx.Preload("TrustedClients")

			err := query.Where(queryParams).First(&whitelist).Error
			if err != nil {
				return nil, err
			}

			return whitelist, nil
		})
	} else {
		getWhitelist, err = database.Read(ks.db.DB, func(rx *gorm.DB) (*types.DBWhitelist, error) {
			var whitelist *types.DBWhitelist
			query := rx

			err := query.Where(queryParams).First(&whitelist).Error
			if err != nil {
				return nil, err
			}

			return whitelist, nil
		})
	}

	if err != nil {
		return nil, fmt.Errorf("database query failed: %v", err)
	}

	response := &v1.GetWhitelistResponse{
		Whitelist: getWhitelist.Proto(),
	}

	return response, nil
}
func (ks *Kritis3m_Scale) GetHardwareConfig(ctx context.Context, req *v1.GetHardwareConfigRequest) (*v1.GetHardwareConfigResponse, error) {
	var includeAssociated bool = false
	var getAll bool = false

	if req.Queryoptions != nil {
		includeAssociated = req.Queryoptions.IncludeAssociatons
		getAll = req.Queryoptions.GetAll
	}

	queryParams, err := structToMap(req.Hardwareconfig)
	if err != nil {
		return nil, fmt.Errorf("failed to convert struct to map: %v", err)
	}

	var getHardwareConfig []*types.HardwareConfig
	if includeAssociated {
		getHardwareConfig, err = database.Read(ks.db.DB, func(rx *gorm.DB) ([]*types.HardwareConfig, error) {
			var hardwareConfig []*types.HardwareConfig
			query := rx.Preload("AssociatedEntities") // Adjust the preload associations as needed
			if getAll {
				err := query.Find(&hardwareConfig).Error
				if err != nil {
					return nil, err
				}
			} else {
				err := query.Where(queryParams).First(&hardwareConfig).Error
				if err != nil {
					return nil, err
				}
			}
			return hardwareConfig, nil
		})
	} else {
		getHardwareConfig, err = database.Read(ks.db.DB, func(rx *gorm.DB) ([]*types.HardwareConfig, error) {
			var hardwareConfig []*types.HardwareConfig
			query := rx
			if getAll {
				err := query.Find(&hardwareConfig).Error
				if err != nil {
					return nil, err
				}
			} else {
				err := query.Where(queryParams).First(&hardwareConfig).Error
				if err != nil {
					return nil, err
				}
			}
			return hardwareConfig, nil
		})
	}

	if err != nil {
		return nil, fmt.Errorf("database query failed: %v", err)
	}
	var returnvalue []*v1.HardwareConfig
	for _, hwConfig := range getHardwareConfig {
		v1_hwConfig := hwConfig.Proto()
		returnvalue = append(returnvalue, v1_hwConfig)
	}

	response := &v1.GetHardwareConfigResponse{
		Hardwareconfig: returnvalue,
	}

	return response, nil
}
func (ks *Kritis3m_Scale) GetTrustedClients(ctx context.Context, req *v1.GetTrustedClientsRequest) (*v1.GetTrustedClientsResponse, error) {
	return nil, nil
}
func (ks *Kritis3m_Scale) GetEndpointConfig(ctx context.Context, req *v1.GetEndpointConfigRequest) (*v1.GetEndpointConfigResponse, error) {
}
func (ks *Kritis3m_Scale) GetIdentity(ctx context.Context, req *v1.GetIdentityRequest) (*v1.GetIdentityResponse, error) {
}

func (ks *Kritis3m_Scale) AddSelectedConfigurations(ctx context.Context, req *v1.AddSelectedConfigurationsRequest) (*v1.AddSelectedConfigurationsResponse, error) {
}
func (ks *Kritis3m_Scale) UpdateSelectedConfigurations(ctx context.Context, req *v1.UpdateSelectedConfigurationsRequest) (*v1.UpdateSelectedConfigurationsResponse, error) {
}
func (ks *Kritis3m_Scale) DeleteSelectedConfigurations(ctx context.Context, req *v1.DeleteSelectedConfigurationsRequest) (*v1.DeleteSelectedConfigurationsResponse, error) {
}
func (ks *Kritis3m_Scale) GetListSelectedConfigurations(ctx context.Context, req *empty.Empty) (*v1.GetListSelectedConfigurationsResponse, error) {
}
func (ks *Kritis3m_Scale) AddNode(ctx context.Context, req *v1.AddNodeRequest) (*v1.AddNodeResponse, error) {
}
func (ks *Kritis3m_Scale) UpdateNode(ctx context.Context, req *v1.UpdateNodeRequest) (*v1.UpdateNodeResponse, error) {
}
func (ks *Kritis3m_Scale) DeleteNode(ctx context.Context, req *v1.DeleteNodeRequest) (*v1.DeleteNodeResponse, error) {
}
func (ks *Kritis3m_Scale) GetListNode(ctx context.Context, req *empty.Empty) (*v1.GetListNodeResponse, error) {
}
func (ks *Kritis3m_Scale) AddConfig(ctx context.Context, req *v1.AddConfigRequest) (*v1.AddConfigResponse, error) {
}
func (ks *Kritis3m_Scale) UpdateConfig(ctx context.Context, req *v1.UpdateConfigRequest) (*v1.UpdateConfigResponse, error) {
}
func (ks *Kritis3m_Scale) DeleteConfig(ctx context.Context, req *v1.DeleteConfigRequest) (*v1.DeleteConfigResponse, error) {
}
func (ks *Kritis3m_Scale) GetListConfig(ctx context.Context, req *empty.Empty) (*v1.GetListConfigResponse, error) {
}
func (ks *Kritis3m_Scale) AddApplication(ctx context.Context, req *v1.AddApplicationRequest) (*v1.AddApplicationResponse, error) {
}
func (ks *Kritis3m_Scale) UpdateApplication(ctx context.Context, req *v1.UpdateApplicationRequest) (*v1.UpdateApplicationResponse, error) {
}
func (ks *Kritis3m_Scale) DeleteApplication(ctx context.Context, req *v1.DeleteApplicationRequest) (*v1.DeleteApplicationResponse, error) {
}
func (ks *Kritis3m_Scale) GetListApplicaiton(ctx context.Context, req *empty.Empty) (*v1.GetListApplicationResponse, error) {
}
func (ks *Kritis3m_Scale) AddWhitelist(ctx context.Context, req *v1.AddWhitelistRequest) (*v1.AddWhitelistResponse, error) {
}
func (ks *Kritis3m_Scale) UpdateWhitelist(ctx context.Context, req *v1.UpdateWhitelistRequest) (*v1.UpdateWhitelistResponse, error) {
}
func (ks *Kritis3m_Scale) DeleteWhitelist(ctx context.Context, req *v1.DeleteWhitelistRequest) (*v1.DeleteWhitelistResponse, error) {
}
func (ks *Kritis3m_Scale) GetListWhitelist(ctx context.Context, req *empty.Empty) (*v1.GetListWhitelistResponse, error) {
}
func (ks *Kritis3m_Scale) AddHardwareConfig(ctx context.Context, req *v1.AddHardwareConfigRequest) (*v1.AddHardwareConfigResponse, error) {
}
func (ks *Kritis3m_Scale) UpdateHardwareConfig(ctx context.Context, req *v1.UpdateHardwareConfigRequest) (*v1.UpdateHardwareConfigResponse, error) {
}
func (ks *Kritis3m_Scale) DeleteHardwareConfig(ctx context.Context, req *v1.DeleteHardwareConfigRequest) (*v1.DeleteHardwareConfigResponse, error) {
}
func (ks *Kritis3m_Scale) GetListHardwareConfig(ctx context.Context, req *empty.Empty) (*v1.GetListHardwareConfigResponse, error) {
}
func (ks *Kritis3m_Scale) AddTrustedClients(ctx context.Context, req *v1.AddTrustedClientsRequest) (*v1.AddTrustedClientsResponse, error) {
}
func (ks *Kritis3m_Scale) UpdateTrustedClients(ctx context.Context, req *v1.UpdateTrustedClientsRequest) (*v1.UpdateTrustedClientsResponse, error) {
}
func (ks *Kritis3m_Scale) DeleteTrustedClients(ctx context.Context, req *v1.DeleteTrustedClientsRequest) (*v1.DeleteHardwareConfigResponse, error) {
}
func (ks *Kritis3m_Scale) GetListTrustedClients(ctx context.Context, req *empty.Empty) (*v1.GetListTrustedClientsResponse, error) {
}
func (ks *Kritis3m_Scale) AddEndpointConfig(ctx context.Context, req *v1.AddEndpointConfigRequest) (*v1.AddEndpointConfigResponse, error) {
}
func (ks *Kritis3m_Scale) UpdateEndpointConfig(ctx context.Context, req *v1.UpdateEndpointConfigRequest) (*v1.UpdateEndpointConfigResponse, error) {
}
func (ks *Kritis3m_Scale) DeleteEndpointConfig(ctx context.Context, req *v1.DeleteEndpointConfigRequest) (*v1.DeleteEndpointConfigResponse, error) {
}
func (ks *Kritis3m_Scale) GetListEndpointConfig(ctx context.Context, req *empty.Empty) (*v1.GetListEndpointConfigResponse, error) {
}
func (ks *Kritis3m_Scale) AddIdentity(ctx context.Context, req *v1.AddIdentityRequest) (*v1.AddIdentityResponse, error) {
}
func (ks *Kritis3m_Scale) UpdateIdentity(ctx context.Context, req *v1.UpdateIdentityRequest) (*v1.UpdateIdentityResponse, error) {
}
func (ks *Kritis3m_Scale) DeleteIdentity(ctx context.Context, req *v1.DeleteIdentityRequest) (*v1.DeleteIdentityResponse, error) {
}
func (ks *Kritis3m_Scale) GetListIdentity(ctx context.Context, req *empty.Empty) (*v1.GetListIdentityResponse, error) {
}
func (ks *Kritis3m_Scale) mustEmbedUnimplementedNorthboundServiceServer() {}
