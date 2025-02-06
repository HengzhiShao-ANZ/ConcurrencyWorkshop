package arithmetics

import (
	"context"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"github.com/romangurevitch/concurrencyworkshop/internal/pattern/fanoutin"
)

func SequentialSum(inputSize int) int {
	sum := 0
	for i := 1; i <= inputSize; i++ {
		sum += process(i)
	}
	return sum
}

// Example squareNonNegative function that squares non-negative integer.
func squareNonNegative(_ context.Context, value int) (int, error) {
	return process(value), nil
}

// ParallelSum implement this method.
func ParallelSum(inputSize int) int {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var sum atomic.Int64
	var jobs []fanoutin.Job[int]
	for i := 1; i <= inputSize; i++ {
		jobs = append(jobs, fanoutin.Job[int]{ID: i, Value: i})
	}

	// Fan out
	results := fanoutin.FanOut(ctx, jobs, squareNonNegative)

	wg := sync.WaitGroup{}
	// Fan in
	for result := range results {
		wg.Add(1)
		go func() {
			defer wg.Done()
			select {
			case <-ctx.Done():
				return
			default:
				if result.Err != nil {
					slog.Error("Error processing job", "jobID", result.Job.ID, "error", result.Err)
					cancel()
					return
				}
				slog.Info("Result for job", "jobID", result.Job.ID, "result", result.Value)
				sum.Add(int64(result.Value))
			}
		}()
	}

	wg.Wait()
	return int(sum.Load())
}

func process(num int) int {
	time.Sleep(100 * time.Millisecond) // simulate processing time
	return num * num
}
