# Changelog

## [v0.9.0](https://github.com/monkescience/reference-service-go/compare/v0.8.2...v0.9.0) (2026-04-24)

### ⚠ BREAKING CHANGES

- switch to hexagonal layout under internal/core ([4f47ae2](https://github.com/monkescience/reference-service-go/commit/4f47ae2e64495864d8477983f030f50cce2aa738))
### Features

- **observability:** wire opentelemetry tracing ([d498519](https://github.com/monkescience/reference-service-go/commit/d498519923de6f6898764c3db5a7ae4d52533fc7))
### Bug Fixes

- **deps:** update module go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp to v0.68.0 (#94) ([19d74a7](https://github.com/monkescience/reference-service-go/commit/19d74a7fd05cf49139ac5321a62586c51317ae17))
- **deps:** update module github.com/jackc/pgx/v5 to v5.9.2 [security] (#92) ([340c1b0](https://github.com/monkescience/reference-service-go/commit/340c1b0c1cb468f030d7eeb3aa05f0badf68cacb))

## [v0.8.2](https://github.com/monkescience/reference-service-go/compare/v0.8.1...v0.8.2) (2026-04-10)

### Features

- **ids:** use UUIDv7 for import and catch IDs ([a11081b](https://github.com/monkescience/reference-service-go/commit/a11081bd7804c38a19365d7837e3002c0f180350))
### Bug Fixes

- **build:** standardize service build flags ([903b3ed](https://github.com/monkescience/reference-service-go/commit/903b3ed50bae032502b8109f27a36ed3818ac449))

## [v0.8.1](https://github.com/monkescience/reference-service-go/compare/v0.8.0...v0.8.1) (2026-04-09)

### Bug Fixes

- **generate:** align oapi-codegen toolchain dependencies ([de90c10](https://github.com/monkescience/reference-service-go/commit/de90c107ad7338f96aa65f22c6ad514eb93f2f34))
- **deps:** update module github.com/oapi-codegen/runtime to v1.4.0 (#76) ([ac9012d](https://github.com/monkescience/reference-service-go/commit/ac9012d2659930984148deec33bb8f235aa8f537))
- **deps:** update module github.com/monkescience/vital to v0.4.0 (#74) ([2370a1c](https://github.com/monkescience/reference-service-go/commit/2370a1c7b622956243be2ea338d1e8a52eafcf6c))
- **deps:** update module github.com/getkin/kin-openapi to v0.135.0 (#84) ([b26b6c1](https://github.com/monkescience/reference-service-go/commit/b26b6c15a86d119abeeac8f73c61c1f008150989))
- **deps:** update module github.com/oapi-codegen/oapi-codegen/v2 to v2.6.0 (#75) ([3328793](https://github.com/monkescience/reference-service-go/commit/3328793c3ae58d6ff66dc964cc63c70c53c0f66a))
- **deps:** update module github.com/getkin/kin-openapi to v0.134.0 (#73) ([e721a36](https://github.com/monkescience/reference-service-go/commit/e721a365fd36f547f6aa219b3a35fd00bbe0eaa4))
- **chart:** separate image tags from chart metadata ([ddab73c](https://github.com/monkescience/reference-service-go/commit/ddab73c85048865aa33fb0734c852efe6575ba16))

## [v0.8.0](https://github.com/monkescience/reference-service-go/compare/0.7.0...v0.8.0) (2026-04-09)

### ⚠ BREAKING CHANGES

- **api:** consolidate reference API and persist catches ([73a934c](https://github.com/monkescience/reference-service-go/commit/73a934c1ad1be555525eefa003d02ece44a70a36))
### Features

- **api:** consolidate reference API and persist catches ([73a934c](https://github.com/monkescience/reference-service-go/commit/73a934c1ad1be555525eefa003d02ece44a70a36))
- **config:** add server port config and align test tooling ([29d16b0](https://github.com/monkescience/reference-service-go/commit/29d16b048f576a05b831f158dd6b29f0fc699129))
- decouple migrations from service startup for K8s PreSync hook ([581f8ca](https://github.com/monkescience/reference-service-go/commit/581f8caae7bb86e15dc05db1e3b43124db41c409))
- add PokeAPI import, Pokeball gacha, and PostgreSQL persistence ([896cf65](https://github.com/monkescience/reference-service-go/commit/896cf65cbd8cfd16592cc026d8ac042e0959f80a))
- add Pokemon import API with OpenAPI-first design ([3c9dda6](https://github.com/monkescience/reference-service-go/commit/3c9dda6e780a33adfc3b4594a93a3584ff59c3d6))
- **ci:** migrate from release-please to yeet ([2f87154](https://github.com/monkescience/reference-service-go/commit/2f87154627f72fa20bcd3735120b88f17baee66b))
### Bug Fixes

- **tests:** clean up postgres before TestMain exit ([30cee5e](https://github.com/monkescience/reference-service-go/commit/30cee5ebc529ab1cde5c4176712c01bd984da695))

## [0.7.0](https://github.com/monkescience/reference-service-go/compare/0.6.0...0.7.0) (2025-12-12)


### Features

* **build,ci:** add Mage build tool, linting config, and initial tests ([6965c8a](https://github.com/monkescience/reference-service-go/commit/6965c8af85aa3a39aeef8292d3439707a0ed389c))
* **order:** add country support in order ID and enhance API responses ([ba759c3](https://github.com/monkescience/reference-service-go/commit/ba759c364b78d6b7fbc2229b840ab186f5bb89cb))
* **ui:** add dark mode toggle, improve styling consistency, and optimize responsiveness ([5c8e6e2](https://github.com/monkescience/reference-service-go/commit/5c8e6e2a4945f75ee5a1df323807f70d692624c5))
* **ui:** add instance dashboard with HTMX and dynamic tile updates ([b0394a8](https://github.com/monkescience/reference-service-go/commit/b0394a8dc302b31cd5b1ba70e64fa091becc9940))

## [0.6.0](https://github.com/monkescience/reference-service-go/compare/0.5.2...0.6.0) (2025-11-15)


### Features

* **release:** add force-release trigger file ([4f3f4d7](https://github.com/monkescience/reference-service-go/commit/4f3f4d745bd38bfc56228a125cf5e53c34d50797))

## [0.5.2](https://github.com/monkescience/reference-service-go/compare/0.5.1...0.5.2) (2025-11-08)


### Bug Fixes

* **ci:** update Helm push target to use HELM_REPOSITORY ([df41bdb](https://github.com/monkescience/reference-service-go/commit/df41bdb0b88cb04a33ad3c6980ce928bb26b10de))

## [0.5.1](https://github.com/monkescience/reference-service-go/compare/0.5.0...0.5.1) (2025-11-08)


### Bug Fixes

* **deps:** update module github.com/getkin/kin-openapi to v0.133.0 ([#33](https://github.com/monkescience/reference-service-go/issues/33)) ([7750519](https://github.com/monkescience/reference-service-go/commit/7750519e5ef468b5b641ed1a38d6cc6c2466fbdf))
* **deps:** update module github.com/go-chi/chi/v5 to v5.2.3 ([#32](https://github.com/monkescience/reference-service-go/issues/32)) ([21fa9d3](https://github.com/monkescience/reference-service-go/commit/21fa9d3e95c3c2e0d45f120244bc0eb430273ac2))
* **deps:** update module github.com/oapi-codegen/oapi-codegen/v2 to v2.5.0 ([#23](https://github.com/monkescience/reference-service-go/issues/23)) ([1c84735](https://github.com/monkescience/reference-service-go/commit/1c84735ba7a65717b1f04edbef3db208c8e05e78))
* **deps:** update module github.com/oapi-codegen/oapi-codegen/v2 to v2.5.1 ([#47](https://github.com/monkescience/reference-service-go/issues/47)) ([4fa4268](https://github.com/monkescience/reference-service-go/commit/4fa42684b29b2299a1ae2b0fac32b56796369adb))
* **deps:** update module github.com/oapi-codegen/runtime to v1.1.2 ([#22](https://github.com/monkescience/reference-service-go/issues/22)) ([69f0efa](https://github.com/monkescience/reference-service-go/commit/69f0efaaea9dba08661de8077083c691c5854eaf))
* **deps:** update module github.com/prometheus/client_golang to v1.23.0 ([#25](https://github.com/monkescience/reference-service-go/issues/25)) ([e89651e](https://github.com/monkescience/reference-service-go/commit/e89651e4de6a64298b4442033e10d9864a641ba3))
* **deps:** update module github.com/prometheus/client_golang to v1.23.1 ([#35](https://github.com/monkescience/reference-service-go/issues/35)) ([19a9f67](https://github.com/monkescience/reference-service-go/commit/19a9f677f9d709d55fedd22e1292e6ae01e5f852))
* **deps:** update module github.com/prometheus/client_golang to v1.23.2 ([#37](https://github.com/monkescience/reference-service-go/issues/37)) ([8725711](https://github.com/monkescience/reference-service-go/commit/8725711419e02e4904d45894156961f2a0fdfc08))

## [0.5.0](https://github.com/monkescience/reference-service-go/compare/0.4.0...0.5.0) (2025-07-11)


### Features

* **ci:** add secrets for GitOps app private key in eu-central-1-prod deployment ([841a47b](https://github.com/monkescience/reference-service-go/commit/841a47b3515c3e00ecf0d89f0164c4e64bd2464f))
* **ci:** add support for GitOps app private key in workflows ([aa6d615](https://github.com/monkescience/reference-service-go/commit/aa6d615ff8a199b9a92aab23876be834c3c8237e))

## [0.4.0](https://github.com/monkescience/reference-service-go/compare/0.3.0...0.4.0) (2025-07-11)


### Features

* add Dockerfile and go.sum for service containerization ([173a832](https://github.com/monkescience/reference-service-go/commit/173a8320ba25a016e2ea45aac4a1c34c2ca09a4a))
* add initial service ([d7c230b](https://github.com/monkescience/reference-service-go/commit/d7c230b997b08b0944437bbdf8ac292f8c324f6b))
* add renovate config and github action setup ([47d0d76](https://github.com/monkescience/reference-service-go/commit/47d0d766e176b588c76abd18b5f4008caad43b6c))
* **build:** add multi-platform build support to Dockerfile and update build-sync workflow ([7e2607b](https://github.com/monkescience/reference-service-go/commit/7e2607be7eed098ff36030b027adf0160d2a451d))
* **ci:** add GitHub Actions workflow for building Docker image and syncing GitOps ([8dd68c9](https://github.com/monkescience/reference-service-go/commit/8dd68c9179016c373d847ffbee8587d59834c7fd))
* **ci:** add GitHub App token generation to release-please workflow ([1825f89](https://github.com/monkescience/reference-service-go/commit/1825f895adf109b790123405fccd85096fa62092))
* **ci:** add kustomize validation and release-please workflows ([bbd6195](https://github.com/monkescience/reference-service-go/commit/bbd6195741aa65b1e8857343a2db08ffe8ab24f8))
* **ci:** grant package write permissions in release-please workflow ([bd55243](https://github.com/monkescience/reference-service-go/commit/bd55243f669a380be0a8376de7c6a42b5c6cbd87))
* **ci:** modularize deployment workflows and add eu-central-1-prod deployment ([eca4a2c](https://github.com/monkescience/reference-service-go/commit/eca4a2cc0388b411bc5acb460cfc2c4276e4c768))
* **ci:** rename build.yaml to build-sync.yaml and update Kustomize download URL ([e413ea2](https://github.com/monkescience/reference-service-go/commit/e413ea209fb82eee60e24217cb315aa1224a8e79))
* **ci:** update Kustomize download method to use executable directly ([1a150ff](https://github.com/monkescience/reference-service-go/commit/1a150ff2f52e57bdaf141d7f181ba032d1ed5354))
* **ci:** update Kustomize version format and change download method to use tar.gz ([c8e1547](https://github.com/monkescience/reference-service-go/commit/c8e1547fe2e5a4024dd50186d81c75b42be764ed))
* **deployment:** add Kubernetes deployment, service, HPA, and kustomization for reference service ([480c000](https://github.com/monkescience/reference-service-go/commit/480c000dd2c0d4e1a2b2c491689e98b53e83270b))
* **kustomization:** add PodDisruptionBudget and update deployment configuration ([cbc4627](https://github.com/monkescience/reference-service-go/commit/cbc46277800775c8487bf36e0b5a6e2edbe333b2))
* **kustomization:** update apiVersion and restructure resource patches in kustomization.yaml ([2b04902](https://github.com/monkescience/reference-service-go/commit/2b049021e5c149d87207fcbaa55e4b0ec2e941ca))
* **kustomization:** update labels and image tags for deployment resources ([e170a3a](https://github.com/monkescience/reference-service-go/commit/e170a3a08775d9871adc035463697cb10756c35a))
* **metrics:** add HTTP response time metrics and health check API ([454dcba](https://github.com/monkescience/reference-service-go/commit/454dcba865664dfe68344a44fa00f148a9081d88))
* **namespace:** add ambient mode label to reference-service-go namespace ([5bc5e3a](https://github.com/monkescience/reference-service-go/commit/5bc5e3aa9df73a9ec3679a7d813d05a2f174eea5))
* update renovate config to use semantic commits ([9a6f371](https://github.com/monkescience/reference-service-go/commit/9a6f371f269cac1ae9106c8d1281583708e36cf8))


### Bug Fixes

* **build-sync:** enable push trigger for main branch ([9bc5401](https://github.com/monkescience/reference-service-go/commit/9bc540154858a58ac54e34d7657780f69cd9709f))
* **build-sync:** set git user configuration for GitOps app ([2b61f08](https://github.com/monkescience/reference-service-go/commit/2b61f088da64e0b64e5c1c89c21c48f83903db99))
* **build-sync:** update app-id reference to use vars and clean up git config commands ([17f854b](https://github.com/monkescience/reference-service-go/commit/17f854b6d0c8fdcb5ef50975f1c7801cde600429))
* **ci:** correct extraction command for Kustomize tar.gz file ([12346f8](https://github.com/monkescience/reference-service-go/commit/12346f8358ac44b1dcd8ca349dd5354963553ae0))
* **deployment:** remove cpu limit and increase memeory limit ([c0c9fb0](https://github.com/monkescience/reference-service-go/commit/c0c9fb05a48f0efe67402a42861b3ae670b30292))
* **deps:** update module github.com/getkin/kin-openapi to v0.131.0 [security] ([#2](https://github.com/monkescience/reference-service-go/issues/2)) ([5c7ddfa](https://github.com/monkescience/reference-service-go/commit/5c7ddfa43cd5fdc352eb4a6bec860ea1b7521a51))
* **deps:** update module github.com/getkin/kin-openapi to v0.132.0 ([#6](https://github.com/monkescience/reference-service-go/issues/6)) ([5512cbd](https://github.com/monkescience/reference-service-go/commit/5512cbd3268b37be628618146ec19c94d03f0957))
* **deps:** update module github.com/go-chi/chi/v5 to v5.2.1 ([#4](https://github.com/monkescience/reference-service-go/issues/4)) ([945a6c9](https://github.com/monkescience/reference-service-go/commit/945a6c9b3a7061bb7e6224857f6164cde26d1426))
* **deps:** update module github.com/go-chi/chi/v5 to v5.2.2 [security] ([#13](https://github.com/monkescience/reference-service-go/issues/13)) ([1e8690e](https://github.com/monkescience/reference-service-go/commit/1e8690e41e8c5cf1875ab386271d90faad03f77b))
* **deps:** update module github.com/oklog/ulid/v2 to v2.1.1 ([#7](https://github.com/monkescience/reference-service-go/issues/7)) ([2545f7e](https://github.com/monkescience/reference-service-go/commit/2545f7edd01f6eefc8cf36c842730fff8b90f791))
* **kustomization:** set includeSelectors to false ([1b1d9a3](https://github.com/monkescience/reference-service-go/commit/1b1d9a30d6edc00955851473e3208e68d07aaacf))
* **kustomization:** update label addition to include --without-selector flag ([0fe363b](https://github.com/monkescience/reference-service-go/commit/0fe363b7d86f95db09efd86a1c479fd4ad8e7886))
* **pdb:** reduce minAvailable to 50% for reference-service-go ([e5d3263](https://github.com/monkescience/reference-service-go/commit/e5d32630a7b8975dee61dbc55cfff925f60af859))
* **release-please:** disable inclusion of 'v' prefix in tags ([0af3286](https://github.com/monkescience/reference-service-go/commit/0af32864102e4f4b276a9606bd0f8928f2735e7b))
* **release-please:** remove draft mode from configuration ([b48a308](https://github.com/monkescience/reference-service-go/commit/b48a3082488cf50cb2ed7f1e7d6b793ee81abb50))

## [0.3.0](https://github.com/monkescience/reference-service-go/compare/v0.2.0...v0.3.0) (2025-07-11)


### Features

* add Dockerfile and go.sum for service containerization ([173a832](https://github.com/monkescience/reference-service-go/commit/173a8320ba25a016e2ea45aac4a1c34c2ca09a4a))
* add initial service ([d7c230b](https://github.com/monkescience/reference-service-go/commit/d7c230b997b08b0944437bbdf8ac292f8c324f6b))
* add renovate config and github action setup ([47d0d76](https://github.com/monkescience/reference-service-go/commit/47d0d766e176b588c76abd18b5f4008caad43b6c))
* **build:** add multi-platform build support to Dockerfile and update build-sync workflow ([7e2607b](https://github.com/monkescience/reference-service-go/commit/7e2607be7eed098ff36030b027adf0160d2a451d))
* **ci:** add GitHub Actions workflow for building Docker image and syncing GitOps ([8dd68c9](https://github.com/monkescience/reference-service-go/commit/8dd68c9179016c373d847ffbee8587d59834c7fd))
* **ci:** add GitHub App token generation to release-please workflow ([1825f89](https://github.com/monkescience/reference-service-go/commit/1825f895adf109b790123405fccd85096fa62092))
* **ci:** add kustomize validation and release-please workflows ([bbd6195](https://github.com/monkescience/reference-service-go/commit/bbd6195741aa65b1e8857343a2db08ffe8ab24f8))
* **ci:** grant package write permissions in release-please workflow ([bd55243](https://github.com/monkescience/reference-service-go/commit/bd55243f669a380be0a8376de7c6a42b5c6cbd87))
* **ci:** modularize deployment workflows and add eu-central-1-prod deployment ([eca4a2c](https://github.com/monkescience/reference-service-go/commit/eca4a2cc0388b411bc5acb460cfc2c4276e4c768))
* **ci:** rename build.yaml to build-sync.yaml and update Kustomize download URL ([e413ea2](https://github.com/monkescience/reference-service-go/commit/e413ea209fb82eee60e24217cb315aa1224a8e79))
* **ci:** update Kustomize download method to use executable directly ([1a150ff](https://github.com/monkescience/reference-service-go/commit/1a150ff2f52e57bdaf141d7f181ba032d1ed5354))
* **ci:** update Kustomize version format and change download method to use tar.gz ([c8e1547](https://github.com/monkescience/reference-service-go/commit/c8e1547fe2e5a4024dd50186d81c75b42be764ed))
* **deployment:** add Kubernetes deployment, service, HPA, and kustomization for reference service ([480c000](https://github.com/monkescience/reference-service-go/commit/480c000dd2c0d4e1a2b2c491689e98b53e83270b))
* **kustomization:** add PodDisruptionBudget and update deployment configuration ([cbc4627](https://github.com/monkescience/reference-service-go/commit/cbc46277800775c8487bf36e0b5a6e2edbe333b2))
* **kustomization:** update apiVersion and restructure resource patches in kustomization.yaml ([2b04902](https://github.com/monkescience/reference-service-go/commit/2b049021e5c149d87207fcbaa55e4b0ec2e941ca))
* **kustomization:** update labels and image tags for deployment resources ([e170a3a](https://github.com/monkescience/reference-service-go/commit/e170a3a08775d9871adc035463697cb10756c35a))
* **metrics:** add HTTP response time metrics and health check API ([454dcba](https://github.com/monkescience/reference-service-go/commit/454dcba865664dfe68344a44fa00f148a9081d88))
* **namespace:** add ambient mode label to reference-service-go namespace ([5bc5e3a](https://github.com/monkescience/reference-service-go/commit/5bc5e3aa9df73a9ec3679a7d813d05a2f174eea5))
* update renovate config to use semantic commits ([9a6f371](https://github.com/monkescience/reference-service-go/commit/9a6f371f269cac1ae9106c8d1281583708e36cf8))


### Bug Fixes

* **build-sync:** enable push trigger for main branch ([9bc5401](https://github.com/monkescience/reference-service-go/commit/9bc540154858a58ac54e34d7657780f69cd9709f))
* **build-sync:** set git user configuration for GitOps app ([2b61f08](https://github.com/monkescience/reference-service-go/commit/2b61f088da64e0b64e5c1c89c21c48f83903db99))
* **build-sync:** update app-id reference to use vars and clean up git config commands ([17f854b](https://github.com/monkescience/reference-service-go/commit/17f854b6d0c8fdcb5ef50975f1c7801cde600429))
* **ci:** correct extraction command for Kustomize tar.gz file ([12346f8](https://github.com/monkescience/reference-service-go/commit/12346f8358ac44b1dcd8ca349dd5354963553ae0))
* **deployment:** remove cpu limit and increase memeory limit ([c0c9fb0](https://github.com/monkescience/reference-service-go/commit/c0c9fb05a48f0efe67402a42861b3ae670b30292))
* **deps:** update module github.com/getkin/kin-openapi to v0.131.0 [security] ([#2](https://github.com/monkescience/reference-service-go/issues/2)) ([5c7ddfa](https://github.com/monkescience/reference-service-go/commit/5c7ddfa43cd5fdc352eb4a6bec860ea1b7521a51))
* **deps:** update module github.com/getkin/kin-openapi to v0.132.0 ([#6](https://github.com/monkescience/reference-service-go/issues/6)) ([5512cbd](https://github.com/monkescience/reference-service-go/commit/5512cbd3268b37be628618146ec19c94d03f0957))
* **deps:** update module github.com/go-chi/chi/v5 to v5.2.1 ([#4](https://github.com/monkescience/reference-service-go/issues/4)) ([945a6c9](https://github.com/monkescience/reference-service-go/commit/945a6c9b3a7061bb7e6224857f6164cde26d1426))
* **deps:** update module github.com/go-chi/chi/v5 to v5.2.2 [security] ([#13](https://github.com/monkescience/reference-service-go/issues/13)) ([1e8690e](https://github.com/monkescience/reference-service-go/commit/1e8690e41e8c5cf1875ab386271d90faad03f77b))
* **deps:** update module github.com/oklog/ulid/v2 to v2.1.1 ([#7](https://github.com/monkescience/reference-service-go/issues/7)) ([2545f7e](https://github.com/monkescience/reference-service-go/commit/2545f7edd01f6eefc8cf36c842730fff8b90f791))
* **kustomization:** set includeSelectors to false ([1b1d9a3](https://github.com/monkescience/reference-service-go/commit/1b1d9a30d6edc00955851473e3208e68d07aaacf))
* **kustomization:** update label addition to include --without-selector flag ([0fe363b](https://github.com/monkescience/reference-service-go/commit/0fe363b7d86f95db09efd86a1c479fd4ad8e7886))
* **pdb:** reduce minAvailable to 50% for reference-service-go ([e5d3263](https://github.com/monkescience/reference-service-go/commit/e5d32630a7b8975dee61dbc55cfff925f60af859))
* **release-please:** remove draft mode from configuration ([b48a308](https://github.com/monkescience/reference-service-go/commit/b48a3082488cf50cb2ed7f1e7d6b793ee81abb50))

## [0.2.0](https://github.com/monkescience/reference-service-go/compare/v0.1.0...v0.2.0) (2025-07-11)


### Features

* add Dockerfile and go.sum for service containerization ([173a832](https://github.com/monkescience/reference-service-go/commit/173a8320ba25a016e2ea45aac4a1c34c2ca09a4a))
* add initial service ([d7c230b](https://github.com/monkescience/reference-service-go/commit/d7c230b997b08b0944437bbdf8ac292f8c324f6b))
* add renovate config and github action setup ([47d0d76](https://github.com/monkescience/reference-service-go/commit/47d0d766e176b588c76abd18b5f4008caad43b6c))
* **build:** add multi-platform build support to Dockerfile and update build-sync workflow ([7e2607b](https://github.com/monkescience/reference-service-go/commit/7e2607be7eed098ff36030b027adf0160d2a451d))
* **ci:** add GitHub Actions workflow for building Docker image and syncing GitOps ([8dd68c9](https://github.com/monkescience/reference-service-go/commit/8dd68c9179016c373d847ffbee8587d59834c7fd))
* **ci:** add GitHub App token generation to release-please workflow ([1825f89](https://github.com/monkescience/reference-service-go/commit/1825f895adf109b790123405fccd85096fa62092))
* **ci:** add kustomize validation and release-please workflows ([bbd6195](https://github.com/monkescience/reference-service-go/commit/bbd6195741aa65b1e8857343a2db08ffe8ab24f8))
* **ci:** modularize deployment workflows and add eu-central-1-prod deployment ([eca4a2c](https://github.com/monkescience/reference-service-go/commit/eca4a2cc0388b411bc5acb460cfc2c4276e4c768))
* **ci:** rename build.yaml to build-sync.yaml and update Kustomize download URL ([e413ea2](https://github.com/monkescience/reference-service-go/commit/e413ea209fb82eee60e24217cb315aa1224a8e79))
* **ci:** update Kustomize download method to use executable directly ([1a150ff](https://github.com/monkescience/reference-service-go/commit/1a150ff2f52e57bdaf141d7f181ba032d1ed5354))
* **ci:** update Kustomize version format and change download method to use tar.gz ([c8e1547](https://github.com/monkescience/reference-service-go/commit/c8e1547fe2e5a4024dd50186d81c75b42be764ed))
* **deployment:** add Kubernetes deployment, service, HPA, and kustomization for reference service ([480c000](https://github.com/monkescience/reference-service-go/commit/480c000dd2c0d4e1a2b2c491689e98b53e83270b))
* **kustomization:** add PodDisruptionBudget and update deployment configuration ([cbc4627](https://github.com/monkescience/reference-service-go/commit/cbc46277800775c8487bf36e0b5a6e2edbe333b2))
* **kustomization:** update apiVersion and restructure resource patches in kustomization.yaml ([2b04902](https://github.com/monkescience/reference-service-go/commit/2b049021e5c149d87207fcbaa55e4b0ec2e941ca))
* **kustomization:** update labels and image tags for deployment resources ([e170a3a](https://github.com/monkescience/reference-service-go/commit/e170a3a08775d9871adc035463697cb10756c35a))
* **metrics:** add HTTP response time metrics and health check API ([454dcba](https://github.com/monkescience/reference-service-go/commit/454dcba865664dfe68344a44fa00f148a9081d88))
* **namespace:** add ambient mode label to reference-service-go namespace ([5bc5e3a](https://github.com/monkescience/reference-service-go/commit/5bc5e3aa9df73a9ec3679a7d813d05a2f174eea5))
* update renovate config to use semantic commits ([9a6f371](https://github.com/monkescience/reference-service-go/commit/9a6f371f269cac1ae9106c8d1281583708e36cf8))


### Bug Fixes

* **build-sync:** enable push trigger for main branch ([9bc5401](https://github.com/monkescience/reference-service-go/commit/9bc540154858a58ac54e34d7657780f69cd9709f))
* **build-sync:** set git user configuration for GitOps app ([2b61f08](https://github.com/monkescience/reference-service-go/commit/2b61f088da64e0b64e5c1c89c21c48f83903db99))
* **build-sync:** update app-id reference to use vars and clean up git config commands ([17f854b](https://github.com/monkescience/reference-service-go/commit/17f854b6d0c8fdcb5ef50975f1c7801cde600429))
* **ci:** correct extraction command for Kustomize tar.gz file ([12346f8](https://github.com/monkescience/reference-service-go/commit/12346f8358ac44b1dcd8ca349dd5354963553ae0))
* **deployment:** remove cpu limit and increase memeory limit ([c0c9fb0](https://github.com/monkescience/reference-service-go/commit/c0c9fb05a48f0efe67402a42861b3ae670b30292))
* **deps:** update module github.com/getkin/kin-openapi to v0.131.0 [security] ([#2](https://github.com/monkescience/reference-service-go/issues/2)) ([5c7ddfa](https://github.com/monkescience/reference-service-go/commit/5c7ddfa43cd5fdc352eb4a6bec860ea1b7521a51))
* **deps:** update module github.com/getkin/kin-openapi to v0.132.0 ([#6](https://github.com/monkescience/reference-service-go/issues/6)) ([5512cbd](https://github.com/monkescience/reference-service-go/commit/5512cbd3268b37be628618146ec19c94d03f0957))
* **deps:** update module github.com/go-chi/chi/v5 to v5.2.1 ([#4](https://github.com/monkescience/reference-service-go/issues/4)) ([945a6c9](https://github.com/monkescience/reference-service-go/commit/945a6c9b3a7061bb7e6224857f6164cde26d1426))
* **deps:** update module github.com/go-chi/chi/v5 to v5.2.2 [security] ([#13](https://github.com/monkescience/reference-service-go/issues/13)) ([1e8690e](https://github.com/monkescience/reference-service-go/commit/1e8690e41e8c5cf1875ab386271d90faad03f77b))
* **deps:** update module github.com/oklog/ulid/v2 to v2.1.1 ([#7](https://github.com/monkescience/reference-service-go/issues/7)) ([2545f7e](https://github.com/monkescience/reference-service-go/commit/2545f7edd01f6eefc8cf36c842730fff8b90f791))
* **kustomization:** set includeSelectors to false ([1b1d9a3](https://github.com/monkescience/reference-service-go/commit/1b1d9a30d6edc00955851473e3208e68d07aaacf))
* **kustomization:** update label addition to include --without-selector flag ([0fe363b](https://github.com/monkescience/reference-service-go/commit/0fe363b7d86f95db09efd86a1c479fd4ad8e7886))
* **pdb:** reduce minAvailable to 50% for reference-service-go ([e5d3263](https://github.com/monkescience/reference-service-go/commit/e5d32630a7b8975dee61dbc55cfff925f60af859))
