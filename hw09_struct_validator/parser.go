package hw09structvalidator

import (
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func checkIn(options []string, fieldValue string) (isIn bool) {
	for _, value := range options {
		if value == fieldValue {
			isIn = true
			break
		}
	}

	return isIn
}

func parseTagIn(options string, fieldValue string) (err error) {
	log.Println(options, fieldValue)
	sliceOptions := strings.Split(options, ",")

	if !checkIn(sliceOptions, fieldValue) {
		err = fmt.Errorf(ErrorIn.Error(), fieldValue, sliceOptions)
	}

	return err
}

func checkLen(field reflect.Value, arg string) (err error) {
	lenString, err := strconv.Atoi(arg)
	if err != nil {
		err = fmt.Errorf(ErrorInvalidArg.Error(), arg, "len")
	} else if field.Len() != lenString {
		err = fmt.Errorf(ErrorLen.Error(), field.Len(), arg)
	}

	return err
}

func checkRegexp(field reflect.Value, arg string) (err error) {
	math, err := regexp.MatchString(arg, field.String())
	if err != nil {
		err = fmt.Errorf(ErrorInvalidArg.Error(), arg, "regexp")
	} else if !math {
		err = fmt.Errorf(ErrorRegexp.Error(), field, arg)
	}

	return err
}

func parserString(field reflect.Value, param, arg string) (err error) {
	switch param {
	case "len":
		err = checkLen(field, arg)
	case "regexp":
		err = checkRegexp(field, arg)
	case "in":
		err = parseTagIn(arg, field.String())
	default:
		err = fmt.Errorf(ErrorInvalidParam.Error(), param)
	}
	log.Println("Error in validate tag string", err)

	return err
}

func minMax(option string, bound int64, field reflect.Value) (err error) {
	if option == "min" {
		if bound > field.Int() {
			err = fmt.Errorf(ErrorMin.Error(), field, bound)
		}
	} else if option == "max" {
		if bound < field.Int() {
			err = fmt.Errorf(ErrorMax.Error(), field, bound)
		}
	}

	return err
}

func parserInt(field reflect.Value, param, arg string) (err error) {
	number, err := strconv.ParseInt(arg, 10, 64)
	if err != nil && param != "in" {
		return fmt.Errorf(ErrorInvalidArg.Error(), arg, param)
	}

	switch param {
	case "min", "max":
		err = minMax(param, number, field)
	case "in":
		err = parseTagIn(arg, strconv.Itoa(int(field.Int())))
	default:
		err = fmt.Errorf(ErrorInvalidParam.Error(), param)
	}
	log.Println("Error in validate tag int", err)
	return err
}

func validateSlice(validate string, field reflect.Value, fieldName string) (err error) {
	sliceErrors := make(ValidationErrors, 0)

	for index := 0; index < field.Len(); index++ {
		value := field.Index(index)
		errorList := manager(value, validate, fieldName)
		sliceErrors.append(fieldName, errorList...)
		log.Println(errorList, len(errorList))
	}

	log.Printf("ERROR SLICE 1 %q\n", sliceErrors)

	if len(sliceErrors) != 0 {
		err = sliceErrors
	}

	log.Printf("ERROR SLICE 2 %q\n", sliceErrors)

	return err
}

func validateStruct(structure reflect.Value) (resultError error) {
	sliceErrors := make(ValidationErrors, 0, structure.NumField())

	for indexField := 0; indexField < structure.NumField(); indexField++ {
		fieldValue := structure.Field(indexField)
		fieldType := structure.Type().Field(indexField)
		tag, ok := fieldType.Tag.Lookup("validate")

		if !fieldType.IsExported() || !ok {
			continue
		}

		manageErr := manager(fieldValue, tag, fieldType.Name)
		log.Println("ERROR IN MANAGER", sliceErrors)
		sliceErrors.append(fieldType.Name, manageErr...)
		log.Println("ALL ERRORS IN SLICE", sliceErrors)
	}

	if len(sliceErrors) != 0 {
		resultError = sliceErrors
	}
	log.Println("RESULT ERROR IN SLICE", resultError)

	return resultError
}
