package hashelper

import (
	"context"
	"hash"
	"io"
)

type CallBack func(ctx context.Context, summed int64)

func Sum(ctx context.Context, r io.Reader, bufferSize int64, h hash.Hash, f CallBack) ([]byte, int64, error) {
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

			if f != nil {
				f(ctx, summed)
			}
		}
	}
}

func SumString(s string, h hash.Hash) ([]byte, int, error) {
	n, err := io.WriteString(h, s)
	if err != nil {
		return nil, 0, err
	}

	checksum := h.Sum(nil)
	return checksum, n, nil
}
