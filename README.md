# midec [![Go Reference](https://pkg.go.dev/badge/github.com/sapphi-red/midec.svg)](https://pkg.go.dev/github.com/sapphi-red/midec) [![CI](https://github.com/sapphi-red/midec/actions/workflows/main.yaml/badge.svg)](https://github.com/sapphi-red/midec/actions/workflows/main.yaml) [![codecov](https://codecov.io/gh/sapphi-red/midec/branch/main/graph/badge.svg?token=H9T7BGUQ7V)](https://codecov.io/gh/sapphi-red/midec)

Golang **M**ulti-**i**mage **de**te**c**tor.
cf. Animated GIF, APNG, Animated WebP, Animated HEIF / AVIF.

## Usage
```go
package main 

import (
	"fmt"
	"os"

	"github.com/sapphi-red/midec"
	_ "github.com/sapphi-red/midec/gif" // import this to detect Animated GIF
	// _ "github.com/sapphi-red/midec/png" // import this to detect APNG
	// _ "github.com/sapphi-red/midec/webp" // import this to detect Animated WebP
	// _ "github.com/sapphi-red/midec/isobmff" // import this to detect Animated HEIF / AVIF
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

## Benchmarks
Comparison with using `image/gif` package's `gif.decodeAll`. See code for [`bench_test.go`](https://github.com/sapphi-red/midec/blob/main/bench_test.go).
```text
goos: windows
goarch: amd64
pkg: github.com/sapphi-red/midec
cpu: AMD Ryzen 7 3700X 8-Core Processor
BenchmarkGIF_ImageGIF-16            2406            451885 ns/op          497565 B/op           1435 allocs/op      
BenchmarkGIF_Midec-16             100472             10931 ns/op            5008 B/op             36 allocs/op      
BenchmarkPNG_Midec-16             183244              6957 ns/op            5008 B/op             13 allocs/op      
BenchmarkWebP_Midec-16            142461              9733 ns/op            5040 B/op             20 allocs/op      
BenchmarkHEIFAVIF_Midec-16        126554              8688 ns/op            5136 B/op             44 allocs/op      
PASS
ok      github.com/sapphi-red/midec     24.975s
```
