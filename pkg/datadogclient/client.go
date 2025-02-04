// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-2021 Datadog, Inc.

package datadogclient

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/go-logr/logr"

	"github.com/DataDog/datadog-operator/pkg/config"

	datadogapiclientv1 "github.com/DataDog/datadog-api-client-go/api/v1/datadog"
	apicommon "github.com/DataDog/datadog-operator/apis/datadoghq/common"
)

const prefix = "https://api."

// DatadogClient contains the Datadog API Client and Authentication context.
type DatadogClient struct {
	Client *datadogapiclientv1.APIClient
	Auth   context.Context
}

// InitDatadogClient initializes the Datadog API Client and establishes credentials.
func InitDatadogClient(logger logr.Logger, creds config.Creds) (DatadogClient, error) {
	if creds.APIKey == "" || creds.AppKey == "" {
		return DatadogClient{}, errors.New("error obtaining API key and/or app key")
	}

	// Initialize the official Datadog V1 API client.
	authV1 := context.WithValue(
		context.Background(),
		datadogapiclientv1.ContextAPIKeys,
		map[string]datadogapiclientv1.APIKey{
			"apiKeyAuth": {
				Key: creds.APIKey,
			},
			"appKeyAuth": {
				Key: creds.AppKey,
			},
		},
	)
	configV1 := datadogapiclientv1.NewConfiguration()

	apiURL := ""
	if os.Getenv(config.DDURLEnvVar) != "" {
		apiURL = os.Getenv(config.DDURLEnvVar)
	} else if site := os.Getenv(apicommon.DDSite); site != "" {
		apiURL = prefix + strings.TrimSpace(site)
	}

	if apiURL != "" {
		logger.Info("Got API URL for DatadogOperator controller", "URL", apiURL)
		parsedAPIURL, parseErr := url.Parse(apiURL)
		if parseErr != nil {
			return DatadogClient{}, fmt.Errorf(`invalid API URL : %w`, parseErr)
		}
		if parsedAPIURL.Host == "" || parsedAPIURL.Scheme == "" {
			return DatadogClient{}, fmt.Errorf(`missing protocol or host : %s`, apiURL)
		}
		// If API URL is passed, set and use the API name and protocol on ServerIndex{1}.
		authV1 = context.WithValue(authV1, datadogapiclientv1.ContextServerIndex, 1)
		authV1 = context.WithValue(authV1, datadogapiclientv1.ContextServerVariables, map[string]string{
			"name":     parsedAPIURL.Host,
			"protocol": parsedAPIURL.Scheme,
		})
	}

	client := datadogapiclientv1.NewAPIClient(configV1)

	return DatadogClient{Client: client, Auth: authV1}, nil
}
