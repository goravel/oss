package oss

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gookit/color"
	"github.com/goravel/framework/contracts/config"
	"github.com/goravel/framework/contracts/filesystem"
	contractshttp "github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/support/carbon"
	"github.com/goravel/framework/support/str"
)

/*
 * Oss OSS
 * Document: https://help.aliyun.com/document_detail/32144.html
 */

type Oss struct {
	ctx            context.Context
	config         config.Config
	disk           string
	instance       *oss.Client
	bucketInstance *oss.Bucket
	bucket         string
	url            string
	endpoint       string
}

func NewOss(ctx context.Context, config config.Config, disk string) (*Oss, error) {
	accessKeyId := config.GetString(fmt.Sprintf("filesystems.disks.%s.key", disk))
	accessKeySecret := config.GetString(fmt.Sprintf("filesystems.disks.%s.secret", disk))
	bucket := config.GetString(fmt.Sprintf("filesystems.disks.%s.bucket", disk))
	url := config.GetString(fmt.Sprintf("filesystems.disks.%s.url", disk))
	endpoint := config.GetString(fmt.Sprintf("filesystems.disks.%s.endpoint", disk))
	if accessKeyId == "" || accessKeySecret == "" || bucket == "" || url == "" || endpoint == "" {
		return nil, fmt.Errorf("please set %s configuration first", disk)
	}

	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("init %s disk error: %s", disk, err)
	}

	bucketInstance, err := client.Bucket(bucket)
	if err != nil {
		return nil, fmt.Errorf("init %s bucket of %s disk error: %s", bucket, disk, err)
	}

	return &Oss{
		ctx:            ctx,
		config:         config,
		disk:           disk,
		instance:       client,
		bucketInstance: bucketInstance,
		bucket:         bucket,
		url:            url,
		endpoint:       endpoint,
	}, nil
}

func (r *Oss) AllDirectories(path string) ([]string, error) {
	var directories []string
	validPath := validPath(path)
	lsRes, err := r.bucketInstance.ListObjectsV2(oss.MaxKeys(MaxFileNum), oss.Prefix(validPath), oss.Delimiter("/"))
	if err != nil {
		return nil, err
	}

	for _, commonPrefix := range lsRes.CommonPrefixes {
		directories = append(directories, strings.ReplaceAll(commonPrefix, validPath, ""))
		subDirectories, err := r.AllDirectories(commonPrefix)
		if err != nil {
			return nil, err
		}
		for _, subDirectory := range subDirectories {
			if strings.HasSuffix(subDirectory, "/") {
				directories = append(directories, strings.ReplaceAll(commonPrefix+subDirectory, validPath, ""))
			}
		}
	}

	return directories, nil
}

func (r *Oss) AllFiles(path string) ([]string, error) {
	var files []string
	validPath := validPath(path)
	lsRes, err := r.bucketInstance.ListObjectsV2(oss.MaxKeys(MaxFileNum), oss.Prefix(validPath))
	if err != nil {
		return nil, err
	}
	for _, object := range lsRes.Objects {
		if !strings.HasSuffix(object.Key, "/") {
			files = append(files, strings.ReplaceAll(object.Key, validPath, ""))
		}
	}

	return files, nil
}

func (r *Oss) Copy(originFile, targetFile string) error {
	if _, err := r.bucketInstance.CopyObject(originFile, targetFile); err != nil {
		return err
	}

	return nil
}

func (r *Oss) Delete(files ...string) error {
	_, err := r.bucketInstance.DeleteObjects(files)
	if err != nil {
		return err
	}

	return nil
}

func (r *Oss) DeleteDirectory(directory string) error {
	if !strings.HasSuffix(directory, "/") {
		directory += "/"
	}

	marker := oss.Marker("")
	prefix := oss.Prefix(directory)
	for {
		lor, err := r.bucketInstance.ListObjects(marker, prefix)
		if err != nil {
			return err
		}
		if len(lor.Objects) == 0 {
			return nil
		}

		var objects []string
		for _, object := range lor.Objects {
			objects = append(objects, object.Key)
		}

		if _, err := r.bucketInstance.DeleteObjects(objects, oss.DeleteObjectsQuiet(true)); err != nil {
			return err
		}

		prefix = oss.Prefix(lor.Prefix)
		marker = oss.Marker(lor.NextMarker)
		if !lor.IsTruncated {
			break
		}
	}

	return nil
}

func (r *Oss) Directories(path string) ([]string, error) {
	var directories []string
	validPath := validPath(path)
	lsRes, err := r.bucketInstance.ListObjectsV2(oss.MaxKeys(MaxFileNum), oss.Prefix(validPath), oss.Delimiter("/"))
	if err != nil {
		return nil, err
	}

	for _, commonPrefix := range lsRes.CommonPrefixes {
		directory := strings.ReplaceAll(commonPrefix, validPath, "")
		if directory != "" {
			directories = append(directories, directory)
		}
	}

	return directories, nil
}

func (r *Oss) Exists(file string) bool {
	exist, err := r.bucketInstance.IsObjectExist(file)
	if err != nil {
		return false
	}

	return exist
}

func (r *Oss) Files(path string) ([]string, error) {
	var files []string
	validPath := validPath(path)
	lsRes, err := r.bucketInstance.ListObjectsV2(oss.MaxKeys(MaxFileNum), oss.Prefix(validPath), oss.Delimiter("/"))
	if err != nil {
		return nil, err
	}
	for _, object := range lsRes.Objects {
		file := strings.ReplaceAll(object.Key, validPath, "")
		if file != "" {
			files = append(files, file)
		}
	}

	return files, nil
}

func (r *Oss) Get(file string) (string, error) {
	data, err := r.GetBytes(file)

	return string(data), err
}

func (r *Oss) GetBytes(file string) ([]byte, error) {
	res, err := r.bucketInstance.GetObject(file)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	data, err := io.ReadAll(res)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (r *Oss) LastModified(file string) (time.Time, error) {
	headers, err := r.bucketInstance.GetObjectDetailedMeta(file)
	if err != nil {
		return time.Time{}, err
	}

	lastModified, err := http.ParseTime(headers.Get("Last-Modified"))
	if err != nil {
		return time.Time{}, err
	}

	l, err := time.LoadLocation(r.config.GetString("app.timezone"))
	if err != nil {
		return time.Time{}, err
	}

	return lastModified.In(l), nil
}

func (r *Oss) MakeDirectory(directory string) error {
	if !strings.HasSuffix(directory, "/") {
		directory += "/"
	}

	return r.bucketInstance.PutObject(directory, bytes.NewReader([]byte("")))
}

func (r *Oss) MimeType(file string) (string, error) {
	headers, err := r.bucketInstance.GetObjectDetailedMeta(file)
	if err != nil {
		return "", err
	}

	return headers.Get("Content-Type"), nil
}

func (r *Oss) Missing(file string) bool {
	return !r.Exists(file)
}

func (r *Oss) Move(oldFile, newFile string) error {
	if err := r.Copy(oldFile, newFile); err != nil {
		return err
	}

	return r.Delete(oldFile)
}

func (r *Oss) Path(file string) string {
	return file
}

func (r *Oss) Put(file string, content string) error {
	// If the file is created in a folder directly, we can't check if the folder exists.
	// So we need to create the top folder first.
	if err := r.makeDirectories(file); err != nil {
		return err
	}

	tempFile, err := r.tempFile(content)
	defer os.Remove(tempFile.Name())
	if err != nil {
		return err
	}

	return r.bucketInstance.PutObjectFromFile(file, tempFile.Name())
}

func (r *Oss) PutFile(filePath string, source filesystem.File) (string, error) {
	return r.PutFileAs(filePath, source, str.Random(40))
}

func (r *Oss) PutFileAs(filePath string, source filesystem.File, name string) (string, error) {
	fullPath, err := fullPathOfFile(filePath, source, name)
	if err != nil {
		return "", err
	}

	// If the file is created in a folder directly, we can't check if the folder exists.
	// So we need to create the top folder first.
	if err := r.makeDirectories(str.Of(filePath).Finish("/").String()); err != nil {
		return "", err
	}

	if err := r.bucketInstance.PutObjectFromFile(fullPath, source.File()); err != nil {
		return "", err
	}

	return fullPath, nil
}

func (r *Oss) Size(file string) (int64, error) {
	props, err := r.bucketInstance.GetObjectDetailedMeta(file)
	if err != nil {
		return 0, err
	}

	lens := props["Content-Length"]
	if len(lens) == 0 {
		return 0, nil
	}

	contentLengthInt, err := strconv.ParseInt(lens[0], 10, 64)
	if err != nil {
		return 0, err
	}

	return contentLengthInt, nil
}

func (r *Oss) TemporaryUrl(file string, t time.Time) (string, error) {
	signedURL, err := r.bucketInstance.SignURL(file, oss.HTTPGet, int64(t.Sub(carbon.Now().StdTime()).Seconds()))
	if err != nil {
		return "", err
	}

	return signedURL, nil
}

func (r *Oss) WithContext(ctx context.Context) filesystem.Driver {
	if httpCtx, ok := ctx.(contractshttp.Context); ok {
		ctx = httpCtx.Context()
	}

	driver, err := NewOss(ctx, r.config, r.disk)
	if err != nil {
		color.Redf("init %s error: %+v\n", r.disk, err)

		return nil
	}

	return driver
}

func (r *Oss) Url(file string) string {
	return r.url + "/" + file
}

func (r *Oss) makeDirectories(path string) error {
	folders := strings.Split(path, "/")
	for i := 1; i < len(folders); i++ {
		folder := strings.Join(folders[:i], "/")
		if err := r.MakeDirectory(folder); err != nil {
			return err
		}
	}

	return nil
}

func (r *Oss) tempFile(content string) (*os.File, error) {
	tempFile, err := os.CreateTemp(os.TempDir(), "goravel-")
	if err != nil {
		return nil, err
	}

	if _, err := tempFile.WriteString(content); err != nil {
		return nil, err
	}

	return tempFile, nil
}
