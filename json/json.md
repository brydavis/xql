http://play.golang.org/p/tIuEypeLSS


```go

// http://play.golang.org/p/WxLgaHYLMi

package main

import "fmt"
import "encoding/json"

func main() {
	var i interface{}
	raw := []byte(`[{"num":6.13,"strs":["a","b"],"bool":true,"flts":[7.0,9.2]}]`)

	json.Unmarshal(raw, &i)

	switch data := i.(type) {
	case []interface{}:
		for _, row := range data {
			m := row.(map[string]interface{})
			DecodeMapJSON(m)
		}
	case map[string]interface{}:
		DecodeMapJSON(data)
	}

	fmt.Println(i)
}

func DecodeMapJSON(data map[string]interface{}) {
	for key, val := range data {
		fmt.Println(key, val)
		switch vval := val.(type) {
		case int, int16, int32, int64:
			fmt.Printf("%d is an integer\n", vval)
		case float32, float64:
			fmt.Printf("%f is a float\n", vval)
		case bool:
			fmt.Printf("%t is a bool\n", vval)
		case string:
			fmt.Printf("%s is a string\n", vval)
		case []interface{}:
			switch vval[0].(type) {
			case int, int16, int32, int64:
				fmt.Printf("%d is an array of integers\n", vval)
			case float32, float64:
				fmt.Printf("%f is an array of floats\n", vval)
			case bool:
				fmt.Printf("%t is an array of bools\n", vval)
			case string:
				fmt.Printf("%s is an array of strings\n", vval)
			case []interface{}:
				fmt.Printf("%v is an array of slices\n", vval)
			}

			//fmt.Printf("%v is a slice\n", vval)
		}
	}

}



```