package config

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/elastic/observability-test-environments/tools/oblt-cli/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	// AllFlag flag name.
	AllFlag = "all"
	// ConfigFlag flag name.
	ConfigFlag = "config"
	// ClusterNameFlag flag name.
	ClusterNameFlag = "cluster-name"
	// RemoteClusterFlag flag name.
	RemoteClusterFlag = "remote-cluster"
	// ClusterNamePrefixFlag flag name.
	ClusterNamePrefixFlag = "cluster-name-prefix"
	// ClusterNameSuffixFlag flag name.
	ClusterNameSuffixFlag = "cluster-name-suffix"
	// StackVersionFlag flag name.
	StackVersionFlag = "stack-version"
	// DryRunFlag flag name.
	DryRunFlag = "dry-run"
	// SlackChannelFlag slack channel flag name.
	SlackChannelFlag = "slack-channel"
	// UsernameFlag Username flag name.
	UsernameFlag = "username"
	// ForceFlag force flag name.
	ForceFlag = "force"
	// InteractiveFlag interactive flag name.
	InteractiveFlag = "interactive"
	// TemplateNameFlag flag name.
	TemplateNameFlag = "template"
	// TemplateFileFlag flag name.
	TemplateFileFlag = "template-file"
	// ParametersFlag flag name.
	ParametersFlag = "parameters"
	// EsDockerImageFlag flag name.
	EsDockerImageFlag = "elasticsearch-docker-image"
	// ElasticAgentDockerImageFlag flag name.
	ElasticAgentDockerImageFlag = "elastic-agent-docker-image"
	// KbDockerImageFlag flag name.
	KbDockerImageFlag = "kibana-docker-image"
	// ApmDockerImageFlag flag name.
	ApmDockerImageFlag = "apm-docker-image"
	// DockerImageVersionFlag flag name.
	DockerImageVersionFlag = "docker-image-version"
	// GitHttpModeFlag flag name.
	GitHttpModeFlag = "git-http-mode"
	// ExperimentalFlag flag name.
	ExperimentalFlag = "experimental"
	// VerboseFlag flag name.
	VerboseFlag = "verbose"
	// OutputFileFlag flag name.
	OutputFileFlag = "output-file"
	// UrlFlag flag name.
	UrlFlag = "url"
	// PasswordFlag flag name.
	PasswordFlag = "password"
	// ApiKeyFlag flag name.
	ApiKeyFlag = "api-key"
	// RecipesFlag flag name.
	RecipesFlag = "recipes"
	// RecipesElasticsearchFlag flag name.
	RecipesElasticsearchFlag = "recipes-elasticsearch"
	// RecipesKibanaFlag flag name.
	RecipesKibanaFlag = "recipes-kibana"
	// ListFlag flag name.
	ListFlag = "list"
	// IgnoreCertificatesFlag flag name.
	IgnoreCertificatesFlag = "ignore-certificates"
	// GithubTokenFlag Github Token flag name.
	GithubTokenFlag = "github-token"
	// WaitFlag flag name.
	WaitFlag = "wait"
	// IsReleaseFlag flag name.
	IsReleaseFlag = "release"
	// RepoFlag flag name.
	RepoFlag = "repo"
	// CommitFlag flag name.
	CommitFlag = "commit"
	// IssueFlag flag name.
	IssueFlag = "issue"
	// PullRequestFlag flag name.
	PullRequestFlag = "pull-request"
	// CommentIdFlag flag name.
	CommentIdFlag = "comment-id"
	// BootstrapFolderFlag flag name.
	BootstrapFolderFlag = "bootstrap-folder"
	// RecipesApmFlag flag name.
	RecipesApmFlag = "recipes-apm"
	// RecipesFleetFlag flag name.
	RecipesFleetFlag = "recipes-fleet"
	// ParametersFileFlag flag name.
	ParametersFileFlag = "parameters-file"
	// BranchFlag flag name.
	BranchFlag = "branch"
	// IntegrationFlag flag name.
	IntegrationFlag = "integration"
	// IntegrationNameFlag flag name.
	ContainerNameFlag = "container"
	// PackageFolderFlag flag name.
	PackageFolderFlag = "package-folder"
	// PortsFlag flag name.
	PortsFlag = "ports"
	// RepositoryFlag flag name.
	RepositoryFlag = "repository"
	// OutputFolderFlag flag name.
	OutputFolderFlag = "output-folder"
	// DisableBannerFlag flag name.
	DisableBannerFlag = "disable-banner"
	// WipeupFlag flag name.
	WipeupFlag = "wipeup"
	// FilterFlag flag name.
	FilterFlag = "filter"
	// TypeFlag flag name.
	TypeFlag = "type"
	// EnvironmentFlag flag name.
	EnvironmentFlag = "environment"
	// SaveConfigFlag flag name.
	SaveConfigFlag = "save-config"
	// ProjectTypeFlag flag name.
	ProjectTypeFlag = "project-type"
	// ElasticsearchDockerImageFlag flag name.
	ElasticsearchDockerImageFlag = "elasticsearch-docker-image"
	// KibanaDockerImageFlag flag name.
	KibanaDockerImageFlag = "kibana-docker-image"
	// FleetDockerImageFlag flag name.
	FleetDockerImageFlag = "fleet-docker-image"
	// TargetFlag flag name.
	TargetFlag = "target"
	// BaseFlag flag name.
	BaseFlag = "base"
	// BodyFlag flag name.
	BodyFlag = "body"
	// HeadFlag flag name.
	HeadFlag = "head"
	// LabelsFlag flag name.
	LabelsFlag = "labels"
	// TitleFlag flag name.
	TitleFlag = "title"
	// KibanaYamlPathFlag flag name.
	KibanaYamlPathFlag = "kibana-yaml-path"
	// KibanaSrcPathFlag flag name.
	KibanaSrcPathFlag = "kibana-src"
	// VPNFlag flag name.
	VPNFlag = "vpn"
	// MKILoginFlag flag name.
	MKILoginFlag = "mki-login"
	// ConfigFileFlag flag name.
	ConfigFileFlag = "config-file"
)

const (
	// configDir config directory.
	configDir = ".oblt-cli"
	// configFile config file name.
	configFile = "config.yaml"
)

type ObltConfiguration struct {
	// ConfigFile is the configuration file path.
	ConfigFile string
	// SlackChannel user Slack member ID.
	SlackChannel string
	// Username it is the username used for naming resources.
	Username string
	// GitHttpMode is true if the git repository is cloned using HTTPS.
	GitHttpMode bool
	// Verbose is true if verbose output is enabled.
	Verbose bool
}

// CfgFile path to the configuration file passes by command line.
var CfgFile string

// RepoBranch is the branch to use for the observability-test-environments repository.
var RepoBranch string = "main"

// Initialise reads in config file, running callbacks after configuration is ready
func Initialise(cfgFileFallback string, callbacks ...func()) {
	isConfigured := false
	if CfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(CfgFile)
	} else {
		viper.SetConfigFile(cfgFileFallback)

		cfgDir := filepath.Dir(cfgFileFallback)
		os.MkdirAll(cfgDir, 0777)
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		logger.Debugf("Using config file: %s", viper.ConfigFileUsed())
		isConfigured = true
	}

	if isConfigured {
		for _, callback := range callbacks {
			callback()
		}
	}
}

// DefaultFile return the default configuration file path.
func DefaultFile() (response string) {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	return filepath.Join(home, configDir, configFile)
}

// ForUser return a unique configuration file path for the given user.
func ForUser(userName string) (response string) {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	return filepath.Join(home, configDir, userName, configFile)
}

// NewObltConfig creates a new ObltConfig.
func NewObltConfig(configFile, username, slackChannel string, gitHttpMode, verbose bool) (config ObltConfiguration) {
	return ObltConfiguration{
		ConfigFile:   configFile,
		SlackChannel: slackChannel,
		Username:     username,
		GitHttpMode:  gitHttpMode,
		Verbose:      verbose,
	}
}

// NewObltConfigFromViper creates a new ObltConfig from the configuration file using Viper library.
// The Viper configuration is not thread safe.
func NewObltConfigFromViper() (config ObltConfiguration) {
	return ObltConfiguration{
		ConfigFile:   viper.ConfigFileUsed(),
		SlackChannel: viper.GetString(SlackChannelFlag),
		Username:     viper.GetString(UsernameFlag),
		GitHttpMode:  viper.GetBool(GitHttpModeFlag),
		Verbose:      viper.GetBool(VerboseFlag),
	}
}

// IsGitMode if the SSH mode is selected for Git operations (SSH/HTTP)
func (o *ObltConfiguration) IsSSHMode() bool {
	return !o.GitHttpMode
}

// GetDir Returns the absolute path to the configuration folder.
func (o *ObltConfiguration) GetDir() (cfgDir string) {
	cfgDir = filepath.Dir(o.ConfigFile)
	if !fileExists(cfgDir) {
		log.Fatalf(`config dir does not exist at "%s". Please create it by running the "configure" command of this tool:

oblt-cli configure --help

`, o.ConfigFile)
	}

	return
}

// fileExists checks if the given file path exists in the file system
func fileExists(file string) bool {
	_, err := os.Stat(file)

	return !os.IsNotExist(err)
}

// ValidateAlphanumeric It checks if the value match the rules alphanumeric, ASCII and lowercase.
func ValidateAlphanumeric(value string) (err error) {
	if !govalidator.IsASCII(value) || !govalidator.IsLowerCase(value) || !govalidator.IsAlphanumeric(value) {
		err = fmt.Errorf("'%s' should contains only the following characters [a-z0-9]", value)
	}

	return err
}

// ValidateUsername It checks if the username is valid.
func ValidateUsername(value string) (err error) {
	err = ValidateDNSLike(value)
	if err != nil {
		err = fmt.Errorf("invalid username %v", err)
	}
	return err
}

// ValidateDNSLike It checks if the value match the rules valid DNS name.
func ValidateDNSLike(value string) (err error) {
	if !govalidator.IsDNSName(value) {
		err = fmt.Errorf("'%s' should be a valid FQDN name", value)
	}
	return err
}

// ValidateNames It checks if the value match the rules valid DNS name, no '_', and lowercase.
func ValidateNames(value string) (err error) {
	if !govalidator.IsDNSName(value) || strings.Contains(value, "_") || !govalidator.IsLowerCase(value) {
		err = fmt.Errorf("'%s' should be a valid name, these are the only accepted characters [a-z][0-9][.-]", value)
	}
	return err
}

// ValidateLicenseType It checks if the license type is valid.
func ValidateLicenseType(licenseType string) (err error) {
	switch licenseType {
	case "release", "dev", "orchestration", "orchestration-dev":
	default:
		err = fmt.Errorf("'%s' is not a valid license type, please use one of the following values [release, dev, orchestration, orchestration-dev]", licenseType)
	}
	return err
}

// ValidateProjectType It checks if the project type is valid.
func ValidateProjectType(projectType string) (err error) {
	switch projectType {
	case "elasticsearch", "observability", "security", "":
	default:
		err = fmt.Errorf("'%s' is not a valid project type, please use one of the following values [elasticsearch, observability, security]", projectType)
	}
	return err
}

// ValidateTarget It checks if the target is valid.
func ValidateTarget(target string) (err error) {
	switch target {
	case "production", "staging", "qa", "":
	default:
		err = fmt.Errorf("'%s' is not a valid target, please use one of the following values [production, staging, qa]", target)
	}
	return err
}

// ValidateDockerImage It checks if the docker image is valid.
func ValidateDockerImage(dockerImage string) (err error) {
	if dockerImage != "" {
		var validExp = regexp.MustCompile(`^([a-z0-9-]+(?:[.-][a-z0-9-]+)*\/)?[a-z0-9-\/]+(?:[.-][a-z0-9-]+)*(:[a-zA-Z0-9.-]+)?$`)
		if !validExp.MatchString(dockerImage) {
			err = fmt.Errorf("'%s' is not a valid Docker image name", dockerImage)
		}
	}
	return err
}

// ValidateClusterName It checks if the clusterName is valid.
func ValidateClusterName(clusterName string) (err error) {
	if clusterName != "" {
		err = errors.Join(
			ValidateNames(clusterName),
			ValidateLength(clusterName, 1, 63),
		)
	}
	return err
}

// ValidatePrefix It checks if the prefix is valid.
func ValidatePrefix(prefix string) (err error) {
	if prefix != "" {
		err = errors.Join(
			ValidateNames(prefix),
			ValidateLength(prefix, 1, 32),
		)
	}
	return err
}

// ValidateSuffix It checks if the suffix is valid.
func ValidateSuffix(suffix string) (err error) {
	return ValidatePrefix(suffix)
}

// ValidateSlackChannel It checks if the slack channel is valid.
func ValidateSlackChannel(slackChannel string) (err error) {
	var validExp = regexp.MustCompile(`^[@#][A-Za-z0-9-]{1,79}$`)
	if !validExp.MatchString(slackChannel) {
		err = fmt.Errorf("the provided Slack channel '%s' is invalid; please use the '#' prefix for Slack channels (e.g., #channelname) and the '@' prefix for Slack user IDs (e.g., @userID)", slackChannel)
	}
	return err
}

// ValidateLength It checks if the value length is between min and max.
func ValidateLength(value string, min, max int) (err error) {
	if len(value) < min || len(value) > max {
		err = fmt.Errorf("'%s' length must be between %d and %d", value, min, max)
	}
	return err
}

// ValidateSemVer It checks if the value is a valid semantic version.
func ValidateSemVer(value string) (err error) {
	if !govalidator.IsSemver(value) {
		err = fmt.Errorf("'%s' is not a valid semantic version", value)
	}
	return err
}

// seed Returns a random string of lowercase characters of the length specified.
func Seed(length int) (seed string) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	characters := []rune("abcdefghijklmnopqrstuvwxyz")
	seedTmp := make([]rune, length)
	limit := len(characters)
	for i := range seedTmp {
		seedTmp[i] = characters[r.Intn(limit)]
	}
	return string(seedTmp)
}

func CheckErr(msg interface{}) {
	if msg != nil {
		logger.Errorf("%v", msg)
		os.Exit(1)
	}
}
