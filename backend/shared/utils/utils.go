package utils

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func AsPointer[T any](e T) *T {
	return &e
}

func NilOrPointer[T any](e *T) *T {
	if e == nil {
		return nil
	}

	return e
}

func NilTimeToRpc(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}

	return timestamppb.New(*t)
}
