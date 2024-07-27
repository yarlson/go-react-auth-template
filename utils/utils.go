package utils

import (
	"math/rand"
	"time"
)

func RetryWithBackoff(operation func() error, maxRetries int) error {
	var err error
	for i := 0; i < maxRetries; i++ {
		err = operation()
		if err == nil {
			return nil
		}
		if i == maxRetries-1 {
			break
		}
		time.Sleep(time.Duration(1<<uint(i))*time.Second + time.Duration(rand.Intn(1000))*time.Millisecond)
	}

	return err
}
