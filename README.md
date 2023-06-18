# OSS

An OSS disk driver for facades.Storage of Goravel.

## Version

| goravel/oss  | goravel/framework    |
| ----------   | --------------       |
| v1.0.0       | v1.12.0             |

## Install

1. Add package

```
go get -u github.com/goravel/oss
```

2. Register service provider

```
// config/app.go
import "github.com/goravel/oss"

"providers": []foundation.ServiceProvider{
    ...
    &oss.ServiceProvider{},
}
```

3. Add oss disk to `config/filesystems.go` file

```
// config/filesystems.go
import (
    "github.com/goravel/framework/contracts/filesystem"
    ossfacades "github.com/goravel/oss/facades"
)

"disks": map[string]any{
    ...
    "oss": map[string]any{
        "driver": "custom",
        "key":      config.Env("ALIYUN_ACCESS_KEY_ID"),
        "secret":   config.Env("ALIYUN_ACCESS_KEY_SECRET"),
        "bucket":   config.Env("ALIYUN_BUCKET"),
        "url":      config.Env("ALIYUN_URL"),
        "endpoint": config.Env("ALIYUN_ENDPOINT"),
        "via": func() (filesystem.Driver, error) {
            return ossfacades.Oss("oss"), nil // The `oss` value is the `disks` key
        },
    },
}
```

## Testing

Run command below to run test(fill your owner oss configuration):

```
ALIYUN_ACCESS_KEY_ID= ALIYUN_ACCESS_KEY_SECRET= ALIYUN_BUCKET= ALIYUN_URL= ALIYUN_ENDPOINT= go test ./...
```
