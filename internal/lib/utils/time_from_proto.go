package utils

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func ProtoToTime(ts *timestamppb.Timestamp) time.Time {
	if ts == nil {
		return time.Time{}
	}
	return ts.AsTime()
}
