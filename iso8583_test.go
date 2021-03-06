// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package iso8583

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestISO struct {
	F2  *Llnumeric `field:"2" length:"19"`
	F3  *Numeric   `field:"3" length:"6"`
	F4  *Numeric   `field:"4" length:"12"`
	F7  *Numeric   `field:"7" length:"10"`
	F11 *Numeric   `field:"11" length:"6"`
	F12 *Numeric   `field:"12" length:"6"`
	F13 *Numeric   `field:"13" length:"4"`
	F14 *Numeric   `field:"14" length:"4"`
	// BCD encoding with right-aligned value with odd length (for ex. "643" as [6 67] == "0643")
	F19  *Numeric      `field:"19" length:"3" encode:"rbcd"`
	F22  *Numeric      `field:"22" length:"3"`
	F25  *Numeric      `field:"25" length:"2"`
	F32  *Llnumeric    `field:"32" length:"11"`
	F35  *Llnumeric    `field:"35" length:"37"`
	F37  *Alphanumeric `field:"37" length:"12"`
	F39  *Alphanumeric `field:"39" length:"2"`
	F41  *Alphanumeric `field:"41" length:"8"`
	F42  *Alphanumeric `field:"42" length:"15"`
	F43  *Alphanumeric `field:"43" length:"40"`
	F49  *Numeric      `field:"49" length:"3" encode:"bcd"`
	F52  *Binary       `field:"52" length:"8"`
	F53  *Numeric      `field:"53" length:"16"`
	F120 *Lllnumeric   `field:"120" length:"999"`
}

type TestISO2 struct {
	F2  *Llnumeric    `field:"2" length:"19" encode:"bcd,rbcd"`
	F3  *Numeric      `field:"3" length:"6" encode:"bcd"`
	F4  *Numeric      `field:"4" length:"12" encode:"ascii"`
	F7  *Numeric      `field:"7" length:"10" encode:"bcd"`
	F11 *Numeric      `field:"11" length:"6" encode:"bcd"`
	F12 *Numeric      `field:"12" length:"6" encode:"bcd"`
	F13 *Numeric      `field:"13" length:"4" encode:"lbcd"`
	F14 *Numeric      `field:"14" length:"4" encode:"lbcd"`
	F19 *Numeric      `field:"19" length:"3" encode:"rbcd"`
	F22 *Numeric      `field:"22" length:"3" encode:"rbcd"`
	F25 *Numeric      `field:"25" length:"2" encode:"bcd"`
	F26 *Numeric      `field:"26" length:"2" encode:"bcd"`
	F28 *Alphanumeric `field:"28" length:"9"`
	F32 *Llnumeric    `field:"32" length:"11" encode:"bcd,rbcd"`
	F35 *Llnumeric    `field:"35" length:"37" encode:"rbcd,ascii"`
	F37 *Alphanumeric `field:"37" length:"12"`
	F39 *Alphanumeric `field:"39" length:"2"`
	F41 *Alphanumeric `field:"41" length:"8"`
	F42 *Alphanumeric `field:"42" length:"15"`
	F43 *Alphanumeric `field:"43" length:"40"`
	F45 *Llnumeric    `field:"45" length:"75" encode:"ascii,bcd"`
	F49 *Numeric      `field:"49" length:"3" encode:"rbcd"`
	F52 *Binary       `field:"52" length:"8"`
	F53 *Numeric      `field:"53" length:"16" encode:"bcd"`
	F54 *Llvar        `field:"54" length:"255" encode:"ascii,ascii"`
	F55 *Llvar        `field:"55" length:"255" encode:"bcd,ascii"`
	F56 *Lllvar       `field:"56" length:"255" encode:"bcd,ascii"`
	F57 *Lllvar       `field:"57" length:"255" encode:"rbcd,ascii"`
	F58 *Lllvar       `field:"58" length:"255" encode:"ascii,ascii"`
	F59 *Llvar        `field:"59" length:"255" encode:"rbcd,ascii"`
	F60 *Lllnumeric   `field:"60" length:"999" encode:"bcd,ascii"`
	F61 *Lllnumeric   `field:"60" length:"999" encode:"bcd,rbcd"`
	F63 *Lllnumeric   `field:"63" length:"999" encode:"rbcd,bcd"`
	F64 *Binary       `field:"64" length:"32"`
}

func TestEncode(t *testing.T) {
	data := &TestISO{
		F2:   NewLlnumeric("4276555555555555"),
		F3:   NewNumeric("000000"),
		F4:   NewNumeric("000000077700"),
		F7:   NewNumeric("0701111844"),
		F11:  NewNumeric("000123"),
		F12:  NewNumeric("131844"),
		F13:  NewNumeric("0701"),
		F14:  NewNumeric("1902"),
		F19:  NewNumeric("643"),
		F22:  NewNumeric("901"),
		F25:  NewNumeric("02"),
		F32:  NewLlnumeric("123456"),
		F35:  NewLlnumeric("4276555555555555=12345678901234567890"),
		F37:  NewAlphanumeric("987654321001"),
		F39:  NewAlphanumeric(""),
		F41:  NewAlphanumeric("00000321"),
		F42:  NewAlphanumeric("120000000000034"),
		F43:  NewAlphanumeric("Test text"),
		F49:  NewNumeric("643"),
		F52:  NewBinary([]byte{1, 2, 3, 4, 5, 6, 7, 8}),
		F53:  NewNumeric("1234000000000000"),
		F120: NewLllnumeric("Another test text"),
	}

	iso := Message{"0100", ASCII, true, false, data}

	res, err := iso.Bytes()

	if err != nil {
		t.Error("ISO Encode error:", err)
	}

	sample := []byte{48, 49, 48, 48, 242, 60, 36, 129, 40, 224, 152, 0, 0, 0, 0, 0, 0, 0, 1, 0, 49, 54, 52, 50, 55, 54, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 55, 55, 55, 48, 48, 48, 55, 48, 49, 49, 49, 49, 56, 52, 52, 48, 48, 48, 49, 50, 51, 49, 51, 49, 56, 52, 52, 48, 55, 48, 49, 49, 57, 48, 50, 6, 67, 57, 48, 49, 48, 50, 48, 54, 49, 50, 51, 52, 53, 54, 51, 55, 52, 50, 55, 54, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 61, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 57, 56, 55, 54, 53, 52, 51, 50, 49, 48, 48, 49, 48, 48, 48, 48, 48, 51, 50, 49, 49, 50, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 51, 52, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 84, 101, 115, 116, 32, 116, 101, 120, 116, 100, 48, 1, 2, 3, 4, 5, 6, 7, 8, 49, 50, 51, 52, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 49, 55, 65, 110, 111, 116, 104, 101, 114, 32, 116, 101, 115, 116, 32, 116, 101, 120, 116}

	if !bytes.Equal(res, sample) {
		t.Error("ISO Encode error!")
	}
}

func TestDecode(t *testing.T) {

	input := []byte{48, 49, 48, 48, 242, 60, 36, 129, 40, 224, 152, 0, 0, 0, 0, 0, 0, 0, 1, 0, 49, 54, 52, 50, 55, 54, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 55, 55, 55, 48, 48, 48, 55, 48, 49, 49, 49, 49, 56, 52, 52, 48, 48, 48, 49, 50, 51, 49, 51, 49, 56, 52, 52, 48, 55, 48, 49, 49, 57, 48, 50, 6, 67, 57, 48, 49, 48, 50, 48, 54, 49, 50, 51, 52, 53, 54, 51, 55, 52, 50, 55, 54, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 61, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 57, 56, 55, 54, 53, 52, 51, 50, 49, 48, 48, 49, 48, 48, 48, 48, 48, 51, 50, 49, 49, 50, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 51, 52, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 84, 101, 115, 116, 32, 116, 101, 120, 116, 100, 48, 1, 2, 3, 4, 5, 6, 7, 8, 49, 50, 51, 52, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 49, 55, 65, 110, 111, 116, 104, 101, 114, 32, 116, 101, 115, 116, 32, 116, 101, 120, 116}

	// init empty iso message struct
	iso := Message{"", ASCII, true, false, newDataIso()}

	// parse data from bytes to iso struct
	err := iso.Load(input)

	if err != nil {
		t.Error("ISO Decode error:", err)
	}

	resultFields := iso.Data.(*TestISO)

	// check BCD numeric values length
	assert.Equal(t, 3, len(resultFields.F19.Value))
	assert.Equal(t, 3, len(resultFields.F49.Value))

	// check values for BCD (lBCD) and rBCD
	assert.Equal(t, "643", resultFields.F19.Value)
	assert.Equal(t, "643", resultFields.F49.Value)

	var res []byte

	// set second bitmap because field 120 in struct (need if more than 63 fields in message)
	iso.SecondBitmap = true

	// before encode add "0" to left of F19 for testing rBCD encoding
	iso.Data.(*TestISO).F19.Value = "0" + iso.Data.(*TestISO).F19.Value

	// encode iso struct to bytes
	res, err = iso.Bytes()

	if err != nil {
		t.Error("ISO Encode error:", err)
	}

	// parse data from bytes to iso struct to test Bytes() function
	err = iso.Load(res)

	if err != nil {
		t.Error(err)
	}

	// set field 120 value to empty string
	iso.Data.(*TestISO).F120.Value = ""

	iso.SecondBitmap = false

	// encode iso struct to bytes
	res, err = iso.Bytes()

	if err != nil {
		t.Error("ISO Encode error:", err)
	}

	// parse data from bytes to iso struct to test Bytes() function
	err = iso.Load(res)

	if err != nil {
		t.Error(err)
	}

	sample := []byte{48, 49, 48, 48, 114, 60, 36, 129, 40, 224, 152, 0, 49, 54, 52, 50, 55, 54, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 55, 55, 55, 48, 48, 48, 55, 48, 49, 49, 49, 49, 56, 52, 52, 48, 48, 48, 49, 50, 51, 49, 51, 49, 56, 52, 52, 48, 55, 48, 49, 49, 57, 48, 50, 6, 67, 57, 48, 49, 48, 50, 48, 54, 49, 50, 51, 52, 53, 54, 51, 55, 52, 50, 55, 54, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 61, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 57, 56, 55, 54, 53, 52, 51, 50, 49, 48, 48, 49, 48, 48, 48, 48, 48, 51, 50, 49, 49, 50, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 51, 52, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 84, 101, 115, 116, 32, 116, 101, 120, 116, 100, 48, 1, 2, 3, 4, 5, 6, 7, 8, 49, 50, 51, 52, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48}

	if !bytes.Equal(res, sample) {
		t.Error("ISO Encode error!")
	}
}

func TestEncodeASCIIBitmap(t *testing.T) {
	data := &TestISO{
		F2:   NewLlnumeric("4276555555555555"),
		F3:   NewNumeric("000000"),
		F4:   NewNumeric("000000077700"),
		F7:   NewNumeric("0701111844"),
		F11:  NewNumeric("000123"),
		F12:  NewNumeric("131844"),
		F13:  NewNumeric("0701"),
		F14:  NewNumeric("1902"),
		F19:  NewNumeric("643"),
		F22:  NewNumeric("901"),
		F25:  NewNumeric("02"),
		F32:  NewLlnumeric("123456"),
		F35:  NewLlnumeric("4276555555555555=12345678901234567890"),
		F37:  NewAlphanumeric("987654321001"),
		F39:  NewAlphanumeric(""),
		F41:  NewAlphanumeric("00000321"),
		F42:  NewAlphanumeric("120000000000034"),
		F43:  NewAlphanumeric("Test text"),
		F49:  NewNumeric("643"),
		F52:  NewBinary([]byte{1, 2, 3, 4, 5, 6, 7, 8}),
		F53:  NewNumeric("1234000000000000"),
		F120: NewLllnumeric("Another test text"),
	}

	iso := Message{"0100", ASCII, true, true, data}

	res, err := iso.Bytes()

	if err != nil {
		t.Error("ISO Encode error:", err)
	}

	expected := "0100F23C248128E098000000000000000100164276555555555555000000000000077700070111184400012313184407011902\x06C9010206123456374276555555555555=1234567890123456789098765432100100000321120000000000034                               Test textd0\x01\x02\x03\x04\x05\x06\x07\x081234000000000000017Another test text"
	if string(res) != expected {
		t.Error("ISO Encode error!")
	}
}
func TestDecodeASCIIBitmap(t *testing.T) {
	input := []byte("0100F23C248128E098000000000000000100164276555555555555000000000000077700070111184400012313184407011902\x06C9010206123456374276555555555555=1234567890123456789098765432100100000321120000000000034                               Test textd0\x01\x02\x03\x04\x05\x06\x07\x081234000000000000017Another test text")

	iso := Message{"", ASCII, true, true, newDataIso()}
	err := iso.Load(input)

	assert.NoError(t, err, "ISO Decode error:")

	resultFields := iso.Data.(*TestISO)

	assert.Equal(t, "4276555555555555", resultFields.F2.Value)
	assert.Equal(t, "000000", resultFields.F3.Value)
	assert.Equal(t, "000000077700", resultFields.F4.Value)
	assert.Equal(t, "0701111844", resultFields.F7.Value)
	assert.Equal(t, "000123", resultFields.F11.Value)
	assert.Equal(t, "131844", resultFields.F12.Value)
	assert.Equal(t, "0701", resultFields.F13.Value)
	assert.Equal(t, "1902", resultFields.F14.Value)
	assert.Equal(t, "643", resultFields.F19.Value)
	assert.Equal(t, "901", resultFields.F22.Value)
	assert.Equal(t, "02", resultFields.F25.Value)
	assert.Equal(t, "123456", resultFields.F32.Value)
	assert.Equal(t, "4276555555555555=12345678901234567890", resultFields.F35.Value)
	assert.Equal(t, "987654321001", resultFields.F37.Value)
	assert.Equal(t, "", resultFields.F39.Value)
	assert.Equal(t, "00000321", resultFields.F41.Value)
	assert.Equal(t, "120000000000034", resultFields.F42.Value)
	assert.Equal(t, "                               Test text", resultFields.F43.Value)
	assert.Equal(t, "643", resultFields.F49.Value)
	assert.Equal(t, []byte{1, 2, 3, 4, 5, 6, 7, 8}, resultFields.F52.Value)
	assert.Equal(t, "1234000000000000", resultFields.F53.Value)
	assert.Equal(t, "Another test text", resultFields.F120.Value)
}

func TestEncodeDecode(t *testing.T) {

	data := &TestISO2{
		F2:  NewLlnumeric("4276555555555555"),
		F3:  NewNumeric("000000"),
		F4:  NewNumeric("000000077700"),
		F7:  NewNumeric("0701111844"),
		F11: NewNumeric("123"),
		F12: NewNumeric("131844"),
		F13: NewNumeric("0701"),
		F14: NewNumeric("1902"),
		F19: NewNumeric("643"),
		F22: NewNumeric("901"),
		F25: NewNumeric("02"),
		F28: NewAlphanumeric("abcd12345"),
		F32: NewLlnumeric("123456"),
		F35: NewLlnumeric("4276555555555555=12345678901234567890"),
		F37: NewAlphanumeric("987654321001"),
		F39: NewAlphanumeric("00"),
		F41: NewAlphanumeric("00000321"),
		F42: NewAlphanumeric("120000000000034"),
		F43: NewAlphanumeric("Test text"),
		F45: NewLlnumeric("1230abc"),
		F49: NewNumeric("643"),
		F52: NewBinary([]byte{1, 2, 3, 4, 5}),
		F53: NewNumeric("1234000000000000"),
		F54: NewLlvar([]byte{7, 8, 56, 71, 35}),
		F55: NewLlvar([]byte{0, 1, 2, 5, 51, 47, 45, 32, 158}),
		F56: NewLllvar([]byte("test data1")),
		F57: NewLllvar([]byte("test data2")),
		F58: NewLllvar([]byte("test data3")),
		F59: NewLlvar([]byte("test data4")),
		F60: NewLllnumeric("123456789"),
		F61: NewLllnumeric("abcdef"),
		F63: NewLllnumeric("123abc456ef7890"),
	}

	iso := NewMessage("0110", data)

	res, err := iso.Bytes()

	if err != nil {
		t.Error("ISO Encode error:", err)
	}

	iso2 := NewMessage("0110", data)

	err = iso2.Load(res)

	if err != nil {
		t.Error("ISO Encode error:", err)
	}

	// check data after encode/decode
	assert.Equal(t, iso, iso2)
}

func TestFieldNumericEncodeErrors(t *testing.T) {

	type test1 struct {
		F2 *Numeric `field:"2" length:"6" encode:"test"`
	}

	data1 := &test1{
		F2: NewNumeric("123456"),
	}

	iso := NewMessage("0110", data1)

	_, err := iso.Bytes()

	assert.EqualError(t, err, "invalid encoder")

	type test2 struct {
		F2 *Numeric `field:"2"`
	}

	data2 := &test2{
		F2: NewNumeric("123456"),
	}

	iso = NewMessage("0110", data2)

	_, err = iso.Bytes()

	assert.EqualError(t, err, "missing length")

	type test3 struct {
		F2 *Numeric `field:"2" length:"3"`
	}

	data3 := &test3{
		F2: NewNumeric("123456"),
	}

	iso = NewMessage("0110", data3)

	_, err = iso.Bytes()

	assert.EqualError(t, err, "length of value is longer than definition; type=Numeric, def_len=3, len=6")
}

func TestFieldAlphanumericEncodeErrors(t *testing.T) {

	type test1 struct {
		F2 *Alphanumeric `field:"2"`
	}

	data1 := &test1{
		F2: NewAlphanumeric("abcdef"),
	}

	iso := NewMessage("0110", data1)

	_, err := iso.Bytes()

	assert.EqualError(t, err, "missing length")

	type test2 struct {
		F2 *Alphanumeric `field:"2" length:"3"`
	}

	data2 := &test2{
		F2: NewAlphanumeric("abcdef"),
	}

	iso = NewMessage("0110", data2)

	_, err = iso.Bytes()

	assert.EqualError(t, err, "length of value is longer than definition; type=Alphanumeric, def_len=3, len=6")
}

func TestFieldBinaryEncodeErrors(t *testing.T) {

	type test1 struct {
		F2 *Binary `field:"2"`
	}

	data1 := &test1{
		F2: NewBinary([]byte("abcdef")),
	}

	iso := NewMessage("0110", data1)

	_, err := iso.Bytes()

	assert.EqualError(t, err, "missing length")

	type test2 struct {
		F2 *Binary `field:"2" length:"3"`
	}

	data2 := &test2{
		F2: NewBinary([]byte("abcdef")),
	}

	iso = NewMessage("0110", data2)

	_, err = iso.Bytes()

	assert.EqualError(t, err, "length of value is longer than definition; type=Binary, def_len=3, len=6")
}

func TestFieldLlnumericEncodeErrors(t *testing.T) {

	type test1 struct {
		F2 *Llnumeric `field:"2" length:"6" encode:"test"`
	}

	data1 := &test1{
		F2: NewLlnumeric("123456"),
	}

	iso := NewMessage("0110", data1)

	_, err := iso.Bytes()

	assert.EqualError(t, err, "invalid encoder")

	type test2 struct {
		F2 *Llnumeric `field:"2" length:"3"`
	}

	data2 := &test2{
		F2: NewLlnumeric("123456"),
	}

	iso = NewMessage("0110", data2)

	_, err = iso.Bytes()

	assert.EqualError(t, err, "length of value is longer than definition; type=Llnumeric, def_len=3, len=6")

	type test3 struct {
		F2 *Llnumeric `field:"2" encode:"ascii,ascii"`
	}

	data3 := &test3{
		F2: NewLlnumeric(string(bytes.Repeat([]byte("a"), 100))),
	}

	iso = NewMessage("0110", data3)

	_, err = iso.Bytes()

	assert.EqualError(t, err, "invalid length head")

	type test4 struct {
		F2 *Llnumeric `field:"2" length:"100" encode:"bcd,ascii"`
	}

	data4 := &test4{
		F2: NewLlnumeric(string(bytes.Repeat([]byte("a"), 100))),
	}

	iso = NewMessage("0110", data4)

	_, err = iso.Bytes()

	assert.EqualError(t, err, "invalid length head")

	type test5 struct {
		F2 *Llnumeric `field:"2" length:"6" encode:"test,ascii"`
	}

	data5 := &test5{
		F2: NewLlnumeric("123456"),
	}

	iso = NewMessage("0110", data5)

	_, err = iso.Bytes()

	assert.EqualError(t, err, "invalid length encoder")
}

func TestFieldLllnumericEncodeErrors(t *testing.T) {

	type test1 struct {
		F2 *Lllnumeric `field:"2" length:"6" encode:"test"`
	}

	data1 := &test1{
		F2: NewLllnumeric("123456"),
	}

	iso := NewMessage("0110", data1)

	_, err := iso.Bytes()

	assert.EqualError(t, err, "invalid encoder")

	type test2 struct {
		F2 *Lllnumeric `field:"2" length:"3"`
	}

	data2 := &test2{
		F2: NewLllnumeric("123456"),
	}

	iso = NewMessage("0110", data2)

	_, err = iso.Bytes()

	assert.EqualError(t, err, "length of value is longer than definition; type=Lllnumeric, def_len=3, len=6")

	type test3 struct {
		F2 *Lllnumeric `field:"2" encode:"ascii,ascii"`
	}

	data3 := &test3{
		F2: NewLllnumeric(string(bytes.Repeat([]byte("a"), 1000))),
	}

	iso = NewMessage("0110", data3)

	_, err = iso.Bytes()

	assert.EqualError(t, err, "invalid length head")

	type test4 struct {
		F2 *Lllnumeric `field:"2" length:"1000" encode:"bcd,ascii"`
	}

	data4 := &test4{
		F2: NewLllnumeric(string(bytes.Repeat([]byte("a"), 1000))),
	}

	iso = NewMessage("0110", data4)

	_, err = iso.Bytes()

	assert.EqualError(t, err, "invalid length head")

	type test5 struct {
		F2 *Lllnumeric `field:"2" length:"6" encode:"test,ascii"`
	}

	data5 := &test5{
		F2: NewLllnumeric("123456"),
	}

	iso = NewMessage("0110", data5)

	_, err = iso.Bytes()

	assert.EqualError(t, err, "invalid length encoder")
}

func TestFieldLlvarEncodeErrors(t *testing.T) {

	type test1 struct {
		F2 *Llvar `field:"2" length:"6" encode:"test"`
	}

	data1 := &test1{
		F2: NewLlvar([]byte("123456")),
	}

	iso := NewMessage("0110", data1)

	_, err := iso.Bytes()

	assert.EqualError(t, err, "invalid encoder")

	type test2 struct {
		F2 *Llvar `field:"2" length:"3"`
	}

	data2 := &test2{
		F2: NewLlvar([]byte("123456")),
	}

	iso = NewMessage("0110", data2)

	_, err = iso.Bytes()

	assert.EqualError(t, err, "length of value is longer than definition; type=Llvar, def_len=3, len=6")

	type test3 struct {
		F2 *Llvar `field:"2" encode:"ascii,ascii"`
	}

	data3 := &test3{
		F2: NewLlvar(bytes.Repeat([]byte("a"), 100)),
	}

	iso = NewMessage("0110", data3)

	_, err = iso.Bytes()

	assert.EqualError(t, err, "invalid length head")

	type test4 struct {
		F2 *Llvar `field:"2" length:"100" encode:"bcd,ascii"`
	}

	data4 := &test4{
		F2: NewLlvar(bytes.Repeat([]byte("a"), 100)),
	}

	iso = NewMessage("0110", data4)

	_, err = iso.Bytes()

	assert.EqualError(t, err, "invalid length head")

	type test5 struct {
		F2 *Llvar `field:"2" length:"6" encode:"test,ascii"`
	}

	data5 := &test5{
		F2: NewLlvar([]byte("123456")),
	}

	iso = NewMessage("0110", data5)

	_, err = iso.Bytes()

	assert.EqualError(t, err, "invalid length encoder")
}

func TestFieldLllvarEncodeErrors(t *testing.T) {

	type test1 struct {
		F2 *Lllvar `field:"2" length:"6" encode:"test"`
	}

	data1 := &test1{
		F2: NewLllvar([]byte("123456")),
	}

	iso := NewMessage("0110", data1)

	_, err := iso.Bytes()

	assert.EqualError(t, err, "invalid encoder")

	type test2 struct {
		F2 *Lllvar `field:"2" length:"3"`
	}

	data2 := &test2{
		F2: NewLllvar([]byte("123456")),
	}

	iso = NewMessage("0110", data2)

	_, err = iso.Bytes()

	assert.EqualError(t, err, "length of value is longer than definition; type=Lllvar, def_len=3, len=6")

	type test3 struct {
		F2 *Lllvar `field:"2" encode:"ascii,ascii"`
	}

	data3 := &test3{
		F2: NewLllvar(bytes.Repeat([]byte("a"), 1000)),
	}

	iso = NewMessage("0110", data3)

	_, err = iso.Bytes()

	assert.EqualError(t, err, "invalid length head")

	type test4 struct {
		F2 *Lllvar `field:"2" length:"1000" encode:"bcd,ascii"`
	}

	data4 := &test4{
		F2: NewLllvar(bytes.Repeat([]byte("a"), 1000)),
	}

	iso = NewMessage("0110", data4)

	_, err = iso.Bytes()

	assert.EqualError(t, err, "invalid length head")

	type test5 struct {
		F2 *Lllvar `field:"2" length:"6" encode:"test,ascii"`
	}

	data5 := &test5{
		F2: NewLllvar([]byte("123456")),
	}

	iso = NewMessage("0110", data5)

	_, err = iso.Bytes()

	assert.EqualError(t, err, "invalid length encoder")
}

func TestFieldNumericDecodeErrors(t *testing.T) {
	type test1 struct {
		F2 *Numeric `field:"2" length:"10" encode:"ascii"`
	}

	data1 := &test1{
		F2: NewNumeric("123456"),
	}

	iso := NewMessage("0110", data1)

	isoBytes, err := iso.Bytes()

	assert.Empty(t, err)

	isoBytes = isoBytes[0 : len(isoBytes)-1]

	err = iso.Load(isoBytes)

	assert.EqualError(t, err, "field 2: bad raw data")

	type test2 struct {
		F2 *Numeric `field:"2" length:"10" encode:"bcd"`
	}

	data2 := &test2{
		F2: NewNumeric("123456"),
	}

	iso = NewMessage("0110", data2)

	isoBytes, err = iso.Bytes()

	assert.Empty(t, err)

	isoBytes = isoBytes[0 : len(isoBytes)-1]

	err = iso.Load(isoBytes)

	assert.EqualError(t, err, "field 2: bad raw data")

	type test3 struct {
		F2 *Numeric `field:"2" length:"10" encode:"rbcd"`
	}

	data3 := &test3{
		F2: NewNumeric("123456"),
	}

	iso = NewMessage("0110", data3)

	isoBytes, err = iso.Bytes()

	assert.Empty(t, err)

	isoBytes = isoBytes[0 : len(isoBytes)-1]

	err = iso.Load(isoBytes)

	assert.EqualError(t, err, "field 2: bad raw data")

	type test4 struct {
		F2 *Numeric `field:"2" encode:"rbcd"`
	}

	data4 := &test4{
		F2: NewNumeric(""),
	}

	iso = NewMessage("0110", data4)

	err = iso.Load(isoBytes)

	assert.EqualError(t, err, "field 2: missing length")

	type test5 struct {
		F2 *Numeric `field:"2" length:"10" encode:"test"`
	}

	data5 := &test5{
		F2: NewNumeric(""),
	}

	iso = NewMessage("0110", data5)

	err = iso.Load(isoBytes)

	assert.EqualError(t, err, "field 2: invalid encoder")
}

func TestFieldLlnumericDecodeErrors(t *testing.T) {
	type test1 struct {
		F2 *Llnumeric `field:"2" length:"10" encode:"ascii"`
	}

	data1 := &test1{
		F2: NewLlnumeric("123456"),
	}

	iso := NewMessage("0110", data1)

	isoBytes, err := iso.Bytes()

	assert.Empty(t, err)

	isoBytes = isoBytes[0 : len(isoBytes)-1]

	err = iso.Load(isoBytes)

	assert.EqualError(t, err, "field 2: bad raw data")

	type test2 struct {
		F2 *Llnumeric `field:"2" length:"10" encode:"bcd"`
	}

	data2 := &test2{
		F2: NewLlnumeric("123456"),
	}

	iso = NewMessage("0110", data2)

	isoBytes, err = iso.Bytes()

	assert.Empty(t, err)

	isoBytes = isoBytes[0 : len(isoBytes)-1]

	err = iso.Load(isoBytes)

	assert.EqualError(t, err, "field 2: bad raw data")

	type test3 struct {
		F2 *Llnumeric `field:"2" length:"10" encode:"rbcd"`
	}

	data3 := &test3{
		F2: NewLlnumeric("123456"),
	}

	iso = NewMessage("0110", data3)

	isoBytes, err = iso.Bytes()

	assert.Empty(t, err)

	isoBytes = isoBytes[0 : len(isoBytes)-1]

	err = iso.Load(isoBytes)

	assert.EqualError(t, err, "field 2: bad raw data")

	type test4 struct {
		F2 *Llnumeric `field:"2" length:"10" encode:"test"`
	}

	data4 := &test4{
		F2: NewLlnumeric(""),
	}

	iso = NewMessage("0110", data4)

	err = iso.Load(isoBytes)

	assert.EqualError(t, err, "field 2: invalid encoder")

	type test5 struct {
		F2 *Llnumeric `field:"2" length:"10" encode:"bcd,ascii"`
	}

	data5 := &test5{
		F2: NewLlnumeric("543210"),
	}

	iso = NewMessage("0110", data5)

	isoBytes, err = iso.Bytes()

	assert.Empty(t, err)

	isoBytes[12] = 123

	err = iso.Load(isoBytes)

	assert.EqualError(t, err, "field 2: parse length head failed: {")

	type test6 struct {
		F2 *Llnumeric `field:"2" length:"10" encode:"rbcd,ascii"`
	}

	data6 := &test6{
		F2: NewLlnumeric(""),
	}

	iso = NewMessage("0110", data6)

	err = iso.Load(isoBytes)

	assert.EqualError(t, err, "field 2: parse length head failed: {")

	type test7 struct {
		F2 *Llnumeric `field:"2" length:"10" encode:"ascii,ascii"`
	}

	data7 := &test7{
		F2: NewLlnumeric("543210"),
	}

	iso = NewMessage("0110", data7)

	isoBytes, err = iso.Bytes()

	assert.Empty(t, err)

	isoBytes[12] = 123
	isoBytes[13] = 125

	err = iso.Load(isoBytes)

	assert.EqualError(t, err, "field 2: parse length head failed: {}")

	type test8 struct {
		F2 *Llnumeric `field:"2" length:"10" encode:"test,ascii"`
	}

	data8 := &test8{
		F2: NewLlnumeric(""),
	}

	iso = NewMessage("0110", data8)

	err = iso.Load(isoBytes)

	assert.EqualError(t, err, "field 2: invalid length encoder")
}

func TestFieldLllnumericDecodeErrors(t *testing.T) {
	type test1 struct {
		F2 *Lllnumeric `field:"2" length:"10" encode:"ascii"`
	}

	data1 := &test1{
		F2: NewLllnumeric("123456"),
	}

	iso := NewMessage("0110", data1)

	isoBytes, err := iso.Bytes()

	assert.Empty(t, err)

	isoBytes = isoBytes[0 : len(isoBytes)-1]

	err = iso.Load(isoBytes)

	assert.EqualError(t, err, "field 2: bad raw data")

	type test2 struct {
		F2 *Lllnumeric `field:"2" length:"10" encode:"bcd"`
	}

	data2 := &test2{
		F2: NewLllnumeric("123456"),
	}

	iso = NewMessage("0110", data2)

	isoBytes, err = iso.Bytes()

	assert.Empty(t, err)

	isoBytes = isoBytes[0 : len(isoBytes)-1]

	err = iso.Load(isoBytes)

	assert.EqualError(t, err, "field 2: bad raw data")

	type test3 struct {
		F2 *Lllnumeric `field:"2" length:"10" encode:"rbcd"`
	}

	data3 := &test3{
		F2: NewLllnumeric("123456"),
	}

	iso = NewMessage("0110", data3)

	isoBytes, err = iso.Bytes()

	assert.Empty(t, err)

	isoBytes = isoBytes[0 : len(isoBytes)-1]

	err = iso.Load(isoBytes)

	assert.EqualError(t, err, "field 2: bad raw data")

	type test4 struct {
		F2 *Lllnumeric `field:"2" length:"10" encode:"test"`
	}

	data4 := &test4{
		F2: NewLllnumeric(""),
	}

	iso = NewMessage("0110", data4)

	err = iso.Load(isoBytes)

	assert.EqualError(t, err, "field 2: invalid encoder")

	type test5 struct {
		F2 *Lllnumeric `field:"2" length:"10" encode:"bcd,ascii"`
	}

	data5 := &test5{
		F2: NewLllnumeric("543210"),
	}

	iso = NewMessage("0110", data5)

	isoBytes, err = iso.Bytes()

	assert.Empty(t, err)

	isoBytes[12] = 123
	isoBytes[13] = 125

	err = iso.Load(isoBytes)

	assert.EqualError(t, err, "field 2: parse length head failed: {}")

	type test6 struct {
		F2 *Lllnumeric `field:"2" length:"10" encode:"rbcd,ascii"`
	}

	data6 := &test6{
		F2: NewLllnumeric(""),
	}

	iso = NewMessage("0110", data6)

	err = iso.Load(isoBytes)

	assert.EqualError(t, err, "field 2: parse length head failed: {}")

	type test7 struct {
		F2 *Lllnumeric `field:"2" length:"10" encode:"ascii,ascii"`
	}

	data7 := &test7{
		F2: NewLllnumeric("543210"),
	}

	iso = NewMessage("0110", data7)

	isoBytes, err = iso.Bytes()

	assert.Empty(t, err)

	isoBytes[12] = 123
	isoBytes[13] = 124
	isoBytes[14] = 125

	err = iso.Load(isoBytes)

	assert.EqualError(t, err, "field 2: parse length head failed: {|}")

	type test8 struct {
		F2 *Lllnumeric `field:"2" length:"10" encode:"test,ascii"`
	}

	data8 := &test8{
		F2: NewLllnumeric(""),
	}

	iso = NewMessage("0110", data8)

	err = iso.Load(isoBytes)

	assert.EqualError(t, err, "field 2: invalid length encoder")
}

func TestFieldLlvarDecodeErrors(t *testing.T) {
	type test1 struct {
		F2 *Llvar `field:"2" length:"10" encode:"ascii"`
	}

	data1 := &test1{
		F2: NewLlvar([]byte("123456")),
	}

	iso := NewMessage("0110", data1)

	isoBytes, err := iso.Bytes()

	assert.Empty(t, err)

	isoBytes = isoBytes[0 : len(isoBytes)-1]

	err = iso.Load(isoBytes)

	assert.EqualError(t, err, "field 2: bad raw data")

	type test2 struct {
		F2 *Llvar `field:"2" length:"10" encode:"bcd,ascii"`
	}

	data2 := &test2{
		F2: NewLlvar([]byte("123456")),
	}

	iso = NewMessage("0110", data2)

	isoBytes, err = iso.Bytes()

	assert.Empty(t, err)

	isoBytes = isoBytes[0 : len(isoBytes)-1]

	err = iso.Load(isoBytes)

	assert.EqualError(t, err, "field 2: bad raw data")

	type test3 struct {
		F2 *Llvar `field:"2" length:"10" encode:"rbcd,ascii"`
	}

	data3 := &test3{
		F2: NewLlvar([]byte("123456")),
	}

	iso = NewMessage("0110", data3)

	isoBytes, err = iso.Bytes()

	assert.Empty(t, err)

	isoBytes2 := isoBytes[0 : len(isoBytes)-1]

	err = iso.Load(isoBytes2)

	assert.EqualError(t, err, "field 2: bad raw data")

	type test4 struct {
		F2 *Llvar `field:"2" length:"10" encode:"rbcd,test"`
	}

	data4 := &test4{
		F2: NewLlvar(nil),
	}

	iso = NewMessage("0110", data4)

	err = iso.Load(isoBytes)

	assert.EqualError(t, err, "field 2: invalid encoder")

	type test5 struct {
		F2 *Llvar `field:"2" length:"10" encode:"bcd,ascii"`
	}

	data5 := &test5{
		F2: NewLlvar([]byte("543210")),
	}

	iso = NewMessage("0110", data5)

	isoBytes, err = iso.Bytes()

	assert.Empty(t, err)

	isoBytes[12] = 123

	err = iso.Load(isoBytes)

	assert.EqualError(t, err, "field 2: parse length head failed: {")

	type test6 struct {
		F2 *Llvar `field:"2" length:"10" encode:"rbcd,ascii"`
	}

	data6 := &test6{
		F2: NewLlvar(nil),
	}

	iso = NewMessage("0110", data6)

	err = iso.Load(isoBytes)

	assert.EqualError(t, err, "field 2: parse length head failed: {")

	type test7 struct {
		F2 *Llvar `field:"2" length:"10" encode:"ascii,ascii"`
	}

	data7 := &test7{
		F2: NewLlvar([]byte("543210")),
	}

	iso = NewMessage("0110", data7)

	isoBytes, err = iso.Bytes()

	assert.Empty(t, err)

	isoBytes[12] = 123
	isoBytes[13] = 125

	err = iso.Load(isoBytes)

	assert.EqualError(t, err, "field 2: parse length head failed: {}")

	type test8 struct {
		F2 *Llvar `field:"2" length:"10" encode:"test,ascii"`
	}

	data8 := &test8{
		F2: NewLlvar(nil),
	}

	iso = NewMessage("0110", data8)

	err = iso.Load(isoBytes)

	assert.EqualError(t, err, "field 2: invalid length encoder")
}

func TestFieldLllvarDecodeErrors(t *testing.T) {
	type test1 struct {
		F2 *Lllvar `field:"2" length:"10" encode:"ascii"`
	}

	data1 := &test1{
		F2: NewLllvar([]byte("123456")),
	}

	iso := NewMessage("0110", data1)

	isoBytes, err := iso.Bytes()

	assert.Empty(t, err)

	isoBytes = isoBytes[0 : len(isoBytes)-1]

	err = iso.Load(isoBytes)

	assert.EqualError(t, err, "field 2: bad raw data")

	type test2 struct {
		F2 *Lllvar `field:"2" length:"10" encode:"bcd,ascii"`
	}

	data2 := &test2{
		F2: NewLllvar([]byte("123456")),
	}

	iso = NewMessage("0110", data2)

	isoBytes, err = iso.Bytes()

	assert.Empty(t, err)

	isoBytes = isoBytes[0 : len(isoBytes)-1]

	err = iso.Load(isoBytes)

	assert.EqualError(t, err, "field 2: bad raw data")

	type test3 struct {
		F2 *Lllvar `field:"2" length:"10" encode:"rbcd,ascii"`
	}

	data3 := &test3{
		F2: NewLllvar([]byte("123456")),
	}

	iso = NewMessage("0110", data3)

	isoBytes, err = iso.Bytes()

	assert.Empty(t, err)

	isoBytes2 := isoBytes[0 : len(isoBytes)-1]

	err = iso.Load(isoBytes2)

	assert.EqualError(t, err, "field 2: bad raw data")

	type test4 struct {
		F2 *Lllvar `field:"2" length:"10" encode:"rbcd,test"`
	}

	data4 := &test4{
		F2: NewLllvar(nil),
	}

	iso = NewMessage("0110", data4)

	err = iso.Load(isoBytes)

	assert.EqualError(t, err, "field 2: invalid encoder")

	type test5 struct {
		F2 *Lllvar `field:"2" length:"10" encode:"bcd,ascii"`
	}

	data5 := &test5{
		F2: NewLllvar([]byte("543210")),
	}

	iso = NewMessage("0110", data5)

	isoBytes, err = iso.Bytes()

	assert.Empty(t, err)

	isoBytes[12] = 123
	isoBytes[13] = 125

	err = iso.Load(isoBytes)

	assert.EqualError(t, err, "field 2: parse length head failed: {}")

	type test6 struct {
		F2 *Lllvar `field:"2" length:"10" encode:"rbcd,ascii"`
	}

	data6 := &test6{
		F2: NewLllvar(nil),
	}

	iso = NewMessage("0110", data6)

	err = iso.Load(isoBytes)

	assert.EqualError(t, err, "field 2: parse length head failed: {}")

	type test7 struct {
		F2 *Lllvar `field:"2" length:"10" encode:"ascii,ascii"`
	}

	data7 := &test7{
		F2: NewLllvar([]byte("543210")),
	}

	iso = NewMessage("0110", data7)

	isoBytes, err = iso.Bytes()

	assert.Empty(t, err)

	isoBytes[12] = 123
	isoBytes[13] = 124
	isoBytes[14] = 125

	err = iso.Load(isoBytes)

	assert.EqualError(t, err, "field 2: parse length head failed: {|}")

	type test8 struct {
		F2 *Lllvar `field:"2" length:"10" encode:"test,ascii"`
	}

	data8 := &test8{
		F2: NewLllvar(nil),
	}

	iso = NewMessage("0110", data8)

	err = iso.Load(isoBytes)

	assert.EqualError(t, err, "field 2: invalid length encoder")
}

func TestFieldAlphanumericDecodeErrors(t *testing.T) {
	type test1 struct {
		F2 *Alphanumeric `field:"2" length:"10"`
	}

	data1 := &test1{
		F2: NewAlphanumeric("123456"),
	}

	iso := NewMessage("0110", data1)

	isoBytes, err := iso.Bytes()

	assert.Empty(t, err)

	isoBytes = isoBytes[0 : len(isoBytes)-1]

	err = iso.Load(isoBytes)

	assert.EqualError(t, err, "field 2: bad raw data")

	type test2 struct {
		F2 *Alphanumeric `field:"2"`
	}

	data2 := &test2{
		F2: NewAlphanumeric("123456"),
	}

	iso = NewMessage("0110", data2)

	err = iso.Load(isoBytes)

	assert.EqualError(t, err, "field 2: missing length")
}

func TestFieldBinaryDecodeErrors(t *testing.T) {
	type test1 struct {
		F2 *Binary `field:"2" length:"10"`
	}

	data1 := &test1{
		F2: NewBinary([]byte("123456")),
	}

	iso := NewMessage("0110", data1)

	isoBytes, err := iso.Bytes()

	assert.Empty(t, err)

	isoBytes = isoBytes[0 : len(isoBytes)-1]

	err = iso.Load(isoBytes)

	assert.EqualError(t, err, "field 2: bad raw data")

	type test2 struct {
		F2 *Binary `field:"2"`
	}

	data2 := &test2{
		F2: NewBinary([]byte("123456")),
	}

	iso = NewMessage("0110", data2)

	err = iso.Load(isoBytes)

	assert.EqualError(t, err, "field 2: missing length")
}

func TestParserErrors(t *testing.T) {

	parser := Parser{}

	err := parser.Register("0100", nil)

	assert.EqualError(t, err, "Critical error:reflect: call of reflect.Value.Type on zero Value")

	err = parser.Register("1", newDataIso())

	assert.EqualError(t, err, "MTI must be a 4 digit numeric field")

	_, err = parser.Parse([]byte{0})

	assert.EqualError(t, err, "bad MTI raw data")

	parser.MtiEncode = BCD

	_, err = parser.Parse([]byte{1, 2})

	assert.EqualError(t, err, "no template registered for MTI: 0102")

	parser.MtiEncode = 10

	_, err = parser.Parse([]byte{1, 2, 3, 4})

	assert.EqualError(t, err, "invalid encode type")

	parser.MtiEncode = ASCII

	input := []byte{48, 49, 48, 48, 242, 60, 36, 129, 40, 224, 152, 0, 0, 0, 0, 0, 0, 0, 1, 0, 49, 54, 52, 50, 55, 54, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 55, 55, 55, 48, 48, 48, 55, 48, 49, 49, 49, 49, 56, 52, 52, 48, 48, 48, 49, 50, 51, 49, 51, 49, 56, 52, 52, 48, 55, 48, 49, 49, 57, 48, 50, 6, 67, 57, 48, 49, 48, 50, 48, 54, 49, 50, 51, 52, 53, 54, 51, 55, 52, 50, 55, 54, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 61, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 57, 56, 55, 54, 53, 52, 51, 50, 49, 48, 48, 49, 48, 48, 48, 48, 48, 51, 50, 49, 49, 50, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 51, 52, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 84, 101, 115, 116, 32, 116, 101, 120, 116, 100, 48, 1, 2, 3, 4, 5, 6, 7, 8, 49, 50, 51, 52, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 49, 55, 65, 110, 111, 116, 104, 101, 114, 32, 116, 101, 115, 116, 32, 116, 101, 120, 116}

	err = parser.Register("0100", newDataIso())
	assert.NoError(t, err)

	_, err = parser.Parse(input[0:23])

	assert.EqualError(t, err, "field 2: bad raw data")

	parser.messages["0100"] = nil

	_, err = parser.Parse(input)

	assert.EqualError(t, err, "Critical error:reflect: New(nil)")
}

func TestParser(t *testing.T) {

	input := []byte{48, 49, 48, 48, 242, 60, 36, 129, 40, 224, 152, 0, 0, 0, 0, 0, 0, 0, 1, 0, 49, 54, 52, 50, 55, 54, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 55, 55, 55, 48, 48, 48, 55, 48, 49, 49, 49, 49, 56, 52, 52, 48, 48, 48, 49, 50, 51, 49, 51, 49, 56, 52, 52, 48, 55, 48, 49, 49, 57, 48, 50, 6, 67, 57, 48, 49, 48, 50, 48, 54, 49, 50, 51, 52, 53, 54, 51, 55, 52, 50, 55, 54, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 61, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 57, 56, 55, 54, 53, 52, 51, 50, 49, 48, 48, 49, 48, 48, 48, 48, 48, 51, 50, 49, 49, 50, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 51, 52, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 84, 101, 115, 116, 32, 116, 101, 120, 116, 100, 48, 1, 2, 3, 4, 5, 6, 7, 8, 49, 50, 51, 52, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 49, 55, 65, 110, 111, 116, 104, 101, 114, 32, 116, 101, 115, 116, 32, 116, 101, 120, 116}

	parser := Parser{}

	err := parser.Register("0100", newDataIso())

	assert.Equal(t, nil, err)
	// parse data from bytes to iso struct
	// parse data from bytes to iso struct
	iso, err := parser.Parse(input)

	if err != nil {
		t.Error("ISO Decode error:", err)
	}

	resultFields := iso.Data.(*TestISO)

	// check BCD numeric values length
	assert.Equal(t, 3, len(resultFields.F19.Value))
	assert.Equal(t, 3, len(resultFields.F49.Value))

	// check values for BCD (lBCD) and rBCD
	assert.Equal(t, "643", resultFields.F19.Value)
	assert.Equal(t, "643", resultFields.F49.Value)

	var res []byte

	// set second bitmap because field 120 in struct (need if more than 63 fields in message)
	iso.SecondBitmap = true

	// before encode add "0" to left of F19 for testing rBCD encoding
	iso.Data.(*TestISO).F19.Value = "0" + iso.Data.(*TestISO).F19.Value

	// encode iso struct to bytes
	res, err = iso.Bytes()

	if err != nil {
		t.Error("ISO Encode error:", err)
	}

	// parse data from bytes to iso struct to test Bytes() function
	err = iso.Load(res)

	if err != nil {
		t.Error(err)
	}

	// set field 120 value to empty string
	iso.Data.(*TestISO).F120.Value = ""

	iso.SecondBitmap = false

	// encode iso struct to bytes
	res, err = iso.Bytes()

	if err != nil {
		t.Error("ISO Encode error:", err)
	}

	// parse data from bytes to iso struct to test Bytes() function
	err = iso.Load(res)

	if err != nil {
		t.Error(err)
	}

	sample := []byte{48, 49, 48, 48, 114, 60, 36, 129, 40, 224, 152, 0, 49, 54, 52, 50, 55, 54, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 55, 55, 55, 48, 48, 48, 55, 48, 49, 49, 49, 49, 56, 52, 52, 48, 48, 48, 49, 50, 51, 49, 51, 49, 56, 52, 52, 48, 55, 48, 49, 49, 57, 48, 50, 6, 67, 57, 48, 49, 48, 50, 48, 54, 49, 50, 51, 52, 53, 54, 51, 55, 52, 50, 55, 54, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 61, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 57, 56, 55, 54, 53, 52, 51, 50, 49, 48, 48, 49, 48, 48, 48, 48, 48, 51, 50, 49, 49, 50, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 51, 52, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 84, 101, 115, 116, 32, 116, 101, 120, 116, 100, 48, 1, 2, 3, 4, 5, 6, 7, 8, 49, 50, 51, 52, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48}

	if !bytes.Equal(res, sample) {
		t.Error("ISO Encode error!")
	}
}

func TestMessage(t *testing.T) {
	type TestIso struct {
		TestISO
		AB *Llnumeric `field:"ab" length:"19"`
	}

	iso := Message{"", ASCII, true, false, TestIso{*newDataIso(), NewLlnumeric("")}}

	input := []byte{48, 49, 48, 48, 114, 60, 36, 129, 40, 224, 152, 0, 49, 54, 52, 50, 55, 54, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 55, 55, 55, 48, 48, 48, 55, 48, 49, 49, 49, 49, 56, 52, 52, 48, 48, 48, 49, 50, 51, 49, 51, 49, 56, 52, 52, 48, 55, 48, 49, 49, 57, 48, 50, 6, 67, 57, 48, 49, 48, 50, 48, 54, 49, 50, 51, 52, 53, 54, 51, 55, 52, 50, 55, 54, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 53, 61, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 57, 56, 55, 54, 53, 52, 51, 50, 49, 48, 48, 49, 48, 48, 48, 48, 48, 51, 50, 49, 49, 50, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 51, 52, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 84, 101, 115, 116, 32, 116, 101, 120, 116, 100, 48, 1, 2, 3, 4, 5, 6, 7, 8, 49, 50, 51, 52, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48}

	err := iso.Load(input)

	assert.EqualError(t, err, "Critical error:value of field must be numeric")

	type TestIso2 struct {
		F2 *Llnumeric `field:"2" length:"19"`
	}

	iso = Message{"", ASCII, true, false, TestIso2{}}

	err = iso.Load(input)

	assert.EqualError(t, err, "field 2 not defined")

}

func TestMTIError(t *testing.T) {
	data := &TestISO{
		F2: NewLlnumeric("4276555555555555"),
	}

	iso := Message{"01000", ASCII, true, false, data}

	_, err := iso.Bytes()

	assert.EqualError(t, err, "MTI is invalid")

	iso.Mti = "abca"

	_, err = iso.Bytes()

	assert.EqualError(t, err, "MTI is invalid")

	iso.Mti = ""

	_, err = iso.Bytes()

	assert.EqualError(t, err, "MTI is required")

	iso = Message{"0100", BCD, true, false, data}

	res, err := iso.Bytes()

	assert.Empty(t, err)

	iso = Message{"", BCD, true, false, data}

	err = iso.Load(res[0:1])

	assert.EqualError(t, err, "bad MTI raw data")

	iso.Mti = "abca"

	_, err = iso.Bytes()

	assert.EqualError(t, err, "MTI is invalid")

}

func TestMTIBCD(t *testing.T) {
	data := &TestISO{
		F2: NewLlnumeric("4276555555555555"),
	}

	iso := Message{"0100", BCD, true, false, data}

	res, err := iso.Bytes()

	assert.Empty(t, err)

	iso2 := Message{"0100", BCD, true, false, data}

	err = iso2.Load(res)

	assert.Empty(t, err)

	assert.Equal(t, iso, iso2)
}

func TestParseFieldsErrors(t *testing.T) {
	type test1 struct {
		F2 *Llnumeric `field:"abc" length:"19"`
	}

	data1 := &test1{
		F2: NewLlnumeric("4276555555555555"),
	}

	iso := Message{"0100", BCD, true, false, data1}

	_, err := iso.Bytes()

	assert.EqualError(t, err, "Critical error:value of field must be numeric")

	type test2 struct {
		F2 *Llnumeric `field:"2" length:"abc"`
	}

	data2 := &test2{
		F2: NewLlnumeric("4276555555555555"),
	}

	iso = Message{"0100", BCD, true, false, data2}

	_, err = iso.Bytes()

	assert.EqualError(t, err, "Critical error:value of length must be numeric")

	type test3 struct {
		F2 string `field:"2" length:"2"`
	}

	data3 := &test3{
		F2: string("123abc"),
	}

	iso = Message{"0100", BCD, true, false, data3}

	_, err = iso.Bytes()

	assert.EqualError(t, err, "Critical error:field must be Iso8583Type")

	iso = Message{"0100", BCD, true, false, nil}

	_, err = iso.Bytes()

	assert.EqualError(t, err, "Critical error:data must be a struct")
}

// newDataIso creates DataIso
func newDataIso() *TestISO {
	return &TestISO{
		F2:   NewLlnumeric(""),
		F3:   NewNumeric(""),
		F4:   NewNumeric(""),
		F7:   NewNumeric(""),
		F11:  NewNumeric(""),
		F12:  NewNumeric(""),
		F13:  NewNumeric(""),
		F14:  NewNumeric(""),
		F19:  NewNumeric(""),
		F22:  NewNumeric(""),
		F25:  NewNumeric(""),
		F32:  NewLlnumeric(""),
		F35:  NewLlnumeric(""),
		F37:  NewAlphanumeric(""),
		F39:  NewAlphanumeric(""),
		F41:  NewAlphanumeric(""),
		F42:  NewAlphanumeric(""),
		F43:  NewAlphanumeric(""),
		F49:  NewNumeric(""),
		F52:  NewBinary(nil),
		F53:  NewNumeric(""),
		F120: NewLllnumeric(""),
	}
}

func TestWindows1252(t *testing.T) {
	type testIso struct {
		F2 *Llvar        `field:"2" length:"10" encode:"ascii"`
		F3 *Lllvar       `field:"3" length:"999" encode:"ascii"`
		F4 *Alphanumeric `field:"4" length:"10" encode:"ascii"`
		F5 *L8var        `field:"5" length:"99999999" encode:"ascii"`
	}

	data := &testIso{
		F2: NewLlvar([]byte("garçon!")),
		F3: NewLllvar([]byte("coração")),
		F4: NewAlphanumeric("solução"),
		F5: NewL8var([]byte("bota mais feijão ai meu irmão")),
	}
	iso := Message{
		Mti:          "0800",
		MtiEncode:    ASCII,
		SecondBitmap: true,
		ASCIIBitmap:  true,
		Data:         data,
	}

	result, err := iso.Bytes()
	assert.NoError(t, err)
	expected := "0800F800000000000000000000000000000007gar\xe7on!007cora\xe7\xe3o   solu\xe7\xe3o00000029bota mais feij\xe3o ai meu irm\xe3o"
	assert.Equal(t, expected, fmt.Sprintf("%s", result))

	emptyData := &testIso{
		F2: NewLlvar([]byte("")),
		F3: NewLllvar([]byte("")),
		F4: NewAlphanumeric(""),
		F5: NewL8var([]byte("")),
	}
	iso = Message{
		Mti:          "",
		MtiEncode:    ASCII,
		SecondBitmap: true,
		ASCIIBitmap:  true,
		Data:         emptyData,
	}
	err = iso.Load(result)
	assert.NoError(t, err)
	resultFields := iso.Data.(*testIso)
	assert.Equal(t, resultFields.F2.Value, []byte("gar\xe7on!"))
	assert.Equal(t, resultFields.F3.Value, []byte("cora\xe7\xe3o"))
	assert.Equal(t, resultFields.F4.Value, "   solu\xe7\xe3o")
	assert.Equal(t, resultFields.F5.Value, []byte("bota mais feij\xe3o ai meu irm\xe3o"))
}
