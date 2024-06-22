package util

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func getCurrentTimestamp() *timestamppb.Timestamp {
	return timestamppb.New(time.Now())
}
