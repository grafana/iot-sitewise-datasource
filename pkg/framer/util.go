package framer

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iotsitewise"
	"github.com/grafana/iot-sitewise-datasource/pkg/util"
)

func getTime(ts *iotsitewise.TimeInNanos) time.Time {
	sec := *ts.TimeInSeconds

	if nanos := ts.OffsetInNanos; nanos != nil {
		return time.Unix(sec, *nanos)
	}

	return time.Unix(sec, 0)
}

func isPropertyDataTypeDefined(dataType string) bool {
	return dataType == "BOOLEAN" || dataType == "DOUBLE" || dataType == "INTEGER" || dataType == "STRING"
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

func getPropertyVariantValueType(variant *iotsitewise.Variant) string {

	if val := variant.BooleanValue; val != nil {
		return "BOOLEAN"
	}

	if val := variant.DoubleValue; val != nil {
		return "DOUBLE"
	}

	if val := variant.IntegerValue; val != nil {
		return "INTEGER"
	}

	if val := variant.StringValue; val != nil {
		return "STRING"
	}

	return ""
}

func getErrorDescription(details *iotsitewise.ErrorDetails) (*string, error) {

	if details == nil {
		return nil, nil
	}

	jb, err := serialize(*details)
	if err != nil {
		return nil, err
	}
	return aws.String(string(jb)), nil
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
