package viron

import (
	"log"
	"reflect"
	"strconv"
)

type printableField struct {
	Name  string
	Value interface{}
}

func Print(env interface{}, logger *log.Logger) {
	t := reflect.TypeOf(env)
	v := reflect.ValueOf(env)
	maxFieldNameLength := 0

	var fields []printableField

	for i := 0; i < v.NumField(); i++ {
		fieldType := t.Field(i)
		fieldValue := v.Field(i)
		if fieldValue.CanInterface() {
			if x := len(fieldType.Name); x > maxFieldNameLength {
				maxFieldNameLength = x
			}
			fields = append(fields, printableField{
				Name:  fieldType.Name,
				Value: fieldValue.Interface(),
			})
		}
	}

	for _, field := range fields {
		logger.Printf("%-"+strconv.Itoa(maxFieldNameLength)+"s -> %+v", field.Name, field.Value)
	}
}
