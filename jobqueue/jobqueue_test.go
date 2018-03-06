package jobqueue

import (
	"testing"
)

func TestQueue(test *testing.T) {
	res := []int{}

	queue := New()
	queue.Job(func() {
		res = append(res, 1)
	})
	queue.Job(func() {
		res = append(res, 2)
	})

	queue.Run()

	if len(res) != 2 || res[0] != 1 || res[1] != 2 {
		test.Errorf("Failed")
	}
}
