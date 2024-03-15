//go:build wireinject
// +build wireinject

//go:generate go run github.com/google/wire/cmd/wire
package wiring

import (
	"github.com/google/wire"
	"github.com/maxuanquang/idm/internal/configs"
	"github.com/maxuanquang/idm/internal/dataaccess"
	"github.com/maxuanquang/idm/internal/handler"
	"github.com/maxuanquang/idm/internal/handler/grpc"
	"github.com/maxuanquang/idm/internal/logic"
)

var WireSet = wire.NewSet(
	configs.WireSet,
	dataaccess.WireSet,
	handler.WireSet,
	logic.WireSet,
)

func InitializeGRPCServer(configFilePath configs.ConfigFilePath) (grpc.Server, func(), error) {
	wire.Build(WireSet)

	return nil, nil, nil
}
