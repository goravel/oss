package main

import (
	"os"

	"github.com/goravel/framework/packages"
	"github.com/goravel/framework/packages/match"
	"github.com/goravel/framework/packages/modify"
	"github.com/goravel/framework/support/env"
	"github.com/goravel/framework/support/path"
)

func main() {
	setup := packages.Setup(os.Args)
	config := `map[string]any{
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

	appConfigPath := path.Config("app.go")
	filesystemsConfigPath := path.Config("filesystems.go")
	moduleImport := setup.Paths().Module().Import()
	ossServiceProvider := "&oss.ServiceProvider{}"
	filesystemContract := "github.com/goravel/framework/contracts/filesystem"
	ossFacades := "github.com/goravel/oss/facades"
	filesystemsDisksConfig := match.Config("filesystems.disks")
	filesystemsConfig := match.Config("filesystems")

	setup.Install(
		// Add oss service provider to app.go if not using bootstrap setup
		modify.When(func(_ map[string]any) bool {
			return !env.IsBootstrapSetup()
		}, modify.GoFile(appConfigPath).
			Find(match.Imports()).Modify(modify.AddImport(moduleImport)).
			Find(match.Providers()).Modify(modify.Register(ossServiceProvider))),

		// Add oss service provider to providers.go if using bootstrap setup
		modify.When(func(_ map[string]any) bool {
			return env.IsBootstrapSetup()
		}, modify.RegisterProvider(moduleImport, ossServiceProvider)),

		// Add oss disk to filesystems.go
		modify.GoFile(filesystemsConfigPath).Find(match.Imports()).Modify(
			modify.AddImport(filesystemContract),
			modify.AddImport(ossFacades, "ossfacades"),
		).
			Find(filesystemsDisksConfig).Modify(modify.AddConfig("oss", config)).
			Find(filesystemsConfig).Modify(modify.AddConfig("default", `"oss"`)),
	).Uninstall(
		// Remove oss disk from filesystems.go
		modify.WhenFileExists(filesystemsConfigPath, modify.GoFile(filesystemsConfigPath).
			Find(filesystemsConfig).Modify(modify.AddConfig("default", `"local"`)).
			Find(filesystemsDisksConfig).Modify(modify.RemoveConfig("oss")).
			Find(match.Imports()).Modify(
			modify.RemoveImport(filesystemContract),
			modify.RemoveImport(ossFacades, "ossfacades"),
		)),

		// Remove oss service provider from app.go if not using bootstrap setup
		modify.When(func(_ map[string]any) bool {
			return !env.IsBootstrapSetup()
		}, modify.GoFile(appConfigPath).
			Find(match.Providers()).Modify(modify.Unregister(ossServiceProvider)).
			Find(match.Imports()).Modify(modify.RemoveImport(moduleImport))),

		// Remove oss service provider from providers.go if using bootstrap setup
		modify.When(func(_ map[string]any) bool {
			return env.IsBootstrapSetup()
		}, modify.UnregisterProvider(moduleImport, ossServiceProvider)),
	).Execute()
}
