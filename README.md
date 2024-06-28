# corekit

```go
	// CreateBucket creates a bucket.
	CreateBucket(ctx context.Context, bucket, region string) error
	// DeleteBucket delete a bucket.
	DeleteBucket(ctx context.Context, bucket string) error
	// ListBuckets list all buckets owned by this authenticated user.
	ListBuckets(ctx context.Context) ([]Bucket, error)
	// SetBucketPolicy sets the bucket policy.
	SetBucketPolicy(ctx context.Context, bucket, policy string) error
	// URL returns the url.
	URL(ctx context.Context, metadata Metadata) (string, error)
	// PutObject puts the object to the oss.
	PutObject(ctx context.Context, obj *Object) (*Object, error)
	// PutObjectReader puts the object reader to the oss.
	PutObjectReader(ctx context.Context, obj *ObjectReader) (*ObjectReader, error)
	// GetObject get the object from the oss.
	GetObject(ctx context.Context, metadata Metadata) (*Object, error)
	// GetObjectReader get the object reader from oss
	GetObjectReader(ctx context.Context, metadata Metadata) (*ObjectReader, error)
	// DeleteObject deletes the object.
	DeleteObject(ctx context.Context, metadata Metadata) error
	// ListObject list all object owned by this authenticated user.
	ListObject(ctx context.Context, metadata Metadata) ([]Object, error)
	// Clear clears the object.
	Clear(ctx context.Context, metadata Metadata) error
	// Download downloads the object.
	Download(ctx context.Context, metadata Metadata) ([]byte, error)
	// Size returns the size.
	Size(ctx context.Context, metadata Metadata) (int64, error)
```
