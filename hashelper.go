package hashelper

import (
	"context"
	"crypto"
	"errors"
	"hash"
	"io"
)

type CallBack func(ctx context.Context, summed int64)

var (
	ErrNoHashFuncs = errors.New("no hash functions")
)

func sum(ctx context.Context, r io.Reader, bufferSize int64, cb CallBack, hashes ...hash.Hash) ([][]byte, int64, error) {
	var (
		summed    int64
		writers   []io.Writer
		checksums [][]byte
	)

	buf := make([]byte, bufferSize)

	for _, h := range hashes {
		writers = append(writers, h)
	}

	w := io.MultiWriter(writers...)

	for {
		select {
		case <-ctx.Done():
			return nil, summed, ctx.Err()
		default:
			n, err := io.CopyBuffer(w, r, buf)
			if err != nil {
				return nil, 0, err
			}

			if n == 0 {
				for _, h := range hashes {
					checksums = append(checksums, h.Sum(nil))
				}
				return checksums, summed, nil
			}

			summed += n

			if cb != nil {
				cb(ctx, summed)
			}
		}
	}
}

func Sum(ctx context.Context, r io.Reader, bufferSize int64, cb CallBack, hashFuncs ...crypto.Hash) ([][]byte, int64, error) {
	var hashes []hash.Hash

	for _, h := range hashFuncs {
		if !h.Available() {
			return nil, 0, ErrNoHashFuncs
		}
		hashes = append(hashes, h.New())
	}

	return sum(ctx, r, bufferSize, cb, hashes...)
}

func sumString(s string, hashes ...hash.Hash) ([][]byte, int, error) {
	var (
		writers   []io.Writer
		checksums [][]byte
	)

	for _, h := range hashes {
		writers = append(writers, h)
	}

	w := io.MultiWriter(writers...)

	n, err := io.WriteString(w, s)
	if err != nil {
		return nil, 0, err
	}

	for _, h := range hashes {
		checksums = append(checksums, h.Sum(nil))
	}
	return checksums, n, nil
}

func SumString(s string, hashFuncs ...crypto.Hash) ([][]byte, int, error) {
	var hashes []hash.Hash

	for _, h := range hashFuncs {
		if !h.Available() {
			return nil, 0, ErrNoHashFuncs
		}
		hashes = append(hashes, h.New())
	}

	return sumString(s, hashes...)
}
