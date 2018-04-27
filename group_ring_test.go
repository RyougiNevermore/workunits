package workunits

import (
	"testing"
	"fmt"
	"time"
)

type ringUnit struct {

}

func (u *ringUnit) Process()  {
	fmt.Println(time.Now())
}

func TestNewRingWorkerGroup(t *testing.T) {
	group := NewRingWorkerGroup(4, 64, 64)
	group.Start()
	for i := 0 ; i < 10 ; i ++ {
		group.Send(&ringUnit{})
	}
	group.Close()
	group.Sync()
}
