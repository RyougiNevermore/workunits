package workunits

import (
	"testing"
	"fmt"
)

type aunit struct {
	i int
}

func (u *aunit) Process() {
	fmt.Printf("[%v] process. \n", u.i)
}

func TestNewArrayBufferedWorkerGroup(t *testing.T) {
	var err error
	group := NewArrayBufferedWorkerGroup(8)
	err = group.Start()
	if err != nil {
		t.Errorf("start failed, %v", err)
		t.FailNow()
		return
	}
	for i := 0; i < 20; i++ {
		err = group.Send(&aunit{i: i})
		if err != nil {
			t.Errorf("send failed, %v", err)
			t.FailNow()
			return
		}
	}
	err = group.Close()
	if err != nil {
		t.Errorf("close failed, %v", err)
		t.FailNow()
		return
	}
	err = group.Sync()
	if err != nil {
		t.Errorf("sync failed, %v", err)
		t.FailNow()
		return
	}
}

