package sip

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const EnvTag = "env"
const SipTag = "sip"

type Key string

func (key *Key) ToEnvVar() string {
	return strings.ReplaceAll(strings.ToUpper(string(*key)), ".", "_")
}

func Fill(i any, prefix string) error {
	if err := LoadFiles(i); err != nil {
		return err
	}

	flat := map[string]reflect.Value{}
	if err := flatten(reflect.ValueOf(i), prefix, flat); err != nil {
		return err
	}

	for tag, value := range flat {
		if err := setValue(Key(tag), value); err != nil {
			return err
		}
	}
	return nil
}

func flatten(v reflect.Value, prefix string, flat map[string]reflect.Value) error {
	if v.Kind() == reflect.Interface && v.IsNil() {
		return nil
	}

	if v.Kind() == reflect.Pointer || v.Kind() == reflect.Interface {
		e := v.Elem()
		if e.Kind() == reflect.Invalid && v.CanSet() {
			v.Set(reflect.New(v.Type().Elem()))
			return flatten(v, prefix, flat)
		} else if e.Kind() == reflect.Invalid {
			return nil
		} else {
			return flatten(e, prefix, flat)
		}
	}

	return flattenFields(v, prefix, flat)
}

func flattenFields(v reflect.Value, prefix string, flat map[string]reflect.Value) error {
	for i := 0; i < v.NumField(); i++ {
		value := v.Field(i)
		field := v.Type().Field(i)

		env, ok := field.Tag.Lookup(EnvTag)
		if ok && !IsNestedType(value) {
			flat[env] = value
		} else {
			sip := field.Tag.Get(SipTag)

			DuplicateTagIfNotExist(field, SipTag, "json")
			DuplicateTagIfNotExist(field, SipTag, "yaml")

			var newPrefix string
			if sip != "" {
				newPrefix = prefix + "." + sip
			}

			if IsNestedType(value) {
				return flatten(value, newPrefix, flat)
			} else if value.Kind() == reflect.Invalid {
				continue
			}

			flat[newPrefix] = value
		}
	}

	return nil
}

func setValue(key Key, value reflect.Value) error { // NOSONAR go:S3776
	for _, provider := range Providers {
		v, ok := provider(key)
		if !ok {
			continue // Provider does not provide a value for the key
		}

		switch value.Kind() {
		case reflect.String:
			value.SetString(v)
		case reflect.Bool:
			if b, err := strconv.ParseBool(v); err != nil {
				return err
			} else {
				value.SetBool(b)
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if i, err := strconv.ParseInt(v, 10, 64); err != nil {
				return err
			} else {
				value.SetInt(i)
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if i, err := strconv.ParseUint(v, 10, 64); err != nil {
				return err
			} else {
				value.SetUint(i)
			}
		case reflect.Float32, reflect.Float64:
			if f, err := strconv.ParseFloat(v, 64); err != nil {
				return err
			} else {
				value.SetFloat(f)
			}
		case reflect.Complex64, reflect.Complex128:
			if c, err := strconv.ParseComplex(v, 128); err != nil {
				return err
			} else {
				value.SetComplex(c)
			}
		default:
			return fmt.Errorf("unable to parse %s to type %s (don't know how)", v, value.Kind())
		}
	}

	return nil
}

func DuplicateTagIfNotExist(field reflect.StructField, src string, dst string) {
	_, ok := field.Tag.Lookup(dst)
	if ok {
		return
	}

	v, ok := field.Tag.Lookup(src)
	if !ok {
		return
	}

	field.Tag = reflect.StructTag(fmt.Sprintf("%s:%s=\"%s\"", field.Tag, dst, v))
}

func IsNestedType(value reflect.Value) bool {
	return value.Kind() == reflect.Pointer || value.Kind() == reflect.Struct || value.Kind() == reflect.Interface
}
