package main // persons

import "fmt"

type Persons struct {
	Likes []string

	Female bool
	Nums   []float64

	Friends []map[string]interface{}

	First  string
	Last   string
	Age    float64
	School string
}

func main() {
	n := NewStruct()
	fmt.Println(n)
}

func NewStruct() persons {
	var j persons
	return j
}
