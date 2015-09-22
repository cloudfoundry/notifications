package viron_test

import (
	"fmt"
	"math"
	"os"
	"reflect"

	"github.com/ryanmoran/viron"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Parse", func() {
	var env *Environment

	BeforeEach(func() {
		env = &Environment{}
		t := reflect.TypeOf(env).Elem()
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			name := field.Tag.Get("env")
			os.Setenv(name, "")
		}
		os.Setenv("REQUIRED", "banana")
	})

	Context("bool values", func() {
		It("parses valid bool values", func() {
			validTrues := []string{"1", "t", "T", "true", "TRUE", "True"}
			validFalses := []string{"0", "f", "F", "false", "FALSE", "False"}

			for _, value := range validTrues {
				os.Setenv("BOOL", value)
				err := viron.Parse(env)
				Expect(err).NotTo(HaveOccurred())
				Expect(env.Bool).To(BeTrue())
			}

			for _, value := range validFalses {
				os.Setenv("BOOL", value)
				err := viron.Parse(env)
				Expect(err).NotTo(HaveOccurred())
				Expect(env.Bool).To(BeFalse())
			}
		})

		It("returns an error when the bool value cannot be parsed", func() {
			os.Setenv("BOOL", "banana")

			err := viron.Parse(env)
			Expect(err).To(Equal(viron.ParseError{
				Name:  "BOOL",
				Value: "banana",
				Kind:  "bool",
			}))
		})
	})

	Context("string values", func() {
		It("parses string values", func() {
			os.Setenv("STRING", "banana")

			err := viron.Parse(env)
			Expect(err).NotTo(HaveOccurred())
			Expect(env.String).To(Equal("banana"))
		})
	})

	Context("integer values", func() {
		It("parses int values", func() {
			os.Setenv("INT", fmt.Sprintf("%d", MaxInt))

			err := viron.Parse(env)
			Expect(err).NotTo(HaveOccurred())
			Expect(env.Int).To(Equal(MaxInt))
		})

		It("parses int8 values", func() {
			os.Setenv("INT8", fmt.Sprintf("%d", MaxInt8))

			err := viron.Parse(env)
			Expect(err).NotTo(HaveOccurred())
			Expect(env.Int8).To(Equal(MaxInt8))
		})

		It("parses int16 values", func() {
			os.Setenv("INT16", fmt.Sprintf("%d", MaxInt16))

			err := viron.Parse(env)
			Expect(err).NotTo(HaveOccurred())
			Expect(env.Int16).To(Equal(MaxInt16))
		})

		It("parses int32 values", func() {
			os.Setenv("INT32", fmt.Sprintf("%d", MaxInt32))

			err := viron.Parse(env)
			Expect(err).NotTo(HaveOccurred())
			Expect(env.Int32).To(Equal(MaxInt32))
		})

		It("parses int64 values", func() {
			os.Setenv("INT64", fmt.Sprintf("%d", MaxInt64))

			err := viron.Parse(env)
			Expect(err).NotTo(HaveOccurred())
			Expect(env.Int64).To(Equal(MaxInt64))
		})

		It("returns an error when the int value cannot be parsed", func() {
			os.Setenv("INT16", "banana")

			err := viron.Parse(env)
			Expect(err).To(Equal(viron.ParseError{
				Name:  "INT16",
				Value: "banana",
				Kind:  "int16",
			}))
		})
	})

	Context("unsigned integer values", func() {
		It("parses uint values", func() {
			os.Setenv("UINT", fmt.Sprintf("%d", MaxUint))

			err := viron.Parse(env)
			Expect(err).NotTo(HaveOccurred())
			Expect(env.Uint).To(Equal(MaxUint))
		})

		It("parses uint8 values", func() {
			os.Setenv("UINT8", fmt.Sprintf("%d", MaxUint8))

			err := viron.Parse(env)
			Expect(err).NotTo(HaveOccurred())
			Expect(env.Uint8).To(Equal(MaxUint8))
		})

		It("parses uint16 values", func() {
			os.Setenv("UINT16", fmt.Sprintf("%d", MaxUint16))

			err := viron.Parse(env)
			Expect(err).NotTo(HaveOccurred())
			Expect(env.Uint16).To(Equal(MaxUint16))
		})

		It("parses uint32 values", func() {
			os.Setenv("UINT32", fmt.Sprintf("%d", MaxUint32))

			err := viron.Parse(env)
			Expect(err).NotTo(HaveOccurred())
			Expect(env.Uint32).To(Equal(MaxUint32))
		})

		It("parses uint64 values", func() {
			os.Setenv("UINT64", fmt.Sprintf("%d", MaxUint64))

			err := viron.Parse(env)
			Expect(err).NotTo(HaveOccurred())
			Expect(env.Uint64).To(Equal(MaxUint64))
		})

		It("parses uintptr values", func() {
			os.Setenv("UINTPTR", fmt.Sprintf("%d", MaxUintptr))

			err := viron.Parse(env)
			Expect(err).NotTo(HaveOccurred())
			Expect(env.Uintptr).To(Equal(MaxUintptr))
		})

		It("returns an error when the uint value cannot be parsed", func() {
			os.Setenv("UINT32", "banana")

			err := viron.Parse(env)
			Expect(err).To(Equal(viron.ParseError{
				Name:  "UINT32",
				Value: "banana",
				Kind:  "uint32",
			}))
		})
	})

	Context("float values", func() {
		It("parses float32 values", func() {
			os.Setenv("FLOAT32", fmt.Sprintf("%f", math.MaxFloat32))

			err := viron.Parse(env)
			Expect(err).NotTo(HaveOccurred())
			Expect(env.Float32).To(Equal(float32(math.MaxFloat32)))
		})

		It("parses float64 values", func() {
			os.Setenv("FLOAT64", fmt.Sprintf("%f", math.MaxFloat64))

			err := viron.Parse(env)
			Expect(err).NotTo(HaveOccurred())
			Expect(env.Float64).To(Equal(math.MaxFloat64))
		})

		It("returns an error when the float value cannot be parsed", func() {
			os.Setenv("FLOAT64", "banana")

			err := viron.Parse(env)
			Expect(err).To(Equal(viron.ParseError{
				Name:  "FLOAT64",
				Value: "banana",
				Kind:  "float64",
			}))
		})
	})

	Context("struct fields", func() {
		It("parses the variable as JSON", func() {
			os.Setenv("JSON", `{"space":"RGB", "point":{"r":98, "g":218, "b":255}}`)

			err := viron.Parse(env)
			Expect(err).NotTo(HaveOccurred())

			Expect(env.JSON.Space).To(Equal("RGB"))
			Expect(env.JSON.Point.R).To(Equal(98))
			Expect(env.JSON.Point.G).To(Equal(218))
			Expect(env.JSON.Point.B).To(Equal(255))
		})

		It("returns an error when the JSON cannot be parsed", func() {
			os.Setenv("JSON", "banana")

			err := viron.Parse(env)
			Expect(err).To(Equal(viron.ParseError{
				Name:  "JSON",
				Value: "banana",
				Kind:  "struct",
			}))
		})
	})

	Context("byte slice fields", func() {
		It("parses the variable from a string to a byte slice", func() {
			os.Setenv("BYTESLICE", "SOMETHING")

			err := viron.Parse(env)
			Expect(err).NotTo(HaveOccurred())
			Expect(env.ByteSlice).To(Equal([]byte("SOMETHING")))
		})
	})

	Context("when the environment passed-in is not a non-zero pointer", func() {
		It("returns an InvalidArgumentError", func() {
			err := viron.Parse(7)
			Expect(err).To(Equal(viron.InvalidArgumentError{
				Value: 7,
			}))
		})

		It("returns an InvalidArgumentError", func() {
			var actualEnv *Environment
			var expectedEnv *Environment
			err := viron.Parse(actualEnv)
			Expect(err).To(Equal(viron.InvalidArgumentError{
				Value: expectedEnv,
			}))
		})
	})

	Context("non-required values", func() {
		It("leaves unassigned", func() {
			emptyEnv := Environment{
				Required: "banana",
				Default:  "banana",
			}

			err := viron.Parse(env)
			Expect(err).NotTo(HaveOccurred())
			Expect(*env).To(Equal(emptyEnv))
		})
	})

	Context("unexported fields", func() {
		It("ignores them", func() {
			os.Setenv("UNEXPORTED", "banana")

			err := viron.Parse(env)
			Expect(err).NotTo(HaveOccurred())
			Expect(env.unexported).To(Equal(""))
		})
	})

	Context("non-tagged fields", func() {
		It("ignores them", func() {
			err := viron.Parse(env)
			Expect(err).NotTo(HaveOccurred())
			Expect(env.NonTagged).To(Equal(struct{}{}))
		})
	})

	Context("required fields", func() {
		It("returns an error if the variable is missing", func() {
			os.Setenv("REQUIRED", "")

			err := viron.Parse(env)
			Expect(err).To(Equal(viron.RequiredFieldError{
				Name: "REQUIRED",
			}))
		})
	})

	Context("default fields", func() {
		It("uses the default value if none can be found", func() {
			os.Setenv("DEFAULT", "")

			err := viron.Parse(env)
			Expect(err).NotTo(HaveOccurred())
			Expect(env.Default).To(Equal("banana"))
		})
	})
})
