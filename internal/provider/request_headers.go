package provider

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/deploymenttheory/go-api-http-client/httpclient"
)

var (
	errDuplicateCustomHeader = errors.New("custom header names must be unique, ignoring case")
	errBasicProxyHeader      = errors.New("a custom Authorization header requires auth_method oauth2; Jamf basic authentication needs Authorization for its token request")
)

type customHeaderModel struct {
	Name  string `tfsdk:"name"`
	Value string `tfsdk:"value"`
}

func requestHeaderConfig(headers []customHeaderModel, tokenHeader, authMethod string) (httpclient.RequestHeaderConfig, error) {
	result := httpclient.RequestHeaderConfig{BearerTokenHeader: tokenHeader, CustomHeaders: make(map[string]string, len(headers))}
	for _, header := range headers {
		name := http.CanonicalHeaderKey(header.Name)
		if _, exists := result.CustomHeaders[name]; exists {
			return httpclient.RequestHeaderConfig{}, fmt.Errorf("%w: %s", errDuplicateCustomHeader, name)
		}
		if authMethod == "basic" && strings.EqualFold(name, "Authorization") {
			return httpclient.RequestHeaderConfig{}, errBasicProxyHeader
		}
		result.CustomHeaders[name] = header.Value
	}
	return result, nil
}
