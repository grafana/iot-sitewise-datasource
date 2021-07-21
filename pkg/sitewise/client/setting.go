package client

import (
	"encoding/json"
	"fmt"

	"github.com/grafana/grafana-aws-sdk/pkg/awsds"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

type AWSSiteWiseDataSourceSetting struct {
	awsds.AWSDatasourceSettings
	Cert         string `json:"-"`
	EdgeAuthMode string `json:"edgeAuthMode"`
	EdgeAuthUser string `json:"edgeAuthUser"`
	EdgeAuthPass string `json:"-"`
}

func (s *AWSSiteWiseDataSourceSetting) Load(config backend.DataSourceInstanceSettings) error {
	if config.JSONData != nil && len(config.JSONData) > 1 {
		if err := json.Unmarshal(config.JSONData, s); err != nil {
			return fmt.Errorf("could not unmarshal DatasourceSettings json: %w", err)
		}
	}

	if s.Region == "default" || s.Region == "" {
		s.Region = s.DefaultRegion
	}

	if s.Profile == "" {
		s.Profile = config.Database // legacy support (only for cloudwatch?)
	}

	s.AccessKey = config.DecryptedSecureJSONData["accessKey"]
	s.SecretKey = config.DecryptedSecureJSONData["secretKey"]
	s.Cert = config.DecryptedSecureJSONData["cert"]
	s.EdgeAuthPass = config.DecryptedSecureJSONData["edgeAuthPass"]
	return nil
}

func (s *AWSSiteWiseDataSourceSetting) toAWSDatasourceSettings() awsds.AWSDatasourceSettings {
	cfg := awsds.AWSDatasourceSettings{
		Profile:       s.Profile,
		Region:        s.Region,
		AuthType:      s.AuthType,
		AssumeRoleARN: s.AssumeRoleARN,
		ExternalID:    s.ExternalID,
		Endpoint:      s.Endpoint,
		DefaultRegion: s.DefaultRegion,
		AccessKey:     s.AccessKey,
		SecretKey:     s.SecretKey,
	}

	return cfg
}
