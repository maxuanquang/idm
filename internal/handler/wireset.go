package handler

import (
	"github.com/google/wire"
	"github.com/maxuanquang/idm/internal/handler/consumer"
	"github.com/maxuanquang/idm/internal/handler/grpc"
	"github.com/maxuanquang/idm/internal/handler/http"
)

var WireSet = wire.NewSet(
	grpc.WireSet,
	http.WireSet,
	consumer.WireSet,
)
