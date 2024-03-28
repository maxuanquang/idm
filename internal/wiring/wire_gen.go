// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package wiring

import (
	"github.com/google/wire"
	"github.com/maxuanquang/idm/internal/app"
	"github.com/maxuanquang/idm/internal/configs"
	"github.com/maxuanquang/idm/internal/dataaccess"
	"github.com/maxuanquang/idm/internal/dataaccess/cache"
	"github.com/maxuanquang/idm/internal/dataaccess/database"
	"github.com/maxuanquang/idm/internal/dataaccess/file"
	consumer2 "github.com/maxuanquang/idm/internal/dataaccess/mq/consumer"
	"github.com/maxuanquang/idm/internal/dataaccess/mq/producer"
	"github.com/maxuanquang/idm/internal/handler"
	"github.com/maxuanquang/idm/internal/handler/consumer"
	"github.com/maxuanquang/idm/internal/handler/grpc"
	"github.com/maxuanquang/idm/internal/handler/http"
	"github.com/maxuanquang/idm/internal/logic"
	"github.com/maxuanquang/idm/internal/utils"
)

// Injectors from wire.go:

func InitializeAppServer(configFilePath configs.ConfigFilePath) (app.Server, func(), error) {
	config, err := configs.NewConfig(configFilePath)
	if err != nil {
		return app.Server{}, nil, err
	}
	configsGRPC := config.GRPC
	configsDatabase := config.Database
	databaseDatabase, cleanup, err := database.InitializeDB(configsDatabase)
	if err != nil {
		return app.Server{}, nil, err
	}
	log := config.Log
	logger, cleanup2, err := utils.InitializeLogger(log)
	if err != nil {
		cleanup()
		return app.Server{}, nil, err
	}
	accountDataAccessor := database.NewAccountDataAccessor(databaseDatabase, logger)
	accountPasswordDataAccessor := database.NewAccountPasswordDataAccessor(databaseDatabase, logger)
	hashLogic := logic.NewHashLogic()
	tokenPublicKeyDataAccessor, err := database.NewTokenPublicKeyDataAccessor(databaseDatabase, logger)
	if err != nil {
		cleanup2()
		cleanup()
		return app.Server{}, nil, err
	}
	auth := config.Auth
	configsCache := config.Cache
	client, err := cache.NewClient(configsCache, logger)
	if err != nil {
		cleanup2()
		cleanup()
		return app.Server{}, nil, err
	}
	tokenPublicKey, err := cache.NewTokenPublicKey(client)
	if err != nil {
		cleanup2()
		cleanup()
		return app.Server{}, nil, err
	}
	tokenLogic, err := logic.NewTokenLogic(accountDataAccessor, tokenPublicKeyDataAccessor, logger, auth, tokenPublicKey)
	if err != nil {
		cleanup2()
		cleanup()
		return app.Server{}, nil, err
	}
	takenAccountName, err := cache.NewTakenAccountName(client)
	if err != nil {
		cleanup2()
		cleanup()
		return app.Server{}, nil, err
	}
	accountLogic := logic.NewAccountLogic(databaseDatabase, accountDataAccessor, accountPasswordDataAccessor, hashLogic, tokenLogic, takenAccountName, logger)
	downloadTaskDataAccessor := database.NewDownloadTaskDataAccessor(databaseDatabase, logger)
	mq := config.MQ
	producerClient, err := producer.NewClient(mq, logger)
	if err != nil {
		cleanup2()
		cleanup()
		return app.Server{}, nil, err
	}
	downloadTaskCreatedProducer, err := producer.NewDownloadTaskCreatedProducer(producerClient, logger)
	if err != nil {
		cleanup2()
		cleanup()
		return app.Server{}, nil, err
	}
	download := config.Download
	fileClient, err := file.NewClient(download, logger)
	if err != nil {
		cleanup2()
		cleanup()
		return app.Server{}, nil, err
	}
	downloadTaskLogic, err := logic.NewDownloadTaskLogic(tokenLogic, downloadTaskDataAccessor, downloadTaskCreatedProducer, fileClient, databaseDatabase, logger)
	if err != nil {
		cleanup2()
		cleanup()
		return app.Server{}, nil, err
	}
	idmServiceServer := grpc.NewHandler(accountLogic, downloadTaskLogic)
	server := grpc.NewServer(configsGRPC, idmServiceServer)
	configsHTTP := config.HTTP
	httpServer := http.NewServer(configsHTTP, configsGRPC, auth, logger)
	downloadTaskCreatedHandler, err := consumer.NewDownloadTaskCreatedHandler(downloadTaskLogic, logger)
	if err != nil {
		cleanup2()
		cleanup()
		return app.Server{}, nil, err
	}
	consumerConsumer, err := consumer2.NewConsumer(mq, logger)
	if err != nil {
		cleanup2()
		cleanup()
		return app.Server{}, nil, err
	}
	rootConsumer := consumer.NewRootConsumer(downloadTaskCreatedHandler, consumerConsumer, logger)
	appServer, err := app.NewServer(server, httpServer, rootConsumer, logger)
	if err != nil {
		cleanup2()
		cleanup()
		return app.Server{}, nil, err
	}
	return appServer, func() {
		cleanup2()
		cleanup()
	}, nil
}

// wire.go:

var WireSet = wire.NewSet(configs.WireSet, dataaccess.WireSet, handler.WireSet, logic.WireSet, utils.WireSet, app.WireSet)
