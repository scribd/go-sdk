# [0.6.0](https://git.lo/microservices/sdk/go-sdk/compare/v0.5.0...v0.6.0) (2020-09-10)


### Features

* Inject the Metrics in the context using the context/metrics package ([fd03e2b](https://git.lo/microservices/sdk/go-sdk/commit/fd03e2bff9f9d974fa1c14fda76fdb6777a19eb0))
* Introduce the gRPC metrics interceptors ([c84f618](https://git.lo/microservices/sdk/go-sdk/commit/c84f6185bb6b8c2c14dc9e64b9daa94fbe2e18ad))
* Remove the Metrics context key as now managed by the context/metrics package ([80a1d68](https://git.lo/microservices/sdk/go-sdk/commit/80a1d686595e6fa004c50f9ccacb61c21055495d))

# [0.5.0](https://git.lo/microservices/sdk/go-sdk/compare/v0.4.0...v0.5.0) (2020-09-09)


### Features

* Add package context/logger to manage the logger in the context ([ecba066](https://git.lo/microservices/sdk/go-sdk/commit/ecba06656e3ad8e1c08a2632ec1ffe90a3242d6c))
* Add the gRPC dependencies ([2a84eef](https://git.lo/microservices/sdk/go-sdk/commit/2a84eeff3ec9a46b5fa37fd5503340193dbcc61b))
* Inject the Logger in the context using the sdklogcontext package ([7c8db72](https://git.lo/microservices/sdk/go-sdk/commit/7c8db72ad78afc207012a7c9b03a1da1044b810c))
* Introduce the gRPC logging interceptors ([f0369e8](https://git.lo/microservices/sdk/go-sdk/commit/f0369e81aba0845040edfad041b6d42a62e0d5d2))
* Remove the Logger context key as now managed by the context/logger package ([d93eb99](https://git.lo/microservices/sdk/go-sdk/commit/d93eb99c2abd1dab12dfdba4fce06ca87f3c2bf5))

# [0.4.0](https://git.lo/microservices/sdk/go-sdk/compare/v0.3.0...v0.4.0) (2020-08-07)


### Features

* datadog-go is now a direct dependency ([0b5849c](https://git.lo/microservices/sdk/go-sdk/commit/0b5849c8088c726b2df1742dfff5ba6b3046a0d5))
* Introduce the metrics package ([f352e72](https://git.lo/microservices/sdk/go-sdk/commit/f352e722950b22b40fa59ef9567f753bac717f77))
* Move the `Metrics` Context key in the package `contextkeys` ([281d91c](https://git.lo/microservices/sdk/go-sdk/commit/281d91c1f16bfa3a8bc3eabee610710003f0b06b))

# [0.3.0](https://git.lo/microservices/sdk/go-sdk/compare/v0.2.0...v0.3.0) (2020-07-27)


### Features

* Add DataDog profiler ([9452058](https://git.lo/microservices/sdk/go-sdk/commit/94520581911c3d69205f212a7e9af0615ae4b32b))

# [0.2.0](https://git.lo/microservices/sdk/go-sdk/compare/v0.1.0...v0.2.0) (2020-07-02)


### Features

* Environment not directly configurable but read from APP_ENV ([f0e6b99](https://git.lo/microservices/sdk/go-sdk/commit/f0e6b999a5bb830ee8cee9572fb0aa01de0d54df))
* Extend the Tracking configuration to expose additional options ([e85d10a](https://git.lo/microservices/sdk/go-sdk/commit/e85d10ac22845e26de684f5de926acecb7bd2333))
* Implement a "native" Sentry hook using the official library ([62def35](https://git.lo/microservices/sdk/go-sdk/commit/62def354b6ff29fd5fca3bab6041d4f4c6572caf))
* Read APP_VERSION and APP_SERVER_NAME from ENV instead config ([8b9bf5d](https://git.lo/microservices/sdk/go-sdk/commit/8b9bf5d22b62a207bf79695d40ce0790005ced43))
* Remove the unsupported Timeout option from the Sentry configuration ([21cc31a](https://git.lo/microservices/sdk/go-sdk/commit/21cc31ab25af512b42ba94a6f8ce1334174f739e))
* Replace LogrusSentry with the official SentriGo (go.mod) ([f63df74](https://git.lo/microservices/sdk/go-sdk/commit/f63df74988fc602ac6880e3a21633404904eddda))

# [0.1.0](https://git.lo/microservices/sdk/go-sdk/compare/v0.0.1...v0.1.0) (2020-05-28)


### Bug Fixes

* Fix entrypoint for release container ([e2f14cc](https://git.lo/microservices/sdk/go-sdk/commit/e2f14cc2e8501b236ab849ecf192a0f5c738eceb))


### Features

* Add .releaserc.yml ([5003f35](https://git.lo/microservices/sdk/go-sdk/commit/5003f35f5adcbab2b8f56c1fe44b55eb1d399238))
* Add release stage ([337b09e](https://git.lo/microservices/sdk/go-sdk/commit/337b09e8de38d67700b9e8b8e3b50aa683fba875))
* Update version in version package during release ([2045421](https://git.lo/microservices/sdk/go-sdk/commit/204542103478087b547e52b0da88f3e1748e5abe))

# CHANGELOG

<!--- next entry here -->

## 0.0.1
2019-10-29

### Features

- Add version package (2b331d73a8ae50ddbbc525b04220c9e24b207f45)
