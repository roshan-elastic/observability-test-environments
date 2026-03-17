const fs = require("fs");
const yaml = require("js-yaml");
const path = require("path");
const core = require("@actions/core");

const getAllFiles = (dir) =>
  fs.readdirSync(dir).reduce((files, file) => {
    const name = path.join(dir, file);
    const isDirectory = fs.statSync(name).isDirectory();
    return isDirectory ? [...files, ...getAllFiles(name)] : [...files, path.resolve(name)];
  }, []);

const getClusterConfigFiles = (clusterName, fileName, directoryPath) => {
  var files = [];
  if (clusterName != "") {
    core.debug(`Reading all files in ${directoryPath}`);
    files = getAllFiles(directoryPath);
  } else {
    core.debug("Reading single file");
    files = [fileName];
  }
  return files;
};

const run = (clusterName, fileName, directoryPath, validateConfig) => {
  core.info("CLUSTER_NAME: " + clusterName);
  core.info("CLUSTER_CONFIG_FILE: " + fileName);

  if (fileName == "" && clusterName == "") {
    const errorMessage =
      "Either CLUSTER_NAME or CLUSTER_CONFIG_FILE must be set";
    core.setFailed(errorMessage);
    throw new Error(errorMessage);
  }

  var files = getClusterConfigFiles(clusterName, fileName, directoryPath);
  core.debug("Found files: " + JSON.stringify(files));
  for (const file of files) {
    if (file.endsWith(".yml") || file.endsWith(".yaml")) {
      core.debug("Reading file: " + file);
      const fileContent = fs.readFileSync(file, "utf8");
      const config = yaml.load(fileContent);
      const clusterNameFile = config.cluster_name?.toString();

      if (
        (fileName == file ) ||
        (clusterName != "" && clusterNameFile == clusterName)
      ) {
        setOutputs(config, clusterName);
        core.setOutput("clusterConfigFile", file);
        if (validateConfig) {
          validate(config);
        }
        break;
      }
    }
  }
};

function setOutputs(config, clusterName) {
  const stack = config.stack || {};
  const stack_mode = stack.mode || "ess";
  const stack_mode_obj = stack[stack_mode] || {};
  const elasticsearch = stack_mode_obj.elasticsearch || {};
  const kibana = stack_mode_obj.kibana || {};
  const stack_version = stack.version;
  const elasticsearch_version = stack_version;
  const kibana_version = stack_version;
  const apm_version = stack_version;
  oblt_managed = true
  if (`${config.oblt_managed}` == 'false') {
    oblt_managed = false
  }
  golden_cluster = false
  if (`${config.golden_cluster}` == 'true') {
    golden_cluster = true
  }
  grab_cluster_info = false
  if (`${config.grab_cluster_info}` == 'true') {
    grab_cluster_info = true
  }


  core.debug("Found cluster: " + clusterName);
  core.setOutput("cluster_name", config.cluster_name);
  core.setOutput("template_name", config.template_name);
  core.setOutput("slack_channel", config.slack_channel);
  core.setOutput("stack_version", stack_version);
  core.setOutput("oblt_username", config.oblt_username);
  core.setOutput('oblt_managed', oblt_managed);
  core.setOutput("golden_cluster", golden_cluster);
  core.setOutput("expire_date", config.expire_date);
  core.setOutput("grab_cluster_info", grab_cluster_info);

  core.setOutput("update_channel", stack.update_channel);
  core.setOutput("update_mode", stack.update_mode || "update");
  core.setOutput("update_schedule", stack.update_schedule || "daily");

  const elasticsearch_image = elasticsearch.image || "";
  core.setOutput("elasticsearch_image", elasticsearch_image);
  core.setOutput("kibana_image", kibana.image || "");
  core.setOutput("elasticsearch_version", elasticsearch_version);
  core.setOutput("kibana_version", kibana_version);
  core.setOutput("apm_version", apm_version);

  core.setOutput("remote_cluster", stack_mode_obj.ccs_remote_cluster || "");

  stack_build = elasticsearch_image.split(":")[1] || stack_version;
  core.setOutput("stack_build", stack_build);

  const gitops = config.gitops || {};
  gitOpsEnabled = false
  if (`${gitops.enabled}` == 'true') {
    gitOpsEnabled = true
  }
  core.setOutput('gitops', gitOpsEnabled );
  if (gitOpsEnabled) {
    core.setOutput('comment_id', gitops.comment_id || '');
    core.setOutput('commit', gitops.commit || '');
    core.setOutput('issue', gitops.issue || '');
    core.setOutput('pull_request', gitops.pull_request || '');
    core.setOutput('repository', gitops.repository || '');
  }
}

function validate(config) {
  validateClusterName(config);
  if (config.stack !== undefined) {
    validateStackVersion(config);
    validateDockerImages(config);
  }
}

function validateClusterName(config) {
  if (config.cluster_name === undefined || config.cluster_name === "" || config.cluster_name === null || config.cluster_name == "") {
    const errorMessage = "cluster_name is required";
    core.setFailed(errorMessage);
    throw new Error(errorMessage);
  }
  if (config.cluster_name?.toString()?.length > 63) {
    const errorMessage = "cluster_name must be shorter than 63 characters";
    core.setFailed(errorMessage);
    throw new Error(errorMessage);
  }
  if (config.cluster_name?.toString()?.match(/^([a-z][a-z0-9-._]+)$/) == null) {
    const errorMessage =
      "cluster_name must be lowercase alphanumeric characters, hyphens and dots (FQDN)";
    core.setFailed(errorMessage);
    throw new Error(errorMessage);
  }
}

function validateStackVersion(config) {
  if (config?.stack?.version !== undefined
      && config?.stack?.version !== "") {
    if (config.stack.version?.toString().match(/^([0-9]+\.[0-9]+\.[0-9]+[a-zA-Z0-9-._]*)$/) == null) {
      const errorMessage = "stack.version must be a valid version number";
      core.setFailed(errorMessage);
      throw new Error(errorMessage);
    }
  }
}

function validateDockerImages(config) {
  stack_mode = config.stack.mode || "ess";
  if (stack_mode.match(/([a-z]+)/) == null) {
    const errorMessage = "stack.mode must be lowercase alphanumeric characters";
    core.setFailed(errorMessage);
    throw new Error(errorMessage);
  }
  validateEsDockerImage(config, stack_mode);
  validateKibanaDockerImage(config, stack_mode);
}

function validateEsDockerImage(config, stack_mode) {
  if (config?.stack[stack_mode]?.elasticsearch?.image !== undefined) {
    if (config.stack[stack_mode].elasticsearch.image?.toString().match(/^([a-z0-9]+(\.[a-z0-9-]+)*(\/[a-z0-9-]+)*:([a-zA-Z0-9-_.]+))$/) == null) {
      const errorMessage = "stack.elasticsearch.image must be a valid docker image";
      core.setFailed(errorMessage);
      throw new Error(errorMessage);
    }
  }
}

function validateKibanaDockerImage(config, stack_mode) {
  if (config?.stack[stack_mode]?.kibana?.image !== undefined) {
    if (config.stack[stack_mode].kibana.image?.toString().match(/^([a-z0-9]+(\.[a-z0-9-]+)*(\/[a-z0-9-]+)*:([a-zA-Z0-9-_.]+))$/) == null) {
      const errorMessage = "stack.kibana.image must be a valid docker image";
      core.setFailed(errorMessage);
      throw new Error(errorMessage);
    }
  }
}


exports = module.exports = {
  run,
};
