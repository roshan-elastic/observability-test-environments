# Changelog

Versions upper than 2.5.0 have their change log at [GitHub Releases](https://github.com/elastic/observability-test-environments/releases)

## 2.5.0 (15/09/2022)

#### 🚀 Enhancements

-  chore: review heartbeat config [#2493](https://github.com/elastic/observability-test-environments/pull/2493)
-  Add opbeans php [#2506](https://github.com/elastic/observability-test-environments/pull/2506)

#### 🐛 Bug Fixes

-  Add opbeans php [#2533](https://github.com/elastic/observability-test-environments/pull/2533)
-  fix: reimplement the Docker images pull and push process [#2530](https://github.com/elastic/observability-test-environments/pull/2530)
-  fix: avoid remove endoflines in config files [#2528](https://github.com/elastic/observability-test-environments/pull/2528)
-  fix: change authproxy secrets paths [#2527](https://github.com/elastic/observability-test-environments/pull/2527)
-  fix: ignore ess images for 7.x [#2525](https://github.com/elastic/observability-test-environments/pull/2525)
-  fix: increase matrix workers [#2524](https://github.com/elastic/observability-test-environments/pull/2524)
-  fix: increase the disk space on gobld workers [#2523](https://github.com/elastic/observability-test-environments/pull/2523)
-  fix: workaround to pull and push Elasticsearch Docker images. [#2522](https://github.com/elastic/observability-test-environments/pull/2522)
- [**area:dx**] fix: change the way we import the Docker images tgz [#2517](https://github.com/elastic/observability-test-environments/pull/2517)

#### 🙈 No user affected

-  chore: bump release cluster to 8.4.1 [#2526](https://github.com/elastic/observability-test-environments/pull/2526)
-  chore: remove Helm Chart Elastic Stack support [#2494](https://github.com/elastic/observability-test-environments/pull/2494)

---

## 2.4.0 (25/07/2022)

#### 🚀 Enhancements

-  feat: Apm mutating webhook [#2465](https://github.com/elastic/observability-test-environments/pull/2465)
- [**area:dx**] feat: support for destroying a cluster [#2441](https://github.com/elastic/observability-test-environments/pull/2441)
- [**area:dx**] feat: add a command to create ESS clusters [#2408](https://github.com/elastic/observability-test-environments/pull/2408)

#### 🙈 No user affected

- [**area:dx**] chore: initialise oblt-repository just once for the CMDs in the cli [#2456](https://github.com/elastic/observability-test-environments/pull/2456)

---

## 2.3.1 (07/07/2022)

#### 🐛 Bug Fixes

- [**area:dx**] fix: do not override flag variables with env variables [#2401](https://github.com/elastic/observability-test-environments/pull/2401)

#### 📚 Documentation

- [**area:dx**] docs: user guide for the slack bot [#2403](https://github.com/elastic/observability-test-environments/pull/2403)

---

## 2.3.0 (06/07/2022)

#### 🚀 Enhancements

-  feat: bump versions [#2337](https://github.com/elastic/observability-test-environments/pull/2337)
-  feat: backup secrets before update [#2336](https://github.com/elastic/observability-test-environments/pull/2336)
-  Test if CCS developer cluster can be created after update [#2329](https://github.com/elastic/observability-test-environments/pull/2329)
- [**area:dx**] feat: bootstrap commands [#2307](https://github.com/elastic/observability-test-environments/pull/2307)
-  feat: script to generate the script to push licenses [#2303](https://github.com/elastic/observability-test-environments/pull/2303)
-  feat: add a Synthetics monitor to monitoring-oblt [#2300](https://github.com/elastic/observability-test-environments/pull/2300)
-  chore: Mix of issues [#2299](https://github.com/elastic/observability-test-environments/pull/2299)
-  feat: backup vault secrets [#2297](https://github.com/elastic/observability-test-environments/pull/2297)
-  fix: improve stack version script [#2296](https://github.com/elastic/observability-test-environments/pull/2296)
- [**area:dx**] feat: retrieve cluster secrets from oblt tools [#2385](https://github.com/elastic/observability-test-environments/pull/2385)

#### 🐛 Bug Fixes

-  fix: pull the repo before make the push [#2338](https://github.com/elastic/observability-test-environments/pull/2338)
-  fix: bump Helm charts versions [#2328](https://github.com/elastic/observability-test-environments/pull/2328)
-  fix: ECE deployment [#2325](https://github.com/elastic/observability-test-environments/pull/2325)
- [**area:dx**] chore: remove dry run from the deployment descriptor [#2316](https://github.com/elastic/observability-test-environments/pull/2316)
-  fix: set autoscale to true [#2302](https://github.com/elastic/observability-test-environments/pull/2302)
-  fix: remove nonexistent function [#2301](https://github.com/elastic/observability-test-environments/pull/2301)

#### 🙈 No user affected

-  chore: remove tag from grafana configuration [#2340](https://github.com/elastic/observability-test-environments/pull/2340)
-  chore: template for update issues [#2332](https://github.com/elastic/observability-test-environments/pull/2332)
- [**area:dx**] feat: do not remove the system-status repository each time [#2282](https://github.com/elastic/observability-test-environments/pull/2282)

---

## 2.2.0 (04/05/2022)

#### 🚀 Enhancements

- [**area:dx**][**requested-by:Ingest Management**] fix: add monitoring cluster to the CCS Slack form [#2275](https://github.com/elastic/observability-test-environments/pull/2275)
- [**area:dx**] chore: refactor output for cluster templates command [#2272](https://github.com/elastic/observability-test-environments/pull/2272)

#### 🙈 No user affected

- [**area:ci**][**area:dx**] chore: report cobertura for oblt-cli [#2242](https://github.com/elastic/observability-test-environments/pull/2242)

---

## 2.1.0 (25/04/2022)

#### 🚀 Enhancements

- [**area:dx**] feat: support showing commit ID and URL in messages [#2245](https://github.com/elastic/observability-test-environments/pull/2245)
- [**area:dx**] feat: add a slack command to list templates [#2235](https://github.com/elastic/observability-test-environments/pull/2235)

#### 🙈 No user affected

- [**area:dx**] chore: release slack bot docker images on git tags [#2250](https://github.com/elastic/observability-test-environments/pull/2250)
- [**area:dx**] chore(slack): improve message for CCS creation [#2243](https://github.com/elastic/observability-test-environments/pull/2243)
- [**area:dx**] fix: use project version for slackbot image [#2236](https://github.com/elastic/observability-test-environments/pull/2236)
- [**area:dx**] chore: always log current version if current is equals to latest [#2231](https://github.com/elastic/observability-test-environments/pull/2231)

---

## 2.0.31 (19/04/2022)

#### 🚀 Enhancements

- [**area:dx**] feat: add a command to self-update the tool [#2223](https://github.com/elastic/observability-test-environments/pull/2223)

#### 🐛 Bug Fixes

- [**area:ci**] fix: do not increment strings in groovy [#2229](https://github.com/elastic/observability-test-environments/pull/2229)
- [**area:dx**] chore: regen blob on releases [#2224](https://github.com/elastic/observability-test-environments/pull/2224)

---

## 2.0.29 (19/04/2022)

#### 🐛 Bug Fixes

- [**area:dx**] fix: proper base dir to create config dir [#2217](https://github.com/elastic/observability-test-environments/pull/2217)

---

## 2.0.28 (13/04/2022)

#### 🐛 Bug Fixes

- [**area:dx**] fix: create default config file's parent dir [#2213](https://github.com/elastic/observability-test-environments/pull/2213)

---

## 2.0.27 (13/04/2022)

#### 🚀 Enhancements

- [**area:dx**] feat: add a search by remote cluster command [#2209](https://github.com/elastic/observability-test-environments/pull/2209)
- [**area:dx**] feat: add cluster versions command [#2205](https://github.com/elastic/observability-test-environments/pull/2205)

---

## 2.0.26 (05/04/2022)
*No changelog for this release.*

---

## 2.0.25 (04/04/2022)
*No changelog for this release.*

---

## 2.0.24 (04/02/2022)
*No changelog for this release.*

---

## 2.0.15 (10/11/2021)
*No changelog for this release.*

---

## 2.0.23 (10/11/2021)
*No changelog for this release.*

---

## 2.0.14 (04/11/2021)
*No changelog for this release.*

---

## 2.0.22 (04/11/2021)
*No changelog for this release.*

---

## 2.0.13 (27/10/2021)
*No changelog for this release.*

---

## 2.0.21 (27/10/2021)
*No changelog for this release.*

---

## 2.0.12 (27/10/2021)
*No changelog for this release.*

---

## 2.0.20 (27/10/2021)
*No changelog for this release.*

---

## 2.0.10 (11/10/2021)
*No changelog for this release.*

---

## 2.0.11 (11/10/2021)
*No changelog for this release.*

---

## 2.0.9 (04/10/2021)
*No changelog for this release.*

---

## 2.0.8 (04/10/2021)
*No changelog for this release.*

---

## 2.0.7 (30/07/2021)
*No changelog for this release.*

---

## 2.0.6 (30/07/2021)

#### 🚀 Enhancements

-  feat: set default stack version across files [#594](https://github.com/elastic/observability-test-environments/pull/594)

---

## 2.0.5 (16/06/2020)
*No changelog for this release.*

---

## 2.0.4 (10/06/2020)
*No changelog for this release.*

---

## 2.0.3 (09/06/2020)
*No changelog for this release.*

---

## 2.0.2 (14/04/2020)
*No changelog for this release.*

---

## 2.0.1 (31/03/2020)
*No changelog for this release.*

---

## v1.0.0 (31/03/2020)

#### 🚀 Enhancements

-  feat: Apm mutating webhook [#2465](https://github.com/elastic/observability-test-environments/pull/2465)
-  feat: bump versions [#2337](https://github.com/elastic/observability-test-environments/pull/2337)
-  chore: review heartbeat config [#2493](https://github.com/elastic/observability-test-environments/pull/2493)
-  feat: backup secrets before update [#2336](https://github.com/elastic/observability-test-environments/pull/2336)
-  Test if CCS developer cluster can be created after update [#2329](https://github.com/elastic/observability-test-environments/pull/2329)
-  Add opbeans php [#2506](https://github.com/elastic/observability-test-environments/pull/2506)
-  feat: set default stack version across files [#594](https://github.com/elastic/observability-test-environments/pull/594)
- [**area:dx**] feat: bootstrap commands [#2307](https://github.com/elastic/observability-test-environments/pull/2307)
-  feat: script to generate the script to push licenses [#2303](https://github.com/elastic/observability-test-environments/pull/2303)
-  feat: add a Synthetics monitor to monitoring-oblt [#2300](https://github.com/elastic/observability-test-environments/pull/2300)
-  chore: Mix of issues [#2299](https://github.com/elastic/observability-test-environments/pull/2299)
-  feat: backup vault secrets [#2297](https://github.com/elastic/observability-test-environments/pull/2297)
-  fix: improve stack version script [#2296](https://github.com/elastic/observability-test-environments/pull/2296)
- [**area:dx**] feat: support for destroying a cluster [#2441](https://github.com/elastic/observability-test-environments/pull/2441)
- [**area:dx**] feat: add a command to create ESS clusters [#2408](https://github.com/elastic/observability-test-environments/pull/2408)
- [**area:dx**] feat: retrieve cluster secrets from oblt tools [#2385](https://github.com/elastic/observability-test-environments/pull/2385)
- [**area:dx**][**requested-by:Ingest Management**] fix: add monitoring cluster to the CCS Slack form [#2275](https://github.com/elastic/observability-test-environments/pull/2275)
- [**area:dx**] chore: refactor output for cluster templates command [#2272](https://github.com/elastic/observability-test-environments/pull/2272)
- [**area:dx**] feat: support showing commit ID and URL in messages [#2245](https://github.com/elastic/observability-test-environments/pull/2245)
- [**area:dx**] feat: add a slack command to list templates [#2235](https://github.com/elastic/observability-test-environments/pull/2235)
- [**area:dx**] feat: add a command to self-update the tool [#2223](https://github.com/elastic/observability-test-environments/pull/2223)
- [**area:dx**] feat: add a search by remote cluster command [#2209](https://github.com/elastic/observability-test-environments/pull/2209)
- [**area:dx**] feat: add cluster versions command [#2205](https://github.com/elastic/observability-test-environments/pull/2205)

#### 🐛 Bug Fixes

-  Add opbeans php [#2533](https://github.com/elastic/observability-test-environments/pull/2533)
-  fix: pull the repo before make the push [#2338](https://github.com/elastic/observability-test-environments/pull/2338)
-  fix: bump Helm charts versions [#2328](https://github.com/elastic/observability-test-environments/pull/2328)
-  fix: reimplement the Docker images pull and push process [#2530](https://github.com/elastic/observability-test-environments/pull/2530)
-  fix: avoid remove endoflines in config files [#2528](https://github.com/elastic/observability-test-environments/pull/2528)
-  fix: change authproxy secrets paths [#2527](https://github.com/elastic/observability-test-environments/pull/2527)
-  fix: ignore ess images for 7.x [#2525](https://github.com/elastic/observability-test-environments/pull/2525)
-  fix: increase matrix workers [#2524](https://github.com/elastic/observability-test-environments/pull/2524)
-  fix: increase the disk space on gobld workers [#2523](https://github.com/elastic/observability-test-environments/pull/2523)
-  fix: workaround to pull and push Elasticsearch Docker images. [#2522](https://github.com/elastic/observability-test-environments/pull/2522)
- [**area:dx**] fix: change the way we import the Docker images tgz [#2517](https://github.com/elastic/observability-test-environments/pull/2517)
-  fix: ECE deployment [#2325](https://github.com/elastic/observability-test-environments/pull/2325)
- [**area:dx**] chore: remove dry run from the deployment descriptor [#2316](https://github.com/elastic/observability-test-environments/pull/2316)
-  fix: set autoscale to true [#2302](https://github.com/elastic/observability-test-environments/pull/2302)
-  fix: remove nonexistent function [#2301](https://github.com/elastic/observability-test-environments/pull/2301)
- [**area:dx**] fix: do not override flag variables with env variables [#2401](https://github.com/elastic/observability-test-environments/pull/2401)
- [**area:ci**] fix: do not increment strings in groovy [#2229](https://github.com/elastic/observability-test-environments/pull/2229)
- [**area:dx**] chore: regen blob on releases [#2224](https://github.com/elastic/observability-test-environments/pull/2224)
- [**area:dx**] fix: proper base dir to create config dir [#2217](https://github.com/elastic/observability-test-environments/pull/2217)
- [**area:dx**] fix: create default config file's parent dir [#2213](https://github.com/elastic/observability-test-environments/pull/2213)

#### 📚 Documentation

- [**area:dx**] docs: user guide for the slack bot [#2403](https://github.com/elastic/observability-test-environments/pull/2403)

#### 🙈 No user affected

-  chore: remove tag from grafana configuration [#2340](https://github.com/elastic/observability-test-environments/pull/2340)
-  chore: template for update issues [#2332](https://github.com/elastic/observability-test-environments/pull/2332)
-  chore: bump release cluster to 8.4.1 [#2526](https://github.com/elastic/observability-test-environments/pull/2526)
-  chore: remove Helm Chart Elastic Stack support [#2494](https://github.com/elastic/observability-test-environments/pull/2494)
- [**area:dx**] chore: initialise oblt-repository just once for the CMDs in the cli [#2456](https://github.com/elastic/observability-test-environments/pull/2456)
- [**area:dx**] feat: do not remove the system-status repository each time [#2282](https://github.com/elastic/observability-test-environments/pull/2282)
- [**area:dx**] chore: release slack bot docker images on git tags [#2250](https://github.com/elastic/observability-test-environments/pull/2250)
- [**area:dx**] chore(slack): improve message for CCS creation [#2243](https://github.com/elastic/observability-test-environments/pull/2243)
- [**area:ci**][**area:dx**] chore: report cobertura for oblt-cli [#2242](https://github.com/elastic/observability-test-environments/pull/2242)
- [**area:dx**] fix: use project version for slackbot image [#2236](https://github.com/elastic/observability-test-environments/pull/2236)
- [**area:dx**] chore: always log current version if current is equals to latest [#2231](https://github.com/elastic/observability-test-environments/pull/2231)

---

## 2.0.30 (01/01/1970)

#### 🚀 Enhancements

-  feat: Apm mutating webhook [#2465](https://github.com/elastic/observability-test-environments/pull/2465)
-  feat: bump versions [#2337](https://github.com/elastic/observability-test-environments/pull/2337)
-  chore: review heartbeat config [#2493](https://github.com/elastic/observability-test-environments/pull/2493)
-  feat: backup secrets before update [#2336](https://github.com/elastic/observability-test-environments/pull/2336)
-  Test if CCS developer cluster can be created after update [#2329](https://github.com/elastic/observability-test-environments/pull/2329)
-  Add opbeans php [#2506](https://github.com/elastic/observability-test-environments/pull/2506)
-  feat: set default stack version across files [#594](https://github.com/elastic/observability-test-environments/pull/594)
- [**area:dx**] feat: bootstrap commands [#2307](https://github.com/elastic/observability-test-environments/pull/2307)
-  feat: script to generate the script to push licenses [#2303](https://github.com/elastic/observability-test-environments/pull/2303)
-  feat: add a Synthetics monitor to monitoring-oblt [#2300](https://github.com/elastic/observability-test-environments/pull/2300)
-  chore: Mix of issues [#2299](https://github.com/elastic/observability-test-environments/pull/2299)
-  feat: backup vault secrets [#2297](https://github.com/elastic/observability-test-environments/pull/2297)
-  fix: improve stack version script [#2296](https://github.com/elastic/observability-test-environments/pull/2296)
- [**area:dx**] feat: support for destroying a cluster [#2441](https://github.com/elastic/observability-test-environments/pull/2441)
- [**area:dx**] feat: add a command to create ESS clusters [#2408](https://github.com/elastic/observability-test-environments/pull/2408)
- [**area:dx**] feat: retrieve cluster secrets from oblt tools [#2385](https://github.com/elastic/observability-test-environments/pull/2385)
- [**area:dx**][**requested-by:Ingest Management**] fix: add monitoring cluster to the CCS Slack form [#2275](https://github.com/elastic/observability-test-environments/pull/2275)
- [**area:dx**] chore: refactor output for cluster templates command [#2272](https://github.com/elastic/observability-test-environments/pull/2272)
- [**area:dx**] feat: support showing commit ID and URL in messages [#2245](https://github.com/elastic/observability-test-environments/pull/2245)
- [**area:dx**] feat: add a slack command to list templates [#2235](https://github.com/elastic/observability-test-environments/pull/2235)
- [**area:dx**] feat: add a command to self-update the tool [#2223](https://github.com/elastic/observability-test-environments/pull/2223)
- [**area:dx**] feat: add a search by remote cluster command [#2209](https://github.com/elastic/observability-test-environments/pull/2209)
- [**area:dx**] feat: add cluster versions command [#2205](https://github.com/elastic/observability-test-environments/pull/2205)

#### 🐛 Bug Fixes

-  Add opbeans php [#2533](https://github.com/elastic/observability-test-environments/pull/2533)
-  fix: pull the repo before make the push [#2338](https://github.com/elastic/observability-test-environments/pull/2338)
-  fix: bump Helm charts versions [#2328](https://github.com/elastic/observability-test-environments/pull/2328)
-  fix: reimplement the Docker images pull and push process [#2530](https://github.com/elastic/observability-test-environments/pull/2530)
-  fix: avoid remove endoflines in config files [#2528](https://github.com/elastic/observability-test-environments/pull/2528)
-  fix: change authproxy secrets paths [#2527](https://github.com/elastic/observability-test-environments/pull/2527)
-  fix: ignore ess images for 7.x [#2525](https://github.com/elastic/observability-test-environments/pull/2525)
-  fix: increase matrix workers [#2524](https://github.com/elastic/observability-test-environments/pull/2524)
-  fix: increase the disk space on gobld workers [#2523](https://github.com/elastic/observability-test-environments/pull/2523)
-  fix: workaround to pull and push Elasticsearch Docker images. [#2522](https://github.com/elastic/observability-test-environments/pull/2522)
- [**area:dx**] fix: change the way we import the Docker images tgz [#2517](https://github.com/elastic/observability-test-environments/pull/2517)
-  fix: ECE deployment [#2325](https://github.com/elastic/observability-test-environments/pull/2325)
- [**area:dx**] chore: remove dry run from the deployment descriptor [#2316](https://github.com/elastic/observability-test-environments/pull/2316)
-  fix: set autoscale to true [#2302](https://github.com/elastic/observability-test-environments/pull/2302)
-  fix: remove nonexistent function [#2301](https://github.com/elastic/observability-test-environments/pull/2301)
- [**area:dx**] fix: do not override flag variables with env variables [#2401](https://github.com/elastic/observability-test-environments/pull/2401)
- [**area:ci**] fix: do not increment strings in groovy [#2229](https://github.com/elastic/observability-test-environments/pull/2229)
- [**area:dx**] chore: regen blob on releases [#2224](https://github.com/elastic/observability-test-environments/pull/2224)
- [**area:dx**] fix: proper base dir to create config dir [#2217](https://github.com/elastic/observability-test-environments/pull/2217)
- [**area:dx**] fix: create default config file's parent dir [#2213](https://github.com/elastic/observability-test-environments/pull/2213)

#### 📚 Documentation

- [**area:dx**] docs: user guide for the slack bot [#2403](https://github.com/elastic/observability-test-environments/pull/2403)

#### 🙈 No user affected

-  chore: remove tag from grafana configuration [#2340](https://github.com/elastic/observability-test-environments/pull/2340)
-  chore: template for update issues [#2332](https://github.com/elastic/observability-test-environments/pull/2332)
-  chore: bump release cluster to 8.4.1 [#2526](https://github.com/elastic/observability-test-environments/pull/2526)
-  chore: remove Helm Chart Elastic Stack support [#2494](https://github.com/elastic/observability-test-environments/pull/2494)
- [**area:dx**] chore: initialise oblt-repository just once for the CMDs in the cli [#2456](https://github.com/elastic/observability-test-environments/pull/2456)
- [**area:dx**] feat: do not remove the system-status repository each time [#2282](https://github.com/elastic/observability-test-environments/pull/2282)
- [**area:dx**] chore: release slack bot docker images on git tags [#2250](https://github.com/elastic/observability-test-environments/pull/2250)
- [**area:dx**] chore(slack): improve message for CCS creation [#2243](https://github.com/elastic/observability-test-environments/pull/2243)
- [**area:ci**][**area:dx**] chore: report cobertura for oblt-cli [#2242](https://github.com/elastic/observability-test-environments/pull/2242)
- [**area:dx**] fix: use project version for slackbot image [#2236](https://github.com/elastic/observability-test-environments/pull/2236)
- [**area:dx**] chore: always log current version if current is equals to latest [#2231](https://github.com/elastic/observability-test-environments/pull/2231)

---

## 2.0.19 (01/01/1970)

#### 🚀 Enhancements

-  feat: Apm mutating webhook [#2465](https://github.com/elastic/observability-test-environments/pull/2465)
-  feat: bump versions [#2337](https://github.com/elastic/observability-test-environments/pull/2337)
-  chore: review heartbeat config [#2493](https://github.com/elastic/observability-test-environments/pull/2493)
-  feat: backup secrets before update [#2336](https://github.com/elastic/observability-test-environments/pull/2336)
-  Test if CCS developer cluster can be created after update [#2329](https://github.com/elastic/observability-test-environments/pull/2329)
-  Add opbeans php [#2506](https://github.com/elastic/observability-test-environments/pull/2506)
-  feat: set default stack version across files [#594](https://github.com/elastic/observability-test-environments/pull/594)
- [**area:dx**] feat: bootstrap commands [#2307](https://github.com/elastic/observability-test-environments/pull/2307)
-  feat: script to generate the script to push licenses [#2303](https://github.com/elastic/observability-test-environments/pull/2303)
-  feat: add a Synthetics monitor to monitoring-oblt [#2300](https://github.com/elastic/observability-test-environments/pull/2300)
-  chore: Mix of issues [#2299](https://github.com/elastic/observability-test-environments/pull/2299)
-  feat: backup vault secrets [#2297](https://github.com/elastic/observability-test-environments/pull/2297)
-  fix: improve stack version script [#2296](https://github.com/elastic/observability-test-environments/pull/2296)
- [**area:dx**] feat: support for destroying a cluster [#2441](https://github.com/elastic/observability-test-environments/pull/2441)
- [**area:dx**] feat: add a command to create ESS clusters [#2408](https://github.com/elastic/observability-test-environments/pull/2408)
- [**area:dx**] feat: retrieve cluster secrets from oblt tools [#2385](https://github.com/elastic/observability-test-environments/pull/2385)
- [**area:dx**][**requested-by:Ingest Management**] fix: add monitoring cluster to the CCS Slack form [#2275](https://github.com/elastic/observability-test-environments/pull/2275)
- [**area:dx**] chore: refactor output for cluster templates command [#2272](https://github.com/elastic/observability-test-environments/pull/2272)
- [**area:dx**] feat: support showing commit ID and URL in messages [#2245](https://github.com/elastic/observability-test-environments/pull/2245)
- [**area:dx**] feat: add a slack command to list templates [#2235](https://github.com/elastic/observability-test-environments/pull/2235)
- [**area:dx**] feat: add a command to self-update the tool [#2223](https://github.com/elastic/observability-test-environments/pull/2223)
- [**area:dx**] feat: add a search by remote cluster command [#2209](https://github.com/elastic/observability-test-environments/pull/2209)
- [**area:dx**] feat: add cluster versions command [#2205](https://github.com/elastic/observability-test-environments/pull/2205)

#### 🐛 Bug Fixes

-  Add opbeans php [#2533](https://github.com/elastic/observability-test-environments/pull/2533)
-  fix: pull the repo before make the push [#2338](https://github.com/elastic/observability-test-environments/pull/2338)
-  fix: bump Helm charts versions [#2328](https://github.com/elastic/observability-test-environments/pull/2328)
-  fix: reimplement the Docker images pull and push process [#2530](https://github.com/elastic/observability-test-environments/pull/2530)
-  fix: avoid remove endoflines in config files [#2528](https://github.com/elastic/observability-test-environments/pull/2528)
-  fix: change authproxy secrets paths [#2527](https://github.com/elastic/observability-test-environments/pull/2527)
-  fix: ignore ess images for 7.x [#2525](https://github.com/elastic/observability-test-environments/pull/2525)
-  fix: increase matrix workers [#2524](https://github.com/elastic/observability-test-environments/pull/2524)
-  fix: increase the disk space on gobld workers [#2523](https://github.com/elastic/observability-test-environments/pull/2523)
-  fix: workaround to pull and push Elasticsearch Docker images. [#2522](https://github.com/elastic/observability-test-environments/pull/2522)
- [**area:dx**] fix: change the way we import the Docker images tgz [#2517](https://github.com/elastic/observability-test-environments/pull/2517)
-  fix: ECE deployment [#2325](https://github.com/elastic/observability-test-environments/pull/2325)
- [**area:dx**] chore: remove dry run from the deployment descriptor [#2316](https://github.com/elastic/observability-test-environments/pull/2316)
-  fix: set autoscale to true [#2302](https://github.com/elastic/observability-test-environments/pull/2302)
-  fix: remove nonexistent function [#2301](https://github.com/elastic/observability-test-environments/pull/2301)
- [**area:dx**] fix: do not override flag variables with env variables [#2401](https://github.com/elastic/observability-test-environments/pull/2401)
- [**area:ci**] fix: do not increment strings in groovy [#2229](https://github.com/elastic/observability-test-environments/pull/2229)
- [**area:dx**] chore: regen blob on releases [#2224](https://github.com/elastic/observability-test-environments/pull/2224)
- [**area:dx**] fix: proper base dir to create config dir [#2217](https://github.com/elastic/observability-test-environments/pull/2217)
- [**area:dx**] fix: create default config file's parent dir [#2213](https://github.com/elastic/observability-test-environments/pull/2213)

#### 📚 Documentation

- [**area:dx**] docs: user guide for the slack bot [#2403](https://github.com/elastic/observability-test-environments/pull/2403)

#### 🙈 No user affected

-  chore: remove tag from grafana configuration [#2340](https://github.com/elastic/observability-test-environments/pull/2340)
-  chore: template for update issues [#2332](https://github.com/elastic/observability-test-environments/pull/2332)
-  chore: bump release cluster to 8.4.1 [#2526](https://github.com/elastic/observability-test-environments/pull/2526)
-  chore: remove Helm Chart Elastic Stack support [#2494](https://github.com/elastic/observability-test-environments/pull/2494)
- [**area:dx**] chore: initialise oblt-repository just once for the CMDs in the cli [#2456](https://github.com/elastic/observability-test-environments/pull/2456)
- [**area:dx**] feat: do not remove the system-status repository each time [#2282](https://github.com/elastic/observability-test-environments/pull/2282)
- [**area:dx**] chore: release slack bot docker images on git tags [#2250](https://github.com/elastic/observability-test-environments/pull/2250)
- [**area:dx**] chore(slack): improve message for CCS creation [#2243](https://github.com/elastic/observability-test-environments/pull/2243)
- [**area:ci**][**area:dx**] chore: report cobertura for oblt-cli [#2242](https://github.com/elastic/observability-test-environments/pull/2242)
- [**area:dx**] fix: use project version for slackbot image [#2236](https://github.com/elastic/observability-test-environments/pull/2236)
- [**area:dx**] chore: always log current version if current is equals to latest [#2231](https://github.com/elastic/observability-test-environments/pull/2231)

---

## 2.0.18 (01/01/1970)

#### 🚀 Enhancements

-  feat: Apm mutating webhook [#2465](https://github.com/elastic/observability-test-environments/pull/2465)
-  feat: bump versions [#2337](https://github.com/elastic/observability-test-environments/pull/2337)
-  chore: review heartbeat config [#2493](https://github.com/elastic/observability-test-environments/pull/2493)
-  feat: backup secrets before update [#2336](https://github.com/elastic/observability-test-environments/pull/2336)
-  Test if CCS developer cluster can be created after update [#2329](https://github.com/elastic/observability-test-environments/pull/2329)
-  Add opbeans php [#2506](https://github.com/elastic/observability-test-environments/pull/2506)
-  feat: set default stack version across files [#594](https://github.com/elastic/observability-test-environments/pull/594)
- [**area:dx**] feat: bootstrap commands [#2307](https://github.com/elastic/observability-test-environments/pull/2307)
-  feat: script to generate the script to push licenses [#2303](https://github.com/elastic/observability-test-environments/pull/2303)
-  feat: add a Synthetics monitor to monitoring-oblt [#2300](https://github.com/elastic/observability-test-environments/pull/2300)
-  chore: Mix of issues [#2299](https://github.com/elastic/observability-test-environments/pull/2299)
-  feat: backup vault secrets [#2297](https://github.com/elastic/observability-test-environments/pull/2297)
-  fix: improve stack version script [#2296](https://github.com/elastic/observability-test-environments/pull/2296)
- [**area:dx**] feat: support for destroying a cluster [#2441](https://github.com/elastic/observability-test-environments/pull/2441)
- [**area:dx**] feat: add a command to create ESS clusters [#2408](https://github.com/elastic/observability-test-environments/pull/2408)
- [**area:dx**] feat: retrieve cluster secrets from oblt tools [#2385](https://github.com/elastic/observability-test-environments/pull/2385)
- [**area:dx**][**requested-by:Ingest Management**] fix: add monitoring cluster to the CCS Slack form [#2275](https://github.com/elastic/observability-test-environments/pull/2275)
- [**area:dx**] chore: refactor output for cluster templates command [#2272](https://github.com/elastic/observability-test-environments/pull/2272)
- [**area:dx**] feat: support showing commit ID and URL in messages [#2245](https://github.com/elastic/observability-test-environments/pull/2245)
- [**area:dx**] feat: add a slack command to list templates [#2235](https://github.com/elastic/observability-test-environments/pull/2235)
- [**area:dx**] feat: add a command to self-update the tool [#2223](https://github.com/elastic/observability-test-environments/pull/2223)
- [**area:dx**] feat: add a search by remote cluster command [#2209](https://github.com/elastic/observability-test-environments/pull/2209)
- [**area:dx**] feat: add cluster versions command [#2205](https://github.com/elastic/observability-test-environments/pull/2205)

#### 🐛 Bug Fixes

-  Add opbeans php [#2533](https://github.com/elastic/observability-test-environments/pull/2533)
-  fix: pull the repo before make the push [#2338](https://github.com/elastic/observability-test-environments/pull/2338)
-  fix: bump Helm charts versions [#2328](https://github.com/elastic/observability-test-environments/pull/2328)
-  fix: reimplement the Docker images pull and push process [#2530](https://github.com/elastic/observability-test-environments/pull/2530)
-  fix: avoid remove endoflines in config files [#2528](https://github.com/elastic/observability-test-environments/pull/2528)
-  fix: change authproxy secrets paths [#2527](https://github.com/elastic/observability-test-environments/pull/2527)
-  fix: ignore ess images for 7.x [#2525](https://github.com/elastic/observability-test-environments/pull/2525)
-  fix: increase matrix workers [#2524](https://github.com/elastic/observability-test-environments/pull/2524)
-  fix: increase the disk space on gobld workers [#2523](https://github.com/elastic/observability-test-environments/pull/2523)
-  fix: workaround to pull and push Elasticsearch Docker images. [#2522](https://github.com/elastic/observability-test-environments/pull/2522)
- [**area:dx**] fix: change the way we import the Docker images tgz [#2517](https://github.com/elastic/observability-test-environments/pull/2517)
-  fix: ECE deployment [#2325](https://github.com/elastic/observability-test-environments/pull/2325)
- [**area:dx**] chore: remove dry run from the deployment descriptor [#2316](https://github.com/elastic/observability-test-environments/pull/2316)
-  fix: set autoscale to true [#2302](https://github.com/elastic/observability-test-environments/pull/2302)
-  fix: remove nonexistent function [#2301](https://github.com/elastic/observability-test-environments/pull/2301)
- [**area:dx**] fix: do not override flag variables with env variables [#2401](https://github.com/elastic/observability-test-environments/pull/2401)
- [**area:ci**] fix: do not increment strings in groovy [#2229](https://github.com/elastic/observability-test-environments/pull/2229)
- [**area:dx**] chore: regen blob on releases [#2224](https://github.com/elastic/observability-test-environments/pull/2224)
- [**area:dx**] fix: proper base dir to create config dir [#2217](https://github.com/elastic/observability-test-environments/pull/2217)
- [**area:dx**] fix: create default config file's parent dir [#2213](https://github.com/elastic/observability-test-environments/pull/2213)

#### 📚 Documentation

- [**area:dx**] docs: user guide for the slack bot [#2403](https://github.com/elastic/observability-test-environments/pull/2403)

#### 🙈 No user affected

-  chore: remove tag from grafana configuration [#2340](https://github.com/elastic/observability-test-environments/pull/2340)
-  chore: template for update issues [#2332](https://github.com/elastic/observability-test-environments/pull/2332)
-  chore: bump release cluster to 8.4.1 [#2526](https://github.com/elastic/observability-test-environments/pull/2526)
-  chore: remove Helm Chart Elastic Stack support [#2494](https://github.com/elastic/observability-test-environments/pull/2494)
- [**area:dx**] chore: initialise oblt-repository just once for the CMDs in the cli [#2456](https://github.com/elastic/observability-test-environments/pull/2456)
- [**area:dx**] feat: do not remove the system-status repository each time [#2282](https://github.com/elastic/observability-test-environments/pull/2282)
- [**area:dx**] chore: release slack bot docker images on git tags [#2250](https://github.com/elastic/observability-test-environments/pull/2250)
- [**area:dx**] chore(slack): improve message for CCS creation [#2243](https://github.com/elastic/observability-test-environments/pull/2243)
- [**area:ci**][**area:dx**] chore: report cobertura for oblt-cli [#2242](https://github.com/elastic/observability-test-environments/pull/2242)
- [**area:dx**] fix: use project version for slackbot image [#2236](https://github.com/elastic/observability-test-environments/pull/2236)
- [**area:dx**] chore: always log current version if current is equals to latest [#2231](https://github.com/elastic/observability-test-environments/pull/2231)

---

## 2.0.17 (01/01/1970)

#### 🚀 Enhancements

-  feat: Apm mutating webhook [#2465](https://github.com/elastic/observability-test-environments/pull/2465)
-  feat: bump versions [#2337](https://github.com/elastic/observability-test-environments/pull/2337)
-  chore: review heartbeat config [#2493](https://github.com/elastic/observability-test-environments/pull/2493)
-  feat: backup secrets before update [#2336](https://github.com/elastic/observability-test-environments/pull/2336)
-  Test if CCS developer cluster can be created after update [#2329](https://github.com/elastic/observability-test-environments/pull/2329)
-  Add opbeans php [#2506](https://github.com/elastic/observability-test-environments/pull/2506)
-  feat: set default stack version across files [#594](https://github.com/elastic/observability-test-environments/pull/594)
- [**area:dx**] feat: bootstrap commands [#2307](https://github.com/elastic/observability-test-environments/pull/2307)
-  feat: script to generate the script to push licenses [#2303](https://github.com/elastic/observability-test-environments/pull/2303)
-  feat: add a Synthetics monitor to monitoring-oblt [#2300](https://github.com/elastic/observability-test-environments/pull/2300)
-  chore: Mix of issues [#2299](https://github.com/elastic/observability-test-environments/pull/2299)
-  feat: backup vault secrets [#2297](https://github.com/elastic/observability-test-environments/pull/2297)
-  fix: improve stack version script [#2296](https://github.com/elastic/observability-test-environments/pull/2296)
- [**area:dx**] feat: support for destroying a cluster [#2441](https://github.com/elastic/observability-test-environments/pull/2441)
- [**area:dx**] feat: add a command to create ESS clusters [#2408](https://github.com/elastic/observability-test-environments/pull/2408)
- [**area:dx**] feat: retrieve cluster secrets from oblt tools [#2385](https://github.com/elastic/observability-test-environments/pull/2385)
- [**area:dx**][**requested-by:Ingest Management**] fix: add monitoring cluster to the CCS Slack form [#2275](https://github.com/elastic/observability-test-environments/pull/2275)
- [**area:dx**] chore: refactor output for cluster templates command [#2272](https://github.com/elastic/observability-test-environments/pull/2272)
- [**area:dx**] feat: support showing commit ID and URL in messages [#2245](https://github.com/elastic/observability-test-environments/pull/2245)
- [**area:dx**] feat: add a slack command to list templates [#2235](https://github.com/elastic/observability-test-environments/pull/2235)
- [**area:dx**] feat: add a command to self-update the tool [#2223](https://github.com/elastic/observability-test-environments/pull/2223)
- [**area:dx**] feat: add a search by remote cluster command [#2209](https://github.com/elastic/observability-test-environments/pull/2209)
- [**area:dx**] feat: add cluster versions command [#2205](https://github.com/elastic/observability-test-environments/pull/2205)

#### 🐛 Bug Fixes

-  Add opbeans php [#2533](https://github.com/elastic/observability-test-environments/pull/2533)
-  fix: pull the repo before make the push [#2338](https://github.com/elastic/observability-test-environments/pull/2338)
-  fix: bump Helm charts versions [#2328](https://github.com/elastic/observability-test-environments/pull/2328)
-  fix: reimplement the Docker images pull and push process [#2530](https://github.com/elastic/observability-test-environments/pull/2530)
-  fix: avoid remove endoflines in config files [#2528](https://github.com/elastic/observability-test-environments/pull/2528)
-  fix: change authproxy secrets paths [#2527](https://github.com/elastic/observability-test-environments/pull/2527)
-  fix: ignore ess images for 7.x [#2525](https://github.com/elastic/observability-test-environments/pull/2525)
-  fix: increase matrix workers [#2524](https://github.com/elastic/observability-test-environments/pull/2524)
-  fix: increase the disk space on gobld workers [#2523](https://github.com/elastic/observability-test-environments/pull/2523)
-  fix: workaround to pull and push Elasticsearch Docker images. [#2522](https://github.com/elastic/observability-test-environments/pull/2522)
- [**area:dx**] fix: change the way we import the Docker images tgz [#2517](https://github.com/elastic/observability-test-environments/pull/2517)
-  fix: ECE deployment [#2325](https://github.com/elastic/observability-test-environments/pull/2325)
- [**area:dx**] chore: remove dry run from the deployment descriptor [#2316](https://github.com/elastic/observability-test-environments/pull/2316)
-  fix: set autoscale to true [#2302](https://github.com/elastic/observability-test-environments/pull/2302)
-  fix: remove nonexistent function [#2301](https://github.com/elastic/observability-test-environments/pull/2301)
- [**area:dx**] fix: do not override flag variables with env variables [#2401](https://github.com/elastic/observability-test-environments/pull/2401)
- [**area:ci**] fix: do not increment strings in groovy [#2229](https://github.com/elastic/observability-test-environments/pull/2229)
- [**area:dx**] chore: regen blob on releases [#2224](https://github.com/elastic/observability-test-environments/pull/2224)
- [**area:dx**] fix: proper base dir to create config dir [#2217](https://github.com/elastic/observability-test-environments/pull/2217)
- [**area:dx**] fix: create default config file's parent dir [#2213](https://github.com/elastic/observability-test-environments/pull/2213)

#### 📚 Documentation

- [**area:dx**] docs: user guide for the slack bot [#2403](https://github.com/elastic/observability-test-environments/pull/2403)

#### 🙈 No user affected

-  chore: remove tag from grafana configuration [#2340](https://github.com/elastic/observability-test-environments/pull/2340)
-  chore: template for update issues [#2332](https://github.com/elastic/observability-test-environments/pull/2332)
-  chore: bump release cluster to 8.4.1 [#2526](https://github.com/elastic/observability-test-environments/pull/2526)
-  chore: remove Helm Chart Elastic Stack support [#2494](https://github.com/elastic/observability-test-environments/pull/2494)
- [**area:dx**] chore: initialise oblt-repository just once for the CMDs in the cli [#2456](https://github.com/elastic/observability-test-environments/pull/2456)
- [**area:dx**] feat: do not remove the system-status repository each time [#2282](https://github.com/elastic/observability-test-environments/pull/2282)
- [**area:dx**] chore: release slack bot docker images on git tags [#2250](https://github.com/elastic/observability-test-environments/pull/2250)
- [**area:dx**] chore(slack): improve message for CCS creation [#2243](https://github.com/elastic/observability-test-environments/pull/2243)
- [**area:ci**][**area:dx**] chore: report cobertura for oblt-cli [#2242](https://github.com/elastic/observability-test-environments/pull/2242)
- [**area:dx**] fix: use project version for slackbot image [#2236](https://github.com/elastic/observability-test-environments/pull/2236)
- [**area:dx**] chore: always log current version if current is equals to latest [#2231](https://github.com/elastic/observability-test-environments/pull/2231)

---

## 2.0.16 (01/01/1970)

#### 🚀 Enhancements

-  feat: Apm mutating webhook [#2465](https://github.com/elastic/observability-test-environments/pull/2465)
-  feat: bump versions [#2337](https://github.com/elastic/observability-test-environments/pull/2337)
-  chore: review heartbeat config [#2493](https://github.com/elastic/observability-test-environments/pull/2493)
-  feat: backup secrets before update [#2336](https://github.com/elastic/observability-test-environments/pull/2336)
-  Test if CCS developer cluster can be created after update [#2329](https://github.com/elastic/observability-test-environments/pull/2329)
-  Add opbeans php [#2506](https://github.com/elastic/observability-test-environments/pull/2506)
-  feat: set default stack version across files [#594](https://github.com/elastic/observability-test-environments/pull/594)
- [**area:dx**] feat: bootstrap commands [#2307](https://github.com/elastic/observability-test-environments/pull/2307)
-  feat: script to generate the script to push licenses [#2303](https://github.com/elastic/observability-test-environments/pull/2303)
-  feat: add a Synthetics monitor to monitoring-oblt [#2300](https://github.com/elastic/observability-test-environments/pull/2300)
-  chore: Mix of issues [#2299](https://github.com/elastic/observability-test-environments/pull/2299)
-  feat: backup vault secrets [#2297](https://github.com/elastic/observability-test-environments/pull/2297)
-  fix: improve stack version script [#2296](https://github.com/elastic/observability-test-environments/pull/2296)
- [**area:dx**] feat: support for destroying a cluster [#2441](https://github.com/elastic/observability-test-environments/pull/2441)
- [**area:dx**] feat: add a command to create ESS clusters [#2408](https://github.com/elastic/observability-test-environments/pull/2408)
- [**area:dx**] feat: retrieve cluster secrets from oblt tools [#2385](https://github.com/elastic/observability-test-environments/pull/2385)
- [**area:dx**][**requested-by:Ingest Management**] fix: add monitoring cluster to the CCS Slack form [#2275](https://github.com/elastic/observability-test-environments/pull/2275)
- [**area:dx**] chore: refactor output for cluster templates command [#2272](https://github.com/elastic/observability-test-environments/pull/2272)
- [**area:dx**] feat: support showing commit ID and URL in messages [#2245](https://github.com/elastic/observability-test-environments/pull/2245)
- [**area:dx**] feat: add a slack command to list templates [#2235](https://github.com/elastic/observability-test-environments/pull/2235)
- [**area:dx**] feat: add a command to self-update the tool [#2223](https://github.com/elastic/observability-test-environments/pull/2223)
- [**area:dx**] feat: add a search by remote cluster command [#2209](https://github.com/elastic/observability-test-environments/pull/2209)
- [**area:dx**] feat: add cluster versions command [#2205](https://github.com/elastic/observability-test-environments/pull/2205)

#### 🐛 Bug Fixes

-  Add opbeans php [#2533](https://github.com/elastic/observability-test-environments/pull/2533)
-  fix: pull the repo before make the push [#2338](https://github.com/elastic/observability-test-environments/pull/2338)
-  fix: bump Helm charts versions [#2328](https://github.com/elastic/observability-test-environments/pull/2328)
-  fix: reimplement the Docker images pull and push process [#2530](https://github.com/elastic/observability-test-environments/pull/2530)
-  fix: avoid remove endoflines in config files [#2528](https://github.com/elastic/observability-test-environments/pull/2528)
-  fix: change authproxy secrets paths [#2527](https://github.com/elastic/observability-test-environments/pull/2527)
-  fix: ignore ess images for 7.x [#2525](https://github.com/elastic/observability-test-environments/pull/2525)
-  fix: increase matrix workers [#2524](https://github.com/elastic/observability-test-environments/pull/2524)
-  fix: increase the disk space on gobld workers [#2523](https://github.com/elastic/observability-test-environments/pull/2523)
-  fix: workaround to pull and push Elasticsearch Docker images. [#2522](https://github.com/elastic/observability-test-environments/pull/2522)
- [**area:dx**] fix: change the way we import the Docker images tgz [#2517](https://github.com/elastic/observability-test-environments/pull/2517)
-  fix: ECE deployment [#2325](https://github.com/elastic/observability-test-environments/pull/2325)
- [**area:dx**] chore: remove dry run from the deployment descriptor [#2316](https://github.com/elastic/observability-test-environments/pull/2316)
-  fix: set autoscale to true [#2302](https://github.com/elastic/observability-test-environments/pull/2302)
-  fix: remove nonexistent function [#2301](https://github.com/elastic/observability-test-environments/pull/2301)
- [**area:dx**] fix: do not override flag variables with env variables [#2401](https://github.com/elastic/observability-test-environments/pull/2401)
- [**area:ci**] fix: do not increment strings in groovy [#2229](https://github.com/elastic/observability-test-environments/pull/2229)
- [**area:dx**] chore: regen blob on releases [#2224](https://github.com/elastic/observability-test-environments/pull/2224)
- [**area:dx**] fix: proper base dir to create config dir [#2217](https://github.com/elastic/observability-test-environments/pull/2217)
- [**area:dx**] fix: create default config file's parent dir [#2213](https://github.com/elastic/observability-test-environments/pull/2213)

#### 📚 Documentation

- [**area:dx**] docs: user guide for the slack bot [#2403](https://github.com/elastic/observability-test-environments/pull/2403)

#### 🙈 No user affected

-  chore: remove tag from grafana configuration [#2340](https://github.com/elastic/observability-test-environments/pull/2340)
-  chore: template for update issues [#2332](https://github.com/elastic/observability-test-environments/pull/2332)
-  chore: bump release cluster to 8.4.1 [#2526](https://github.com/elastic/observability-test-environments/pull/2526)
-  chore: remove Helm Chart Elastic Stack support [#2494](https://github.com/elastic/observability-test-environments/pull/2494)
- [**area:dx**] chore: initialise oblt-repository just once for the CMDs in the cli [#2456](https://github.com/elastic/observability-test-environments/pull/2456)
- [**area:dx**] feat: do not remove the system-status repository each time [#2282](https://github.com/elastic/observability-test-environments/pull/2282)
- [**area:dx**] chore: release slack bot docker images on git tags [#2250](https://github.com/elastic/observability-test-environments/pull/2250)
- [**area:dx**] chore(slack): improve message for CCS creation [#2243](https://github.com/elastic/observability-test-environments/pull/2243)
- [**area:ci**][**area:dx**] chore: report cobertura for oblt-cli [#2242](https://github.com/elastic/observability-test-environments/pull/2242)
- [**area:dx**] fix: use project version for slackbot image [#2236](https://github.com/elastic/observability-test-environments/pull/2236)
- [**area:dx**] chore: always log current version if current is equals to latest [#2231](https://github.com/elastic/observability-test-environments/pull/2231)
