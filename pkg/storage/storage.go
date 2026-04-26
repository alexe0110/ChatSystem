package storage

import "context"

type FileStorage interface {
	Upload(ctx context.Context, fileName string, data []byte) (string, error)
}
