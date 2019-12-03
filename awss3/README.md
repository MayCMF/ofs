# AWS S3

[AWS S3](https://aws.amazon.com/s3/) backend for [MayCMF OSS](https://github.com/MayCMF/ofs)

## Usage

```go
import "github.com/MayCMF/oss/awss3"

func main() {
  storage := s3.New(s3.Config{
    AccessID: "access_id",
    AccessKey: "access_key",
    Region: "region",
    Bucket: "bucket",
    Endpoint: "cdn.example.com",
    ACL: awss3.BucketCannedACLPublicRead,
  })

  // Save a reader interface into storage
  storage.Put("/example.txt", reader)

  // Get file with path
  storage.Get("/example.txt")

  // Get object as io.ReadCloser
  storage.GetStream("/example.txt")

  // Delete file with path
  storage.Delete("/example.txt")

  // List all objects under path
  storage.List("/")

  // Get Public Accessible URL (useful if current file saved privately)
  storage.GetURL("/example.txt")
}
```


