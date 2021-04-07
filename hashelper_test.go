package hashelper_test

import (
	"context"
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"hash"
	"hash/crc32"
	"log"
	"strings"

	"github.com/northbright/hashelper"
)

func ExampleSum() {
	var bufferSize int64 = 1

	// Hash strings.
	strs := []string{
		"",
		"Hello World!",
	}

	for _, str := range strs {
		ctx := context.Background()
		hashes := []hash.Hash{
			md5.New(),
			sha1.New(),
			crc32.NewIEEE(),
		}

		hashFuncNames := []string{
			"MD5",
			"SHA-1",
			"CRC32",
		}

		r := strings.NewReader(str)
		checksums, summed, err := hashelper.Sum(ctx, r, bufferSize, nil, hashes...)
		if err != nil {
			log.Printf("Sum() error: %v", err)
			return
		}

		fmt.Printf("Summed: %v bytes for \"%v\"\n", summed, str)
		for i, checksum := range checksums {
			fmt.Printf("%v: %X\n", hashFuncNames[i], checksum)
		}
	}

	// Output:
	// Summed: 0 bytes for ""
	// MD5: D41D8CD98F00B204E9800998ECF8427E
	// SHA-1: DA39A3EE5E6B4B0D3255BFEF95601890AFD80709
	// CRC32: 00000000
	// Summed: 12 bytes for "Hello World!"
	// MD5: ED076287532E86365E841E92BFC50D8C
	// SHA-1: 2EF7BDE608CE5404E97D5F042F95F89F1C232871
	// CRC32: 1C291CA3
}

func ExampleSumString() {
	// Hash strings.
	strs := []string{
		"",
		"Hello World!",
	}

	hashes := []hash.Hash{
		md5.New(),
		sha1.New(),
		crc32.NewIEEE(),
	}

	hashFuncNames := []string{
		"MD5",
		"SHA-1",
		"CRC32",
	}

	for _, str := range strs {
		checksums, summed, err := hashelper.SumString(str, hashes...)
		if err != nil {
			log.Printf("Sum() error: %v", err)
			return
		}

		fmt.Printf("Summed: %v bytes for \"%v\"\n", summed, str)
		for i, checksum := range checksums {
			fmt.Printf("%v: %X\n", hashFuncNames[i], checksum)
		}
	}

	// Output:
	// Summed: 0 bytes for ""
	// MD5: D41D8CD98F00B204E9800998ECF8427E
	// SHA-1: DA39A3EE5E6B4B0D3255BFEF95601890AFD80709
	// CRC32: 00000000
	// Summed: 12 bytes for "Hello World!"
	// MD5: ED076287532E86365E841E92BFC50D8C
	// SHA-1: 2EF7BDE608CE5404E97D5F042F95F89F1C232871
	// CRC32: 1C291CA3
}
