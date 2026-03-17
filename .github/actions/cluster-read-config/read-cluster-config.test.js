const fs = require('fs');
const yaml = require('js-yaml');
const path = require('path');
const core = require('@actions/core');

// Mock the required dependencies
jest.mock('@actions/core');

// Import the run function
const { run } = require('./read-cluster-config');
const exp = require('constants');

describe('run', () => {
  beforeEach(() => {
    // Clear all mock calls and instances before each test
    jest.clearAllMocks();
  });

  test('should return an error when clusterName and fileName are empty', () => {
    // Arrange
    const clusterName = '';
    const fileName = '';
    const directoryPath = '';

    // Act
    expect(() => run(clusterName, fileName, directoryPath, false)).toThrow('Either CLUSTER_NAME or CLUSTER_CONFIG_FILE must be set');
    expect(core.setFailed).toHaveBeenCalledWith('Either CLUSTER_NAME or CLUSTER_CONFIG_FILE must be set');
  });

  test('should set outputs when clusterName is empty and fileName is not empty', () => {
    // Arrange
    const clusterName = '';
    const fileName = './data/user/config.yaml';
    const directoryPath = '';

    // Act
    run(clusterName, fileName, directoryPath, false);

    // Assert
    expect(core.setFailed).not.toHaveBeenCalled();
    expect(core.setOutput).toHaveBeenCalledWith('cluster_name', 'cluster-1');
    expect(core.setOutput).toHaveBeenCalledWith('template_name', 'template-1');
    expect(core.setOutput).toHaveBeenCalledWith('slack_channel', '#slack-channel-1');
    expect(core.setOutput).toHaveBeenCalledWith('stack_version', '8.15.0');
    expect(core.setOutput).toHaveBeenCalledWith('oblt_username', 'oblt-username-1');
    expect(core.setOutput).toHaveBeenCalledWith('golden_cluster', true);
    expect(core.setOutput).toHaveBeenCalledWith('expire_date', '2023-01-01');
    expect(core.setOutput).toHaveBeenCalledWith('grab_cluster_info', false);
    expect(core.setOutput).toHaveBeenCalledWith('update_channel', 'beta');
    expect(core.setOutput).toHaveBeenCalledWith('update_mode', 'recreate');
    expect(core.setOutput).toHaveBeenCalledWith('update_schedule', 'daily');
    expect(core.setOutput).toHaveBeenCalledWith('elasticsearch_image', 'docker.elastic.co/elasticsearch/elasticsearch:8.15.0-build');
    expect(core.setOutput).toHaveBeenCalledWith('kibana_image', 'docker.elastic.co/kibana/kibana:8.15.0-build');
    expect(core.setOutput).toHaveBeenCalledWith('elasticsearch_version', '8.15.0');
    expect(core.setOutput).toHaveBeenCalledWith('kibana_version', '8.15.0');
    expect(core.setOutput).toHaveBeenCalledWith('apm_version', '8.15.0');
    expect(core.setOutput).toHaveBeenCalledWith('remote_cluster', 'test-cluster');
    expect(core.setOutput).toHaveBeenCalledWith('stack_build', '8.15.0-build');
    expect(core.setOutput).toHaveBeenCalledWith('gitops', true);
    expect(core.setOutput).toHaveBeenCalledWith('comment_id', 123456789);
    expect(core.setOutput).toHaveBeenCalledWith('commit', 'commit-hash-1');
    expect(core.setOutput).toHaveBeenCalledWith('issue', 'issue-1');
    expect(core.setOutput).toHaveBeenCalledWith('pull_request', 'pull-request-1');
    expect(core.setOutput).toHaveBeenCalledWith('repository', 'repository-1');
  });

  test('Find cluster by name', () => {
    // Arrange
    const clusterName = 'bar';
    const fileName = '';
    const directoryPath = './data';

    // Act
    run(clusterName, fileName, directoryPath, true);
    // Assert
    expect(core.setFailed).not.toHaveBeenCalled();
    expect(core.setOutput).toHaveBeenCalledWith('cluster_name', 'bar');
    expect(core.setOutput).toHaveBeenCalledWith('clusterConfigFile', `${path.resolve(directoryPath)}/user/config.yml`);
  });

  test('should set outputs when clusterName is not empty and fileName is empty', () => {
    // Arrange
    const clusterName = 'cluster-1';
    const fileName = '';
    const directoryPath = './data';

    // Act
    run(clusterName, fileName, directoryPath, false);

    // Assert
    expect(core.setFailed).not.toHaveBeenCalled();
    expect(core.setOutput).toHaveBeenCalledWith('cluster_name', 'cluster-1');
    expect(core.setOutput).toHaveBeenCalledWith('template_name', 'template-1');
    expect(core.setOutput).toHaveBeenCalledWith('slack_channel', '#slack-channel-1');
    expect(core.setOutput).toHaveBeenCalledWith('stack_version', '8.15.0');
    expect(core.setOutput).toHaveBeenCalledWith('oblt_username', 'oblt-username-1');
    expect(core.setOutput).toHaveBeenCalledWith('golden_cluster', true);
    expect(core.setOutput).toHaveBeenCalledWith('expire_date', '2023-01-01');
    expect(core.setOutput).toHaveBeenCalledWith('grab_cluster_info', false);
    expect(core.setOutput).toHaveBeenCalledWith('update_channel', 'beta');
    expect(core.setOutput).toHaveBeenCalledWith('update_mode', 'recreate');
    expect(core.setOutput).toHaveBeenCalledWith('update_schedule', 'daily');
    expect(core.setOutput).toHaveBeenCalledWith('elasticsearch_image', 'docker.elastic.co/elasticsearch/elasticsearch:8.15.0-build');
    expect(core.setOutput).toHaveBeenCalledWith('kibana_image', 'docker.elastic.co/kibana/kibana:8.15.0-build');
    expect(core.setOutput).toHaveBeenCalledWith('elasticsearch_version', '8.15.0');
    expect(core.setOutput).toHaveBeenCalledWith('kibana_version', '8.15.0');
    expect(core.setOutput).toHaveBeenCalledWith('apm_version', '8.15.0');
    expect(core.setOutput).toHaveBeenCalledWith('remote_cluster', 'test-cluster');
    expect(core.setOutput).toHaveBeenCalledWith('stack_build', '8.15.0-build');
    expect(core.setOutput).toHaveBeenCalledWith('gitops', true);
    expect(core.setOutput).toHaveBeenCalledWith('comment_id', 123456789);
    expect(core.setOutput).toHaveBeenCalledWith('commit', 'commit-hash-1');
    expect(core.setOutput).toHaveBeenCalledWith('issue', 'issue-1');
    expect(core.setOutput).toHaveBeenCalledWith('pull_request', 'pull-request-1');
    expect(core.setOutput).toHaveBeenCalledWith('repository', 'repository-1');
  });

  test('should not fail on a basic cluster config', () => {
    // Arrange
    const clusterName = 'bar';
    const fileName = '';
    const directoryPath = './data';

    // Act
    run(clusterName, fileName, directoryPath, false);
    expect(core.setFailed).not.toHaveBeenCalled();
  });

  test('Cluster name must be sort than 63 characters', () => {
    const clusterName = '';
    const fileName = './data/validation/cluster_name_63.yml';
    const directoryPath = '';

    // Act
    expect(() => run(clusterName, fileName, directoryPath, true)).toThrow('cluster_name must be shorter than 63 characters');
    expect(core.setFailed).toHaveBeenCalledWith('cluster_name must be shorter than 63 characters');
  });

  test('Cluster name must start with a letter', () => {
    const clusterName = '';
    const fileName = './data/validation/cluster_name_fqdn.yml';
    const directoryPath = '';

    // Act
    expect(() => run(clusterName, fileName, directoryPath, true)).toThrow('cluster_name must be lowercase alphanumeric characters, hyphens and dots (FQDN)');
    expect(core.setFailed).toHaveBeenCalledWith('cluster_name must be lowercase alphanumeric characters, hyphens and dots (FQDN)');
  });

  test('Cluster name must not be empty', () => {
    const clusterName = '';
    const fileName = './data/validation/cluster_name_empty.yml';
    const directoryPath = '';

    // Act
    expect(() => run(clusterName, fileName, directoryPath, true)).toThrow('cluster_name is required');
    expect(core.setFailed).toHaveBeenCalledWith('cluster_name is required');
  });

  test('Cluster name must not be empty 2', () => {
    const clusterName = '';
    const fileName = './data/validation/cluster_name_empty_2.yml';
    const directoryPath = '';

    // Act
    expect(() => run(clusterName, fileName, directoryPath, true)).toThrow('cluster_name is required');
    expect(core.setFailed).toHaveBeenCalledWith('cluster_name is required');
  });

  test('Cluster name must not be empty 3', () => {
    const clusterName = '';
    const fileName = './data/validation/cluster_name_empty_3.yml';
    const directoryPath = '';

    // Act
    expect(() => run(clusterName, fileName, directoryPath, true)).toThrow('cluster_name is required');
    expect(core.setFailed).toHaveBeenCalledWith('cluster_name is required');
  });

  test('Cluster name must not contain capitals', () => {
    const clusterName = '';
    const fileName = './data/validation/cluster_name_invalid_1.yml';
    const directoryPath = '';

    // Act
    expect(() => run(clusterName, fileName, directoryPath, true)).toThrow('cluster_name must be lowercase alphanumeric characters, hyphens and dots (FQDN)');
    expect(core.setFailed).toHaveBeenCalledWith('cluster_name must be lowercase alphanumeric characters, hyphens and dots (FQDN)');
  });

  test('Docker image valid', () => {
    const clusterName = '';
    const fileName = './data/validation/docker_image_valid.yml';
    const directoryPath = '';

    // Act
    run(clusterName, fileName, directoryPath, true);
    expect(core.setFailed).not.toHaveBeenCalled();
  });

  test('Docker image invalid 1', () => {
    const clusterName = '';
    const fileName = './data/validation/docker_image_invalid_1.yml';
    const directoryPath = '';

    // Act
    expect(() => run(clusterName, fileName, directoryPath, true)).toThrow('stack.elasticsearch.image must be a valid docker image');
    expect(core.setFailed).toHaveBeenCalledWith('stack.elasticsearch.image must be a valid docker image');
  });

  test('Docker image invalid 2', () => {
    const clusterName = '';
    const fileName = './data/validation/docker_image_invalid_1.yml';
    const directoryPath = '';

    // Act
    expect(() => run(clusterName, fileName, directoryPath, true)).toThrow('stack.elasticsearch.image must be a valid docker image');
    expect(core.setFailed).toHaveBeenCalledWith('stack.elasticsearch.image must be a valid docker image');
  });

  test('Stack version release', () => {
    const clusterName = '';
    const fileName = './data/validation/stack_version_release.yml';
    const directoryPath = '';

    // Act
    run(clusterName, fileName, directoryPath, true);
    expect(core.setFailed).not.toHaveBeenCalled();
  });

  test('Stack version SNAPSHOT', () => {
    const clusterName = '';
    const fileName = './data/validation/stack_version_snapshot.yml';
    const directoryPath = '';

    // Act
    run(clusterName, fileName, directoryPath, true);
    expect(core.setFailed).not.toHaveBeenCalled();
  });


  test('Stack version invalid 1', () => {
    const clusterName = '';
    const fileName = './data/validation/stack_version_invalid_1.yml';
    const directoryPath = '';

    // Act
    expect(() => run(clusterName, fileName, directoryPath, true)).toThrow('stack.version must be a valid version number');
    expect(core.setFailed).toHaveBeenCalledWith('stack.version must be a valid version number');
  });

  test('Stack version invalid 2', () => {
    const clusterName = '';
    const fileName = './data/validation/stack_version_invalid_2.yml';
    const directoryPath = '';

    // Act
    expect(() => run(clusterName, fileName, directoryPath, true)).toThrow('stack.version must be a valid version number');
    expect(core.setFailed).toHaveBeenCalledWith('stack.version must be a valid version number');
  });

  test('Stack version invalid 3', () => {
    const clusterName = '';
    const fileName = './data/validation/stack_version_invalid_3.yml';
    const directoryPath = '';

    // Act
    expect(() => run(clusterName, fileName, directoryPath, true)).toThrow('stack.version must be a valid version number');
    expect(core.setFailed).toHaveBeenCalledWith('stack.version must be a valid version number');
  });

});
