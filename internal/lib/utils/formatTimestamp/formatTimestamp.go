package formatTimestamp

import (
	"google.golang.org/protobuf/types/known/timestamppb"
)

func FormatTimestamp(ts *timestamppb.Timestamp) string {
	if ts == nil {
		return ""
	}
	return ts.AsTime().Format("2006-01-02 15:04:05")
}
