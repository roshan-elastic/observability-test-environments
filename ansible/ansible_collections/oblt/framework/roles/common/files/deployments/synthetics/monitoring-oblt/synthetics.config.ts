import type { SyntheticsConfig } from '@elastic/synthetics';

// verify that all environment variables are defined before running synthetics
const requiredEnvVars = [
  'KIBANA_URL',
  'KIBANA_USERNAME',
  'KIBANA_PASSWORD',
  'SYNTHETICS_PROJECT_ID',
  'SYNTHETICS_KIBANA_URL',
  'SYNTHETICS_LOCALE',
  'SYNTHETICS_TIMEZONE'
];
requiredEnvVars.forEach((envVar) => {
  if (!process.env[envVar]) {
    console.log(`Environment variable ${envVar} is not defined`)
    throw new Error(`Environment variable ${envVar} is not defined`);
  }
});

export default env => {
  const config: SyntheticsConfig = {
    params: {
      url: process.env.KIBANA_URL,
      username: process.env.KIBANA_USERNAME,
      password: process.env.KIBANA_PASSWORD,
      projectId: process.env.SYNTHETICS_PROJECT_ID,
    },
    playwrightOptions: {
      ignoreHTTPSErrors: false,
      locale: process.env.SYNTHETICS_LOCALE,
      timezoneId: process.env.SYNTHETICS_TIMEZONE,
    },
    monitor: {
      schedule: 10,
      locations: ['us_east'],
      privateLocations: [],
    },
    project: {
      id: process.env.SYNTHETICS_PROJECT_ID,
      url: process.env.SYNTHETICS_KIBANA_URL,
      space: 'default',
    },
  };
  return config;
};
