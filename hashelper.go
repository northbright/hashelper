package hashelper

import (
	"context"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"hash"
	"hash/crc32"
	"io"
)

const (
	DefBufSize = int64(32 * 1024)
	MaxBufSize = int64(4 * 1024 * 1024 * 1024)
)

var (
	ErrUnSupportedHashFunc = errors.New("unsupported hash function")
	SupportedHashFuncs     = []string{
		"MD5",
		"CRC-32",
		"SHA-1",
		"SHA-256",
		"SHA-512",
	}
)

func GetSupportedHashFuncNames() []string {
	return SupportedHashFuncs
}

func GetHashByName(name string) (hash.Hash, error) {

	switch name {
	case "MD5":
		return md5.New(), nil
	case "CRC-32":
		return crc32.NewIEEE(), nil
	case "SHA-1":
		return sha1.New(), nil
	case "SHA-256":
		return sha256.New(), nil
	case "SHA-512":
		return sha512.New(), nil
	default:
		return nil, ErrUnSupportedHashFunc
	}
}

func updateBufferSize(bufferSize int64) int64 {
	switch {
	case bufferSize <= 0:
		return DefBufSize
	case bufferSize > MaxBufSize:
		return MaxBufSize
	default:
		return bufferSize
	}
}

func SumCtx(
	ctx context.Context,
	r io.Reader,
	hashFuncNames []string,
	bufferSize int64,
) (checksums [][]byte, summed int64, err error) {
	var (
		hashes  []hash.Hash
		writers []io.Writer
	)

	bufferSize = updateBufferSize(bufferSize)
	buf := make([]byte, bufferSize)

	// Get hash.Hash from hash func name.
	for _, name := range hashFuncNames {
		hash, err := GetHashByName(name)
		if err != nil {
			return nil, 0, err
		}
		hashes = append(hashes, hash)
	}

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
		}
	}
}

func SumString(s string, hashFuncNames []string) (checksums [][]byte, summed int, err error) {
	var (
		hashes  []hash.Hash
		writers []io.Writer
	)

	// Get hash.Hash from hash func name.
	for _, name := range hashFuncNames {
		hash, err := GetHashByName(name)
		if err != nil {
			return nil, 0, err
		}
		hashes = append(hashes, hash)
	}

	for _, h := range hashes {
		writers = append(writers, h)
	}

	w := io.MultiWriter(writers...)

	summed, err = io.WriteString(w, s)
	if err != nil {
		return nil, 0, err
	}

	for _, h := range hashes {
		checksums = append(checksums, h.Sum(nil))
	}
	return checksums, summed, nil
}
