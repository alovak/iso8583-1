// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package iso8583

import (
	"errors"
	"fmt"
	"strconv"
)

// A Llnumeric contains numeric value only in non-fix length, contains length in first 2 symbols. It holds numeric
// value as a string. Supportted encoder are ascii, bcd and rbcd. Length is
// required for marshaling and unmarshalling.
type Llnumeric struct {
	Value string
}

// NewLlnumeric create new Llnumeric field
func NewLlnumeric(val string) *Llnumeric {
	return &Llnumeric{val}
}

// IsEmpty check Llnumeric field for empty value
func (l *Llnumeric) IsEmpty() bool {
	return len(l.Value) == 0
}

// Bytes encode Llnumeric field to bytes
func (l *Llnumeric) Bytes(encoder, lenEncoder, length int) ([]byte, error) {
	raw := []byte(l.Value)
	if length != -1 && len(raw) > length {
		return nil, fmt.Errorf(ErrValueTooLong, "Llnumeric", length, len(raw))
	}

	val := raw
	switch encoder {
	case ASCII:
	case BCD:
		val = lbcd(raw)
	case rBCD:
		val = rbcd(raw)
	default:
		return nil, errors.New(ErrInvalidEncoder)
	}

	lenStr := fmt.Sprintf("%02d", len(raw)) // length of digital characters
	contentLen := []byte(lenStr)
	var lenVal []byte
	switch lenEncoder {
	case ASCII:
		lenVal = contentLen
		if len(lenVal) > 2 {
			return nil, errors.New(ErrInvalidLengthHead)
		}
	case rBCD:
		fallthrough
	case BCD:
		lenVal = rbcd(contentLen)
		if len(lenVal) > 1 || len(contentLen) > 3 {
			return nil, errors.New(ErrInvalidLengthHead)
		}
	default:
		return nil, errors.New(ErrInvalidLengthEncoder)
	}
	return append(lenVal, val...), nil
}

// Load decode Llnumeric field from bytes
func (l *Llnumeric) Load(raw []byte, encoder, lenEncoder, length int) (read int, err error) {
	// parse length head:
	var contentLen int
	switch lenEncoder {
	case ASCII:
		read = 2
		contentLen, err = strconv.Atoi(string(raw[:read]))
		if err != nil {
			return 0, errors.New(ErrParseLengthFailed + ": " + string(raw[:2]))
		}
	case rBCD:
		fallthrough
	case BCD:
		read = 1
		contentLen, err = strconv.Atoi(string(bcdr2Ascii(raw[:read], 2)))
		if err != nil {
			return 0, errors.New(ErrParseLengthFailed + ": " + string(raw[0]))
		}
	default:
		return 0, errors.New(ErrInvalidLengthEncoder)
	}

	// parse body:
	switch encoder {
	case ASCII:
		if len(raw) < (read + contentLen) {
			return 0, errors.New(ErrBadRaw)
		}
		l.Value = string(raw[read : read+contentLen])
		read += contentLen
	case rBCD:
		fallthrough
	case BCD:
		bcdLen := (contentLen + 1) / 2
		if len(raw) < (read + bcdLen) {
			return 0, errors.New(ErrBadRaw)
		}
		l.Value = string(bcdl2Ascii(raw[read:read+bcdLen], contentLen))
		read += bcdLen
	default:
		return 0, errors.New(ErrInvalidEncoder)
	}
	return read, nil
}
