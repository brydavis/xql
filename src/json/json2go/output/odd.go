package main // odd

type Odd struct {
	Flts	[]float64

Nested struct {
	N2	[]bool

	N1	float64

}
	Num	float64
	Strs	[]string

	Bln	bool

}
Odd struct {
	Name	string
	Age	float64
	Occupation	string

}


func main() {
	n:= NewStruct()
	fmt.Println(n)
}

func NewStruct() odd {
	var j odd
	return j
}

