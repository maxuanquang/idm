package dataaccess

import (
	"github.com/google/wire"
	"github.com/maxuanquang/idm/internal/dataaccess/database"
)

var WireSet = wire.NewSet(
	database.WireSet,
)