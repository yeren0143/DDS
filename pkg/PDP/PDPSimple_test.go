package PDPSimple

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"testing"
)

func TestListenMultiCast(t *testing.T) {
	fmt.Println("TestListenMultiCast")

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	var wg sync.WaitGroup
	wg.Add(1)

	ctxWithCancel, cancel := context.WithCancel(context.Background())
	go func() {
		wg.Done()
		fmt.Println("begin test listenMultiCast")
		listenMultiCast(ctxWithCancel, "239.255.0.1:7400")
	}()

	wg.Wait()
	<-done
	cancel()

	t.Logf("Test listenMultiCast Done")
}
