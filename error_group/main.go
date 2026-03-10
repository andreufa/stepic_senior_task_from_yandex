package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync/atomic"
	"time"

	"golang.org/x/sync/errgroup"
)

func fitch(ctx context.Context, id int) (int, error) {
	res := rand.Intn(100)
	select {
	case <-ctx.Done():
		return res, ctx.Err()
	default:
		time.Sleep(100 * time.Microsecond)
		if res > 50 {
			return res, nil
		}
		return res, fmt.Errorf("value < 50 %d", res)
	}
}

func multyFetch(ctx context.Context, ips []int) (int64, error) {
	var sum int64

	group, ectx := errgroup.WithContext(ctx)
	group.SetLimit(5)

	for id := range ips {
		id := id
		group.Go(func() (err error) {
			defer func(){
				if e:=recover(); e != nil{
					err = fmt.Errorf("recover panic %s", e)
				}
			}()
			_, err = fitch(ectx, id)
			if err != nil {
				return err
			}
			atomic.AddInt64(&sum, 1)
			return nil
		})
	}
	err := group.Wait()
	if err != nil {
		return sum, err
	}
	return sum, nil
}

func main() {
	ips := make([]int, 10)
	for i := 0; i < 10; i++ {
		ips[i] = i
	}

	ctx := context.Background()
	res, err := multyFetch(ctx, ips)
	if err != nil {
		fmt.Printf("Error!  -> count %d | error %s", res, err)
	} else {
		fmt.Printf("Success  -> count %d | error %s", res, err)
	}
}
