package oss

import (
	"context"

	"github.com/goravel/framework/contracts/binding"
	"github.com/goravel/framework/contracts/foundation"
)

const Binding = "goravel.oss"

var App foundation.Application

type ServiceProvider struct {
}

func (r *ServiceProvider) Relationship() binding.Relationship {
	return binding.Relationship{
		Bindings: []string{
			Binding,
		},
		Dependencies: []string{
			binding.Config,
		},
		ProvideFor: []string{
			binding.Storage,
		},
	}
}

func (r *ServiceProvider) Register(app foundation.Application) {
	App = app

	app.BindWith(Binding, func(app foundation.Application, parameter map[string]any) (any, error) {
		return NewOss(context.Background(), app.MakeConfig(), parameter["disk"].(string))
	})
}

func (r *ServiceProvider) Boot(app foundation.Application) {

}
