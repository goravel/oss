# OSS

An OSS disk driver for `facades.Storage()` of Goravel.

## Version

| goravel/oss | goravel/framework |
|-------------|-------------------|
| v1.3.*      | v1.15.*           |
| v1.2.*      | v1.14.*           |
| v1.1.*      | v1.13.*           |
| v1.0.*      | v1.12.*           |

## Install

Run the command below in your project to install the package automatically:

```
./artisan package:install github.com/goravel/oss
```

Or check [the setup file](./setup/setup.go) to install the package manually.

## Testing

Run command below to run test:

```
ALIYUN_ACCESS_KEY_ID= ALIYUN_ACCESS_KEY_SECRET= ALIYUN_BUCKET= ALIYUN_URL= ALIYUN_ENDPOINT= go test ./...
```
