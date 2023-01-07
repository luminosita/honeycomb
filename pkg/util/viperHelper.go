package util

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

const TAG_NAME = "mapstructure"

type ValueRetriever = func(string) string

func OverrideConfig(vr ValueRetriever, overrideConfigItems map[string]string, c any) error {
	if vr == nil || overrideConfigItems == nil || c == nil {
		//TODO: Externalize
		return errors.New(fmt.Sprintf("Bad request: %+v, %+v", overrideConfigItems, c))
	}

	t := reflect.TypeOf(c).Elem()
	s := reflect.ValueOf(c).Elem()

	for k, v := range overrideConfigItems {
		newValue := vr(k)

		sp := strings.Split(k, ".")
		tagName := sp[len(sp)-1]

		for i := 0; i < t.NumField(); i++ {
			tv, ok := t.Field(i).Tag.Lookup(TAG_NAME)

			if ok && tv == tagName {
				f := s.FieldByName(t.Field(i).Name)
				if f.Kind() == reflect.String && f.CanSet() {
					f.SetString(newValue)
				} else {
					//TODO: Externalize
					fmt.Printf("Wrong config field to override: %s", v)
				}
			}
		}
	}

	return nil
}
