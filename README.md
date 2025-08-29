# realdebrid

Wrapper and CLI for the RealDebrid API

```sh
go get -u github.com/tbdsux/realdebrid-cli/realdebrid
```

## Usage

```go
package main

import (
	"fmt"
	"os"

	"github.com/tbdsux/realdebrid-cli/realdebrid"
)


func main() {
	client := realdebrid.NewClient(os.Getenv("REALDEBRID_API_KEY"))

	user, err := client.GetUser()
	fmt.Println(user, err)
}
```

## [CLI](./rd)
