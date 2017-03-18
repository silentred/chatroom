package chatroom

import "testing"
import "fmt"
import "time"

func TestForRange(t *testing.T) {
	ch := make(chan int, 10)
	go func() {
		for i := range ch {
			fmt.Println(i)
		}
		fmt.Println("end loop")
	}()

	after := time.After(time.Second * 1)
	<-after
	close(ch)
	//ch <- 1 will panic: send on closed channel
	fmt.Println("close chan")
	time.Sleep(time.Second)
}
