# procControl
多进程并发控制

## usage
```go
package main

import (
	"fmt"
	"time"
	
	pc "github.com/mosesyu95/procControl"
)


func main() {
	sema := pc.NewProcControl(10,100)
	for i := 0; i < 100; i++ {
		fmt.Println("start ",i)
		sema.Acquire()
		go func(i int) {
			defer func() {
				sema.Release()
			}()
			time.Sleep(1 * time.Second)
			fmt.Println("print ",i)
		}(i)
	}
	sema.Wait()
}
```
