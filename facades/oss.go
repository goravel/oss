package facades

import (
	"log"

	"github.com/goravel/framework/contracts/filesystem"

	"github.com/goravel/oss"
)

func Oss() filesystem.Driver {
	instance, err := oss.App.Make(oss.Binding)
	if err != nil {
		log.Fatalln(err)
		return nil
	}

	return instance.(*oss.Oss)
}
