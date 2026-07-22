package api

import (
	"testing"

	"github.com/grafana/iot-sitewise-datasource/pkg/models"
	"github.com/grafana/iot-sitewise-datasource/pkg/util"
)

func TestGetNextToken(t *testing.T) {
	const (
		assetID    = "11111111-1111-1111-1111-111111111111"
		propertyID = "22222222-2222-2222-2222-222222222222"
		alias      = "/my/property/alias"
	)
	batchedEntryID := *util.GetEntryIdFromAssetProperty(assetID, propertyID)
	aliasEntryID := *util.GetEntryIdFromPropertyAlias(alias)

	tests := []struct {
		name  string
		query models.BaseQuery
		want  *string
	}{
		{
			name:  "no token returns nil",
			query: models.BaseQuery{},
			want:  nil,
		},
		{
			name:  "single-entry non-empty token is forwarded",
			query: models.BaseQuery{NextToken: "abc123"},
			want:  stringPtr("abc123"),
		},
		{
			name:  "single-entry empty token returns nil",
			query: models.BaseQuery{NextToken: ""},
			want:  nil,
		},
		{
			name: "batched entry with non-empty token is forwarded",
			query: models.BaseQuery{
				AssetPropertyEntries: []models.AssetPropertyEntry{
					{AssetId: assetID, PropertyId: propertyID},
				},
				NextTokens: map[string]string{batchedEntryID: "page2token"},
			},
			want: stringPtr("page2token"),
		},
		{
			// Regression test for https://github.com/grafana/iot-sitewise-datasource/issues/765
			// A map miss (entry that finished paginating on a prior page) yields "".
			// getNextToken must return nil, not aws.String(""), to avoid a 400
			// InvalidRequestException from BatchGetAssetPropertyValueHistory.
			name: "batched entry with map miss returns nil (issue #765)",
			query: models.BaseQuery{
				AssetPropertyEntries: []models.AssetPropertyEntry{
					{AssetId: assetID, PropertyId: propertyID},
				},
				NextTokens: map[string]string{"some-other-entry": "page2token"},
			},
			want: nil,
		},
		{
			name: "batched entry with explicitly empty token returns nil (issue #765)",
			query: models.BaseQuery{
				AssetPropertyEntries: []models.AssetPropertyEntry{
					{AssetId: assetID, PropertyId: propertyID},
				},
				NextTokens: map[string]string{batchedEntryID: ""},
			},
			want: nil,
		},
		{
			name: "batched entry keyed by property alias with non-empty token is forwarded",
			query: models.BaseQuery{
				AssetPropertyEntries: []models.AssetPropertyEntry{
					{PropertyAlias: alias},
				},
				NextTokens: map[string]string{aliasEntryID: "aliasToken"},
			},
			want: stringPtr("aliasToken"),
		},
		{
			name: "batched entry keyed by property alias with map miss returns nil (issue #765)",
			query: models.BaseQuery{
				AssetPropertyEntries: []models.AssetPropertyEntry{
					{PropertyAlias: alias},
				},
				NextTokens: map[string]string{"some-other-entry": "aliasToken"},
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getNextToken(tt.query)
			switch {
			case tt.want == nil && got != nil:
				t.Fatalf("expected nil nextToken, got %q", *got)
			case tt.want != nil && got == nil:
				t.Fatalf("expected nextToken %q, got nil", *tt.want)
			case tt.want != nil && got != nil && *tt.want != *got:
				t.Fatalf("expected nextToken %q, got %q", *tt.want, *got)
			}
		})
	}
}

func stringPtr(s string) *string {
	return &s
}
