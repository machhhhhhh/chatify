package constants

import (
	"net/http"
	"slices"
)

// enum type
type HTTPMethod string

var ValidatorHTTPMethods []HTTPMethod = []HTTPMethod{
	http.MethodGet,
	http.MethodHead,
	http.MethodPut,
	http.MethodPatch,
	http.MethodPost,
	http.MethodDelete,
	http.MethodOptions,
}

func IsHTTPMethodExist(method HTTPMethod) bool {
	return slices.Contains(ValidatorHTTPMethods, method) == true
}

func GetAllHTTPMethods() []string {
	var methods []string
	for i := range ValidatorHTTPMethods {
		methods = append(methods, string(ValidatorHTTPMethods[i]))
	}
	return methods
}
