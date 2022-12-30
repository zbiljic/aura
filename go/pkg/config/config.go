package config

import (
	"errors"
	"os"
	"reflect"
	"strings"

	"github.com/creasty/defaults"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

const appNameKey = "appName"

const (
	jsonTagName  = "json"
	envTagName   = "env"
	aliasTagName = "alias"
)

var envKeyReplacer = strings.NewReplacer(".", "_")

// Load loads configuration into 'config' variable.
func Load(filename, appName string, config interface{}, opts ...viper.DecoderConfigOption) error {
	var v = viper.New()

	if appName != "" {
		v.SetDefault(appNameKey, appName)
		v.SetEnvPrefix(appName)
	}

	v.SetEnvKeyReplacer(envKeyReplacer)
	v.AutomaticEnv()

	if filename != "" {
		v.SetConfigFile(filename)

		if err := v.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				// config file not found; ignore error
			} else {
				return err
			}
		}
	} else {
		err := godotenv.Load()
		if err != nil {
			// ignore error if .env file does not exists
			if !os.IsNotExist(err) {
				return err
			}
		}
	}

	if err := structBindEnv(config, v); err != nil {
		return err
	}

	if err := defaults.Set(config); err != nil {
		return err
	}

	// include default decoder configuration
	opts = append([]viper.DecoderConfigOption{
		func(dc *mapstructure.DecoderConfig) {
			dc.Squash = true
			dc.TagName = jsonTagName
		},
	}, opts...)

	if err := v.Unmarshal(config, opts...); err != nil {
		return err
	}

	if err := ValidateConfig(config); err != nil {
		return err
	}

	return nil
}

// structBindEnv binds environment variables into the Viper instance for
// the given struct.
// It gets the pointer of a struct that is going to holds the variables.
func structBindEnv(structure interface{}, v *viper.Viper, prefix ...string) error {
	inputType := reflect.TypeOf(structure)
	if inputType != nil {
		if inputType.Kind() == reflect.Ptr {
			if inputType.Elem().Kind() == reflect.Struct {
				return bindStruct(reflect.ValueOf(structure).Elem(), v, prefix)
			}
		}
	}

	return errors.New("config: invalid structure")
}

// bindStruct binds a reflected struct fields path to Viper instance.
func bindStruct(s reflect.Value, v *viper.Viper, prefix []string) error {
	for i := 0; i < s.NumField(); i++ {
		fieldName := s.Type().Field(i).Name

		if t, exist := s.Type().Field(i).Tag.Lookup(jsonTagName); exist {
			fieldName = t
		}

		key := strings.Join(append(prefix, fieldName), ".")

		if t, exist := s.Type().Field(i).Tag.Lookup(aliasTagName); exist {
			v.RegisterAlias(key, t)
		}

		if t, exist := s.Type().Field(i).Tag.Lookup(envTagName); exist {
			v.BindEnv(key, t)
		} else if s.Type().Field(i).Type.Kind() == reflect.Struct {
			if s.Type().Field(i).Anonymous {
				// squash the field down if anonymous
				if err := bindStruct(s.Field(i), v, prefix); err != nil {
					return err
				}
			} else {
				if err := bindStruct(s.Field(i), v, append(prefix, fieldName)); err != nil {
					return err
				}
			}
		} else if s.Type().Field(i).Type.Kind() == reflect.Ptr {
			if !s.Field(i).IsZero() && s.Field(i).Elem().Type().Kind() == reflect.Struct {
				if err := bindStruct(s.Field(i).Elem(), v, prefix); err != nil {
					return err
				}
			}
		} else {
			v.BindEnv(key, strings.ToUpper(envKeyReplacer.Replace(key)))
		}
	}

	return nil
}

// ValidateConfig validates configuration based on provided tags.
// It returns first error it encounters.
func ValidateConfig(config interface{}) error {
	validate := validator.New()

	if err := validate.Struct(config); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return validationErrors[0]
	}

	return nil
}
