# realdebrid

Wrapper and CLI for the RealDebrid API


```sh
go get -u github.com/tbdsux/realdebrid-cli/realdebrid
```


## Usage


```go
package main


import (
	"github.com/tbdsux/realdebrid-cli/realdebrid"
	"fmt"
)


func main() {
	client := NewClient(os.Getenv("REALDEBRID_API_KEY"))

	user, err := client.GetUser()
	fmt.Println(user, err)
}
```



## [CLI](./rd)
