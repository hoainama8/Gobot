package main

import (
	"errors" // [THEM_MOI]
	"net/http" // [THEM_MOI]
	"reflect" // [THEM_MOI]
	"strings" // [THEM_MOI]
)

func (app *Application) decodePostForm(r *http.Request, dst any) error { // [THEM_MOI]
	if r.Method != http.MethodPost {
		return errors.New("invalid method")
	}
	if err := r.ParseForm(); err != nil {
		return err
	}

	val := reflect.ValueOf(dst)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return errors.New("dst must be a non-nil pointer")
	}

	elem := val.Elem()
	if elem.Kind() != reflect.Struct {
		return errors.New("dst must point to struct")
	}

	typ := elem.Type()
	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		sf := typ.Field(i)
		formKey := sf.Tag.Get("form")
		if formKey == "" || formKey == "-" {
			continue
		}
		if field.Kind() == reflect.String && field.CanSet() {
			field.SetString(strings.TrimSpace(r.Form.Get(formKey)))
		}
	}
	return nil
}
