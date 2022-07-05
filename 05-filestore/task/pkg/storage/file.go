package storage

import (
	"context"
	"fmt"
	"io"

	"github.com/google/uuid"
	"github.com/minio/minio-go"
	"github.com/pkg/errors"
)

const (
	bucketName = "files"

	// metadata key of file (minio)
	ownerIdMdKey   = "Filestore-Owner-Id"
	encryptedMdKey = "Filestore-Encrypted"
	mdKeyPrefix    = "X-Amz-Meta"
)

type FileOptions struct {
	Name      string
	Size      int64
	OwnerId   uuid.UUID
	Encrypted bool
}

type FileStat struct {
	OwnerId   uuid.UUID
	Encrypted bool
}

type FileStorage interface {
	Update(ctx context.Context, file io.Reader, opts *FileOptions) error
	List() []string
	Stat(name string) (*FileStat, error)
	Get(ctx context.Context, name string) ([]byte, error)
}

type store struct {
	minioClient *minio.Client
}

func NewFileStorage(minioClient *minio.Client) (FileStorage, error) {
	store := &store{
		minioClient: minioClient,
	}

	exists, err := minioClient.BucketExists(bucketName)
	if err != nil {
		return nil, fmt.Errorf("check '%v' bucket existence failed: %v", bucketName, err)
	}
	if exists {
		return store, nil
	}

	if err := minioClient.MakeBucket(bucketName, "us-east-1"); err != nil {
		return nil, fmt.Errorf("make '%v' bucket failed: %v", bucketName, err)
	}
	return store, nil
}

func (s *store) Update(ctx context.Context, file io.Reader, opts *FileOptions) error {
	n, err := s.minioClient.PutObjectWithContext(ctx, bucketName, opts.Name, file, opts.Size, minio.PutObjectOptions{
		UserMetadata: map[string]string{
			ownerIdMdKey:   opts.OwnerId.String(),
			encryptedMdKey: fmt.Sprint(opts.Encrypted),
		},
	})
	if err != nil {
		return fmt.Errorf("put object to minio failed: %v", err)
	}
	if n != opts.Size {
		return fmt.Errorf("object wasn't completely put to minio, want: %v, put: %v", opts.Size, n)
	}

	return nil
}

func (s *store) List() []string {
	filesInfo := make([]string, 0)

	doneCh := make(chan struct{})
	defer close(doneCh)
	for obj := range s.minioClient.ListObjects(bucketName, "", false, doneCh) {
		filesInfo = append(filesInfo, obj.Key)
	}

	return filesInfo
}

func (s *store) Stat(name string) (*FileStat, error) {
	info, err := s.minioClient.StatObject(bucketName, name, minio.StatObjectOptions{})
	if err != nil {
		errResp := minio.ToErrorResponse(err)
		if errResp.Code == "NoSuchKey" {
			return nil, ErrNotExist
		}

		return nil, fmt.Errorf("get object failed: %v", err)
	}

	ownerKey := fmt.Sprintf("%v-%v", mdKeyPrefix, ownerIdMdKey)
	ownerId := uuid.MustParse(info.Metadata.Get(ownerKey))

	encrypted := false
	encryptedKey := fmt.Sprintf("%v-%v", mdKeyPrefix, encryptedMdKey)
	encryptedValue := info.Metadata.Get(encryptedKey)
	switch encryptedValue {
	case "true":
		encrypted = true
	case "false":
		encrypted = false
	default:
		return nil, errors.New("flag \"encrypted\" of user metadata must be true or false")
	}

	return &FileStat{
		OwnerId:   ownerId,
		Encrypted: encrypted,
	}, nil
}

func (s *store) Get(ctx context.Context, name string) ([]byte, error) {
	obj, err := s.minioClient.GetObjectWithContext(ctx, bucketName, name, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("get object failed: %v", err)
	}

	content, err := io.ReadAll(obj)
	if err != nil {
		return nil, fmt.Errorf("read object failed: %v", err)
	}

	return content, nil
}
