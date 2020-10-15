package sitewise

import (
	framerImpl "github.com/grafana/iot-sitewise-datasource/pkg/framer"
	"github.com/grafana/iot-sitewise-datasource/pkg/resource"

	"github.com/grafana/iot-sitewise-datasource/pkg/models"

	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/client"

	"github.com/grafana/iot-sitewise-datasource/pkg/sitewise/framer"
)

func framePropertyValueResponse(query *models.AssetPropertyValueQuery, data framer.FrameData, client client.Client) framer.Framer {

	rp := resource.NewSitewiseResourceProvider(client)
	mp := framerImpl.NewPropertyValueMetaProvider(rp, *query)
	fr := &framerImpl.PropertyValueQueryFramer{
		FrameData:    data,
		MetaProvider: mp,
		Request:      *query,
	}

	return fr
}

//func getRegion(ctx backend.PluginContext, q models.BaseQuery) (*string, error) {
//
//	if q.AwsRegion != "" {
//		return &q.AwsRegion, nil
//	}
//
//	settings, err := gaws.LoadSettings(*ctx.DataSourceInstanceSettings)
//	if err != nil {
//		return nil, err
//	}
//
//	if settings.DefaultRegion == "" {
//		return nil, errors.New("unable to determine aws region for region")
//	}
//
//	return &settings.DefaultRegion, nil
//
//}
