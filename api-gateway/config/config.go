package config

import (
	"log"
	"os"
	"reflect"
	"strconv"
)

type Config struct {
	Port             int    `env:"PORT"                      envDefault:"9090"`
	UserServiceAddr  string `env:"USER_SERVICE_ADDR"         envDefault:"localhost:50051"` // TODO: Change to proper address
	IntersectionAddr string `env:"INTERSECTION_SERVICE_ADDR" envDefault:"localhost:50052"` // TODO: Change to proper address
}

func Load() *Config {
	cfg := &Config{}
	t := reflect.TypeOf(*cfg)
	v := reflect.ValueOf(cfg).Elem()

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
				log.Printf("Invalid bool for field %s: %s (defaulting to %s)", envName, valStr, defaultVal)
				parsed, _ = strconv.ParseBool(defaultVal)
			}
			v.Field(i).SetBool(parsed)
		default:
			log.Fatalf("unsupported field type: %s", field.Type)
		}
	}

	return cfg
}
