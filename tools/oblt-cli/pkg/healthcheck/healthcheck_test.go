package healthcheck

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/gcp"
	"github.com/stretchr/testify/assert"
)

func TestElasticsearchHCOkStatus(t *testing.T) {

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		responseBody := []byte(`{"status": "ok"}`)
		_, _ = w.Write(responseBody)
	}))
	defer mockServer.Close()

	clusterState := gcp.ClusterStateSecret{
		EsSecret: gcp.EsSecret{
			ApiKey:   "key==",
			Url:      mockServer.URL,
			Password: "password",
			Username: "username",
		},
	}

	hcResponse := ElasticsearchHC(clusterState)

	assert.NotNil(t, hcResponse)
	assert.Equal(t, "elasticsearch", hcResponse.Service)
	assert.Equal(t, mockServer.URL, hcResponse.Url)
	assert.Nil(t, hcResponse.Err)
	assert.NotNil(t, hcResponse.Res)
	assert.Equal(t, `{"status": "ok"}`, string(hcResponse.Res))
}

func TestElasticsearchHCServerError(t *testing.T) {
	errorMessage := "Special error message"

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(errorMessage))
	}))

	defer mockServer.Close()

	clusterState := gcp.ClusterStateSecret{
		EsSecret: gcp.EsSecret{
			ApiKey:   "key==",
			Url:      mockServer.URL,
			Password: "password",
			Username: "username",
		},
	}

	hcResponse := ElasticsearchHC(clusterState)

	assert.NotNil(t, hcResponse)
	assert.Equal(t, "elasticsearch", hcResponse.Service)
	assert.Equal(t, mockServer.URL, hcResponse.Url)
	assert.Nil(t, hcResponse.Err)
	assert.Equal(t, errorMessage, string(hcResponse.Res))

}

func TestElasticsearchHCInvalidUrl(t *testing.T) {
	clusterState := gcp.ClusterStateSecret{
		EsSecret: gcp.EsSecret{
			ApiKey:   "key==",
			Url:      "example.com",
			Password: "password",
			Username: "username",
		},
	}

	hcResponse := ElasticsearchHC(clusterState)

	assert.NotNil(t, hcResponse)
	assert.Equal(t, "elasticsearch", hcResponse.Service)
	assert.NotNil(t, hcResponse.Err)
	assert.Contains(t, hcResponse.Err.Error(), "example.com")
}

func TestElasticsearchHCNilCredentials(t *testing.T) {
	clusterState := gcp.ClusterStateSecret{
		EsSecret: gcp.EsSecret{
			ApiKey:   "key==",
			Url:      "",
			Password: "password",
			Username: "username",
		},
	}

	hcResponse := ElasticsearchHC(clusterState)

	assert.Nil(t, hcResponse)
}

func TestElasticsearchHCWhitespaceCredentials(t *testing.T) {
	clusterState := gcp.ClusterStateSecret{
		EsSecret: gcp.EsSecret{
			ApiKey:   "key==",
			Url:      " ",
			Password: "password",
			Username: "username",
		},
	}

	hcResponse := ElasticsearchHC(clusterState)

	assert.Nil(t, hcResponse)
}

func TestResponsesToJsonFunc(t *testing.T) {
	checklist := []*HealthcheckResponse{
		{
			Service: "elasticsearch",
			Url:     "http://example.com/elasticsearch",
			Res:     []byte(`{"status":"ok"}`),
			Err:     nil,
		},
		{
			Service: "kibana",
			Url:     "http://example.com/kibana",
			Res:     []byte(`{"status":"ok"}`),
			Err:     fmt.Errorf("Failed to connect"),
		},
	}

	result := ResponsesToJson(checklist)

	assert.Equal(t, 2, len(result))
	assert.Equal(t, "ok", result["elasticsearch"].(map[string]interface{})["status"])
	assert.Equal(t, "Failed to connect", result["kibana"].(string))
}
