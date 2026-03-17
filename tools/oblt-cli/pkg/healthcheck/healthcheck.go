package healthcheck

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/gcp"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/logger"
)

// Struct to hold the response of a healthcheck
type HealthcheckResponse struct {
	Service string
	Url     string
	Res     []byte
	Err     error
}

// Executes healthcheck for Elasticsearch
func ElasticsearchHC(clusterState gcp.ClusterStateSecret) *HealthcheckResponse {
	esUrl := clusterState.EsSecret.Url
	esPwd := clusterState.EsSecret.Password
	esUser := clusterState.EsSecret.Username
	if IsNotBlank(esUrl) && IsNotBlank(esPwd) && IsNotBlank(esUser) {
		hcResponse := executeHC(esUrl, "elasticsearch", func(req http.Request) {
			auth := fmt.Sprintf("%s:%s", esUser, esPwd)
			basicAuth := base64.StdEncoding.EncodeToString([]byte(auth))
			req.Header.Add("Authorization", "Basic "+basicAuth)
		})
		return hcResponse
	}
	return nil
}

// Executes healthcheck for Kibana
func KibanaHC(clusterState gcp.ClusterStateSecret) *HealthcheckResponse {
	kibanaUrl := clusterState.KibanaSecret.Url
	kibanaPwd := clusterState.KibanaSecret.Password
	kibanaUser := clusterState.KibanaSecret.Username
	if IsNotBlank(kibanaUrl) && IsNotBlank(kibanaPwd) && IsNotBlank(kibanaUser) {
		hcResponse := executeHC(kibanaUrl+"/api/status", "kibana", func(req http.Request) {
			auth := fmt.Sprintf("%s:%s", kibanaUser, kibanaPwd)
			basicAuth := base64.StdEncoding.EncodeToString([]byte(auth))
			req.Header.Add("Authorization", "Basic "+basicAuth)
		})
		return hcResponse
	}
	return nil
}

// Executes healthcheck APM
func ApmHC(clusterState gcp.ClusterStateSecret) *HealthcheckResponse {
	apmUrl := clusterState.ApmSecret.Url
	apmToken := clusterState.ApmSecret.Token
	if IsNotBlank(apmUrl) && IsNotBlank(apmToken) {
		hcResponse := executeHC(apmUrl, "apm", func(req http.Request) {
			req.Header.Add("Authorization", "Bearer "+apmToken)
		})
		return hcResponse
	}
	return nil
}

// Executes an arbitrary healthcheck with specific auth
func executeHC(url string, service string, setAuth func(http.Request)) *HealthcheckResponse {
	hcError := func(err error, msg string) *HealthcheckResponse {
		return &HealthcheckResponse{
			Service: service,
			Url:     url,
			Res:     nil,
			Err:     fmt.Errorf("%s: %w", msg, err),
		}
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return hcError(err, "error creating request")
	}
	setAuth(*req)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return hcError(err, "error sending request")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return hcError(err, "error reading response body")
	}

	return &HealthcheckResponse{
		Service: service,
		Url:     url,
		Res:     body,
		Err:     nil,
	}
}

// Prints a healthcheck response in a pretty way
func PrettyPrintHC(hc *HealthcheckResponse) string {
	result := fmt.Sprintf("Service: %s\n", hc.Service)
	result += fmt.Sprintf("Url: %s\n", hc.Url)
	if hc.Err != nil {
		result += fmt.Sprintf("Error: %s\n", hc.Err.Error())
	} else {
		prettyBody, err := tryPrettyJson(hc.Res)
		if err != nil {
			result += fmt.Sprintf("Response: %s\n", string(hc.Res))
		} else {
			result += fmt.Sprintf("Response: %s\n", string(prettyBody))
		}
	}
	return result
}

// Tries to pretty print a json response
func tryPrettyJson(body []byte) ([]byte, error) {
	var jsonData map[string]interface{}
	err := json.Unmarshal(body, &jsonData)
	if err != nil {
		return nil, err
	}

	prettyJSON, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		return nil, err
	}

	return prettyJSON, nil
}

// saves the results of a healthcheck map for file output
func ResponsesToJson(checklist []*HealthcheckResponse) map[string]interface{} {
	fileReport := make(map[string]interface{})
	for _, hcResponse := range checklist {
		if hcResponse.Err != nil {
			fileReport[hcResponse.Service] = hcResponse.Err.Error()
		} else {
			var jsonData map[string]interface{}
			err := json.Unmarshal(hcResponse.Res, &jsonData)
			if err != nil {
				logger.Infof("Error parsing json response: %s", err)
				fileReport[hcResponse.Service] = string(hcResponse.Res)
			} else {
				fileReport[hcResponse.Service] = jsonData
			}
		}
	}
	return fileReport
}
