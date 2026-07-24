package main

import (
	"reflect"
	"testing"
)

func TestNewAppServicesInitializesEveryService(t *testing.T) {
	services := newAppServices(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)

	value := reflect.ValueOf(services)
	serviceType := value.Type()

	for index := range value.NumField() {
		field := value.Field(index)
		if field.Kind() == reflect.Pointer && field.IsNil() {
			t.Fatalf("expected %s to be initialized", serviceType.Field(index).Name)
		}
	}
}
