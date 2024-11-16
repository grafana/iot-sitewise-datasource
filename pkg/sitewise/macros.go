package sitewise

import (
	"fmt"
	"strings"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data/sqlutil"
	"github.com/grafana/grafana-plugin-sdk-go/experimental/errorsource"
	"github.com/grafana/sqlds/v4"
	"github.com/pkg/errors"
)

var TableColumnsNotFoundError = errors.New("Table name not found in tableColumns")

var tableColumns = map[string][]string{
	"asset": {
		"asset_id", "asset_name", "asset_description", "asset_model_id", "asset_root_id",
	},
	"asset_property": {
		"property_id", "asset_id", "property_name", "property_data_type", "property_alias", "asset_composite_model_id",
	},
	"raw_time_series": {
		"asset_id", "property_id", "property_alias", "event_timestamp", "quality", "boolean_value", "int_value", "double_value", "string_value",
	},
	"latest_value_time_series": {
		"asset_id", "property_id", "property_alias", "event_timestamp", "quality", "boolean_value", "int_value", "double_value", "string_value",
	},
	"precomputed_aggregates": {
		"asset_id", "property_id", "property_alias", "event_timestamp", "resolution", "sum_value", "count_value", "average_value", "maximum_value", "minimum_value", "stdev_value",
	},
}

func extractTableName(query *sqlutil.Query) (string, error) {
	lowerSQL := strings.ToLower(query.RawSQL)
	fromIndex := strings.Index(lowerSQL, "from")
	if fromIndex == -1 {
		return "", errors.New("Missing FROM clause in SQL")
	}

	// Extract the part of the query after the "FROM" clause
	afterFrom := query.RawSQL[fromIndex+len("from"):]

	// Split by spaces and find the first non-empty part
	parts := strings.Fields(afterFrom)
	if len(parts) == 0 {
		return "", errors.New("Table name not found")
	}

	return parts[0], nil
}

func macroSelectAll(query *sqlutil.Query, args []string) (string, error) {
	// find the table name and return all columns
	tableName, err := extractTableName(query)
	if err != nil {
		return "selectAll", TableColumnsNotFoundError
	}
	columns, ok := tableColumns[tableName]
	if !ok {
		return "selectAll", TableColumnsNotFoundError
	}
	return strings.Join(columns, ", "), nil
}

func macroRawTimeFrom(query *sqlutil.Query, args []string) (string, error) {
	return fmt.Sprintf("%d", query.TimeRange.From.UTC().Unix()), nil
}

func macroRawTimeTo(query *sqlutil.Query, args []string) (string, error) {
	return fmt.Sprintf("%d", query.TimeRange.To.UTC().Unix()), nil
}

func macroUnixEpochFilter(query *sqlutil.Query, args []string) (string, error) {
	if len(args) != 1 {
		return "", errorsource.DownstreamError(errors.WithMessagef(sqlds.ErrorBadArgumentCount, "expected one argument"), false)
	}

	var (
		column = args[0]
		from   = query.TimeRange.From.UTC().Unix()
		to     = query.TimeRange.To.UTC().Unix()
	)
	backend.Logger.Debug("macroUnixEpochFilter", "from", query.TimeRange.From, "to", query.TimeRange.To, "from stamp", from, "to stamp", to)

	return fmt.Sprintf("%s >= %d and %s <= %d", column, from, column, to), nil
}

var macros = map[string]sqlutil.MacroFunc{
	"selectAll":       macroSelectAll,
	"rawTimeFrom":     macroRawTimeFrom,
	"rawTimeTo":       macroRawTimeTo,
	"unixEpochFilter": macroUnixEpochFilter,
}

func (s *Datasource) Macros() sqlutil.Macros {
	return macros
}
