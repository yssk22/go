package retry

import (
	"fmt"
	"testing"
	"time"

	"github.com/yssk22/go/x/xtesting/assert"

	"context"
)

func Test_Do_CancelContext(t *testing.T) {
	a := assert.New(t)
	var interval = ConstBackoff(500 * time.Millisecond)
	var i = 0
	var ut = time.Now().Add(3 * time.Second)
	var until = Until(ut)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(1 * time.Second)
		cancel()
	}()
	err := Do(ctx, func(_ context.Context) error {
		i++
		if i < 10 {
			return fmt.Errorf("Need retry")
		}
		return nil
	}, interval, until)
	a.NotNil(err)
	a.EqStr("retry canceled: context canceled", err.Error())
}

func Test_Do_Until(t *testing.T) {
	a := assert.New(t)
	var interval = ConstBackoff(100 * time.Millisecond)
	var i = 0
	var ut = time.Now().Add(3 * time.Second)
	var until = Until(ut)
	err := Do(context.Background(), func(_ context.Context) error {
		i++
		if i < 10 {
			return fmt.Errorf("Need retry")
		}
		return nil
	}, interval, until)

	a.Nil(err)
	a.OK(time.Now().Before(ut))
	a.EqInt(10, i)

	i = 0
	ut = time.Now().Add(500 * time.Millisecond)
	until = Until(ut)
	err = Do(context.Background(), func(_ context.Context) error {
		i++
		if i < 10 {
			return fmt.Errorf("Need retry")
		}
		return nil
	}, interval, until)
	a.NotNil(err)
	a.OK(time.Now().After(ut))
	a.OK(i > 1)
	a.OK(i < 10)
}

func Test_Do_MaxRetries(t *testing.T) {
	a := assert.New(t)
	var interval = ConstBackoff(100 * time.Millisecond)
	var i = 0
	var max = MaxRetries(15)
	err := Do(context.Background(), func(_ context.Context) error {
		i++
		if i < 10 {
			return fmt.Errorf("Need retry")
		}
		return nil
	}, interval, max)

	a.Nil(err)
	a.EqInt(10, i)

	i = 0
	max = MaxRetries(3)
	err = Do(context.Background(), func(context.Context) error {
		i++
		if i < 10 {
			return fmt.Errorf("Need retry")
		}
		return nil
	}, interval, max)
	a.NotNil(err)
	a.EqInt(3, i)

	i = 0
	max = MaxRetries(0)
	err = Do(context.Background(), func(context.Context) error {
		i++
		if i < 10 {
			return fmt.Errorf("Need retry")
		}
		return nil
	}, interval, max)
	a.OK(i == 1)
	a.NotNil(err)
}
