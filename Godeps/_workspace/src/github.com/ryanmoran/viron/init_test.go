package viron_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestViron(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Viron Suite")
}

const (
	MaxUint    = ^uint(0)
	MaxUint8   = ^uint8(0)
	MaxUint16  = ^uint16(0)
	MaxUint32  = ^uint32(0)
	MaxUint64  = ^uint64(0)
	MaxUintptr = ^uintptr(0)
	MaxInt     = int(MaxUint >> 1)
	MaxInt8    = int8(MaxUint8 >> 1)
	MaxInt16   = int16(MaxUint16 >> 1)
	MaxInt32   = int32(MaxUint32 >> 1)
	MaxInt64   = int64(MaxUint64 >> 1)
)

type Environment struct {
	Bool    bool    `env:"BOOL"`
	String  string  `env:"STRING"`
	Int     int     `env:"INT"`
	Int8    int8    `env:"INT8"`
	Int16   int16   `env:"INT16"`
	Int32   int32   `env:"INT32"`
	Int64   int64   `env:"INT64"`
	Uint    uint    `env:"UINT"`
	Uint8   uint8   `env:"UINT8"`
	Uint16  uint16  `env:"UINT16"`
	Uint32  uint32  `env:"UINT32"`
	Uint64  uint64  `env:"UINT64"`
	Uintptr uintptr `env:"UINTPTR"`
	Float32 float32 `env:"FLOAT32"`
	Float64 float64 `env:"FLOAT64"`
	JSON    struct {
		Space string
		Point struct {
			R, G, B int
		}
	} `env:"JSON"`
	ByteSlice  []byte `env:"BYTESLICE"`
	unexported string `env:"UNEXPORTED"`
	NonTagged  struct{}
	Required   string `env:"REQUIRED" env-required:"true"`
	Default    string `env:"DEFAULT" env-default:"banana"`
}
