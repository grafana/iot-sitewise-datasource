package framer

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iotsitewise"
	iotsitewisetypes "github.com/aws/aws-sdk-go-v2/service/iotsitewise/types"

	"github.com/grafana/iot-sitewise-datasource/pkg/util"
)

func getTime(ts *iotsitewisetypes.TimeInNanos) time.Time {
	sec := *ts.TimeInSeconds

	if nanos := ts.OffsetInNanos; nanos != nil {
		return time.Unix(sec, int64(*nanos))
	}

	return time.Unix(sec, 0)
}

func isPropertyDataTypeDefined(dataType iotsitewisetypes.PropertyDataType) bool {
	return dataType == iotsitewisetypes.PropertyDataTypeBoolean ||
		dataType == iotsitewisetypes.PropertyDataTypeDouble ||
		dataType == iotsitewisetypes.PropertyDataTypeInteger ||
		dataType == iotsitewisetypes.PropertyDataTypeString
}

func getPropertyVariantValue(variant *iotsitewisetypes.Variant) interface{} {

	if val := variant.BooleanValue; val != nil {
		return *val
	}

	if val := variant.DoubleValue; val != nil {
		return *val
	}

	if val := variant.IntegerValue; val != nil {
		return int64(*val)
	}

	if val := variant.StringValue; val != nil {
		return *val
	}

	return nil
}

func getPropertyVariantValueType(variant *iotsitewisetypes.Variant) iotsitewisetypes.PropertyDataType {

	if val := variant.BooleanValue; val != nil {
		return iotsitewisetypes.PropertyDataTypeBoolean
	}

	if val := variant.DoubleValue; val != nil {
		return iotsitewisetypes.PropertyDataTypeDouble
	}

	if val := variant.IntegerValue; val != nil {
		return iotsitewisetypes.PropertyDataTypeInteger
	}

	if val := variant.StringValue; val != nil {
		return iotsitewisetypes.PropertyDataTypeString
	}

	return ""
}

func getErrorDescription(details *iotsitewisetypes.ErrorDetails) (*string, error) {

	if details == nil {
		return nil, nil
	}

	jb, err := serialize(*details)
	if err != nil {
		return nil, err
	}
	return aws.String(jb), nil
}

func serialize(item interface{}) (string, error) {
	serialized, err := json.Marshal(item)

	if err != nil {
		return "", err
	}
	return string(serialized), nil
}

func getFrameName(property *iotsitewise.DescribeAssetPropertyOutput) string {
	propertyName := util.GetPropertyName(property)

	if propertyName != "" {
		if *property.AssetName != "" {
			return fmt.Sprintf("%s %s", *property.AssetName, propertyName)
		} else {
			return propertyName
		}
	}

	if *property.AssetName != "" {
		return *property.AssetName
	}

	return ""
}
