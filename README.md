# genius

genius config

## feature

* 大小写不忽略, 保持差异
* 支持几乎所有的配置文件类型

## quick start

```shell
go get github.com/jaronnie/genius@latest
```

```go
package main

import (
	"fmt"
	"github.com/jaronnie/genius"
)

func main() {
	g, err := genius.NewFromRawJSON([]byte(`
{
    "name":"jaronnie",
    "age": 23
}`))
	if err != nil {
		panic(err)
	}

	fmt.Println(g.Get("name"))
}
```