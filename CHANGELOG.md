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
