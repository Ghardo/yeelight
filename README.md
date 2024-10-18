# yeelight
Very simple yeelight control without discover

basic example
```go
package main

import (
	"fmt"
	"log"

	"github.com/Ghardo/yeelight"
)

func main() {
	yl := yeelight.Yeelight{Address: "<local-ip>:55443"}
	yl.SetHexColor("#f59542")
	yl.SetOn()
	fmt.Println(yl.IsOn())

	p := []string{"rgb", "bright"}
	r, err := yl.GetProperties(p)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(r.Result)
}
```


how to use matrix functions
this example is for 1 dot matrix and a spot light
```go
package main

import (
	"log"

	"github.com/Ghardo/yeelight"
)

func main() {
	var err error
	yl := yeelight.Yeelight{Address: "<local-ip>:55443", Persistent: true}
	err = yl.SetDirectMode()
	if err != nil {
		log.Fatal(err)
	}

	// Method with predefined array of hex values
	colors1 := []string{"#FF0000", "#000000", "#000000", "#000000", "#FF0000",
		"#FF0000", "#FF0000", "#000000", "#FF0000", "#FF0000",
		"#FF0000", "#FF0000", "#FF0000", "#FF0000", "#FF0000",
		"#FF0000", "#FFFFFF", "#FF0000", "#FFFFFF", "#FF0000",
		"#000000", "#FF0000", "#FF0000", "#FF0000", "#000000"}
	matrix1 := yeelight.MakeFromHexColors(colors1)
	matrix1 = matrix1.Rotate(90.0)
	
	//matrix1 := yeelight.MakeMatrix("#0000FF", 25)
	//v2 := yeelight.Vector{Row: 1, Column: 2}
	//matrix1.SetHex(v2, "#FF0000")

	matrix2 := yeelight.MakeSpotMatrix("#ff9100")

	m := []yeelight.ColorMatrix{matrix1, matrix2}

	err = yl.SetMatrix(m)

	if err != nil {
		log.Fatal(err)
	}
}

```
