package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

/* TODO
- Generalize API for other formats
- Transform data to actual data types
- Integrate into "xql" package
*/

// var final = ""

type Node struct {
	XMLName xml.Name
	Gender  xml.Attr `xml:"gender,attr"`
	Content []byte   `xml:",innerxml"`
	Nodes   []Node   `xml:",any"`
}

func main() {
	n := Nodify("xml/complex.xml")
	n.Write("json/complex.json", 0700)
	// n.Pretty(0)

	// fmt.Println(final)

}

func Nodify(filename string) Node {
	b, _ := ioutil.ReadFile(filename)

	// buf := bytes.NewBuffer(data)
	reader := strings.NewReader(string(b))

	// dec := xml.NewDecoder(buf)
	dec := xml.NewDecoder(reader)

	var data string
	// var tabs int
	for {
		token, err := dec.Token()
		if err != nil {
			break
		}

		switch t := token.(type) {
		case xml.StartElement:
			// var indent string
			// for i := 0; i < tabs; i++ {
			// 	indent += "\t"
			// }
			// tabs += 1

			data += fmt.Sprintf("<%s>", t.Name.Local)

			// var attrs string
			for _, v := range t.Attr {
				// attrs += fmt.Sprintf("%s\t<%s>%s</%s>\n",  v.Name.Local, v.Value, v.Name.Local)
				data += fmt.Sprintf("<%s>%s</%s>", v.Name.Local, v.Value, v.Name.Local)

			}

			// fmt.Println(attrs)

		case xml.EndElement:
			// tabs -= 1
			// var indent string
			// for i := 0; i < tabs; i++ {
			// 	indent += "\t"
			// }

			data += fmt.Sprintf("</%s>", t.Name.Local)

		case xml.CharData:
			val := strings.Replace(strings.TrimSpace(string(xml.CharData(t))), "\n", "", -1)
			fmt.Println(val)
			if val != "" {
				// var indent string
				// for i := 0; i < tabs; i++ {
				// 	indent += "\t"
				// }

				// data += fmt.Sprintf("%s<value>%s</value>\n", indent, val)
				data += fmt.Sprintf("<value>%s</value>", val)

			}

		default:
			// fmt.Println(t)

		}

	}

	// data, _ := ioutil.ReadFile(filename)

	// buf := bytes.NewBuffer(data)
	reader = strings.NewReader(data)

	// dec := xml.NewDecoder(buf)
	dec = xml.NewDecoder(reader)

	var n Node
	dec.Decode(&n)
	return n

}

func (n Node) Walk() interface{} {
	if len(n.Nodes) < 1 {
		return string(n.Content)
		// final += fmt.Sprintf("<%s>%s</%s>", n.XMLName.Local, string(n.Content), n.XMLName.Local)

	} else {
		x := make(map[string][]interface{})
		for _, v := range n.Nodes {
			x[v.XMLName.Local] = append(x[v.XMLName.Local], v.Walk())

		}

		y := make(map[string]interface{})
		for k, v := range x {

			if len(x[k]) == 1 {
				// final += fmt.Sprintf("<%s>%s</%s>", k, v, k)

				switch v[0].(type) {
				case float32, float64:
					f, _ := strconv.ParseFloat(v[0].(string), 64)
					y[k] = f // v[0].(float64)
				case bool:
					y[k] = v[0].(bool)
				default:
					y[k] = v[0]
				}
			} else {
				y[k] = v
			}
		}
		return y
	}
	var v interface{}
	return v
}

func (n Node) Import() []byte {
	var y []interface{}
	for _, v := range n.Nodes {
		y = append(y, v.Walk())
	}

	b, _ := json.Marshal(y)
	return b
}

func (n Node) Pretty(indent int) {
	var tabs string
	for i := 0; i < indent; i++ {
		tabs += "\t"
	}

	for k, v := range n.Nodes {
		if len(v.Nodes) < 1 {
			fmt.Printf("%s(%d) %s == content (%s)\tAttr => %v\n", tabs, k, strings.Title(strings.ToLower(v.XMLName.Local)), string(v.Content), v.Gender)
		} else {
			fmt.Printf("%s(%d) %s =>\tAttr => %v\n", tabs, k, strings.Title(strings.ToLower(v.XMLName.Local)), v.Gender)
			v.Pretty(indent + 1)
		}
	}
}

func (n Node) Write(name string, perm os.FileMode) {
	j := n.Import()
	ioutil.WriteFile(name, j, perm)
}

func (n Node) Querify() {
	// for k, v := range n.Nodes {
	// 	if len(v.Nodes) < 1 {
	// 		fmt.Printf("%s(%d) %s == content node (%s)\n", k, strings.Title(strings.ToLower(v.XMLName.Local)), string(v.Content))
	// 	} else {
	// 		fmt.Printf("%s(%d) %s == parent node\n", k, strings.Title(strings.ToLower(v.XMLName.Local)))
	// 		v.Pretty()
	// 	}
	// }
}

func Pretty(name string) (string, error) {
	b, err := ioutil.ReadFile(name)
	if err != nil {
		return "", err
	}

	var data interface{}

	if err := json.Unmarshal(b, &data); err != nil {
		return "", err
	}
	b, err = json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}

	return string(b), nil
}
