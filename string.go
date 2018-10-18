package saft

import (
	"fmt"
	"github.com/johan-bolmsjo/errors"
	"net"
	"strconv"
)

// String is surprisingly a string.
type String struct {
	pos LexPos
	V   string // String value
}

// Pos returns positional information useful for context dependent error reporting.
func (s *String) Pos() LexPos {
	return s.pos
}

func (s *String) elemType() elemType {
	return elemTypeString
}

func wrapStrconvIntError(err error, s *String) error {
	// strconv error contain source string
	return errors.Wrap(err, s.pos.String())
}

func wrapStrconvFloatError(err error, s *String) error {
	// strconv error contain source string
	return errors.Wrap(err, s.pos.String())
}

// Parse string as a boolean.
// Returns the parsed value or an error containing positional information.
func (s *String) Bool() (v bool, err error) {
	if v, err = strconv.ParseBool(s.V); err != nil {
		// strconv error contain source string
		err = errors.Wrap(err, s.pos.String())
	}
	return
}

// Parse string as a signed 32 bit integer.
// Returns the parsed value or an error containing positional information.
func (s *String) Int32() (v int32, err error) {
	var t int64
	if t, err = strconv.ParseInt(s.V, 10, 32); err != nil {
		err = wrapStrconvIntError(err, s)
	}
	v = int32(t)
	return
}

// Parse string as an unsigned 32 bit integer.
// Returns the parsed value or an error containing positional information.
func (s *String) Uint32() (v uint32, err error) {
	var t uint64
	if t, err = strconv.ParseUint(s.V, 10, 32); err != nil {
		err = wrapStrconvIntError(err, s)
	}
	v = uint32(t)
	return
}

// Parse string as a signed 64 bit integer.
// Returns the parsed value or an error containing positional information.
func (s *String) Int64() (v int64, err error) {
	if v, err = strconv.ParseInt(s.V, 10, 64); err != nil {
		err = wrapStrconvIntError(err, s)
	}
	return
}

// Parse string as an unsigned 64 bit integer.
// Returns the parsed value or an error containing positional information.
func (s *String) Uint64() (v uint64, err error) {
	if v, err = strconv.ParseUint(s.V, 10, 64); err != nil {
		err = wrapStrconvIntError(err, s)
	}
	return
}

// Parse string as a 32 bit floating-point number.
// Returns the parsed value or an error containing positional information.
func (s *String) Float32() (v float32, err error) {
	var t float64
	if t, err = strconv.ParseFloat(s.V, 32); err != nil {
		err = wrapStrconvFloatError(err, s)
	}
	v = float32(t)
	return
}

// Parse string as a 64 bit floating-point number.
// Returns the parsed value or an error containing positional information.
func (s *String) Float64() (v float64, err error) {
	if v, err = strconv.ParseFloat(s.V, 64); err != nil {
		err = wrapStrconvFloatError(err, s)
	}
	return
}

// Parse string as CIDR notation IP address and prefix length.
// Returns the parsed values or an error containing positional information.
func (s *String) CIDR() (ip net.IP, ipnet *net.IPNet, err error) {
	if ip, ipnet, err = net.ParseCIDR(s.V); err != nil {
		// net parse error contain source string
		err = errors.Wrap(err, s.pos.String())
	}
	return
}

// Parse string as IP address.
// Returns the parsed value or an error containing positional information.
func (s *String) IP() (ip net.IP, err error) {
	if ip = net.ParseIP(s.V); ip == nil {
		err = fmt.Errorf("%s: invalid IP address: %s", &s.pos, s.V)
	}
	return
}

// Parse string as MAC address.
// Returns the parsed value or an error containing positional information.
func (s *String) MAC() (hw net.HardwareAddr, err error) {
	if hw, err = net.ParseMAC(s.V); err != nil {
		// override error to make it similar to IP and CIDR parsing errors.
		err = fmt.Errorf("%s: invalid MAC address: %s", &s.pos, s.V)
	}
	return
}
