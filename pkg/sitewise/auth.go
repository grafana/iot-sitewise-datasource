package sitewise

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/grafana/iot-sitewise-datasource/pkg/models"
)

type EdgeAuthenticator struct {
	Settings models.AWSSiteWiseDataSourceSetting
}

func (a *EdgeAuthenticator) Authorize(ctx context.Context) (models.AuthInfo, error) {
	reqBody := map[string]string{
		"username":      a.Settings.EdgeAuthUser,
		"password":      a.Settings.EdgeAuthPass,
		"authMechanism": a.Settings.EdgeAuthMode,
	}
	reqBodyJson, err := json.Marshal(reqBody)
	if err != nil {
		return models.AuthInfo{}, err
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	authEndpoint := a.Settings.AWSDatasourceSettings.Endpoint + "authenticate"
	resp, err := client.Post(authEndpoint, "application/json", bytes.NewBuffer(reqBodyJson))
	if err != nil {
		return models.AuthInfo{}, err
	}
	fmt.Print(resp)

	res := make(map[string]string)
	json.NewDecoder(resp.Body).Decode(&res)

	timeLayout := time.RFC3339
	sessionExpiryTime, err := time.Parse(timeLayout, res["sessionExpiryTime"])
	if err != nil {
		return models.AuthInfo{}, err
	}

	authInfo := models.AuthInfo{
		Username:          res["username"],
		AccessKeyId:       res["accessKeyId"],
		SecretAccessKey:   res["secretAccessKey"],
		SessionToken:      res["sessionToken"],
		SessionExpiryTime: sessionExpiryTime,
		AuthMechanism:     res["authMechanism"],
	}
	return authInfo, nil
}
