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

var (
	defBufSize             = int64(32 * 1024)
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

func SumCtx(
	ctx context.Context,
	r io.Reader,
	bufferSize int64,
	hashFuncNames []string,
) ([][]byte, int64, error) {
	var (
		summed    int64
		hashes    []hash.Hash
		writers   []io.Writer
		checksums [][]byte
	)

	if bufferSize <= 0 {
		bufferSize = defBufSize
	}

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

func SumString(s string, hashFuncNames []string) ([][]byte, int, error) {
	var (
		hashes    []hash.Hash
		writers   []io.Writer
		checksums [][]byte
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

	n, err := io.WriteString(w, s)
	if err != nil {
		return nil, 0, err
	}

	for _, h := range hashes {
		checksums = append(checksums, h.Sum(nil))
	}
	return checksums, n, nil
}
