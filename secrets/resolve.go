package secrets

import (
	"fmt"
	"reflect"
)

func ResolveAll(sm SecretManager, config interface{}) error {
	v := reflect.ValueOf(config)
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("secrets: config must be a pointer to struct")
	}
	return resolveRecursive(sm, v)
}

func resolveRecursive(sm SecretManager, v reflect.Value) error {
	if !v.IsValid() {
		return nil
	}

	switch v.Kind() {
	case reflect.Ptr:
		if v.IsNil() {
			return nil
		}
		return resolveRecursive(sm, v.Elem())
	case reflect.String:
		str := v.String()
		if IsReference(str) {
			resolved, err := sm.Resolve(str)
			if err != nil {
				return fmt.Errorf("secrets: resolve %s: %w", str, err)
			}
			v.SetString(resolved)
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			if field.CanSet() {
				if err := resolveRecursive(sm, field); err != nil {
					return err
				}
			}
		}
	case reflect.Map:
		if v.Type().Key().Kind() == reflect.String {
			for _, key := range v.MapKeys() {
				elem := v.MapIndex(key)
				if elem.Kind() == reflect.Ptr && elem.IsNil() {
					continue
				}
				newElem := reflect.New(elem.Type()).Elem()
				newElem.Set(elem)
				if err := resolveRecursive(sm, newElem); err != nil {
					return err
				}
				v.SetMapIndex(key, newElem)
			}
		}
	}
	return nil
}
