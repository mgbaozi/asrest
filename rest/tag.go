package rest

import (
	"reflect"
	"strconv"
)

func getTag(index int, key string, model interface{}) string {
	model_type := reflect.ValueOf(model).Elem().Type()
	field := model_type.Field(index)
	return field.Tag.Get(key)
}

func getTags(key string, model interface{}) map[string]int {
	result := make(map[string]int)
	model_type := reflect.ValueOf(model).Elem().Type()
	for i := 0; i < model_type.NumField(); i++ {
		field := model_type.Field(i)
		if param := field.Tag.Get(key); len(param) > 0 {
			result[param] = i
		}
	}
	return result
}

func ConverseType(model interface{}, index int, value string) interface{} {
	model_type := reflect.ValueOf(model).Elem().Type()
	field := model_type.Field(index)
	//field_type := field.Type()
	var result interface{}
	switch reflect.New(field.Type).Elem().Interface().(type) {
	case string:
		result = value
	case int:
		result, _ = strconv.Atoi(value)
	case uint:
		result, _ = strconv.ParseUint(value, 10, 64)
	case int64:
		result, _ = strconv.ParseInt(value, 10, 64)
	case float64:
		result, _ = strconv.ParseFloat(value, 64)
	case bool:
		result, _ = strconv.ParseBool(value)
	default:
		result = nil
	}
	return result
}
