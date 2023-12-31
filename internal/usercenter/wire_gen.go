// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package usercenter

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/superproj/zero/internal/pkg/bootstrap"
	validation2 "github.com/superproj/zero/internal/pkg/validation"
	"github.com/superproj/zero/internal/usercenter/auth"
	"github.com/superproj/zero/internal/usercenter/biz"
	"github.com/superproj/zero/internal/usercenter/server"
	"github.com/superproj/zero/internal/usercenter/service"
	"github.com/superproj/zero/internal/usercenter/store"
	"github.com/superproj/zero/internal/usercenter/validation"
	"github.com/superproj/zero/pkg/db"
	"github.com/superproj/zero/pkg/options"
)

// Injectors from wire.go:

// wireApp builds and returns a Kratos app with the given options.
// It uses the Wire library to automatically generate the dependency injection code.
func wireApp(appInfo bootstrap.AppInfo, config *server.Config, mySQLOptions *db.MySQLOptions, jwtOptions *options.JWTOptions, redisOptions *options.RedisOptions, etcdOptions *options.EtcdOptions, kafkaOptions *options.KafkaOptions) (*kratos.App, func(), error) {
	logger := bootstrap.NewLogger(appInfo)
	registrar := bootstrap.NewEtcdRegistrar(etcdOptions)
	appConfig := bootstrap.AppConfig{
		Info:      appInfo,
		Logger:    logger,
		Registrar: registrar,
	}
	gormDB, err := db.NewMySQL(mySQLOptions)
	if err != nil {
		return nil, nil, err
	}
	datastore := store.NewStore(gormDB)
	authenticator, cleanup, err := NewAuthenticator(jwtOptions, redisOptions)
	if err != nil {
		return nil, nil, err
	}
	secretGetter := newSecretGetter(datastore)
	authn, err := auth.NewAuthn(secretGetter)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	kafkaLogger, err := auth.NewLogger(kafkaOptions)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	authz, err := auth.NewAuthz(gormDB, redisOptions, kafkaLogger)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	authAuth := auth.NewAuth(authn, authz)
	bizBiz := biz.NewBiz(datastore, authenticator, authAuth)
	userCenterService := service.NewUserCenterService(bizBiz)
	customValidator, err := validation.New(datastore)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	validator := validation2.New(customValidator)
	v := server.NewMiddlewares(logger, authenticator, validator)
	httpServer := server.NewHTTPServer(config, userCenterService, authenticator, v)
	grpcServer := server.NewGRPCServer(config, userCenterService, v)
	v2 := server.NewServers(httpServer, grpcServer)
	app := bootstrap.NewApp(appConfig, v2...)
	return app, func() {
		cleanup()
	}, nil
}
