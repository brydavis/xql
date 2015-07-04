package util

import "fmt"

func TypeOf(a interface{}) string {
	return fmt.Sprintf("%T", a)
}
