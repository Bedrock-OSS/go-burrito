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
	burrito.PrintStackTrace = true
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

Setting `burrito.PrintStackTrace` to `true` will also output the stack trace for the error:

```
We failed to do X and we can provide the cause
   [github.com/Bedrock-OSS/go-burrito/burrito.TestReadmeExample2] main_test.go:104
[+]: We failed to do Y and we can provide the cause
   [github.com/Bedrock-OSS/go-burrito/burrito.TestReadmeExample2] main_test.go:103
[+]: This is a root error
   [github.com/Bedrock-OSS/go-burrito/burrito.TestReadmeExample2] main_test.go:102
```

You can also override the default setting for printing the stack trace per error by using the `burrito.ForceStackTrace` function.

### Tags

You can also add tags to errors to help identify them. This is useful for when you want to handle errors differently based on the nature of the error.

```go
package main

import (
	"fmt"
	"github.com/Bedrock-OSS/go-burrito/burrito"
)

const ErrNotFound = "not_found"

func main() {
	err := burrito.WrappedError("File not found")
	err.(*burrito.Error).AddTag(ErrNotFound)
	err = burrito.WrapErrorf(err, "We failed to do Y and we can provide the cause")
	err = burrito.WrapErrorf(err, "We failed to do X and we can provide the cause")
    if err.(*burrito.Error).HasTag(ErrNotFound) {
        fmt.Println("File not found")
    }
}
```
