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

type Node struct {
	XMLName xml.Name
	Content []byte `xml:",innerxml"`
	Nodes   []Node `xml:",any"`
	// HardCopy bool
	// FileName string
}

func main() {
	n := Nodify("../data/hmis.xml")
	n.Write("../data/hmis.json", 0700)
	n.Pretty(0)

}

func Nodify(filename string) (n Node) {
	b, _ := ioutil.ReadFile(filename)
	reader := strings.NewReader(Elementize(xml.NewDecoder(strings.NewReader(string(b)))))

	xml.NewDecoder(reader).Decode(&n)
	return
}

func Elementize(dec *xml.Decoder) (data string) {
	for {
		token, err := dec.Token()
		if err != nil {
			break
		}

		switch t := token.(type) {
		case xml.StartElement:
			data += fmt.Sprintf("<%s>", t.Name.Local)
			for _, v := range t.Attr {
				data += fmt.Sprintf("<%s>%s</%s>", v.Name.Local, v.Value, v.Name.Local)
			}

		case xml.EndElement:
			data += fmt.Sprintf("</%s>", t.Name.Local)

		case xml.CharData:
			val := strings.Replace(strings.TrimSpace(string(xml.CharData(t))), "\n", "", -1)
			if val != "" {
				data += fmt.Sprintf("<value>%s</value>", val)
			}

		default:
			// fmt.Println(t)

		}
	}
	return
}

func (n Node) Walk() interface{} {
	if len(n.Nodes) < 1 {
		return string(n.Content)

	} else {
		x := make(map[string][]interface{})
		for _, v := range n.Nodes {
			x[v.XMLName.Local] = append(x[v.XMLName.Local], v.Walk())

		}

		y := make(map[string]interface{})
		for k, v := range x {
			if len(x[k]) == 1 {
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
			fmt.Printf("%s(%d) %s == content (%s)\n", tabs, k, strings.Title(strings.ToLower(v.XMLName.Local)), string(v.Content))
		} else {
			fmt.Printf("%s(%d) %s =>\n", tabs, k, strings.Title(strings.ToLower(v.XMLName.Local)))
			v.Pretty(indent + 1)
		}
	}
}

func (n *Node) Write(name string, perm os.FileMode) {
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
