package saft_test

import (
	"github.com/johan-bolmsjo/saft"
	"testing"
)

// Go stdlib functions are used to convert strings so extensive testing is not needed.

func TestString_Bool(t *testing.T) {
	const testString = "true"
	v, err := (&saft.String{V: testString}).Bool()
	if v != true || err != nil {
		t.Fatalf(`%q.Bool() = (%v, %q); want (true, nil)`, testString, v, errorString(err))
	}

	v, err = (&saft.String{V: "x"}).Bool()
	errStr, wantErr := errorString(err), "0:0: strconv.ParseBool: parsing \"x\": invalid syntax"
	if v != false || errStr != wantErr {
		t.Fatalf(`"x".Bool() = (%v, %q); want (false, %q)`, v, errStr, wantErr)
	}
}

func TestString_Int32(t *testing.T) {
	const testString = "42"
	v, err := (&saft.String{V: testString}).Int32()
	if v != 42 || err != nil {
		t.Fatalf(`%q.Int32() = (%v, %q); want (true, nil)`, testString, v, errorString(err))
	}

	v, err = (&saft.String{V: "x"}).Int32()
	errStr, wantErr := errorString(err), "0:0: strconv.ParseInt: parsing \"x\": invalid syntax"
	if v != 0 || errStr != wantErr {
		t.Fatalf(`"x".Int32() = (%v, %q); want (false, %q)`, v, errStr, wantErr)
	}
}

func TestString_Uint32(t *testing.T) {
	const testString = "42"
	v, err := (&saft.String{V: testString}).Uint32()
	if v != 42 || err != nil {
		t.Fatalf(`%q.Uint32() = (%v, %q); want (true, nil)`, testString, v, errorString(err))
	}

	v, err = (&saft.String{V: "x"}).Uint32()
	errStr, wantErr := errorString(err), "0:0: strconv.ParseUint: parsing \"x\": invalid syntax"
	if v != 0 || errStr != wantErr {
		t.Fatalf(`"x".Uint32() = (%v, %q); want (false, %q)`, v, errStr, wantErr)
	}
}

func TestString_Int64(t *testing.T) {
	const testString = "42"
	v, err := (&saft.String{V: testString}).Int64()
	if v != 42 || err != nil {
		t.Fatalf(`%q.Int64() = (%v, %q); want (true, nil)`, testString, v, errorString(err))
	}

	v, err = (&saft.String{V: "x"}).Int64()
	errStr, wantErr := errorString(err), "0:0: strconv.ParseInt: parsing \"x\": invalid syntax"
	if v != 0 || errStr != wantErr {
		t.Fatalf(`"x".Int64() = (%v, %q); want (false, %q)`, v, errStr, wantErr)
	}
}

func TestString_Uint64(t *testing.T) {
	const testString = "42"
	v, err := (&saft.String{V: testString}).Uint64()
	if v != 42 || err != nil {
		t.Fatalf(`%q.Uint64() = (%v, %q); want (true, nil)`, testString, v, errorString(err))
	}

	v, err = (&saft.String{V: "x"}).Uint64()
	errStr, wantErr := errorString(err), "0:0: strconv.ParseUint: parsing \"x\": invalid syntax"
	if v != 0 || errStr != wantErr {
		t.Fatalf(`"x".Uint64() = (%v, %q); want (false, %q)`, v, errStr, wantErr)
	}
}

func TestString_Float32(t *testing.T) {
	const testString = "42.0"
	v, err := (&saft.String{V: testString}).Float32()
	if v != 42 || err != nil {
		t.Fatalf(`%q.Float32() = (%v, %q); want (true, nil)`, testString, v, errorString(err))
	}

	v, err = (&saft.String{V: "x"}).Float32()
	errStr, wantErr := errorString(err), "0:0: strconv.ParseFloat: parsing \"x\": invalid syntax"
	if v != 0 || errStr != wantErr {
		t.Fatalf(`"x".Float32() = (%v, %q); want (false, %q)`, v, errStr, wantErr)
	}
}

func TestString_Float64(t *testing.T) {
	const testString = "42.0"
	v, err := (&saft.String{V: testString}).Float64()
	if v != 42 || err != nil {
		t.Fatalf(`%q.Float64() = (%v, %q); want (true, nil)`, testString, v, errorString(err))
	}

	v, err = (&saft.String{V: "x"}).Float64()
	errStr, wantErr := errorString(err), "0:0: strconv.ParseFloat: parsing \"x\": invalid syntax"
	if v != 0 || errStr != wantErr {
		t.Fatalf(`"x".Float64() = (%v, %q); want (false, %q)`, v, errStr, wantErr)
	}
}

func TestString_CIDR(t *testing.T) {
	const testString = "127.0.0.1/8"
	_, _, err := (&saft.String{V: testString}).CIDR()
	if err != nil {
		t.Fatalf(`%q.CIDR() = (_, _, %q); want (_, _, nil)`, testString, errorString(err))
	}

	_, _, err = (&saft.String{V: "x"}).CIDR()
	errStr, wantErr := errorString(err), "0:0: invalid CIDR address: x"
	if errStr != wantErr {
		t.Fatalf(`"x".CIDR() = (_, _, %q); want (_, _, %q)`, errStr, wantErr)
	}
}

func TestString_IP(t *testing.T) {
	const testString = "127.0.0.1"
	_, err := (&saft.String{V: testString}).IP()
	if err != nil {
		t.Fatalf(`%q.IP() = (_, %q); want (_, nil)`, testString, errorString(err))
	}

	_, err = (&saft.String{V: "x"}).IP()
	errStr, wantErr := errorString(err), "0:0: invalid IP address: x"
	if errStr != wantErr {
		t.Fatalf(`"x".IP() = (_, %q); want (_, %q)`, errStr, wantErr)
	}
}

func TestString_MAC(t *testing.T) {
	const testString = "01:02:03:04:05:06"
	_, err := (&saft.String{V: testString}).MAC()
	if err != nil {
		t.Fatalf(`%q.MAC() = (_, %q); want (_, nil)`, testString, errorString(err))
	}

	_, err = (&saft.String{V: "x"}).MAC()
	errStr, wantErr := errorString(err), "0:0: invalid MAC address: x"
	if errStr != wantErr {
		t.Fatalf(`"x".MAC() = (_, %q); want (_, %q)`, errStr, wantErr)
	}
}
