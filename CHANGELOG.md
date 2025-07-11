# Changelog

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
