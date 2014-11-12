package viron_test

import (
	"bytes"
	"log"
	"strings"

	"github.com/ryanmoran/viron"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Print", func() {
	It("prints the configuration object to the logger", func() {
		buffer := bytes.NewBuffer([]byte{})
		logger := log.New(buffer, "", 0)

		viron.Print(Environment{
			Int32: int32(16),
		}, logger)

		lines := strings.Split(buffer.String(), "\n")
		Expect(lines).To(Equal([]string{
			"Bool      -> false",
			"String    -> ",
			"Int       -> 0",
			"Int8      -> 0",
			"Int16     -> 0",
			"Int32     -> 16",
			"Int64     -> 0",
			"Uint      -> 0",
			"Uint8     -> 0",
			"Uint16    -> 0",
			"Uint32    -> 0",
			"Uint64    -> 0",
			"Uintptr   -> 0",
			"Float32   -> 0",
			"Float64   -> 0",
			"JSON      -> {Space: Point:{R:0 G:0 B:0}}",
			"ByteSlice -> []",
			"NonTagged -> {}",
			"Required  -> ",
			"Default   -> ",
			"",
		}))
	})
})
