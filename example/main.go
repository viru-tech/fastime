package main

import (
	"context"
	"fmt"
	"time"

	"github.com/viru-tech/fastime"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t := fastime.New()
	t.StartTimerD(ctx, time.Millisecond*5) // time precision is 5ms

	s1 := t.Now()
	s2 := t.Now()
	s3 := t.Now()
	time.Sleep(time.Second * 2)
	s4 := t.Now()

	time.Sleep(time.Second * 5)

	fmt.Printf("s1=%v\ns2=%v\ns3=%v\ns4=%v\n", s1, s2, s3, s4)

	fmt.Printf("nanonow %v\nnow unixnano %v\nnow add unixnano%v\nnanonow + dur %v\nstring %v\n",
		t.UnixNanoNow(),
		t.Now().Unix(),
		t.Now().Add(time.Second),
		t.UnixNanoNow()+int64(time.Second),
		string(t.FormattedNow()))

	for i := 0; i < 30; i++ {
		time.Sleep(time.Millisecond * 400)
		fmt.Println(t.Now())
	}
}
