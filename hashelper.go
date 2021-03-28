package hashelper

import (
	"context"
	"hash"
	"io"
)

func Sum(ctx context.Context, r io.Reader, bufferSize int64, h hash.Hash) ([]byte, int64, error) {
	var summed int64
	buf := make([]byte, bufferSize)

	for {
		select {
		case <-ctx.Done():
			return h.Sum(nil), summed, ctx.Err()
		default:
			n, err := io.CopyBuffer(h, r, buf)
			if err != nil {
				return nil, 0, err
			}

			if n == 0 {
				return h.Sum(nil), summed, nil
			}

			summed += n
		}
	}
}
