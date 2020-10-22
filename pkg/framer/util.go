package framer

import (
	"encoding/json"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iotsitewise"
)

func getTime(ts *iotsitewise.TimeInNanos) time.Time {
	sec := *ts.TimeInSeconds

	if nanos := ts.OffsetInNanos; nanos != nil {
		return time.Unix(sec, *nanos)
	}

	return time.Unix(sec, 0)
}

func getPropertyVariantValue(variant *iotsitewise.Variant) interface{} {

	if val := variant.BooleanValue; val != nil {
		return *val
	}

	if val := variant.DoubleValue; val != nil {
		return *val
	}

	if val := variant.IntegerValue; val != nil {
		return *val
	}

	if val := variant.StringValue; val != nil {
		return *val
	}

	return nil
}

func getErrorDescription(details *iotsitewise.ErrorDetails) (*string, error) {

	if details == nil {
		return nil, nil
	}

	jb, err := json.Marshal(*details)
	if err != nil {
		return nil, err
	}
	return aws.String(string(jb)), nil
}

func getAssetHierarchies(hierarchies []*iotsitewise.AssetHierarchy) (string, error) {

	serialized, err := json.Marshal(hierarchies)

	if err != nil {
		return "", err
	}
	return string(serialized), nil
}
