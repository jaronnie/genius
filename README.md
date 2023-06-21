# genius

genius config

## feature

* 大小写不忽略, 保持差异
* 支持几乎所有的配置文件类型
* 支持 get array
* 支持 set array
* 支持 append array

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
    "skills": ["golang", "python", "c"]
}`))
	if err != nil {
		panic(err)
	}

	fmt.Println(g.Get("name"))
	fmt.Println(g.Get("skills"))
	fmt.Println(g.Get("skills.0"))
	fmt.Println(g.Set("skills.0", "go"))
}
```