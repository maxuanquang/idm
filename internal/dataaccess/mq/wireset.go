package mq

import (
	"github.com/google/wire"
	"github.com/maxuanquang/idm/internal/dataaccess/mq/producer"
)

var WireSet = wire.NewSet(
	producer.WireSet,
)
