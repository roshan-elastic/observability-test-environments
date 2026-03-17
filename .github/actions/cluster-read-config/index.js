const readClusterConfig = require('./read-cluster-config');
const core = require('@actions/core');

const clusterNameParam = core.getInput('clusterName') || process.env.CLUSTER_NAME || '';
const fileNameParam = core.getInput('fileName') || process.env.CLUSTER_CONFIG_FILE || '';
const directoryPathParam = './environments/users' || process.env.DIRECTORY_PATH || './';
const validateConfig = core.getInput('validateConfig') == true || process.env.VALIDATE || true;

readClusterConfig.run(clusterNameParam, fileNameParam, directoryPathParam, validateConfig);
