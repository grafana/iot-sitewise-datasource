package api

const (
	// Max number of entries from: https://docs.aws.amazon.com/iot-sitewise/latest/APIReference/API_BatchGetAssetPropertyAggregates.html#iotsitewise-BatchGetAssetPropertyAggregates-request-entries
	BatchGetAssetPropertyAggregatesMaxEntries = 16
	// Max results number from: https://docs.aws.amazon.com/iot-sitewise/latest/APIReference/API_BatchGetAssetPropertyValueHistory.html#iotsitewise-BatchGetAssetPropertyValueHistory-request-maxResults
	BatchGetAssetPropertyValueHistoryMaxResults = 20000

	// Max number of entries from: https://docs.aws.amazon.com/iot-sitewise/latest/APIReference/API_BatchGetAssetPropertyValueHistory.html#iotsitewise-BatchGetAssetPropertyValueHistory-request-entries
	BatchGetAssetPropertyValueHistoryMaxEntries = 16
	// Max results number from: https://docs.aws.amazon.com/iot-sitewise/latest/APIReference/API_BatchGetAssetPropertyAggregates.html#iotsitewise-BatchGetAssetPropertyAggregates-request-maxResults
	BatchGetAssetPropertyAggregatesMaxResults = 4000

	// Max number of entries from: https://docs.aws.amazon.com/iot-sitewise/latest/APIReference/API_BatchGetAssetPropertyValue.html#iotsitewise-BatchGetAssetPropertyValue-request-entries
	BatchGetAssetPropertyValueMaxEntries = 128
)
