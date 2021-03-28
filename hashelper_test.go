package hashelper_test

import (
	"context"
	"crypto/md5"
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
		checksum, summed, err := hashelper.Sum(ctx, r, bufferSize, h)
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
