package ratelimit

import "fmt"

type Limiter struct {
	Bucket *Bucket
}

func (r *Limiter) Limit() bool {
	tokenRes := r.Bucket.TakeAvailable(1)
	if tokenRes == 0 {
		fmt.Printf("Reached Rate-Limiting %d \n", r.Bucket.Available())
		return true
	}

	// if tokenRes is not zero, means gRpc request can continue to flow without rate limiting :)
	return false
}
