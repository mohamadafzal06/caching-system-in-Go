package main

import (
	"fmt"
	"time"

	"github.com/mohamadafzal06/cache-in-go/cache"
)

func main() {
	c := cache.New()

	c.Set([]byte("m"), []byte("a"), time.Minute)
	d, err := c.Get([]byte("m"))
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(d))

	b := c.Has([]byte("m"))
	if !b {
		fmt.Println("the is no key")
	}

	err = c.Delete([]byte("m"))
	if err != nil {
		fmt.Println(err.Error())
	}

}
