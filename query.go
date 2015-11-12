package main

// func (base *Base) Select(elements ...string) *Base {
// 	var stmt string
// 	for _, e := range elements {
// 		stmt += e + `, `
// 	}

// 	base.Query += `select ` + stmt[:len(stmt)-2]
// 	return base
// }

// func (base *Base) From(tables ...string) *Base {
// 	var stmt string
// 	for _, t := range tables {
// 		stmt += t + `, `
// 	}

// 	base.Query += "\nfrom " + stmt[:len(stmt)-2]
// 	return base
// }

// func (base *Base) Join(table, join, keys string) *Base {
// 	base.Query += ` A ` + join + ` join ` + table + ` B on A.` + keys + ` = B.` + keys
// 	return base
// }
