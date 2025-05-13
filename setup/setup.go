package main

import (
	"os"

	"github.com/goravel/framework/packages"
	"github.com/goravel/framework/packages/match"
	"github.com/goravel/framework/packages/modify"
	"github.com/goravel/framework/support/path"
)

var config = `map[string]any{
        "driver": "custom",
        "key":      config.Env("ALIYUN_ACCESS_KEY_ID"),
        "secret":   config.Env("ALIYUN_ACCESS_KEY_SECRET"),
        "bucket":   config.Env("ALIYUN_BUCKET"),
        "url":      config.Env("ALIYUN_URL"),
        "endpoint": config.Env("ALIYUN_ENDPOINT"),
        "via": func() (filesystem.Driver, error) {
            return ossfacades.Oss("oss") // The ` + "`oss`" + ` value is the ` + "`disks`" + ` key
        },
    }`

func main() {
	packages.Setup(os.Args).
		Install(
			modify.GoFile(path.Config("app.go")).
				Find(match.Imports()).Modify(modify.AddImport(packages.GetModulePath())).
				Find(match.Providers()).Modify(modify.Register("&oss.ServiceProvider{}")),
			modify.GoFile(path.Config("filesystems.go")).
				Find(match.Imports()).Modify(modify.AddImport("github.com/goravel/framework/contracts/filesystem"), modify.AddImport("github.com/goravel/oss/facades", "ossfacades")).
				Find(match.Config("filesystems.disks")).Modify(modify.AddConfig("oss", config)),
		).
		Uninstall(
			modify.GoFile(path.Config("app.go")).
				Find(match.Imports()).Modify(modify.RemoveImport(packages.GetModulePath())).
				Find(match.Providers()).Modify(modify.Unregister("&oss.ServiceProvider{}")),
			modify.GoFile(path.Config("filesystems.go")).
				Find(match.Config("filesystems.disks")).Modify(modify.RemoveConfig("oss")).
				Find(match.Imports()).Modify(modify.RemoveImport("github.com/goravel/framework/contracts/filesystem"), modify.RemoveImport("github.com/goravel/oss/facades", "ossfacades")),
		).
		Execute()
}
