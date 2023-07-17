## ginx 

#### Example
```go
package main

import (
	"context"
	"github.com/eatmoreapple/ginx"
	"github.com/gin-gonic/gin"
)

type User struct {
	Name string `json:"name" form:"name"`
}

func helloworld(_ context.Context, req User) (any, error) {
	if req.Name == "" {
		req.Name = "hello world"
	}
	return req, nil
}

func main() {
	engine := gin.Default()
	router := ginx.NewRouter(engine)
	router.GET("/", ginx.G(helloworld).JSON())
	engine.Run(":8080")
}
```