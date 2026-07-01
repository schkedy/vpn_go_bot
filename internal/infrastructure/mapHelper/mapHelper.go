package maphelper

import (
	"strconv"
	"strings"
)

func ConvertStrToInterface(s map[string]string) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	for k, v := range s {
		switch getType(k) {
		case "int":
			m[k], err = strconv.Atoi(v)
		case "int8":
			parsed, _ := strconv.ParseInt(v, 10, 8)
			m[k] = int8(parsed)
		}
	}
}

func getType(value string) string {
	if strings.HasSuffix(value, "<{int}>") {
		return "int"
	}
	if strings.HasSuffix(value, "<{int8}>") {
		return "int8"
	}
	if strings.HasSuffix(value, "<{int16}>") {
		return "int16"
	}
	if strings.HasSuffix(value, "<{int32}>") {
		return "int32"
	}
	if strings.HasSuffix(value, "<{rune}>") {
		return "int32"
	}
	if strings.HasSuffix(value, "<{int64}>") {
		return "int64"
	}
	if strings.HasSuffix(value, "<{uint}>") {
		return "uint"
	}
	if strings.HasSuffix(value, "<{uint8}>") {
		return "uint8"
	}
	if strings.HasSuffix(value, "<{byte}>") {
		return "uint8"
	}
	if strings.HasSuffix(value, "<{uint16}>") {
		return "uint16"
	}
	if strings.HasSuffix(value, "<{uint32}>") {
		return "uint32"
	}
	if strings.HasSuffix(value, "<{uint64}>") {
		return "uint64"
	}
	if strings.HasSuffix(value, "<{float}>") {
		return "float"
	}
	if strings.HasSuffix(value, "<{float32}>") {
		return "float32"
	}
	if strings.HasSuffix(value, "<{float64}>") {
		return "float64"
	}
	if strings.HasSuffix(value, "<{complex64}>") {
		return "complex64"
	}
	if strings.HasSuffix(value, "<{complex128}>") {
		return "complex128"
	}
	if strings.HasSuffix(value, "<{string}>") {
		return "string"
	}
	if strings.HasSuffix(value, "<{bool}>") {
		return "bool"
	}

}
