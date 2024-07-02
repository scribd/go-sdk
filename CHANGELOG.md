# [1.36.0](https://github.com/scribd/go-sdk/compare/v1.35.0...v1.36.0) (2024-07-02)


### Features

* Add MySQL driver tracing ([fc38460](https://github.com/scribd/go-sdk/commit/fc38460bc5346162cffd3533d2e64a6d3cc2b4b8))
* Upgrade franz-go to v1.17.0 ([02cc344](https://github.com/scribd/go-sdk/commit/02cc344bd9468e1bf39e90e560f58ee6266ca153))

# [1.35.0](https://github.com/scribd/go-sdk/compare/v1.34.0...v1.35.0) (2024-06-20)


### Features

* Add sentry recovery to kafka consumer ([208a634](https://github.com/scribd/go-sdk/commit/208a634a91a62ea196e1ab13c850e5c7ad216f40))

# [1.34.0](https://github.com/scribd/go-sdk/compare/v1.33.0...v1.34.0) (2024-06-19)


### Features

* Add recovery middleware for GRPC server ([42c0254](https://github.com/scribd/go-sdk/commit/42c02541d18e97b0eac54e57e37825807d6ff3d9))

# [1.33.0](https://github.com/scribd/go-sdk/compare/v1.32.0...v1.33.0) (2024-06-14)


### Features

* Replacing custom database ddtrace with dd-trace-go library ([47a2e39](https://github.com/scribd/go-sdk/commit/47a2e39e62a5ddf1c4da5b5c1873137e73865885))

# [1.32.0](https://github.com/scribd/go-sdk/compare/v1.31.0...v1.32.0) (2024-06-13)


### Features

* Add API to create AWS services based on the Go SDK configuration ([f237347](https://github.com/scribd/go-sdk/commit/f23734734405eaba5af4aa0257ab507dabcac24f))
* Add configuration for AWS services ([326ebb5](https://github.com/scribd/go-sdk/commit/326ebb5a3e8baa7e35f9e1f751407af3c85be7cb))
* Expose AWS configuration via the main configuration entry point ([de0e023](https://github.com/scribd/go-sdk/commit/de0e02372d62082b3f61156d944711b015488e61))

# [1.31.0](https://github.com/scribd/go-sdk/compare/v1.30.0...v1.31.0) (2024-04-24)


### Features

* Create instrumented AWS v2 client ([d74b563](https://github.com/scribd/go-sdk/commit/d74b5633932eba62fbf1930d2bea50ea98ae56c1))
* Use AWS version 2 to get the AWS credentials for AWS MSK ([a8aebc8](https://github.com/scribd/go-sdk/commit/a8aebc86491ad5ba768e7141b452a7ba6eff2f6c))

# [1.30.0](https://github.com/scribd/go-sdk/compare/v1.29.0...v1.30.0) (2024-04-23)


### Features

* Configure redis cache ([b2d9376](https://github.com/scribd/go-sdk/commit/b2d93769c9d5df0ef6d56fcb1683d8b569a93727))
* Create instrumented redis client ([ca9c96e](https://github.com/scribd/go-sdk/commit/ca9c96e4e342745078431966f72a68445d4978b8))
* Propagate cache configration to the root data type ([07f0a59](https://github.com/scribd/go-sdk/commit/07f0a59f26d2109f49742cbd63f48763d944da23))
* Upgrade dd-trace-go and dependant libraries ([5130deb](https://github.com/scribd/go-sdk/commit/5130deb85eb0b124526ab4ff27ea421135f32202))

# [1.29.0](https://github.com/scribd/go-sdk/compare/v1.28.0...v1.29.0) (2024-04-17)


### Features

* Add go-kit Kafka transport ([5afd272](https://github.com/scribd/go-sdk/commit/5afd27290aec1729581b88342c9b5eba1444769b))
* Add PubSub Kafka client ([eba9a3f](https://github.com/scribd/go-sdk/commit/eba9a3f53457ae54f6594ae90921a279505dbfd8))

# [1.28.0](https://github.com/scribd/go-sdk/compare/v1.27.0...v1.28.0) (2024-04-08)


### Features

* Update go directive in go.mod to 1.22 ([4f0bf88](https://github.com/scribd/go-sdk/commit/4f0bf8879859a0a9b70306ea8c3b5369133ee04a))

# [1.27.0](https://github.com/scribd/go-sdk/compare/v1.26.0...v1.27.0) (2024-02-08)


### Features

* Upgrade google/uuid to v1.6.0 ([981427e](https://github.com/scribd/go-sdk/commit/981427e21895dbcbfa3c7ae2514c1aaecb81e03c))
* Upgrafe mage to v1.15.0 ([735df53](https://github.com/scribd/go-sdk/commit/735df5378263af6e3082e147e883c8a3ecf7d1b6))

# [1.26.0](https://github.com/scribd/go-sdk/compare/v1.25.0...v1.26.0) (2024-01-16)


### Features

* Upgrade grpc to v1.60.1 ([bdda913](https://github.com/scribd/go-sdk/commit/bdda913b54f4a48385201df715a14b4ebb057782))

# [1.25.0](https://github.com/scribd/go-sdk/compare/v1.24.0...v1.25.0) (2023-09-18)


### Features

* Update go directive in go.mod to 1.21 ([1014230](https://github.com/scribd/go-sdk/commit/1014230249f8d1f5fd116959cb8b6426be590b0d))

# [1.24.0](https://github.com/scribd/go-sdk/compare/v1.23.0...v1.24.0) (2023-08-23)


### Features

* Add new database connection settings ([f7694fa](https://github.com/scribd/go-sdk/commit/f7694fa0e6944e56347c947d7573d86fc0717d81))

# [1.23.0](https://github.com/scribd/go-sdk/compare/v1.22.0...v1.23.0) (2023-07-11)


### Features

* Update Go to v1.20 ([597094a](https://github.com/scribd/go-sdk/commit/597094ab3c33578b4667e89eebb47c9981b49d5a))

# [1.22.0](https://github.com/scribd/go-sdk/compare/v1.21.0...v1.22.0) (2023-06-16)


### Bug Fixes

* **docker:** Install required dependencies ([53ea80c](https://github.com/scribd/go-sdk/commit/53ea80c8e16649f7232bc40df5d6fbb26a622157))
* **docker:** Use alpine base image ([1f9a1f4](https://github.com/scribd/go-sdk/commit/1f9a1f47ded3635e32aa684bec500e198085a7aa))
* Upgrade golang image to 1.19.9 ([75a8289](https://github.com/scribd/go-sdk/commit/75a82894ccdf05a5e81c56c2625be54c6e4f11e7))


### Features

* **ci:** Upgrade build-push-action to v3 ([e1beb31](https://github.com/scribd/go-sdk/commit/e1beb31ed23a0b487f9073ea809faa7337be7994))
* **ci:** Upgrade cache to v3 ([fa3d6e0](https://github.com/scribd/go-sdk/commit/fa3d6e029e4bc968745b6dd7c914467e2b0fe093))
* **ci:** Upgrade checkout to v3 ([e501d45](https://github.com/scribd/go-sdk/commit/e501d4537e62642e8f56c40f95f6d0ee3571765c))
* **ci:** Upgrade commit-message-checker to v2 ([ec79084](https://github.com/scribd/go-sdk/commit/ec79084402b53954cf112c4da2c7095a6da2152d))
* **ci:** Upgrade configure-aws-credentials to v2 ([cd8cdf6](https://github.com/scribd/go-sdk/commit/cd8cdf6c8619807c2cbb21c004cceae81a51324d))
* **ci:** Upgrade job-notification to v1.1.0 ([7a2f540](https://github.com/scribd/go-sdk/commit/7a2f5406daaca3f632c05aee8d0f7cd22f5bd5e3))
* **ci:** Upgrade setup-buildx-action to v2 ([121b340](https://github.com/scribd/go-sdk/commit/121b34093dd4d632b2693a64f685ae95e0c1238b))
* **ci:** Use job-notification action for release notification ([cc8e269](https://github.com/scribd/go-sdk/commit/cc8e269dd0576cc4781db4d2406feeb0c268f62a))

# [1.21.0](https://github.com/scribd/go-sdk/compare/v1.20.1...v1.21.0) (2023-03-02)


### Features

* Add instrumentation wrapper for kgo.FetchPartition data type ([747dbb0](https://github.com/scribd/go-sdk/commit/747dbb0703899de8c01307cdad455cb6214d52af))
* Configure Kafka subscriber rebalance strategy ([7a05db0](https://github.com/scribd/go-sdk/commit/7a05db015fff4a2639f314362cfebfbd6814abd7))
* Configure max_records for Kafka subscriber ([36d6e73](https://github.com/scribd/go-sdk/commit/36d6e73cc1408ec8a4c1b86aabbdc2adc8104a10))
* Configure number of workers for Kafka consumers ([25af545](https://github.com/scribd/go-sdk/commit/25af54578362bada61c1ca53d83b77e9084511b9))
* Upgrade franz-go to v1.12.1 ([c5a551e](https://github.com/scribd/go-sdk/commit/c5a551ea1dc5e8c6f693b6811e5772228c55e711))

## [1.20.1](https://github.com/scribd/go-sdk/compare/v1.20.0...v1.20.1) (2023-02-22)


### Bug Fixes

* Fix passing data using context in gorm ([d8fc5bd](https://github.com/scribd/go-sdk/commit/d8fc5bd6e93e2e25309b86c8b2a0a72de02ce63d))

# [1.20.0](https://github.com/scribd/go-sdk/compare/v1.19.0...v1.20.0) (2023-02-17)


### Features

* Gorm logger only logs on trace ([19d6969](https://github.com/scribd/go-sdk/commit/19d6969f16581c58be69e6715804a72c758eef2e))
* Update gorm to v2 and replace it's drivers ([8352f54](https://github.com/scribd/go-sdk/commit/8352f54e7b13a93bc39c3c72281c919de4b14d8a))

# [1.19.0](https://github.com/scribd/go-sdk/compare/v1.18.0...v1.19.0) (2023-02-16)


### Features

* Configure auto commit for PubSub ([82c6315](https://github.com/scribd/go-sdk/commit/82c631549bd09efb134ce1d40eb8f55538df5594))

# [1.18.0](https://github.com/scribd/go-sdk/compare/v1.17.1...v1.18.0) (2023-02-06)


### Features

* Update Go to v1.19.5 ([f932500](https://github.com/scribd/go-sdk/commit/f932500bac5c86b4aba9b5da4dec3f652ce472ca))

## [1.17.1](https://github.com/scribd/go-sdk/compare/v1.17.0...v1.17.1) (2023-01-06)


### Bug Fixes

* Finish remaining spans on Fetches iterable Done ([ca0a35b](https://github.com/scribd/go-sdk/commit/ca0a35bc4caf6433225a877946538e7d40ee9e07))

# [1.17.0](https://github.com/scribd/go-sdk/compare/v1.16.1...v1.17.0) (2022-08-08)


### Features

* Set module go version to 1.18 ([6625ccd](https://github.com/scribd/go-sdk/commit/6625ccd65325abceb225b67af5c725af744fae9a))
* Upgrade go version to 1.18.5 in Dockerfile ([1ab8b17](https://github.com/scribd/go-sdk/commit/1ab8b176e98cba1b9a1a0760e17f2bd5280941ba))
* Upgrade golangci-lint to v1.47.3 ([c6ffdd4](https://github.com/scribd/go-sdk/commit/c6ffdd4a8ac1cce140af854ade04992f4127b675))

## [1.16.1](https://github.com/scribd/go-sdk/compare/v1.16.0...v1.16.1) (2022-08-03)


### Bug Fixes

* Fire sentry event fix ([451595c](https://github.com/scribd/go-sdk/commit/451595c5602e24c4be1a8dc01d69ac33dbf76cd3))

## [1.15.4](https://github.com/scribd/go-sdk/compare/v1.15.3...v1.15.4) (2022-07-28)


### Bug Fixes

* Runtime metrics fix ([5758eb7](https://github.com/scribd/go-sdk/commit/5758eb75dc0ad59ae27ac809412d3885cb9f4fa9))

# [1.15.0](https://github.com/scribd/go-sdk/compare/v1.14.0...v1.15.0) (2022-07-06)


### Features

* Pass down service version to the DataDog ([fbe0f6c](https://github.com/scribd/go-sdk/commit/fbe0f6c4a7030062ca4920a53ecd446bed9b515b))
* Upgrade dd-trace-go package version to v1.39.0 ([ef58946](https://github.com/scribd/go-sdk/commit/ef58946ed66388e5b7bf1174fa7962c72ffc39c8))

# [1.14.0](https://github.com/scribd/go-sdk/compare/v1.13.0...v1.14.0) (2022-07-01)


### Features

* Add assumable role and session name configuration settings to AWS MSK IAM SASL ([ef08bd1](https://github.com/scribd/go-sdk/commit/ef08bd15238ad0d9bdc81294ef45994a93410789))

# [1.12.0](https://github.com/scribd/go-sdk/compare/v1.11.0...v1.12.0) (2022-06-22)


### Features

* Add franz-go plugin to log kafka client using go-sdk logger ([ba9000e](https://github.com/scribd/go-sdk/commit/ba9000e101561cbaf1d96647e531f3ddc294178c))
* Add franz-go plugin to publish kafka related metrics ([9ad6f12](https://github.com/scribd/go-sdk/commit/9ad6f1289d9be8399d066975f8f07183a9298fbb))
* Configure PubSub metrics ([d7116b5](https://github.com/scribd/go-sdk/commit/d7116b5597712a23da00804da69e13b0cd46c838))

# [1.11.0](https://github.com/scribd/go-sdk/compare/v1.10.1...v1.11.0) (2022-06-10)


### Features

* Extend PubSub configuration to include TLS and SASL auth ([f723b44](https://github.com/scribd/go-sdk/commit/f723b447cfadbb2e0c06ddc192776c2eea7894d4))
* Trace producers and consumers of the kafka client ([2d4ddf8](https://github.com/scribd/go-sdk/commit/2d4ddf84ecb34f5d8aa960aef6983d74e581c065))

## [1.10.1](https://github.com/scribd/go-sdk/compare/v1.10.0...v1.10.1) (2022-04-14)


### Bug Fixes

* Convert broker urls to string slice explicitly ([0c0713e](https://github.com/scribd/go-sdk/commit/0c0713edad89c53159e0c923c3225b26754fac77))

# [1.10.0](https://github.com/scribd/go-sdk/compare/v1.9.0...v1.10.0) (2022-03-24)


### Features

* Add PubSub configuration ([1576521](https://github.com/scribd/go-sdk/commit/15765218470b17face9b2917aab7c29048b88557))
* Expose PubSub config via main config entrypoint ([9dbb14c](https://github.com/scribd/go-sdk/commit/9dbb14cca2ff7e9da9a669b191d42ec35eccc95c))

# [1.9.0](https://github.com/scribd/go-sdk/compare/v1.8.0...v1.9.0) (2022-02-21)


### Features

* Provide a way to configure profiling code hostpots ([79b9a45](https://github.com/scribd/go-sdk/commit/79b9a45949cd5b3e8a1226abcea0d19b67a0e4a2))
* Upgrade dd-trace-go to v1.36.0 version ([11370a1](https://github.com/scribd/go-sdk/commit/11370a1c8659d41c1447b8709db9c4f3702e65bb))

# [1.8.0](https://github.com/scribd/go-sdk/compare/v1.7.1...v1.8.0) (2022-01-24)


### Features

* Support //go:build lines together with // +build ([4cbdd4d](https://github.com/scribd/go-sdk/commit/4cbdd4df028266015d456ebf2c72186c135f5f60))
* **Dockerfile:** Bump to Go 1.17.6 ([4080b1d](https://github.com/scribd/go-sdk/commit/4080b1df82bca1161b02c32d1f7cff8bced4aeff))
* **go:** Bump to Go 1.17 ([8d5cc40](https://github.com/scribd/go-sdk/commit/8d5cc404f683e18d13fa821f7a06b8e776e252e2))

# [1.7.1](https://github.com/scribd/go-sdk/compare/v1.7.0...v1.7.1) (2022-01-10)


### Security Fixes

* Upgrade sentry-go module to version v0.12.0 ([4bbe60c](https://github.com/scribd/go-sdk/commit/4bbe60ca19d45e607bb77acffd95ce545819cf2c))
* Upgrade viper module to version  v1.10.1 ([2ab5fb5](https://github.com/scribd/go-sdk/commit/2ab5fb5486176b90556426f94a4c8ab38da64b61))

# [1.7.0](https://github.com/scribd/go-sdk/compare/v1.6.0...v1.7.0) (2021-12-16)


### Bug Fixes

* Retrieve request id from the context in logger middleware instead of getting it from the request header ([0d74e33](https://github.com/scribd/go-sdk/commit/0d74e33541ddc8c950f32698189e3712df71ad8a))


### Features

* Add context helpers for request id ([7997f79](https://github.com/scribd/go-sdk/commit/7997f79ecc440bd88854522f35302fba3aac9347))
* Add request ID gRPC interceptors ([5e73a20](https://github.com/scribd/go-sdk/commit/5e73a20ce2acbf979d42af8d66de94df01d7ab0c))
* Pass request id in logger interceptor ([46c8596](https://github.com/scribd/go-sdk/commit/46c859667c3723736d6608b443de3fd5d4a3b5a9))
* Use pkg/context/requestid for request_id middleware ([20e2349](https://github.com/scribd/go-sdk/commit/20e234932c9a5e1400fb9592064f6c59aa734dac))

# [1.6.0](https://github.com/scribd/go-sdk/compare/v1.5.0...v1.6.0) (2021-12-14)


### Bug Fixes

* Create new gorm db instance in database_logger middleware when set logger settings to avoid data races ([5380a87](https://github.com/scribd/go-sdk/commit/5380a8726210fb8908399dc430684b0f10973b84))
* Update database gRPC interceptors to include DB tracing (parity with database middleware) ([616c45d](https://github.com/scribd/go-sdk/commit/616c45dfb8a86f6dc2ad7ec435a5291150716d2c))


### Features

* Add database logging gRPC server interceptors ([ef8ea53](https://github.com/scribd/go-sdk/commit/ef8ea532bbd02f9a9a62c28f96e46132ea0aa93c))
* Add mage command to generate test proto files ([6e2ffc5](https://github.com/scribd/go-sdk/commit/6e2ffc556a5a5b8fd37c517bd8dd18b25567e41c))
* Create test grpc service ([b84d88e](https://github.com/scribd/go-sdk/commit/b84d88eb90d8b4b9a77866c8ce5874ba84a178d1))
* Move gorm logger to logger package ([8fe5b26](https://github.com/scribd/go-sdk/commit/8fe5b26e6344e21f8607c83eb21428a939ce8d2f))
* Upgrade grpc to v1.32.0 ([3c9127d](https://github.com/scribd/go-sdk/commit/3c9127d8e4606c90aa71a2e39967d4102e355bba))

# [1.5.0](https://github.com/scribd/go-sdk/compare/v1.4.1...v1.5.0) (2021-12-06)


### Features

* Add server and client gRPC trace interceptors ([ac66d1b](https://github.com/scribd/go-sdk/commit/ac66d1b2ddda9ff6a43ce20314cd5cbc76b3756d))
* Correlate traces and logs for gRPC ([e2f7f4c](https://github.com/scribd/go-sdk/commit/e2f7f4c2546e9568014b533043735881671ae101))

## [1.4.1](https://github.com/scribd/go-sdk/compare/v1.4.0...v1.4.1) (2021-11-18)


### Bug Fixes

* Downgrade cors library to v1.7.0 ([7faef78](https://github.com/scribd/go-sdk/commit/7faef78108a2a4a93ee01d19a0a114b3b26330a8))

# [1.4.0](https://github.com/scribd/go-sdk/compare/v1.3.0...v1.4.0) (2021-11-17)


### Features

* Update modules dependencies ([e4eaa07](https://github.com/scribd/go-sdk/commit/e4eaa07f60018bd118fd236925708ed068d98c27))

# [1.3.0](https://github.com/scribd/go-sdk/compare/v1.2.3...v1.3.0) (2021-11-01)


### Features

* Add CORS middleware ([74a61f7](https://github.com/scribd/go-sdk/commit/74a61f7a687656d6953b4789d5a34fac54828f5c))
* Add CORS settings to server ([4789fc7](https://github.com/scribd/go-sdk/commit/4789fc7096fdbd710be241971ea8836e86407f15))

## [1.2.3](https://github.com/scribd/go-sdk/compare/v1.2.2...v1.2.3) (2021-11-01)


### Bug Fixes

* **ci:** Do not persist Github token when checking out in Release workflow ([b96c6b0](https://github.com/scribd/go-sdk/commit/b96c6b05249dc451fbe2abb773f32e2cf379974f))

## [1.2.2](https://github.com/scribd/go-sdk/compare/v1.2.1...v1.2.2) (2021-11-01)


### Bug Fixes

* **ci:** Use seperate PAT for release workflow ([7b019a9](https://github.com/scribd/go-sdk/commit/7b019a9dac1cd00a406cfbe1f3ca843e9b05273f))
* **CI:** Use automatic token authentication via GITHUB_TOKEN in actions ([41fd4cf](https://github.com/scribd/go-sdk/commit/41fd4cfc175355726a9df013bda61207049b19ed))

## [1.2.1](https://github.com/scribd/go-sdk/compare/v1.2.0...v1.2.1) (2021-06-09)


### Bug Fixes

* **ci:** Remove CI workflow Slack notification ([7501035](https://github.com/scribd/go-sdk/commit/7501035fe0d5e9fd9663a9691a2d1d2d7435efd0))

# [1.2.0](https://github.com/scribd/go-sdk/compare/v1.1.0...v1.2.0) (2021-05-11)


### Bug Fixes

* Update pkg/version accordingly when releasing new version ([8772dfb](https://github.com/scribd/go-sdk/commit/8772dfbda6791e23945813cbbc814f6af706fe26))


### Features

* **ci:** Add a linter stage in Dockerfile ([cc3ee11](https://github.com/scribd/go-sdk/commit/cc3ee1140711253c8f60c3cbeb3df704c426dccd))
* **ci:** Add commit message check workflow ([b1476d6](https://github.com/scribd/go-sdk/commit/b1476d69446715d9b70703c6ba1711a7d1436e43))
* **ci:** Add datadog metrics for CI pipeline ([73b16a8](https://github.com/scribd/go-sdk/commit/73b16a8c2e642afc4a1aa5d53634df67864ed450))
* **ci:** Add datadog metrics for release pipeline ([7b8e195](https://github.com/scribd/go-sdk/commit/7b8e195d871eff1c82dbaf08cbb906cafd14bd09))
* **ci:** Add velocity metrics for the repo ([f918947](https://github.com/scribd/go-sdk/commit/f918947321aeb8ecd338caa9f7e2ba6acc0b4f90))
* **ci:** Build the development stage from the linter stage ([2b62f8d](https://github.com/scribd/go-sdk/commit/2b62f8d40439ddbb3dd5dbcc29f068c85f204bf0))
* **ci:** Build the SDK before the linter stage in Dockerfile ([c04bdf0](https://github.com/scribd/go-sdk/commit/c04bdf0a4bdd660af4d939927b8877964e7e093e))
* **ci:** Bump golangci-lint to v1.39.0 ([838b1a4](https://github.com/scribd/go-sdk/commit/838b1a406e32b5faa474e277782db80016933181))
* **ci:** Remove the unnecessary build ([9cd7436](https://github.com/scribd/go-sdk/commit/9cd7436240aa7dd66384c9dd3f9ed26235ffc8d1))
* **ci:** Remove unsupported event from commit message check ([e33a978](https://github.com/scribd/go-sdk/commit/e33a978007e5ff0493b9199f79205c9228ff25bd))
* **ci:** Run the fmt checks at the linter stage ([f67348f](https://github.com/scribd/go-sdk/commit/f67348f78596cfcac86d759cdee1934e7b2f710b))
* **ci:** Use mage to run golangci-lint-action in Docker ([8bc71aa](https://github.com/scribd/go-sdk/commit/8bc71aa41cfe5a10c7b388cbfe35d9429c39afd6))
* **Dockerfile:** Bump to Go 1.16.4 ([fdc116e](https://github.com/scribd/go-sdk/commit/fdc116ee138df96f4fc6ff25904de19bdd5d3db9))
* **go:** Bump to Go 1.16.4 ([038076b](https://github.com/scribd/go-sdk/commit/038076b42266ca3a79a1bddd7e98d5be54d493b0))

# [1.2.0](https://github.com/scribd/go-sdk/compare/v1.1.0...v1.2.0) (2021-05-11)


### Bug Fixes

* Update pkg/version accordingly when releasing new version ([8772dfb](https://github.com/scribd/go-sdk/commit/8772dfbda6791e23945813cbbc814f6af706fe26))


### Features

* **ci:** Add a linter stage in Dockerfile ([cc3ee11](https://github.com/scribd/go-sdk/commit/cc3ee1140711253c8f60c3cbeb3df704c426dccd))
* **ci:** Add commit message check workflow ([b1476d6](https://github.com/scribd/go-sdk/commit/b1476d69446715d9b70703c6ba1711a7d1436e43))
* **ci:** Add datadog metrics for CI pipeline ([73b16a8](https://github.com/scribd/go-sdk/commit/73b16a8c2e642afc4a1aa5d53634df67864ed450))
* **ci:** Add datadog metrics for release pipeline ([7b8e195](https://github.com/scribd/go-sdk/commit/7b8e195d871eff1c82dbaf08cbb906cafd14bd09))
* **ci:** Add velocity metrics for the repo ([f918947](https://github.com/scribd/go-sdk/commit/f918947321aeb8ecd338caa9f7e2ba6acc0b4f90))
* **ci:** Build the development stage from the linter stage ([2b62f8d](https://github.com/scribd/go-sdk/commit/2b62f8d40439ddbb3dd5dbcc29f068c85f204bf0))
* **ci:** Build the SDK before the linter stage in Dockerfile ([c04bdf0](https://github.com/scribd/go-sdk/commit/c04bdf0a4bdd660af4d939927b8877964e7e093e))
* **ci:** Bump golangci-lint to v1.39.0 ([838b1a4](https://github.com/scribd/go-sdk/commit/838b1a406e32b5faa474e277782db80016933181))
* **ci:** Remove the unnecessary build ([9cd7436](https://github.com/scribd/go-sdk/commit/9cd7436240aa7dd66384c9dd3f9ed26235ffc8d1))
* **ci:** Remove unsupported event from commit message check ([e33a978](https://github.com/scribd/go-sdk/commit/e33a978007e5ff0493b9199f79205c9228ff25bd))
* **ci:** Run the fmt checks at the linter stage ([f67348f](https://github.com/scribd/go-sdk/commit/f67348f78596cfcac86d759cdee1934e7b2f710b))
* **ci:** Use mage to run golangci-lint-action in Docker ([8bc71aa](https://github.com/scribd/go-sdk/commit/8bc71aa41cfe5a10c7b388cbfe35d9429c39afd6))
* **Dockerfile:** Bump to Go 1.16.4 ([fdc116e](https://github.com/scribd/go-sdk/commit/fdc116ee138df96f4fc6ff25904de19bdd5d3db9))
* **go:** Bump to Go 1.16.4 ([038076b](https://github.com/scribd/go-sdk/commit/038076b42266ca3a79a1bddd7e98d5be54d493b0))

# [1.1.0](https://github.com/scribd/go-sdk/compare/v1.0.0...v1.1.0) (2021-03-02)


### Bug Fixes

* **ci:** Fix CI pipeline by using builder Docker layer as a target and download golangci-lint manually ([dcc1143](https://github.com/scribd/go-sdk/commit/dcc1143d884b37dfb5dde1b34bdfb8936c5100b6))
* **go:** Clean up go.sum ([ed7348a](https://github.com/scribd/go-sdk/commit/ed7348a37b13d395121b551b93973ceff18ea907))
* Add the missing -app suffix to the service name in metrics ([d6f2f0e](https://github.com/scribd/go-sdk/commit/d6f2f0ecfec3ff97e736f00bb65c07034867226c))


### Features

* **ci:** Add the CI Action pipeline ([d1ab240](https://github.com/scribd/go-sdk/commit/d1ab240aff73a328e765e44de007ad91d35158dd))
* **ci:** Add the Release Action pipeline ([b1c3222](https://github.com/scribd/go-sdk/commit/b1c32221ca7e10400a5533359dbff449fd027902))
* **ci:** Cache docker layers ([2329506](https://github.com/scribd/go-sdk/commit/2329506abe0749ff5289c663e0b8e471d83e5db0))

# [1.0.0](https://github.com/scribd/go-sdk/compare/v0.8.0...v1.0.0) (2021-02-04)


### Code Refactoring

* Rename module name to github.com/scribd/go-sdk ([923a8ae](https://github.com/scribd/go-sdk/commit/923a8ae9f8b3f38734ec5d737956ba8fb59cf772))


### BREAKING CHANGES

* Rename module name to github.com/scribd/go-sdk as we are migrating from gitlab to github.

# [0.8.0](https://github.com/scribd/go-sdk/compare/v0.7.0...v0.8.0) (2020-12-16)


### Features

* Bump modules to Go 1.15.6 ([7105080](https://github.com/scribd/go-sdk/commit/7105080f36960e23ab61503b22b2f7a737ed2fcb))
* Bump to Go 1.15.6 ([ac25391](https://github.com/scribd/go-sdk/commit/ac253911f530f2707c4a186d5f97046349378822))

# [0.7.0](https://github.com/scribd/go-sdk/compare/v0.6.0...v0.7.0) (2020-09-10)


### Features

* Inject gorm.DB in the context using the context/database package ([66ef7ba](https://github.com/scribd/go-sdk/commit/66ef7baf99cc0338c817fabe283dfc75456ddfa9))
* Introduce the gRPC database interceptors ([cc73f8c](https://github.com/scribd/go-sdk/commit/cc73f8c6884822c123e6398758bd2657ab132d9d))
* Remove the Database context key as now managed by the context/database package ([2984490](https://github.com/scribd/go-sdk/commit/29844909f6ef5ab418fd18abe8fc95e7b89366bc))

# [0.6.0](https://github.com/scribd/go-sdk/compare/v0.5.0...v0.6.0) (2020-09-10)


### Features

* Inject the Metrics in the context using the context/metrics package ([fd03e2b](https://github.com/scribd/go-sdk/commit/fd03e2bff9f9d974fa1c14fda76fdb6777a19eb0))
* Introduce the gRPC metrics interceptors ([c84f618](https://github.com/scribd/go-sdk/commit/c84f6185bb6b8c2c14dc9e64b9daa94fbe2e18ad))
* Remove the Metrics context key as now managed by the context/metrics package ([80a1d68](https://github.com/scribd/go-sdk/commit/80a1d686595e6fa004c50f9ccacb61c21055495d))

# [0.5.0](https://github.com/scribd/go-sdk/compare/v0.4.0...v0.5.0) (2020-09-09)


### Features

* Add package context/logger to manage the logger in the context ([ecba066](https://github.com/scribd/go-sdk/commit/ecba06656e3ad8e1c08a2632ec1ffe90a3242d6c))
* Add the gRPC dependencies ([2a84eef](https://github.com/scribd/go-sdk/commit/2a84eeff3ec9a46b5fa37fd5503340193dbcc61b))
* Inject the Logger in the context using the sdklogcontext package ([7c8db72](https://github.com/scribd/go-sdk/commit/7c8db72ad78afc207012a7c9b03a1da1044b810c))
* Introduce the gRPC logging interceptors ([f0369e8](https://github.com/scribd/go-sdk/commit/f0369e81aba0845040edfad041b6d42a62e0d5d2))
* Remove the Logger context key as now managed by the context/logger package ([d93eb99](https://github.com/scribd/go-sdk/commit/d93eb99c2abd1dab12dfdba4fce06ca87f3c2bf5))

# [0.4.0](https://github.com/scribd/go-sdk/compare/v0.3.0...v0.4.0) (2020-08-07)


### Features

* datadog-go is now a direct dependency ([0b5849c](https://github.com/scribd/go-sdk/commit/0b5849c8088c726b2df1742dfff5ba6b3046a0d5))
* Introduce the metrics package ([f352e72](https://github.com/scribd/go-sdk/commit/f352e722950b22b40fa59ef9567f753bac717f77))
* Move the `Metrics` Context key in the package `contextkeys` ([281d91c](https://github.com/scribd/go-sdk/commit/281d91c1f16bfa3a8bc3eabee610710003f0b06b))

# [0.3.0](https://github.com/scribd/go-sdk/compare/v0.2.0...v0.3.0) (2020-07-27)


### Features

* Add DataDog profiler ([9452058](https://github.com/scribd/go-sdk/commit/94520581911c3d69205f212a7e9af0615ae4b32b))

# [0.2.0](https://github.com/scribd/go-sdk/compare/v0.1.0...v0.2.0) (2020-07-02)


### Features

* Environment not directly configurable but read from APP_ENV ([f0e6b99](https://github.com/scribd/go-sdk/commit/f0e6b999a5bb830ee8cee9572fb0aa01de0d54df))
* Extend the Tracking configuration to expose additional options ([e85d10a](https://github.com/scribd/go-sdk/commit/e85d10ac22845e26de684f5de926acecb7bd2333))
* Implement a "native" Sentry hook using the official library ([62def35](https://github.com/scribd/go-sdk/commit/62def354b6ff29fd5fca3bab6041d4f4c6572caf))
* Read APP_VERSION and APP_SERVER_NAME from ENV instead config ([8b9bf5d](https://github.com/scribd/go-sdk/commit/8b9bf5d22b62a207bf79695d40ce0790005ced43))
* Remove the unsupported Timeout option from the Sentry configuration ([21cc31a](https://github.com/scribd/go-sdk/commit/21cc31ab25af512b42ba94a6f8ce1334174f739e))
* Replace LogrusSentry with the official SentriGo (go.mod) ([f63df74](https://github.com/scribd/go-sdk/commit/f63df74988fc602ac6880e3a21633404904eddda))

# [0.1.0](https://github.com/scribd/go-sdk/compare/v0.0.1...v0.1.0) (2020-05-28)


### Bug Fixes

* Fix entrypoint for release container ([e2f14cc](https://github.com/scribd/go-sdk/commit/e2f14cc2e8501b236ab849ecf192a0f5c738eceb))


### Features

* Add .releaserc.yml ([5003f35](https://github.com/scribd/go-sdk/commit/5003f35f5adcbab2b8f56c1fe44b55eb1d399238))
* Add release stage ([337b09e](https://github.com/scribd/go-sdk/commit/337b09e8de38d67700b9e8b8e3b50aa683fba875))
* Update version in version package during release ([2045421](https://github.com/scribd/go-sdk/commit/204542103478087b547e52b0da88f3e1748e5abe))

# CHANGELOG

<!--- next entry here -->

## 0.0.1
2019-10-29

### Features

- Add version package (2b331d73a8ae50ddbbc525b04220c9e24b207f45)
