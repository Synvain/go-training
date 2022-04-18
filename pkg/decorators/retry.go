package util

import (
	"fmt"
	"reflect"
	"runtime"
	"time"
)


func getFuncName(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

func Retry(tries int, delay time.Duration, backoff int, f func(attempt int) (retry bool, err error)) error {
	retry, err := false, error(nil)
	for i := 1; tries < 0 || i <= tries; i++ {
		fmt.Printf("Trying to call %v, attempt %d\n", getFuncName(f), i)
		if retry, err = f(i); err == nil || !retry {
			return err
		}
		fmt.Printf("Sleeping %s before retrying\n", delay)
		time.Sleep(delay)
		delay = delay * time.Duration(backoff)
	}
	return fmt.Errorf("Retry failed after %d attempts, last error was: '%s'", tries, err.Error())
}