package utils

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"unicode"
)

type IValidator interface {
	Max(n int) RuleFunc
	Min(n int) RuleFunc
	Regexp(n *regexp.Regexp) RuleFunc
	Required() RuleSet
	Rules(rules ...RuleFunc) []RuleSet
	SetupValidator(data any, fields Fields) *Instance
	Validate() []ErrorMessage
}

type Instance struct {
	data   any
	fields Fields

	ErrorMessages []ErrorMessage
	Checker       bool
}

type RuleSet struct {
	Name         string
	RuleValue    any
	FieldValue   any
	FieldName    any
	MessageFunc  func(RuleSet) string
	ValidateFunc func(RuleSet) bool
}
type (
	RuleFunc     func() RuleSet
	Fields       map[string][]RuleSet
	ErrorMessage struct {
		FieldName string
		Err       error
	}
)

func NewValidator() IValidator {
	return &Instance{}
}

func (v *Instance) SetupValidator(data any, fields Fields) *Instance {
	return &Instance{
		data:   data,
		fields: fields,
	}
}

func (v *Instance) Rules(rules ...RuleFunc) []RuleSet {
	ruleSets := make([]RuleSet, len(rules))
	for i := 0; i < len(ruleSets); i++ {
		ruleSets[i] = rules[i]()
	}
	return ruleSets
}

func (v *Instance) Required() RuleSet {
	return RuleSet{
		Name: "required",
		MessageFunc: func(set RuleSet) string {
			return fmt.Sprintf("%s is a required field", set.FieldName)
		},
		ValidateFunc: func(rule RuleSet) bool {
			str, ok := rule.FieldValue.(string)
			if !ok {
				return false
			}
			return len(str) > 0
		},
	}
}

func (v *Instance) Max(n int) RuleFunc {
	return func() RuleSet {
		return RuleSet{
			Name:      "max",
			RuleValue: n,
			ValidateFunc: func(set RuleSet) bool {
				str, ok := set.FieldValue.(string)
				if !ok {
					return false
				}
				return len(str) <= n
			},
			MessageFunc: func(set RuleSet) string {
				return fmt.Sprintf("%s should be maximum %d characters long", set.FieldName, n)
			},
		}
	}
}

func (v *Instance) Min(n int) RuleFunc {
	return func() RuleSet {
		return RuleSet{
			Name:      "min",
			RuleValue: n,
			ValidateFunc: func(set RuleSet) bool {
				str, ok := set.FieldValue.(string)
				if !ok {
					return false
				}
				return len(str) >= n
			},
			MessageFunc: func(set RuleSet) string {
				return fmt.Sprintf("%s should be at least %d characters long", set.FieldName, n)
			},
		}
	}
}

func (v *Instance) Regexp(rx *regexp.Regexp) RuleFunc {
	return func() RuleSet {
		return RuleSet{
			Name:      "regexp",
			RuleValue: rx,
			ValidateFunc: func(set RuleSet) bool {
				str, ok := set.FieldValue.(string)
				if !ok {
					return false
				}
				return rx.MatchString(str)
			},
			MessageFunc: func(set RuleSet) string {
				return fmt.Sprintf("%s is invalid pattern.", set.FieldName)
			},
		}
	}
}

func (v *Instance) Validate() []ErrorMessage {
	var errorMessage []ErrorMessage

	ok := true
	for fieldName, ruleSets := range v.fields {
		// reflect panics on un-exported variables.
		if !unicode.IsUpper(rune(fieldName[0])) {
			continue
		}
		fieldValue := getFieldValueByName(v.data, fieldName)
		for _, set := range ruleSets {
			set.FieldValue = fieldValue
			set.FieldName = fieldName
			if set.Name == "message" {

				errorMessage = append(errorMessage, ErrorMessage{
					FieldName: fieldName,
					Err:       errors.New(set.RuleValue.(string)),
				})
				continue
			}
			if !set.ValidateFunc(set) {
				msg := set.MessageFunc(set)

				errorMessage = append(errorMessage, ErrorMessage{
					FieldName: fieldName,
					Err:       errors.New(msg),
				})
				ok = false
			}
		}
	}
	v.Checker = ok
	return errorMessage
}

func getFieldValueByName(v any, name string) any {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil
	}
	fieldVal := val.FieldByName(name)
	if !fieldVal.IsValid() {
		return nil
	}
	return fieldVal.Interface()
}
