//go:build wireinject
// +build wireinject

//go:generate go run github.com/google/wire/cmd/wire
package wiring

import (
	"github.com/google/wire"
	"github.com/maxuanquang/idm/internal/configs"
	"github.com/maxuanquang/idm/internal/dataaccess"
	"github.com/maxuanquang/idm/internal/handler"
	"github.com/maxuanquang/idm/internal/logic"
	"github.com/maxuanquang/idm/internal/utils"
	"github.com/maxuanquang/idm/internal/app"
)

var WireSet = wire.NewSet(
	configs.WireSet,
	dataaccess.WireSet,
	handler.WireSet,
	logic.WireSet,
	utils.WireSet,
	app.WireSet,
)

func InitializeStandaloneServer(configFilePath configs.ConfigFilePath) (app.StandaloneServer, func(), error) {
	wire.Build(WireSet)

	return app.StandaloneServer{}, nil, nil
}