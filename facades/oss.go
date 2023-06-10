package facades

import (
	"log"

	"github.com/goravel/framework/contracts/filesystem"

	"github.com/goravel/oss"
)

func Oss(disk string) filesystem.Driver {
	instance, err := oss.App.MakeWith(oss.Binding, map[string]any{"disk": disk})
	if err != nil {
		log.Fatalln(err)
		return nil
	}

	return instance.(*oss.Oss)
}
