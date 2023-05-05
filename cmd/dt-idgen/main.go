package main

import (
	"flag"
	"fmt"
	"github.com/google/uuid"
	"strings"
)

var n int
var idType string

var upper bool

func init() {
	flag.IntVar(&n, "n", 1, "id numbers")
	flag.StringVar(&idType, "type", "uuid", "generated id type")
	flag.BoolVar(&upper, "upper", false, "show the uppercase alphabet")
}

func main() {
	flag.Parse()
	var generator func() string

	switch idType {
	case "uuid":
		generator = UUID
	default:
		fmt.Println("unsupported type")
		return
	}

	for i := 0; i < n; i++ {
		id := generator()
		if upper {
			id = strings.ToUpper(id)
		}
		fmt.Println(id)
	}
}

func UUID() string {
	return uuid.New().String()
}
