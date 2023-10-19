package logger

import "fmt"

func LogFailure(err error, msg string) {
	fmt.Println("NEW ERROR OCCURRED")
	fmt.Printf("error message: %s\n", err.Error())
	fmt.Printf("custome message: %s\n\n\n", msg)
}
