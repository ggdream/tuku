# TuKu

A picture object storage server implemented by Golang that can return preset sizes and custom sizes in real time.

Get Start It!!

```bash
go get -d github.com/ggdream/tuku/cmd/tuku
go install github.com/ggdream/tuku/cmd/tuku
```

```bash
./tuku ./config.yaml
```

# Docs

- [API Document](#api-document)
  1. [Upload Simple File](#upload-simple-file)
  2. [Upload Multiple Files](#upload-multiple-files)
  3. [Get Your File Basic](#get-your-file-basic)
  4. [Get Your File With Custom Size](#get-your-file-with-custom-size)
- [HTTP Response](#http-response)
  - [`code` type](#code-type)
- [Config File Document](#config-file-document)

# API Document

## Upload Simple File

```yaml
method: POST
url: http://127.0.0.1:9779/api/upload
params:
  - name: custom name # Optional
  - file: your file data
```

## Upload Multiple Files

```yaml
method: POST
url: http://127.0.0.1:9779/api/uploads
params:
  - names: custom name list # Optional
  - files: your files data
```

If you use this parameter, you must give each file a custom name!!

## Get Your File Basic

```yaml
method: GET
url: http://127.0.0.1:9779/api/image/{object} # object: your file name
```

## Get Your File With Custom Size

```yaml
method: GET
url: http://127.0.0.1:9779/api/image/{object}?width={w}&height={h} # object: your file name
```

width and height are optional, specify them and return an image of the specified size.

# HTTP Response

The response has three fields, `code` indicates the result code, `data` indicates the received response data, and `message` indicates the server processing result.

```json
{ "code": 0, "data": null, "message": "" }
```

```json
{ "code": 3, "data": null, "message": "upload file failed" }
```

## `code` type

| enum                       | code |
| -------------------------- | ---- |
| TypePerfect                | 0    |
| TypeFileOpenFailed         | 1    |
| TypeFileReadFailed         | 2    |
| TypeFileUploadFailed       | 3    |
| TypeFileCannotFind         | 4    |
| TypeImageDecodeFailed      | 5    |
| TypeImageResizeFailed      | 6    |
| TypeFileFormatNotSupported | 7    |
| TypeParamsMissingErr       | 8    |
| TypeParamsParsingErr       | 9    |
| TypeParamsInvalidErr       | 10   |

# Config File Document

[Click me to get config.yaml sample file](./cmd/tuku/config.yaml)

```yaml
image:
  raw: true # Keep original image
  preset: # Preset Size, {width}*{height}
    - 480*480
    - 1080*1080

http:
  host: 0.0.0.0
  port: 9779 # Service port
  cert: cert.pem
  key: key.pem
  tls: false # Whether to use SSL, and specify the file path in the cert and key parameters

storage:
  type: local # local or minio: local is refers to storing files locally, and minio is refers to storing files in MinIO
  path: ./gallery # Local storage path, enabled when `type` is `local`
  endpoint: 127.0.0.1:9000
  bucket: gallery
  accessKey: test
  secretKey: cQas7im3
  tls: false # Whether to use SSL to connect to MinIO service

dev: true # Whether to start the service in debug mode
```
