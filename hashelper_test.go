package hashelper_test

import (
	"context"
	"crypto/md5"
	"crypto/sha1"
	"fmt"
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

		r := strings.NewReader(str)
		h := md5.New()
		checksum, summed, err := hashelper.Sum(ctx, r, bufferSize, h, nil)
		if err != nil {
			log.Printf("Sum() error: %v", err)
			return
		}

		fmt.Printf("MD5 checksum of \"%v\": %X\n", str, checksum)
		fmt.Printf("Summed: %v\n", summed)
	}

	// Output:
	// MD5 checksum of "": D41D8CD98F00B204E9800998ECF8427E
	// Summed: 0
	// MD5 checksum of "Hello World!": ED076287532E86365E841E92BFC50D8C
	// Summed: 12
}

func ExampleSumString() {
	// Hash strings.
	strs := []string{
		"",
		"Hello World!",
	}

	for _, str := range strs {

		h := sha1.New()
		checksum, summed, err := hashelper.SumString(str, h)
		if err != nil {
			log.Printf("Sum() error: %v", err)
			return
		}

		fmt.Printf("SHA-1 checksum of \"%v\": %X\n", str, checksum)
		fmt.Printf("Summed: %v\n", summed)
	}

	// Output:
	// SHA-1 checksum of "": DA39A3EE5E6B4B0D3255BFEF95601890AFD80709
	// Summed: 0
	// SHA-1 checksum of "Hello World!": 2EF7BDE608CE5404E97D5F042F95F89F1C232871
	// Summed: 12
}
