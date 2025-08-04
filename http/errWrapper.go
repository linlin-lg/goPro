package main

import (
	"fmt"
	"net/http"
	"os"
)

type appHandler func(writer http.ResponseWriter, request *http.Request) error

type UserError string
func (e UserError) Error() string {
	return e.Message()
}
func (e UserError) Message() string {
	return string(e)
}


//type UserError struct {
//	Err string
//	Code int
//}
//func (e UserError) Error() string {
//	return e.Err
//}
//func (e UserError) getCode() int {
//	return e.Code
//}

func ErrWrok(handler appHandler) func(http.ResponseWriter, *http.Request) {

	return func(writer http.ResponseWriter, request *http.Request) {

		defer func() {
			// 发生宕机时，获取panic传递的上下文并打印
			if r := recover(); r != nil {
				fmt.Println("\n err:", r)
				http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		err := handler(writer, request)
		if err == nil {
			return
		}

		// user error 略过

		//system error
		code := http.StatusOK
		switch {
		case os.IsNotExist(err):
			code = http.StatusNotFound
		case os.IsPermission(err):
			code = http.StatusForbidden
		default:
			code = http.StatusInternalServerError
		}
		fmt.Println("\n handling request error = ", err.Error())
		http.Error(writer, err.Error(), code)
	}
}