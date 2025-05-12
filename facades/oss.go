package facades

import (
	"github.com/goravel/framework/contracts/filesystem"

	"github.com/goravel/oss"
)

func Oss(disk string) (filesystem.Driver, error) {
	instance, err := oss.App.MakeWith(oss.Binding, map[string]any{"disk": disk})
	if err != nil {
		return nil, err
	}

	return instance.(*oss.Oss), nil
}
