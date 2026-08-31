package crud

import (
	"errors"
	"reflect"
	"time"
)

func MapSchemaToModel[T any, S any](schema S) (T, error) {
	var model T
	modelVal := reflect.ValueOf(&model).Elem()
	schemaVal := reflect.ValueOf(schema)

	if schemaVal.Kind() == reflect.Ptr {
		if schemaVal.IsNil() {
			return model, errors.New("schema is nil")
		}
		schemaVal = schemaVal.Elem()
	}
	if modelVal.Kind() != reflect.Struct || schemaVal.Kind() != reflect.Struct {
		return model, errors.New("both schema and model must be structs")
	}

	modelType := modelVal.Type()
	for i := 0; i < schemaVal.NumField(); i++ {
		schemaField := schemaVal.Field(i)
		schemaFieldName := schemaVal.Type().Field(i).Name

		_, found := modelType.FieldByName(schemaFieldName)
		if !found {
			continue 
		}
		modelFieldVal := modelVal.FieldByName(schemaFieldName)
		if !modelFieldVal.CanSet() {
			continue
		}

		if err := setFieldValue(modelFieldVal, schemaField); err != nil {
			continue
		}
	}
	return model, nil
}

func setFieldValue(modelField reflect.Value, schemaField reflect.Value) error {
	if !schemaField.IsValid() {
		return errors.New("invalid schema field")
	}
	modelKind := modelField.Kind()
	schemaKind := schemaField.Kind()

	switch modelKind {
	case reflect.String:
		if schemaKind == reflect.String {
			modelField.SetString(schemaField.String())
			return nil
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if schemaKind == reflect.Int || schemaKind == reflect.Int64 {
			modelField.SetInt(schemaField.Int())
			return nil
		}
		if schemaKind == reflect.Float64 {
			modelField.SetInt(int64(schemaField.Float()))
			return nil
		}
		if schemaKind == reflect.Uint {
			modelField.SetInt(int64(schemaField.Uint()))
			return nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if schemaKind == reflect.Uint {
			modelField.SetUint(schemaField.Uint())
			return nil
		}
		if schemaKind == reflect.Int || schemaKind == reflect.Int64 {
			val := schemaField.Int()
			if val >= 0 {
				modelField.SetUint(uint64(val))
				return nil
			}
		}
	case reflect.Bool:
		if schemaKind == reflect.Bool {
			modelField.SetBool(schemaField.Bool())
			return nil
		}
	case reflect.Struct:
		if modelField.Type() == reflect.TypeOf(time.Time{}) {
			if schemaField.Type() == reflect.TypeOf(time.Time{}) {
				modelField.Set(schemaField)
				return nil
			}
		}
	}
	return errors.New("type mismatch or unsupported conversion")
}

func GetNonZeroFields(schema interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	val := reflect.ValueOf(schema)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return result
		}
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return result
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldName := typ.Field(i).Name
		if !field.IsZero() {
			result[fieldName] = field.Interface()
		}
	}
	return result
}

func MapModelToResponse[R any, T any](model T) (R, error) {
	var resp R
	respVal := reflect.ValueOf(&resp).Elem()
	modelVal := reflect.ValueOf(model)

	if modelVal.Kind() == reflect.Ptr {
		if modelVal.IsNil() {
			return resp, errors.New("model is nil")
		}
		modelVal = modelVal.Elem()
	}
	if respVal.Kind() != reflect.Struct || modelVal.Kind() != reflect.Struct {
		return resp, errors.New("both model and response must be structs")
	}

	respType := respVal.Type()
	for i := 0; i < modelVal.NumField(); i++ {
		modelField := modelVal.Field(i)
		modelFieldName := modelVal.Type().Field(i).Name

		_, found := respType.FieldByName(modelFieldName)
		if !found {
			continue 
		}
		respFieldVal := respVal.FieldByName(modelFieldName)
		if !respFieldVal.CanSet() {
			continue
		}

		if err := setFieldValue(respFieldVal, modelField); err != nil {
			continue
		}
	}
	return resp, nil
}

func MapModelsToResponse[R any, T any](models []T) ([]R, error) {
	result := make([]R, len(models))
	for i, m := range models {
		r, err := MapModelToResponse[R](m)
		if err != nil {
			return nil, err
		}
		result[i] = r
	}
	return result, nil
}