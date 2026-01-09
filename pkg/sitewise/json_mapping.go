// Package sitewise contains helpers for post-processing AWS IoT SiteWise query results before returning them to Grafana.
// Responsibilities:
//   - Resolve asset property IDs to readable names
//   - Detect and parse JSON embedded in string fields
//   - Expand JSON attributes into Grafana data frame fields
//   - Normalize diagnostic contribution values for visualization
package sitewise

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/resource"
)

const (
	jsonFieldTimestamp        = "timestamp"
	jsonFieldPrediction       = "prediction"
	jsonFieldPredictionReason = "prediction_reason"
	jsonFieldDiagnostics      = "diagnostics"
	jsonFieldAnomalyScore     = "anomaly_score"
)

var requiredJSONFields = []string{
	jsonFieldTimestamp,
	jsonFieldPrediction,
	jsonFieldPredictionReason,
}

// BuildPropertyNameMap builds a lookup map of assetPropertyID -> assetPropertyName for the given asset.
// Returns an error if:
//   - assetID is empty
//   - the asset model cannot be fetched
//   - no properties are found in the asset model

func BuildPropertyNameMap(
	ctx context.Context,
	resources resource.ResourceProvider,
	assetID string,
) (map[string]string, error) {

	if assetID == "" {
		return nil, fmt.Errorf("BuildPropertyNameMap: assetId is empty")
	}

	modelResp, err := resources.AssetModel(ctx)
	if err != nil {
		return nil, fmt.Errorf("BuildPropertyNameMap: failed to describe asset model: %w", err)
	}

	propertyNameMap := make(map[string]string)

	for _, prop := range modelResp.AssetModelProperties {
		if prop.Id != nil && prop.Name != nil {
			propertyNameMap[*prop.Id] = *prop.Name
		}
	}

	if len(propertyNameMap) == 0 {
		return nil, fmt.Errorf(
			"BuildPropertyNameMap: no properties found for assetId %s",
			assetID,
		)
	}

	return propertyNameMap, nil
}

// requiresJsonParsing determines whether a given query type

func requiresJsonParsing(query models.BaseQuery) bool {
	switch query.QueryType {
	case models.QueryTypePropertyValueHistory,
		models.QueryTypePropertyAggregate,
		models.QueryTypePropertyValue:
		return true
	default:
		return false
	}
}

// ParseJSONFields scans data frames returned from SiteWise queries and dynamically expands JSON-encoded string fields into
// individual Grafana data frame fields.
//
// Behavior:
//   - Only string fields are inspected
//   - Rows that do not resemble JSON objects are skipped
//   - Valid JSON objects are flattened into separate fields
//   - Diagnostic contribution values are normalized to percentages
func ParseJSONFields(
	ctx context.Context,
	frames data.Frames,
	resources resource.ResourceProvider,
	assetID string,
) data.Frames {
	backend.Logger.Info("ParseJSONFields: starting JSON parsing", "assetID", assetID)

	// Build property ID -> readable name map once per request
	propertyNameMap, err := BuildPropertyNameMap(ctx, resources, assetID)
	if err != nil {
		backend.Logger.Error("ParseJSONFields: failed to build property name map", "err", err)
		return frames
	}
	// Resolve asset name for diagnostic field prefixes.
	// Fallback to assetID if the name cannot be fetched.
	assetName := assetID
	assetResp, err := resources.Asset(ctx)
	if err != nil {
		backend.Logger.Warn("ParseJSONFields: failed to fetch asset name, using assetID", "assetID", assetID, "err", err)
	} else if assetResp.AssetName != nil && *assetResp.AssetName != "" {
		assetName = *assetResp.AssetName
	}
	newFrames := data.Frames{}

	for _, frame := range frames {
		newFields := make([]*data.Field, 0, len(frame.Fields))
		jsonParsed := false

		for _, field := range frame.Fields {
			newFields = append(newFields, field)
			// Only string fields can contain embedded JSON payloads
			if field.Type() != data.FieldTypeString || field.Len() == 0 {
				continue
			}

			rowCount := field.Len()
			jsonFields := map[string]*data.Field{}

			for r := 0; r < rowCount; r++ {
				rawStr, ok := field.At(r).(string)
				if !ok || rawStr == "" {
					continue
				}
				// Perform a fast check to ensure the value looks like a JSON object
				if !strings.HasPrefix(rawStr, "{") {
					continue
				}

				var obj map[string]interface{}
				if err := json.Unmarshal([]byte(rawStr), &obj); err != nil {
					// Invalid JSON should not fail the query. Log and skip the corrupted row safely.
					backend.Logger.Warn("ParseJSONFields: corrupted JSON, skipping row", "err", err, "frame", frame.Name, "field", field.Name, "row", r)
					continue
				}
				jsonParsed = true
				for _, req := range requiredJSONFields {
					if _, ok := obj[req]; !ok {
						backend.Logger.Error("ParseJSONFields: missing required field", "jsonField", req, "field", field.Name, "frame", frame.Name, "row", r)
					}
				}

				for key, val := range obj {
					if key == jsonFieldDiagnostics || key == jsonFieldTimestamp || key == jsonFieldAnomalyScore {
						continue
					}
					switch v := val.(type) {
					case float64:
						if _, exists := jsonFields[key]; !exists {
							jsonFields[key] = data.NewField(key, nil, make([]float64, rowCount))
						}
						jsonFields[key].Set(r, v)
					case string:
						if _, exists := jsonFields[key]; !exists {
							jsonFields[key] = data.NewField(key, nil, make([]string, rowCount))
						}
						jsonFields[key].Set(r, v)
					case bool:
						if _, exists := jsonFields[key]; !exists {
							jsonFields[key] = data.NewField(key, nil, make([]bool, rowCount))
						}
						jsonFields[key].Set(r, v)
					}
				}

				if v, ok := obj[jsonFieldAnomalyScore].(float64); ok {
					if _, exists := jsonFields[jsonFieldAnomalyScore]; !exists {
						jsonFields[jsonFieldAnomalyScore] = data.NewField(jsonFieldAnomalyScore, nil, make([]float64, rowCount))
					}
					jsonFields[jsonFieldAnomalyScore].Set(r, v)
				}

				// Diagnostics contain per-property contribution values.
				// These are expanded into separate fields using the format:
				//  contrib_<assetName>_<propertyName>
				// Contribution values are normalized to percentages.
				if diagArr, ok := obj[jsonFieldDiagnostics].([]interface{}); ok {
					contribValues := map[string]float64{}
					for _, item := range diagArr {
						diagObj, ok := item.(map[string]interface{})
						if !ok {
							continue
						}
						rawName, ok := diagObj["name"].(string)
						if !ok || rawName == "" {
							backend.Logger.Warn(
								"ParseJSONFields: diagnostic name is not a valid string",
								"frame", frame.Name,
								"field", field.Name,
								"row", r,
							)
							continue
						}
						parts := strings.Split(rawName, "\\")
						if len(parts) < 2 {
							continue
						}
						propertyID := parts[1]
						readable := propertyID
						if mapped, ok := propertyNameMap[propertyID]; ok {
							readable = mapped
						}
						fieldName := "contrib_" + assetName + "_" + readable
						if _, exists := jsonFields[fieldName]; !exists {
							jsonFields[fieldName] = data.NewField(fieldName, nil, make([]float64, rowCount))
						}
						if v, ok := diagObj["value"].(float64); ok {
							contribValues[fieldName] = v
						}
					}
					// Normalize contribution values so the total equals 100%
					total := 0.0
					for _, v := range contribValues {
						total += v
					}
					if total > 0 {
						for fieldName, rawValue := range contribValues {
							jsonFields[fieldName].Set(r, (rawValue/total)*100.0)
						}
					}
				}
			}

			// Append parsed JSON fields
			for _, f := range jsonFields {
				newFields = append(newFields, f)
			}
		}

		newFrame := data.NewFrame(frame.Name, newFields...)
		newFrame.Meta = frame.Meta
		newFrames = append(newFrames, newFrame)

		if jsonParsed {
			backend.Logger.Info("Parsed JSON in frame", "frame", frame.Name)
		}
	}

	return newFrames
}
