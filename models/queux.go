package models

import (
	"encoding/json"
	"fmt"
)

type Queux []Question

func (q Queux) Len() int                { return len(q) }
func (q *Queux) Enqueue(value Question) { *q = append(*q, value) }
func (q *Queux) Dequeue() {
	if q.Len() == 0 {
		return
	}
	*q = (*q)[1:]
}
func (q *Queux) SwapEnd() {
	if q.Len() > 1 {
		*q = append((*q)[1:], (*q)[0])
	}
}
func (q *Queux) First() *Question {
	if q.Len() != 0 {
		return &(*q)[0]
	}
	return nil
}
func (q *Queux) Marshall() ([]byte, error) {
	if q.Len() != 0 {
		b, err := json.Marshal(*q)
		if err != nil {
			return nil, err
		}
		return b, nil
	}
	return nil, fmt.Errorf("something wrong happened when marshalling")
}
