// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http:// www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

// Package gcp contains the functions related to the Google Cloud Platform.
// This file contains the functions to access the cluster secrets.
package gcp

import (
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/logger"
	"gopkg.in/yaml.v3"
)

var environmentsMatch = map[string]string{
	"production": "pro",
	"pro":        "pro",
	"staging":    "staging",
	"qa":         "qa",
}

// EssCreds is the map of the ESS credentials secrets.
type EsSecret struct {
	ApiKey   string `json:"elasticsearch_apikey" yaml:"elasticsearch_apikey"`
	Url      string `json:"elasticsearch_url" yaml:"elasticsearch_url"`
	Username string `json:"elasticsearch_username" yaml:"elasticsearch_username"`
	Password string `json:"elasticsearch_password" yaml:"elasticsearch_password"`
	Name     string `json:"cluster_name" yaml:"cluster_name"`
}

// NewEsSecret returns a new instance of EsSecret.
func NewEsSecret() EsSecret {
	return EsSecret{
		ApiKey:   "VnVhQ2ZHY0JDZGJrUW0tZTVhT3g6dWkybHAyYXhUTm1zeWFrdzl0dk5udw==",
		Url:      "https://es.example.com",
		Username: "elastic",
		Password: "changeme",
		Name:     "cluster-name",
	}
}

// KibanaSecret is the map of the Kibana secrets.
type KibanaSecret struct {
	Url      string `json:"kibana_url" yaml:"kibana_url"`
	Username string `json:"kibana_username" yaml:"kibana_username"`
	Password string `json:"kibana_password" yaml:"kibana_password"`
}

// NewKibanaSecret returns a new instance of KibanaSecret.
func NewKibanaSecret() KibanaSecret {
	return KibanaSecret{
		Url:      "https://kibana.example.com",
		Username: "elastic",
		Password: "changeme",
	}
}

// ApmSecret is the map of the APM secrets.
type ApmSecret struct {
	Url       string `json:"apm_url" yaml:"apm_url"`
	Token     string `json:"apm_token" yaml:"apm_token"`
	ApmApiKey string `json:"apm_apikey" yaml:"apm_apikey"`
}

// NewApmSecret returns a new instance of ApmSecret.
func NewApmSecret() ApmSecret {
	return ApmSecret{
		Url:   "https://apm.example.com",
		Token: "changeme",
	}
}

// FleetSecret is the map of the Fleet secrets.
type FleetSecret struct {
	Url            string `json:"fleet_url" yaml:"fleet_url"`
	ApiToken       string `json:"fleet_api_token" yaml:"fleet_api_token"`
	EnrollToken    string `json:"fleet_enroll_token" yaml:"fleet_enroll_token"`
	PolicyId       string `json:"fleet_policy_id" yaml:"fleet_policy_id"`
	ServerPolicyId string `json:"fleet_server_policy_id" yaml:"fleet_server_policy_id"`
	ServerToken    string `json:"fleet_server_token" yaml:"fleet_server_token"`
}

// NewFleetSecret returns a new instance of FleetSecret.
func NewFleetSecret() FleetSecret {
	return FleetSecret{
		Url:            "https://fleet.example.com",
		ApiToken:       "changeme",
		EnrollToken:    "changeme",
		PolicyId:       "my-policy",
		ServerPolicyId: "my-server-policy",
		ServerToken:    "changeme",
	}
}

// ESSDeployment is the map of the ESS deployment secrets.
type ESSDeployment struct {
	Id      string `json:"ess_deployment_id" yaml:"ess_deployment_id" `
	CloudId string `json:"ess_cloud_id" yaml:"ess_cloud_id"`
}

// NewESSDeployment returns a new instance of ESSDeployment.
func NewESSDeployment() ESSDeployment {
	return ESSDeployment{
		Id:      "es-id",
		CloudId: "ess-cloud-id",
	}
}

// ESSCredentials is the map of the ESS credentials secrets.
type ESSCredentials struct {
	ApiKey         string `json:"apiKey" yaml:"apiKey"`
	Url            string `json:"ess_endpoint" yaml:"ess_endpoint"`
	Console        string `json:"console" yaml:"console"`
	ConsoleApiKey  string `json:"consoleApiKey" yaml:"consoleApiKey"`
	OrganizationID string `json:"organizationID" yaml:"organizationID"`
	Password       string `json:"password" yaml:"password"`
	Username       string `json:"username" yaml:"username"`
}

// NewESSCredentials returns a new instance of ESSCredentials.
func NewESSCredentials() ESSCredentials {
	return ESSCredentials{
		ApiKey:         "changeme",
		Url:            "https://ess.example.com",
		Console:        "https://console.example.com",
		ConsoleApiKey:  "changeme",
		OrganizationID: "changeme",
		Password:       "changeme",
		Username:       "changeme",
	}
}

// ServerlessCredentials is the map of the serverless credentials secrets.
type ServerlessCredentials struct {
	Endpoint       string `json:"serverless_endpoint" yaml:"serverless_endpoint"`
	ApiKey         string `json:"apiKey" yaml:"apiKey"`
	Console        string `json:"console" yaml:"console"`
	ConsoleApiKey  string `json:"consoleApiKey" yaml:"consoleApiKey"`
	OrganizationID string `json:"organizationID" yaml:"organizationID"`
	Password       string `json:"password" yaml:"password"`
	Username       string `json:"username" yaml:"username"`
}

// NewServerlessCredentials returns a new instance of ServerlessCredentials.
func NewServerlessCredentials() ServerlessCredentials {
	return ServerlessCredentials{
		Endpoint:       "https://serverless.example.com",
		ApiKey:         "changeme",
		Console:        "https://console.example.com",
		ConsoleApiKey:  "changeme",
		OrganizationID: "changeme",
		Password:       "changeme",
		Username:       "changeme",
	}
}

// ServerlessProject is the map of the serverless project secrets.
type ServerlessProject struct {
	ServerlessDeploymentId string `json:"serverless_deployment_id" yaml:"serverless_deployment_id"`
}

// NewServerlessProject returns a new instance of ServerlessProject.
func NewServerlessProject() ServerlessProject {
	return ServerlessProject{
		ServerlessDeploymentId: "serverless-deployment-id",
	}
}

// StackInfo is the map of the stack info secrets.
type StackInfo struct {
	StackMode     string `json:"stack_mode" yaml:"stack_mode"`
	StackVersion  string `json:"stack_version" yaml:"stack_version"`
	StackTarget   string `json:"stack_target" yaml:"stack_target"`
	StackTemplate string `json:"stack_template" yaml:"stack_template"`
}

// NewStackInfo returns a new instance of StackInfo.
func NewStackInfo() StackInfo {
	return StackInfo{
		StackMode:     "esse",
		StackVersion:  "8.11.3",
		StackTarget:   "production",
		StackTemplate: "observability",
	}
}

// K8sInfo is the map of the k8s info secrets.
type K8sInfo struct {
	K8sDefaultNamespace string `json:"k8s_default_namespace" yaml:"k8s_default_namespace"`
	K8sProject          string `json:"k8s_project" yaml:"k8s_project"`
	K8sProvider         string `json:"k8s_provider" yaml:"k8s_provider"`
	K8sRegion           string `json:"k8s_region" yaml:"k8s_region"`
	K8sDomain           string `json:"k8s_domain" yaml:"k8s_domain"`
	IngressAvailable    string `json:"ingress_available" yaml:"ingress_available"`
	IngressPassword     string `json:"ingress_password" yaml:"ingress_password"`
	IngressUsername     string `json:"ingress_username" yaml:"ingress_username"`
	KubeSystemNamespace string `json:"kube_system_namespace" yaml:"kube_system_namespace"`
	LoadBalancerIp      string `json:"load_balancer_ip" yaml:"load_balancer_ip"`
}

// NewK8sInfo returns a new instance of K8sInfo.
func NewK8sInfo() K8sInfo {
	return K8sInfo{
		K8sDefaultNamespace: "default",
		K8sProject:          "my-project",
		K8sProvider:         "gcp",
		K8sRegion:           "us-central1",
		K8sDomain:           "example.com",
		IngressAvailable:    "true",
		IngressPassword:     "changeme",
		IngressUsername:     "changeme",
		KubeSystemNamespace: "kube-system",
		LoadBalancerIp:      "127.0.0.1",
	}
}

type ProfilingSecret struct {
	ProfilingToken string `json:"profiling_token" yaml:"profiling_token"`
	ProfilingUrl   string `json:"profiling_url" yaml:"profiling_url"`
}

// ClusterStateSecret is the map of the cluster state secrets.
type ClusterStateSecret struct {
	ApmSecret         `json:",inline" yaml:",inline"`
	EsSecret          `json:",inline" yaml:",inline"`
	FleetSecret       `json:",inline" yaml:",inline"`
	KibanaSecret      `json:",inline" yaml:",inline"`
	StackInfo         `json:",inline" yaml:",inline"`
	K8sInfo           `json:",inline" yaml:",inline"`
	ServerlessProject `json:",inline" yaml:",inline"`
	ESSDeployment     `json:",inline" yaml:",inline"`
	ProfilingSecret   `json:",inline" yaml:",inline"`
	ClusterCreated    bool `json:"cluster_created" yaml:"cluster_created"`
}

// NewClusterStateSecret returns a new instance of ClusterStateSecret.
func NewClusterStateSecret() ClusterStateSecret {
	return ClusterStateSecret{
		ApmSecret:         NewApmSecret(),
		EsSecret:          NewEsSecret(),
		FleetSecret:       NewFleetSecret(),
		KibanaSecret:      NewKibanaSecret(),
		StackInfo:         NewStackInfo(),
		K8sInfo:           NewK8sInfo(),
		ServerlessProject: NewServerlessProject(),
		ESSDeployment:     NewESSDeployment(),
	}
}

// SecretUnmarsahler is the interface to unmarshal the secrets.
// https://go.dev/doc/tutorial/generics
type SecretUnmarsahler interface {
	ClusterStateSecret | EsSecret | KibanaSecret | ApmSecret | FleetSecret | ESSDeployment | ESSCredentials | ServerlessCredentials | ServerlessProject | StackInfo | K8sInfo
}

// ClusterSecrets is the struct to interact with the cluster secrets.
type ClusterSecrets struct {
	client SecretsManager
	auth   Auth
}

// NewClusterSecrets returns a new instance of ClusterSecrets.
func NewClusterSecrets() (clusterSecrets *ClusterSecrets, err error) {
	return NewClusterSecretsWithProjectId(default_project_id)
}

// NewClusterSecretsWithProjectId returns a new instance of ClusterSecrets with the given project ID.
func NewClusterSecretsWithProjectId(projectId string) (clusterSecrets *ClusterSecrets, err error) {
	auth := NewGcpAuthWithProjectId(projectId)
	auth.Authenticate()
	gcsmClient, err := NewClientWithOpts(auth.Options...)
	if err != nil {
		logger.Errorf("failed to create secretmanager client: %v", err)
	} else {
		clusterSecrets = NewClusterSecretsWithClient(gcsmClient, &auth)
	}
	return clusterSecrets, err
}

// NewClusterSecretsWithClient returns a new instance of ClusterSecrets with the given client.
func NewClusterSecretsWithClient(client SecretsManager, auth Auth) (clusterSecrets *ClusterSecrets) {
	return &ClusterSecrets{
		client: client,
		auth:   auth,
	}
}

// Close closes the secrets manager client.
func (s *ClusterSecrets) Close() error {
	return s.client.Close()
}

// getProjectPath returns the project path.
func (s *ClusterSecrets) getProjectPath() string {
	return "projects/" + s.auth.GetProjectId()
}

// getSecretsRootPath returns the root path for the secrets.
func (s *ClusterSecrets) getSecretsRootPath() string {
	return s.getProjectPath() + "/secrets"
}

// getSecretsClusterBasePath returns the base path for the cluster secrets.
func (s *ClusterSecrets) getSecretsClusterBasePath() string {
	return s.getSecretsRootPath() + "/oblt-clusters_"
}

// getClusterStateSecretPath returns the path for the cluster state secret for the given cluster.
func (s *ClusterSecrets) getClusterStateSecretPath(clusterName string) string {
	return s.getSecretsClusterBasePath() + clusterName + "_cluster-state"
}

// getDeployInfoSecretPath returns the path for the deploy info secret for the given cluster.
func (s *ClusterSecrets) getDeployInfoSecretPath(clusterName string) string {
	return s.getSecretsClusterBasePath() + clusterName + "_deploy-info"
}

// getKibanaYamlSecretPath returns the path for the kibana yaml secret for the given cluster.
func (s *ClusterSecrets) getKibanaYamlSecretPath(clusterName string) string {
	return s.getSecretsClusterBasePath() + clusterName + "_kibana-yml"
}

// getCredentialsSecretPath returns the path for the credentials secret for the given cluster.
func (s *ClusterSecrets) getCredentialsSecretPath(clusterName string) string {
	return s.getSecretsClusterBasePath() + clusterName + "_credentials"
}

// getEnvSecretPath returns the path for the Elastic Stack environment variables secret for the given cluster.
func (s *ClusterSecrets) getEnvSecretPath(clusterName string) string {
	return s.getSecretsClusterBasePath() + clusterName + "_env"
}

// getActivateSecretPath returns the path for the activate secret for the given cluster.
func (s *ClusterSecrets) getActivateSecretPath(clusterName string) string {
	return s.getSecretsClusterBasePath() + clusterName + "_activate"
}

// getESSCredentialsPath returns the path for the ESS credentials secret for the given environment.
func (s *ClusterSecrets) getESSCredentialsPath(environment string) string {
	return s.getSecretsRootPath() + "/elastic-cloud-observability-team-" + environmentsMatch[environment]
}

// getServerlessCredentialsPath returns the path for the serverless credentials secret for the given environment.
func (s *ClusterSecrets) getServerlessCredentialsPath(environment string) string {

	return s.getSecretsRootPath() + "/elastic-cloud-observability-team-" + environmentsMatch[environment]
}

func (s *ClusterSecrets) getLicenseSecretPath(licenseType string) string {
	return s.getSecretsRootPath() + "/elastic-stack-license-" + licenseType
}

// CreateClusterStateSecret creates a new cluster state secret with the given value.
func (s *ClusterSecrets) CreateClusterStateSecret(secret ClusterStateSecret) (err error) {
	return s.client.CreateSecret(s.getClusterStateSecretPath(secret.Name), secret)
}

// ListClusterSecrets returns the list of the cluster secrets.
func (s *ClusterSecrets) ListClusterSecrets(clusterName string, fullPath bool) (secrets []string, err error) {
	secretsPaths, err := s.client.ListSecrets(s.getProjectPath(), "labels.cluster_name="+clusterName)
	if err == nil {
		for _, item := range secretsPaths {
			path := item
			if !fullPath {
				path = s.client.GetSecretName(item)
			}
			secrets = append(secrets, path)
		}
	}
	return secrets, err
}

func unmarshallSecret[T SecretUnmarsahler](client SecretsManager, secretsPath string, value *T) (err error) {
	data, err := client.GetSecret(secretsPath)
	if err == nil {
		err = yaml.Unmarshal([]byte(data), value)
	}
	if err != nil {
		logger.Errorf("failed to read the secret %s: %v", secretsPath, err)
	}
	return err
}

// ReadClusterStateSecret returns the value of the cluster state secret.
func (s *ClusterSecrets) ReadClusterStateSecret(clusterName string) (value ClusterStateSecret, err error) {
	err = unmarshallSecret(s.client, s.getClusterStateSecretPath(clusterName), &value)
	return value, err
}

// ReadEsSecret returns the value of the ES secret.
func (s *ClusterSecrets) ReadEsSecret(clusterName string) (value EsSecret, err error) {
	err = unmarshallSecret(s.client, s.getClusterStateSecretPath(clusterName), &value)
	return value, err
}

// ReadKibanaSecret returns the value of the Kibana secret.
func (s *ClusterSecrets) ReadKibanaSecret(clusterName string) (value KibanaSecret, err error) {
	err = unmarshallSecret(s.client, s.getClusterStateSecretPath(clusterName), &value)
	return value, err
}

// ReadApmSecret returns the value of the APM secret.
func (s *ClusterSecrets) ReadApmSecret(clusterName string) (value ApmSecret, err error) {
	err = unmarshallSecret(s.client, s.getClusterStateSecretPath(clusterName), &value)
	return value, err
}

// ReadFleetSecret returns the value of the Fleet secret.
func (s *ClusterSecrets) ReadFleetSecret(clusterName string) (value FleetSecret, err error) {
	err = unmarshallSecret(s.client, s.getClusterStateSecretPath(clusterName), &value)
	return value, err
}

// ReadStackInfo returns the value of the stack info secret.
func (s *ClusterSecrets) ReadLicense(licenseType string) (value string, err error) {
	return s.client.GetSecret(s.getLicenseSecretPath(licenseType))
}

// ReadStackInfo returns the value of the stack info secret.
func (s *ClusterSecrets) ReadESSDeployment(clusterName string) (value ESSDeployment, err error) {
	err = unmarshallSecret(s.client, s.getClusterStateSecretPath(clusterName), &value)
	return value, err
}

// ReadESSCredentials returns the value of the ESS credentials secret.
func (s *ClusterSecrets) ReadESSCredentials(environment string) (value ESSCredentials, err error) {
	err = unmarshallSecret(s.client, s.getESSCredentialsPath(environment), &value)
	return value, err
}

// ReadServerlessCredentials returns the value of the serverless credentials secret.
func (s *ClusterSecrets) ReadServerlessCredentials(environment string) (value ServerlessCredentials, err error) {
	err = unmarshallSecret(s.client, s.getServerlessCredentialsPath(environment), &value)
	return value, err
}

// ReadDeployInfoSecret returns the value of the deploy info secret.
func (s *ClusterSecrets) ReadDeployInfoSecret(clusterName string) (value string, err error) {
	return s.client.GetSecret(s.getDeployInfoSecretPath(clusterName))
}

// ReadKibanaYamlSecret returns the value of the kibana yaml secret.
func (s *ClusterSecrets) ReadKibanaYamlSecret(clusterName string) (value string, err error) {
	return s.client.GetSecret(s.getKibanaYamlSecretPath(clusterName))
}

// ReadCredentialsSecret returns the value of the credentials secret.
func (s *ClusterSecrets) ReadCredentialsSecret(clusterName string) (value string, err error) {
	return s.client.GetSecret(s.getCredentialsSecretPath(clusterName))
}

// ReadEnvSecret returns the value of the Elastic Stack environment variables secret.
func (s *ClusterSecrets) ReadEnvSecret(clusterName string) (value string, err error) {
	return s.client.GetSecret(s.getEnvSecretPath(clusterName))
}

// ReadActivateSecret returns the value of the activate secret.
func (s *ClusterSecrets) ReadActivateSecret(clusterName string) (value string, err error) {
	return s.client.GetSecret(s.getActivateSecretPath(clusterName))
}
