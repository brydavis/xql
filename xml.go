package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type xmlNode struct {
	XMLName xml.Name
	Content []byte    `xml:",innerxml"`
	Nodes   []xmlNode `xml:",any"`
}

func X2J(filename string) []byte { //(n xmlNode) {
	var n xmlNode
	b, _ := ioutil.ReadFile(filename)
	reader := strings.NewReader(Elementize(xml.NewDecoder(strings.NewReader(string(b)))))

	xml.NewDecoder(reader).Decode(&n)
	var y []interface{}
	for _, v := range n.Nodes {
		y = append(y, v.walk())
	}

	j, _ := json.Marshal(y)
	return j
}

func (n xmlNode) walk() interface{} {
	if len(n.Nodes) < 1 {
		return string(n.Content)
	} else {
		x := make(map[string][]interface{})
		for _, v := range n.Nodes {
			x[v.XMLName.Local] = append(x[v.XMLName.Local], v.walk())
		}

		y := make(map[string]interface{})
		for k, v := range x {
			if vv := v[0]; len(x[k]) == 1 {
				f, errFloat := strconv.ParseFloat(vv.(string), 64)
				i, errInt := strconv.ParseInt(vv.(string), 10, 64)
				b, errBool := strconv.ParseBool(vv.(string))

				switch {
				case errFloat == nil:
					y[k] = f
				case errInt == nil:
					y[k] = i
				case errBool == nil:
					y[k] = b
				default:
					y[k] = vv
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

func Beautify(name string) (string, error) {
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

func Elementize(dec *xml.Decoder) (data string) {
	var parentNode string
	var attrs bool
	for {
		token, err := dec.Token()
		if err != nil {
			break
		}

		switch t := token.(type) {
		case xml.StartElement:
			parentNode = t.Name.Local
			data += fmt.Sprintf("<%s>", t.Name.Local)
			if len(t.Attr) > 0 {
				attrs = true
				// data += "<attr>"
				for _, v := range t.Attr {
					data += fmt.Sprintf("<%s>%s</%s>", v.Name.Local, v.Value, v.Name.Local)
				}
				// data += "</attr>"
			} else {
				attrs = false
			}

		case xml.EndElement:
			data += fmt.Sprintf("</%s>", t.Name.Local)

		case xml.CharData:
			val := strings.Replace(strings.TrimSpace(string(xml.CharData(t))), "\n", "", -1)
			if val != "" {
				// data += fmt.Sprintf("<value>%s</value>", val)
				if attrs {
					data += fmt.Sprintf("<%s>%s</%s>", parentNode, val, parentNode)

				} else {
					data += fmt.Sprintf("%s", val)

				}
			}

		case xml.ProcInst:
			// fmt.Printf("not supported: %v\n", t)

		default:
			fmt.Printf("type not supported: %v\n", t)
		}
	}
	return
}

// func Comma(records [][]string) {
// 	csvfile, err := os.Create("output.csv")
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 	}
// 	defer csvfile.Close()

// 	writer := csv.NewWriter(csvfile)
// 	for _, record := range records {
// 		err := writer.Write(record.([]string))
// 		if err != nil {
// 			fmt.Println("Error:", err)
// 		}
// 	}
// 	writer.Flush()
// }

// func (n Node) Pretty(indent int) {
// 	var tabs string
// 	for i := 0; i < indent; i++ {
// 		tabs += "\t"
// 	}

// 	for k, v := range n.Nodes {
// 		if len(v.Nodes) < 1 {
// 			fmt.Printf("%s(%d) %s == content (%s)\n", tabs, k, strings.Title(strings.ToLower(v.XMLName.Local)), string(v.Content))
// 		} else {
// 			fmt.Printf("%s(%d) %s =>\n", tabs, k, strings.Title(strings.ToLower(v.XMLName.Local)))
// 			v.Pretty(indent + 1)
// 		}
// 	}
// }

// func (n Node) Querify() {
// for k, v := range n.Nodes {
// 	if len(v.Nodes) < 1 {
// 		fmt.Printf("%s(%d) %s == content node (%s)\n", k, strings.Title(strings.ToLower(v.XMLName.Local)), string(v.Content))
// 	} else {
// 		fmt.Printf("%s(%d) %s == parent node\n", k, strings.Title(strings.ToLower(v.XMLName.Local)))
// 		v.Pretty()
// 	}
// }
// }
