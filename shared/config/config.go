package config

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
)

func Load(cfg any) error {
	p := reflect.ValueOf(cfg)
	if p.Kind() != reflect.Ptr || p.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("Load expects a pointer to a struct")
	}

	v := p.Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		envName := field.Tag.Get("env")
		defaultVal := field.Tag.Get("envDefault")
		valStr := os.Getenv(envName)
		if valStr == "" {
			valStr = defaultVal
		}

		switch field.Type.Kind() {
		case reflect.String:
			v.Field(i).SetString(valStr)

		case reflect.Int:
			parsed, err := strconv.Atoi(valStr)
			if err != nil {
				log.Printf("Invalid int for %s: %s (defaulting to %s)", envName, valStr, defaultVal)
				parsed, _ = strconv.Atoi(defaultVal)
			}
			v.Field(i).SetInt(int64(parsed))

		case reflect.Bool:
			parsed, err := strconv.ParseBool(valStr)
			if err != nil {
				log.Printf(
					"Invalid bool for field %s: %s (defaulting to %s)",
					envName,
					valStr,
					defaultVal,
				)
				parsed, _ = strconv.ParseBool(defaultVal)
			}
			v.Field(i).SetBool(parsed)

		case reflect.Float32,
			reflect.Float64, reflect.Uint, reflect.Uint8, reflect.Slice,
			reflect.Map, reflect.Struct, reflect.Ptr, reflect.Interface,
			reflect.Array, reflect.Chan, reflect.Func, reflect.Invalid,
			reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
			reflect.Complex64, reflect.Complex128, reflect.UnsafePointer:
			return fmt.Errorf("unsupported field type: %s", field.Type)

		default:
			return fmt.Errorf("unsupported field type: %s", field.Type)
		}
	}

	return nil
}
