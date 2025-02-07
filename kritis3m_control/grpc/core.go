package grpc

import (
	"context"
	"fmt"

	v1 "github.com/Laboratory-for-Safe-and-Secure-Systems/kritis3m_scale/gen/go/v1"
	"github.com/Laboratory-for-Safe-and-Secure-Systems/kritis3m_scale/kritis3m_control/db"
	"github.com/Laboratory-for-Safe-and-Secure-Systems/kritis3m_scale/kritis3m_control/types"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type NorthboundService struct {
	logger     *zerolog.Logger
	northbound v1.UnimplementedNorthboundServer
	db         *gorm.DB
}

func NewNorthboundService(db *gorm.DB, logger *zerolog.Logger) NorthboundService {
	return NorthboundService{
		db:     db,
		logger: logger,
	}
}

func (northbound NorthboundService) Get(ctx context.Context, req *v1.GetRequest) (*v1.GetResponse, error) {

	id := req.Id
	databasetype := req.Type

	var Entity *v1.Entity
	var err error

	switch databasetype {
	case v1.EntityType_NODE:
		{

			var updateNode types.UpdateNode = types.UpdateNode{ID: uint(id)}
			var node *types.UpdateNode
			node, err = db.Read(northbound.db, func(rx *gorm.DB) (*types.UpdateNode, error) {

				err := rx.First(&updateNode).Error
				if err != nil {
					return nil, err
				}
				return &updateNode, nil
			})
			Entity = &v1.Entity{

				Entity: &v1.Entity_Node{
					Node: node.Proto(),
				},
			}
		}
	case v1.EntityType_HWCONFIG:
		{
			var update_hwconfig types.UpdateHardwareConfig = types.UpdateHardwareConfig{ID: uint(id)}
			var hwconf *types.UpdateHardwareConfig
			hwconf, err = db.Read[*types.UpdateHardwareConfig](northbound.db, func(rx *gorm.DB) (*types.UpdateHardwareConfig, error) {

				err := rx.First(&update_hwconfig).Error
				if err != nil {
					return nil, err
				}
				return &update_hwconfig, nil

			})
			Entity = &v1.Entity{

				Entity: &v1.Entity_Hwconfig{
					Hwconfig: hwconf.Proto(),
				},
			}

		}
	case v1.EntityType_PROXY:
		{
			var update_proxy types.UpdateProxy = types.UpdateProxy{ID: uint(id)}
			var proxy *types.UpdateProxy
			proxy, err = db.Read[*types.UpdateProxy](northbound.db, func(rx *gorm.DB) (*types.UpdateProxy, error) {

				err := rx.First(&update_proxy).Error
				if err != nil {
					return nil, err
				}
				return &update_proxy, nil

			})

			Entity = &v1.Entity{

				Entity: &v1.Entity_Proxy{
					Proxy: proxy.Proto(),
				},
			}

		}
	case v1.EntityType_GROUP:
		{

			var update_group types.UpdateGroup = types.UpdateGroup{ID: uint(id)}
			var group *types.UpdateGroup
			group, err = db.Read[*types.UpdateGroup](northbound.db, func(rx *gorm.DB) (*types.UpdateGroup, error) {

				err := rx.First(&update_group).Error
				if err != nil {
					return nil, err
				}
				return &update_group, nil

			})
			Entity = &v1.Entity{

				Entity: &v1.Entity_Group{
					Group: group.Proto(),
				},
			}

		}
	case v1.EntityType_EPCONFIG:
		{
			var update_epconfig types.UpdateEndpointConfig = types.UpdateEndpointConfig{ID: uint(id)}
			var ep *types.UpdateEndpointConfig
			ep, err = db.Read[*types.UpdateEndpointConfig](northbound.db, func(rx *gorm.DB) (*types.UpdateEndpointConfig, error) {

				err := rx.First(&update_epconfig).Error
				if err != nil {
					return nil, err
				}
				return &update_epconfig, nil

			})

			Entity = &v1.Entity{

				Entity: &v1.Entity_Epconfig{
					Epconfig: ep.Proto(),
				},
			}

		}
	case v1.EntityType_IDENTITY:
		{

			var update_identity types.UpdateIdentity = types.UpdateIdentity{ID: uint(id)}
			var identity *types.UpdateIdentity
			identity, err = db.Read[*types.UpdateIdentity](northbound.db, func(rx *gorm.DB) (*types.UpdateIdentity, error) {

				err := rx.First(&update_identity).Error
				if err != nil {
					return nil, err
				}
				return &update_identity, nil

			})

			Entity = &v1.Entity{

				Entity: &v1.Entity_Identity{
					Identity: identity.Proto(),
				},
			}

		}
	default:
		{
			northbound.logger.Warn().Msg("identity unknown")
			err = fmt.Errorf("wrong entity provided")

		}

	}
	if err != nil {
		return nil, err
	}

	return &v1.GetResponse{Entity: Entity}, nil

}
func (northbound NorthboundService) GetAll(ctx context.Context, req *v1.GetAllRequest) (*v1.GetAllResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAll not implemented")
}
func (northbound NorthboundService) Add(ctx context.Context, req *v1.AddRequest) (*v1.AddResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Add not implemented")
}
func (northbound NorthboundService) Update(ctx context.Context, empty *empty.Empty) (*v1.AddResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (northbound NorthboundService) Delete(ctx context.Context, req *v1.DeleteRequest) (*v1.DeleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}

type CoreService struct {
	logger *zerolog.Logger
	core   v1.UnimplementedSDNcoreServer
	db     *gorm.DB
}

func NewCoreService(db *gorm.DB, logger *zerolog.Logger) CoreService {
	return CoreService{
		db:     db,
		logger: logger,
	}
}
