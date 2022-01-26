package duino

import (
	"crypto/sha1"
	"encoding/hex"
	"strconv"
	"sync"
)

type Job struct {
	// The hash "base", basically a string that is appended before every hash
	Base string
	// The expected result's hash
	Expected string
	// The difficulty of the job
	Difficulty int
}

// Calculates a job with a single thread. This function will return -1 in case a valid answer is not found.
func DoJob(job Job) (int, error) {
	h := sha1.New()
	for i := 0; i < 100*job.Difficulty+1; i++ {
		_, err := h.Write([]byte(job.Base + strconv.Itoa(i)))
		if err != nil {
			return i, err
		}

		sum := hex.EncodeToString(h.Sum(nil))
		if sum == job.Expected {
			return i, nil
		}

		h.Reset()
	}

	return -1, nil
}

func dojobpart(base, expected string, start, end int) (int, error) {
	h := sha1.New()
	for i := start; i < end; i++ {
		_, err := h.Write([]byte(base + strconv.Itoa(i)))
		if err != nil {
			return i, err
		}

		sum := hex.EncodeToString(h.Sum(nil))
		if sum == expected {
			return i, nil
		}

		h.Reset()
	}

	return -1, nil
}

// Calculates a job with multiple threads. This function will return -1 in case a valid answer is not found.
func DoJobMulti(job Job, threads int) (result int, err error) {
	diffPerThread := (100*job.Difficulty + 1) / threads

	resultch := make(chan int, 1)
	errorch := make(chan error, 1)

	var wg sync.WaitGroup

	for i := 0; i < threads; i++ {
		start := diffPerThread * i
		end := start + diffPerThread

		wg.Add(1)
		go func() {
			defer wg.Done()
			result, err := dojobpart(job.Base, job.Expected, start, end)
			if err != nil {
				// I have no clue if this can even error out, but just in case
				errorch <- err
				return
			}

			if result != -1 {
				resultch <- result
			}
		}()
	}

	wg.Wait()
	close(resultch)
	close(errorch)

	return <-resultch, <-errorch
}
