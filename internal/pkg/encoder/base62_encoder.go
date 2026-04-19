package encoder

import (
	"fmt"
	"log/slog"
)

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

type Encoder interface {
	Encode(numbers []int64) ([]string, error)
}

type encoder struct {
}

func NewBase62Encoder() Encoder {
	return &encoder{}
}

func (e *encoder) Encode(numbers []int64) ([]string, error) {
	hashes := make([]string, 0, len(numbers))
	base := int64(len(base62Chars))

	for _, number := range numbers {
		if number < 0 {
			slog.Error("passed negative number in base62 encoder")
			return nil, fmt.Errorf("number must be non-negative: %d", number)
		}

		if number == 0 {
			hashes = append(hashes, string(base62Chars[0]))
			continue
		}

		var res []byte
		n := number

		for n > 0 {
			remainder := n % base
			res = append(res, base62Chars[remainder])
			n /= base
		}

		for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
			res[i], res[j] = res[j], res[i]
		}

		hashes = append(hashes, string(res))
	}

	return hashes, nil
}
