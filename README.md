# go-burrito

Simple, fast, and easy to use library for handling errors.

## Installation

```bash
go get github.com/Bedrock-OSS/go-burrito
```

## Usage

```go
package main

import (
	"fmt"
	"github.com/Bedrock-OSS/go-burrito/burrito"
)

func main() {
	burrito.Debug = true
	err := burrito.WrappedError("This is a root error")
	err = burrito.WrapErrorf(err, "We failed to do Y and we can provide the cause")
	err = burrito.WrapErrorf(err, "We failed to do X and we can provide the cause")
	fmt.Println(err)
}

```

This will output:

```
We failed to do X and we can provide the cause
[+]: We failed to do Y and we can provide the cause
[+]: This is a root error
```

Setting `burrito.Debug` to `true` will also output the stack trace for the error:

```
We failed to do X and we can provide the cause
[+]: We failed to do Y and we can provide the cause
[+]: This is a root error
   [main.main] main.go:10
   [main.main] main.go:11
   [main.main] main.go:12
```