package main

import (
	"context"
	"fmt"
)

func main() {
	fmt.Println("Hello, World! I generate docs")

	ctx := context.Background()

	go Listen(ctx)
	Write(ctx, "Oh", "Test passed")
}
