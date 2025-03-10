package test

import (
	"github.com/grafana/iot-sitewise-datasource/pkg/util"
)

var mockAssetId = "1assetid-aaaa-2222-bbbb-3333cccc4444"
var mockPropertyId = "11propid-aaaa-2222-bbbb-3333cccc4444"
var mockPropertyAlias = "/amazon/renton/1/rpm"
var mockAssetPropertyEntryId = util.GetEntryIdFromAssetProperty(mockAssetId, mockPropertyId)
var mockPropertyAliasEntryId = util.GetEntryIdFromPropertyAlias(mockPropertyAlias)
