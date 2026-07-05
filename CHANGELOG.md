# Changelog

All notable changes to the Pulumi Example Foundation (Go) will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## 1.0.0 (2026-07-05)


### Features

* **0-bootstrap:** add billing.creator, SA impersonation, and bucket IAM ([57c5ce5](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/57c5ce56c20b25ce503986e54d38198d7c725612))
* **0-bootstrap:** add optional Google Workspace group creation ([d15f3cf](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/d15f3cfa9464a0924a3ed4796df222f4f3772e0f))
* **1-org:** achieve full IAM/policy parity with Terraform foundation ([fb62524](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/fb6252484fac8fa23b7c4d84e55875e80433003a))
* **1-org:** achieve full parity with Terraform Enterprise Foundation ([#7](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/7)) ([937c14f](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/937c14fe84d806972c54685abe697955c2ec0bee))
* **2-environments:** add bootstrap stack reference for common_config ([11e3ee6](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/11e3ee63c8a09efb556fc41b222808b4e5e964a1))
* **2-environments:** implement full upstream parity ([#16](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/16)) ([7121723](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/7121723afca2748b18e4d63789e883ccdea3d27a))
* **3-networks:** close networking parity gaps with upstream ([#86](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/86)) ([d302a91](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/d302a914929b3837312774a5c1115adafd319abd))
* **3-networks:** expose VPC flow log options and DNS/firewall logging toggles ([#58](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/58)) ([800f274](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/800f274022575ad8aa177df858c314931bffcccc))
* **4-projects:** full upstream parity — peering, CMEK, VPC-SC, labels, budgets ([9b09420](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/9b09420f04c01790d06efd7d4c77854b9cbd584d))
* achieve full output parity with Terraform reference architecture ([#44](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/44)) ([6d8c47d](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/6d8c47d3492889436d4bcb9474babf83139acd3c))
* achieve parity with upstream 5-app-infra ([#28](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/28)) ([10fa891](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/10fa891fcd399be70fd933d792cbcb8fa398580d))
* add E2E testing infrastructure ([#37](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/37)) ([96fac02](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/96fac0201547b66a94d2b9c6b1892e23baca823b))
* add Pulumi stack configuration templates and documentation for environment-based deployments ([01d03dc](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/01d03dcdfbe3e0e2f8befd4fe927d0820ce9c416))
* add release-please automation for proper versioning ([#32](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/32)) ([91c301f](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/91c301f3f6df04ae5712ef0fd5723ea4617e36a4))
* adopt descriptive upstream-consistent library module names ([#85](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/85)) ([4cf9fc8](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/4cf9fc8a81bcc0354b968cc18fd14705389a19e6))
* **bootstrap:** achieve full parity with Terraform Enterprise Foundation ([#6](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/6)) ([ddae105](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/ddae1056225ff9bca760990e10989c95f38bf4ff))
* complete parity remediation for 4-projects ([#27](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/27)) ([2578ff3](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/2578ff349adc3981e4d83f8794cf52f442116a9d))
* e2e testing ([#39](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/39)) ([1db4eb7](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/1db4eb762efd00fb0f3fb810c3d92af28116b8c2))
* enable random project ID suffixes across all stages ([#4](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/4)) ([03ad067](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/03ad067b2f03b2121c7d5a2e5df6cfd9e37ece0d))
* implement 1-org gaps (billing sink, CAI monitoring, dependency ordering, KMS agent, budget) ([#9](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/9)) ([33bd281](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/33bd2818d960dbcf555e00c138f4aa2f1bbad449))
* initialize Pulumi Go foundation with CI/CD, documentation, and policy library scaffolding ([7f39e0d](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/7f39e0d327055a0930e713465a864e28bf6fd70c))
* integrate all fine-grained library packages ([#43](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/43)) ([5c167c4](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/5c167c4c7237cc230d5fb29755f58c3b4d27355b))
* migrate to restructured pulumi-library go monorepo ([#31](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/31)) ([65c2c43](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/65c2c431379865b627384c9997f1dff55daa5fea))
* **networks:** align SVPC and hub-spoke with upstream terraform foundation ([#21](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/21)) ([50bc0a1](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/50bc0a1dbb15856fd065686bd4d9b683dd990e9d))
* Phase 7 interconnect module parity ([#80](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/80)) ([2672b54](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/2672b54181a75240363899c33bad88042e2619cd))
* Phase 8 Upstream Drift Parity ([#81](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/81)) ([0de4558](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/0de4558b18be96b41052c4bc0bea289f300a0483))
* Plumb network component gaps for svpc and hub-and-spoke ([#26](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/26)) ([d862fa0](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/d862fa0aace4afce27b5d13025cc50b8d03a4479))
* remediate foundation parity gaps in phases 0-2 ([#18](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/18)) ([acb98b9](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/acb98b9b13b9f53bae85b5dbf37b42448929eacf))


### Bug Fixes

* **0-bootstrap:** add missing iamcredentials API ([#15](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/15)) ([3f7c697](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/3f7c697b53270283b41eb7be8afe24aa25088330))
* **0-bootstrap:** align project labels, config, and exports with upstream ([#13](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/13)) ([600a4c5](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/600a4c5455afa3ad2b95feb21b1bfa34f6ea2084))
* **0-bootstrap:** correct KMS region and add KMS prevention parity ([#14](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/14)) ([5792b6f](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/5792b6f186a492bea4c81149c10977e7921ed2e7))
* **0-bootstrap:** enforce group dependencies for IAM bindings ([#70](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/70)) ([b8662ba](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/b8662babee77337148797b2e945b6d3afaca379b))
* **0-bootstrap:** grant billing.admin to billing admins group ([#66](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/66)) ([6660a8c](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/6660a8cf9d2e1a1363961395371a904baa5a70af))
* **0-bootstrap:** map SA output names to match upstream TF ([#67](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/67)) ([0906c8f](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/0906c8f2c34a1c81f81ded96ae658d7a02362778))
* **1-org:** add missing labels, random suffixes, and default corrections ([#12](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/12)) ([18d0576](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/18d0576d93e7ba66c8c672e712c6a8a50d9f48bc))
* **1-org:** align export names with terraform foundation ([#11](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/11)) ([90a0365](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/90a0365b972faef1999068b7a3308056ec5066d4))
* **1-org:** align pulumi 1-org foundation with upstream terraform ([#10](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/10)) ([4a1a35d](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/4a1a35d57754a66c607abeaaaedfae358bb7c649))
* **1-org:** consume bootstrap stack outputs for group IAM ([#60](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/60)) ([8694c04](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/8694c04479ef98e00c7e4d19fadd533e908fb036))
* **1-org:** deterministic Essential Contacts resource naming ([#53](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/53)) ([d04af29](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/d04af29055e82ba43c56130544e00ead7604136a))
* **1-org:** provision KMS org service agent at org scope ([#50](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/50)) ([8d790e9](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/8d790e9ea4768e7fe4c13621fb79ea89f041bfa3))
* **1-org:** refactor CAIMonitoring to use library component ([#24](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/24)) ([7f497dc](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/7f497dc88829c394815dad95ab639faeebfb7a51))
* **1-org:** remove duplicate AccessPolicy creation in policies.go ([#47](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/47)) ([564ff85](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/564ff8586d74aed28607a57a6d482c3073fda176))
* **2-environments:** add shared_network to budget config for parity ([#25](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/25)) ([fef602e](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/fef602e6fa31e835fcb22b6adc4a753e4493c7d3))
* **2-environments:** export assured_workload_id and assured_workload_resources ([#55](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/55)) ([8862c3c](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/8862c3cbc676bedb5bfa152d63aeeb9d4cce2488))
* **2-environments:** remove broken async applyCommonConfig ([babf3b1](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/babf3b1dd7d0333b079ff90aad80a5aab7701e1a))
* **3-networks-hub-and-spoke:** add health-check firewall for transitivity ILBs ([#63](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/63)) ([3ac5af1](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/3ac5af1b414d8289dfa0f92a040eb9d091888507))
* **3-networks-hub-and-spoke:** add tagged internet egress route ([#62](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/62)) ([a28f152](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/a28f152d63a5bd178f3dabd2bfc24030a1b8b232))
* **3-networks-hub-and-spoke:** differentiate per-env spoke CIDRs ([#59](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/59)) ([99f4f96](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/99f4f9664671c53c26bb1792222f6f5e616fb820))
* **3-networks-hub-and-spoke:** populate subnets_secondary_ranges output ([#54](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/54)) ([0b709bb](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/0b709bb048f7b9ea352c67557ecb99e3ca228f5e))
* **3-networks-hub-and-spoke:** require parent_id and fix example configs ([#56](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/56)) ([0a8dc78](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/0a8dc7819beb5a6df36cbd3a62942160abfc4d18))
* **3-networks:** create VPC-SC hub perimeter once in shared stack; gate svpc perimeter on policy source ([8195bac](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/8195bac883f190b7e8a986219114ae7e79d3aafe))
* **3-networks:** default enforce_vpcsc to false (dry-run first, matching TF) ([beb816e](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/beb816e7cfe64728cb96b01a8f82faf16a08027b))
* **3-networks:** wire VPC-SC propagation delay and interconnect policies ([#74](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/74)) ([012282d](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/012282dde3e8b901c3c87b55774518caf0f97b5a))
* **3-networks:** wire VPC-SC unconditionally and add hub perimeter ([#61](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/61)) ([183ba04](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/183ba04033d7ff8716cdfa8e36f1f28b558587c1))
* **4-projects:** add VPC-SC dry run attachment branch ([#73](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/73)) ([a5dbaca](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/a5dbaca77c1f2d10ce58c432fec8ba089961a706))
* **4-projects:** deprivilege default compute SA on sample projects ([1667dd6](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/1667dd6f528195e289bd873107e0e020e38a2e31))
* **4-projects:** remove stale -base suffix from peering VPC target ([#49](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/49)) ([ed1eaed](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/ed1eaedbdc0b4e0f03e47f191cf1b60fd169a543))
* **4-projects:** resolve folder ID from 2-environments stack ([#48](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/48)) ([aa0e31d](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/aa0e31da587aafd2601ad6db7c797818251f6edc))
* **5-app-infra:** address gap report issues for go foundation ([#75](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/75)) ([f9cdce8](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/f9cdce8e46026e9d86eb52c7c17c32de447e2faa))
* **5-app-infra:** address secondary audit parity gaps ([#30](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/30)) ([f6e7caf](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/f6e7cafd15406811756b65e3e44d792f4dae713c))
* **5-app-infra:** remediate review findings ([#29](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/29)) ([1e37b05](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/1e37b054bf466e61ae055e184a8eef23f42d49ff))
* **5-app-infra:** replace active scaffold and use dynamic image reference ([#78](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/78)) ([3764889](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/3764889e91162299d7c4607bc96027792537cdde))
* **5-app-infra:** update library dependencies to latest main ([#76](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/76)) ([b13556d](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/b13556dae325d9342dd0defc69051e67bcb10379))
* **5-app-infra:** wire subnetwork into the Confidential Space instance template ([e5a8475](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/e5a847591e78f956ec8c68e7eb35d959ba1156ec))
* achieve 100% architectural parity ([#46](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/46)) ([28cd39b](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/28cd39b11c0b68ee4bf6e0e338402f483c52f399))
* address ci and foundational gaps ([#77](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/77)) ([6ea6f0c](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/6ea6f0c9ae6560ea339fdf477013a05f266be147))
* Align go foundation exports with Terraform logic ([#45](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/45)) ([b49b6cf](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/b49b6cfa4b3e394cf54b0a03ed4e1bd877111186))
* API and IAM bindings usage of pulumi-library ([#1](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/1)) ([a1b5f36](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/a1b5f36a665931e6f399ca86755865f9ade20936))
* **bootstrap:** fix WIF wildcard binding string formatting ([#71](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/71)) ([b5807ec](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/b5807ec79ac93a4f4e775f2c1a5c9ce34056f15d))
* correct deploy/compile-breakers from gap-fix review ([e7aa871](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/e7aa8717c2c9cefa1bdbdf289fd934faad374978))
* correct deploy/compile-breakers surfaced by gap-fix review ([7ba7f9c](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/7ba7f9ca888e5b96bc7726e80d63a5a71194e2c5))
* **deps:** resolve packages from go module proxy ([#35](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/35)) ([2dd327f](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/2dd327fe31b9a2bdbc7a74fd91a8893b944dae1a))
* **docs:** fix output name references in READMEs ([#64](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/64)) ([97941b7](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/97941b7175a96f892d24312fbb0eab3fb4c259ad))
* e2e folder support ([#40](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/40)) ([f433413](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/f433413781ae534b2afa8a362cdb738192c9c8a2))
* e2e folder support and robust clean script ([#38](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/38)) ([0b5b224](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/0b5b224457dbfe823de1e6d3315455c552e03653))
* **go:** change enable_scc_resources and log_export_storage_versioning defaults to false ([#72](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/72)) ([265e4ae](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/265e4aec6ae6f25fd8392e8f3d56a47b078e5289))
* **go:** update module paths for workspace migration ([00877a4](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/00877a419fa96ed27c8584b6a41cb17ac42a5863))
* **go:** update module paths for workspace migration ([545d669](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/545d66955725fad422744a489c46377694a48f39))
* Phase 6 gap remediations ([#79](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/79)) ([baedada](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/baedadad3aafbddff9efbca3fd9a129e7e926103))
* remove deprecated sourcerepo API and align KMS regions ([#41](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/41)) ([79d1aa8](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/79d1aa8ea4f54c0f8b016c31c42e4a6f33bd5a68))
* resolve phase 3 architectural gaps ([#22](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/22)) ([5d712c6](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/5d712c6dfe46eebb1b976e7e88fd71a4db0768cd))
* **test:** add 1-org config tests ([#65](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/65)) ([bb8b129](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/bb8b1295059cff486b2b8c46402fe7c435a7fbaf))
* **test:** correct validate_e2e.sh KMS location and project filters ([#52](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/52)) ([6129bca](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/6129bca66b49245d88165616b72137428e89d342))
* **test:** update config default assertions in 1-org ([9d97366](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/9d97366f6681b6cfc30661228730164d1d746f04))
* update pulumi-library dependency to include kms region parity ([#42](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/issues/42)) ([42cba4c](https://github.com/VitruvianSoftware/pulumi_go-example-foundation/commit/42cba4cf7525f0999d59a29b353d9a2900ad7094))

## [Unreleased]

### Added

- Initial 6-stage foundation (0-bootstrap through 5-app-infra)
- Shared VPC and Hub-and-Spoke network topologies
- GitHub Actions CI/CD pipeline with Workload Identity Federation
- GitLab CI/CD pipeline alternative
- Comprehensive onboarding guide (`ONBOARDING.md`)
- Pre-flight validation script (`scripts/validate-requirements.sh`)
- Documentation suite: README, CONTRIBUTING, SECURITY, ERRATA, FAQ, GLOSSARY, TROUBLESHOOTING
- CrossGuard policy pack skeleton (`policy-library/`)
- Per-stage Configuration Reference and Outputs tables
- Resource hierarchy change guide (`docs/change_resource_hierarchy.md`)

### Changed

- Migrated shared components to [pulumi-library](https://github.com/VitruvianSoftware/pulumi-library/go)

### Security

- WIF-only authentication (no service account keys stored in CI/CD)
- KMS-encrypted Pulumi state bucket with configurable protection level
- Deletion protection on bootstrap folder, seed project, and CI/CD project
