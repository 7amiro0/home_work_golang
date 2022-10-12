package hw09structvalidator

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
)

var (
	ErrorInvalidArg       = errors.New("invalid arg %v in param %v")
	ErrorInvalidParam     = errors.New("param does not exist: %v")
	ErrorInvalidInputData = errors.New("invalid input data")

	ErrorIn     = errors.New("%v not in multiplicity %v")
	ErrorRegexp = errors.New("%v not equal regexp %v")
	ErrorLen    = errors.New("len %v not equal %v")

	ErrorMin = errors.New("%v smaller then min %v")
	ErrorMax = errors.New("%v bigger then max %v")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (sliceErrors *ValidationErrors) append(name string, errs ...error) {
	for _, err := range errs {
		*sliceErrors = append(*sliceErrors, ValidationError{
			Field: name,
			Err:   err,
		})
	}
}

func (sliceErrors ValidationErrors) Error() string {
	var err string
	for _, value := range sliceErrors {
		err += fmt.Sprintf("%v\n", value.Err.Error())
	}
	return err
}

func checkType(field reflect.Value, validate, fieldName string) (err error) {
	args := strings.Split(validate, ":")

	switch field.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		err = parserInt(field, args[0], args[1])
	case reflect.String:
		err = parserString(field, args[0], args[1])
	case reflect.Slice:
		err = validateSlice(validate, field, fieldName)
	case reflect.Struct:
		if validate == "nested" {
			err = validateStruct(field)
		}
	default:
		err = fmt.Errorf(ErrorInvalidInputData.Error(), field)
	}

	return err
}

func manager(field reflect.Value, tag, fieldName string) []error {
	sliceErrors := make([]error, 0)
	validate := strings.Split(tag, "|")

	for indexTag := 0; indexTag < len(validate); indexTag++ {
		err := checkType(field, validate[indexTag], fieldName)
		log.Println("CURRENT ERROR", err)
		if err != nil {
			sliceErrors = append(sliceErrors, err)
		}
	}
	log.Println("ALL ERROR IN MANAGER", sliceErrors)

	return sliceErrors
}

func Validate(v interface{}) (finaleError error) {
	value := reflect.ValueOf(v)

	if value.Type().Kind() != reflect.Struct {
		finaleError = fmt.Errorf(ErrorInvalidInputData.Error())
	} else {
		finaleError = validateStruct(value)
	}

	log.Println("FINAL ERROR", finaleError)

	return finaleError
}
