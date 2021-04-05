# midec

Golang **M**ulti-**i**mage **de**te**c**tor.
cf. Animated GIF, APNG, Animated WebP, Animated AVIF.

Supports only Animated GIF and APNG for now.

## Usage
```go
package main 

import (
	"fmt"
	"os"

	"github.com/sapphi-red/midec"
	_ "github.com/sapphi-red/midec/gif" // import this to detect Animated GIF
	// _ "github.com/sapphi-red/midec/png" // import this to detect APNG
)

func main() {
	fp, err := os.Open("test.gif")
	if err != nil {
		panic(err)
	}

	isAnimated := midec.IsAnimated(fp)
	fmt.Println(isAnimated)
}
```

## Extension
To add support for other formats, use `midec.RegisterFormat`.
This function is very similar to [`image.RegisterFormat`](https://golang.org/pkg/image/#RegisterFormat).

```go
func init() {
	midec.RegisterFormat("gif", gifHeader, isAnimated)
}
```
