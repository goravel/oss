# oss

A oss disk driver for facades.Storage of Goravel.

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

3. Publish configuration file
dd
```
go run . artisan vendor:publish --package=github.com/goravel/oss
```

4. Fill your oss configuration to `config/oss.go` file

5. Add oss disk to `config/filesystems.go` file

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
        "via": func() (filesystem.Driver, error) {
            return ossfacades.Oss(), nil
        },
    },
}
```

## Testing

Run command below to run test(fill your owner oss configuration):

```
ALIYUN_ACCESS_KEY_ID= ALIYUN_ACCESS_KEY_SECRET= ALIYUN_BUCKET= ALIYUN_URL= ALIYUN_ENDPOINT= go test ./...
```
