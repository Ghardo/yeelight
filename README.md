# yeelight
Very simple yeelight control without discover

```go
package main

import (
	"fmt"
	"log"

	"github.com/Ghardo/yeelight"
)

func main() {
	yl := yeelight.Yeelight{Address: "IP:55443"}
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
