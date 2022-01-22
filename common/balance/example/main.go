package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/smallnest/weighted"
)

func ExampleW2_Next() {
	w := &weighted.RRW{}
	w.Add("a", 5)
	w.Add("b", 2)
	w.Add("c", 3)


	count := make(map[string]int)
	var lock sync.Mutex
	for i := 0; i < 20; i++ {
		i := i

		go func() {
			if i %19 == 0 {
				w.Add("a", 10)
			}

			next := w.Next()
			fmt.Println(next)

			item, ok := next.(string)
			if !ok {
				fmt.Printf("!ok")
				return
			}
			lock.Lock()
			count[item] ++
			lock.Unlock()

			//fmt.Printf("%s ", item)
		}()
	}
	time.Sleep(1*time.Second)

	for k ,v := range count {
		fmt.Printf("k:%v, v: %v\n",k,v)
	}


}

func main()  {
	ExampleW2_Next()
}

