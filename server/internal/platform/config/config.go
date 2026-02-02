package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
)

type Config struct {
	// Database URI to connect to the database
	DatabaseURI string `config:"DATABASE_URI"`
	// Address of the http server
	Port string `config:"PORT"`
}

func Init() (*Config, error) {
	cfg := &Config{}

	v := reflect.ValueOf(cfg).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		ftype := t.Field(i)
		fval := v.Field(i)

		key := ftype.Tag.Get("config")
		if key == "" {
			continue
		}

		value := os.Getenv(key)
		if len(value) == 0 {
			return nil, fmt.Errorf("key %q must have a value", key)
		}

		if err := setField(fval, value); err != nil {
			return nil, err
		}
	}

	return cfg, nil
}

func setField(field reflect.Value, raw string) error {
	if !field.CanSet() {
		return nil
	}

	switch field.Kind() {
	case reflect.String:
		field.SetString(raw)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(raw, 10, field.Type().Bits())
		if err != nil {
			return err
		}
		field.SetInt(i)

	case reflect.Bool:
		b, err := strconv.ParseBool(raw)
		if err != nil {
			return err
		}
		field.SetBool(b)

	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(raw, field.Type().Bits())
		if err != nil {
			return err
		}
		field.SetFloat(f)

	default:
		return fmt.Errorf("unsupported kind: %s", field.Kind())
	}

	return nil
}
