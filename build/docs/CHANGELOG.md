# Changelog

All notable changes to this module will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.0.0.42] - 2023-05-12

- Added private keys for all (except fisherman) actors
- Changed the debug_keybase package to support multiple yaml secrets in one yaml file
- Added full node (non-staked validator)

## [0.0.0.41] - 2023-05-08

- Updated Dockerfiles using outdated go version to 1.19

## [0.0.0.40] - 2023-05-04

- Added `network_id` parameter to the node config files

## [0.0.0.39] - 2023-04-28

- Added a `fisherman_per_session` governance parameter
- Updated the default `blocks_per_session` from `4` to `1`

## [0.0.0.38] - 2023-04-28

- Removed unused `tmpDir` in `debug_keybase` package
- Changed the name of the secret that holds private keys so it will be the same across local and dev networks.
- Moved `regexp.MustCompile` into `init()` to avoid recompiling the regex on every call (requested in previus PR)

## [0.0.0.37] - 2023-04-24

- Adds kubectl to dev Dockerfile image

## [0.0.0.36] - 2023-04-19

- Changed validator DNS names to match new naming convention (again, helm chart was renamed)
- Changed the way `cluster-manager` looks up the new validators to be staked.

## [0.0.0.35] - 2023-04-17

- Removed runtime/configs.Config#UseLibp2p field
- Use pod IP for validator DNS resolution tilt localnet
- Add `LIBP2P_DEBUG` env var

## [0.0.0.34] - 2023-04-14

- Changed LocalNet validators to use the new `pocket-validator` helm chart instead of templating the manifests with `sed`.
- Each validator now has it's own postgres instance (as a helm chart dependency), which allows for clean scale up/down of the validators.
- Cleaned up old manifests and scripts that are no longer needed.
- Changed LocalNet documentation to reflect the new changes.

## [0.0.0.33] - 2023-04-13

- Add persistent txIndexerPath to node configs

## [0.0.0.32] - 2023-04-10

- Adds e2e-tests button to Tiltfile

## [0.0.0.31] - 2023-04-06

- Updated `genesis.json` and `configs.yaml` to reflect pools address changes

## [0.0.0.30] - 2023-03-31

- Include `cluster-manager` to `-dev` flavor of container images.

## [0.0.0.29] - 2023-03-30

- `cluster-manager` now waits for `v1-validator001` to be online AND responsive by checking the `/v1/health` endpoint (dogfooding)
- `cluster-manager` skips auto staking for the validators that are already staked in genesis

## [0.0.0.28] - 2023-03-30

- Update `pacemaker_timeout` from 5 to 10 seconds to make the logging output less noisy during development
- Updated configurations related to postgres connection pooling

## [0.0.0.27] - 2023-03-28

- Silence gopls build error about missing DebugKeybaseBackup

## [0.0.0.26] - 2023-03-28

- Make k8s distribution recommendation more opinionated

## [0.0.0.25] - 2023-03-24

- Introduced a new binary that's used to check if the debug_keybase.bak is up to date with the private-keys.yaml file and to update it if it's not.

## [0.0.0.24] - 2023-03-17

- Added resource limits and PVC for debug client pod.
- Added `procps` to the client pod image.

## [0.0.0.23] - 2023-03-14

- Simplifies the debug CLI tooling by embedding private-keys.yaml manifest
  into the CLI binary when the debug build tag is present.

## [0.0.0.22] - 2023-03-08

- add permissions to get private keys secret on cluster-manager-account service account

## [0.0.0.21] - 2023-03-02

- added HOME environment variable to image build to fix tilt installation

## [0.0.0.20] - 2023-03-01

- replace `consensus_port` with `port` in P2P config
- update default P2P config `port` to from `8080` to `42069`
- add `use_libp2p` field to base config
- add `hostname` field to P2P config

## [0.0.0.19] - 2023-02-28

- Renamed `generic_param` to `service_url` in the config files
- Renamed a few governance parameters to make self explanatory

## [0.0.0.18] - 2023-02-21

- Rename ServiceNode Actor Type Name to Servicer

## [0.0.0.17] - 2023-02-21

- Updated `docker-compose` to allow for editing port mappings via environment variables.

## [0.0.0.16] - 2023-02-17

- Updated genesis to include accounts for all the validators that we can use in LocalNet based on the pre-generated keys in `build/localnet/manifests/private-keys.yaml`
- Updated `docker-compose` to name the deployment as `pocket-v1` instead of `deployments` (default is the containing folder name)
- Introduced the `cluster-manager`, which is a standalone microservice in the K8S LocalNet that takes care of (for now) automatically staking/unstaking nodes that are added/removed from the deployment
- Updated manifests and K8S resources to reflect the new `cluster-manager` addition
- In K8S LocalNet, the `cli-client` now waits for `v1-validator001` since its required for address book sourcing
- Added labels in `Tiltfile` to group resources

## [0.0.0.15] - 2023-02-17

- Added manifests to handle `Roles`, `RoleBindings` and `ServiceAccounts` and referenced them in the `Tiltfile`
- Updated `cli-client.yaml` to bind the `debug-client-account` `ServiceAccount` that has permissions to read the private keys from the `Secret`

## [0.0.0.14] - 2023-02-09

- Updated all `config*.json` files with new `server_mode_enabled` field (for state sync)

## [0.0.0.13] - 2023-02-08

- Fix bug related to installing Tilt in the Docker containers

## [0.0.0.12] - 2023-02-07

- Code formatting by VSCode

## [0.0.0.11] - 2023-02-07

- Added GITHUB_WIKI tags where it was missing

## [0.0.0.10] - 2023-02-06

- Added `genesis_localhost.json`, a copy of `genesis.json` to be used by the localhost instead of a debug docker container

## [0.0.0.9] - 2023-02-06

- Address legacy linter errors from `golangci-lint`

## [0.0.0.8] - 2023-02-06

- Added LocalNet on Kubernetes with tilt.dev

## [0.0.0.7] - 2023-02-04

- Added `--decoration="none"` flag to `reflex`

## [0.0.0.6] - 2023-01-23

- Added pprof feature flag guideline in docker-compose.yml

## [0.0.0.5] - 2023-01-20

- Update the docker-compose and relevant files to automatically load `pgadmin` server configs by binding the appropriate configs

## [0.0.0.4] - 2023-01-14

- Added `max_conns_count`, `min_conns_count`, `max_conn_lifetime`, `max_conn_idle_time` and `health_check_period` to config files

## [0.0.0.3] - 2023-01-11

- Reorder private keys so addresses (retrieved by transforming private keys) to reflect the numbering in LocalNet appropriately. The address for val1, based on config1, will have the lexicographically first address. This makes debugging easier.

## [0.0.0.2] - 2023-01-10

- Removed `BaseConfig` from `configs`
- Centralized `PersistenceGenesisState` and `ConsensusGenesisState` into `GenesisState`
- Removed `is_client_only` since it's set programmatically in the CLI

## [0.0.0.1] - 2022-12-29

- Updated all `config*.json` files with the missing `max_mempool_count` value
- Added `is_client_only` to `config1.json` so Viper knows it can be overridden. The config override is done in the Makefile's `client_connect` target. Setting this can be avoided if we merge the changes in https://github.com/pokt-network/pocket/compare/main...issue/cli-viper-environment-vars-fix

## [0.0.0.0] - 2022-12-22

- Introduced this `CHANGELOG.md`

<!-- GITHUB_WIKI: changelog/build -->
