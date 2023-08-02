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

	app.BindWith(Binding, func(app foundation.Application, parameter map[string]any) (any, error) {
		return NewOss(context.Background(), app.MakeConfig(), parameter["disk"].(string))
	})
}

func (receiver *ServiceProvider) Boot(app foundation.Application) {

}
