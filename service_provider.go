package oss

import (
	"context"

	"github.com/goravel/framework/contracts/foundation"
)

const Binding = "goravel.oss"

var App foundation.Application

type ServiceProvider struct {
}

func (receiver *ServiceProvider) Register(app foundation.Application) {
	App = app

	app.Bind(Binding, func(app foundation.Application) (any, error) {
		return NewOss(context.Background(), app.MakeConfig())
	})
}

func (receiver *ServiceProvider) Boot(app foundation.Application) {
	app.Publishes("github.com/goravel/oss", map[string]string{
		"config/oss.go": app.ConfigPath(""),
	})
}
