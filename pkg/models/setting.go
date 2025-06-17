package models

import (
	"encoding/json"
	"fmt"

	"github.com/grafana/grafana-aws-sdk/pkg/awsds"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

const EDGE_REGION string = "Edge"
const EDGE_AUTH_MODE_DEFAULT string = "default"
const EDGE_AUTH_MODE_LDAP string = "ldap"
const EDGE_AUTH_MODE_LINUX string = "linux"

type AWSSiteWiseDataSourceSetting struct {
	awsds.AWSDatasourceSettings
	Cert         string `json:"-"`
	EdgeAuthMode string `json:"edgeAuthMode"`
	EdgeAuthUser string `json:"edgeAuthUser"`
	EdgeAuthPass string `json:"-"`
}

func (s *AWSSiteWiseDataSourceSetting) Load(config backend.DataSourceInstanceSettings) error {
	if len(config.JSONData) > 1 {
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

	// Make sure to set an auth mode
	if s.Region == EDGE_REGION && s.EdgeAuthMode == "" {
		s.EdgeAuthMode = EDGE_AUTH_MODE_DEFAULT
	}

	s.AccessKey = config.DecryptedSecureJSONData["accessKey"]
	s.SecretKey = config.DecryptedSecureJSONData["secretKey"]
	s.Cert = config.DecryptedSecureJSONData["cert"]
	s.EdgeAuthPass = config.DecryptedSecureJSONData["edgeAuthPass"]
	return nil
}

func (s *AWSSiteWiseDataSourceSetting) Validate() error {
	if s.Region != EDGE_REGION {
		return nil
	}

	if s.Endpoint == "" {
		return fmt.Errorf("edge region requires an explicit endpoint")
	}
	if s.Cert == "" {
		return fmt.Errorf("edge region requires an SSL certificate")
	}

	if s.EdgeAuthMode != EDGE_AUTH_MODE_DEFAULT {
		if s.EdgeAuthUser == "" {
			return fmt.Errorf("missing edge auth user")
		}
		if s.EdgeAuthPass == "" {
			return fmt.Errorf("missing edge auth password")
		}
	}

	return nil
}

func (s *AWSSiteWiseDataSourceSetting) ToAWSDatasourceSettings() awsds.AWSDatasourceSettings {
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
		SessionToken:  s.SessionToken,
	}

	return cfg
}
