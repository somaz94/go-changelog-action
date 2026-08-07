# Changelog

All notable changes to this project will be documented in this file.

## Unreleased (2026-08-07)

### Performance Improvements

- ship a prebuilt multi-arch image instead of building per run ([d19f85c](https://github.com/somaz94/go-changelog-action/commit/d19f85cc789ae2aadc571b2f039fa2da2305483c))

### Continuous Integration

- add a golangci-lint config scoped to defect-finding linters ([4c6f375](https://github.com/somaz94/go-changelog-action/commit/4c6f37529b6e96f4df421be88942e83c62a8a018))

### Contributors

- somaz

<br/>

## [v1.0.10](https://github.com/somaz94/go-changelog-action/compare/v1.0.9...v1.0.10) (2026-07-21)

### Code Refactoring

- hoist buildEntry invariant args into entryBuilder context struct ([5749bad](https://github.com/somaz94/go-changelog-action/commit/5749bad188b14a5ffcb27332b30fd47cbc6d8971))

### Contributors

- somaz

<br/>

## [v1.0.9](https://github.com/somaz94/go-changelog-action/compare/v1.0.8...v1.0.9) (2026-07-21)

### Code Refactoring

- propagate context to the git safe.directory command ([536ac6c](https://github.com/somaz94/go-changelog-action/commit/536ac6ccc50efedb324babdea6a63390c9ba790f))

### Tests

- isolate global git config in tests to stop ~/.gitconfig pollution ([8071ebf](https://github.com/somaz94/go-changelog-action/commit/8071ebfc607420771083ac8e3994b183bc65149d))

### Builds

- **deps:** bump actions/checkout from 6 to 7 (#6) ([#6](https://github.com/somaz94/go-changelog-action/pull/6)) ([85e90c1](https://github.com/somaz94/go-changelog-action/commit/85e90c1a52c256f6489bbda2932cb134cb3c2c18))
- **deps:** bump alpine from 3.23 to 3.24 in the docker-minor group (#5) ([#5](https://github.com/somaz94/go-changelog-action/pull/5)) ([ec5f0eb](https://github.com/somaz94/go-changelog-action/commit/ec5f0eba8e67b8ff70becd52e545f41e07eb8106))

### Continuous Integration

- remove DCO workflow ([45ee5c5](https://github.com/somaz94/go-changelog-action/commit/45ee5c50c86bd88a374cc857045d220f197515b6))
- adopt semantic-pr, labels, lock-threads, PR size, and auto-assign reusables ([b405a1a](https://github.com/somaz94/go-changelog-action/commit/b405a1a4e92bc53e312dfe79c6296ef06bbd1baf))
- use reusable stale-issues workflow ([5002bcf](https://github.com/somaz94/go-changelog-action/commit/5002bcf8aa290acae12a6670e6d705bf800470a9))
- use reusable issue-greeting workflow ([22cc5cd](https://github.com/somaz94/go-changelog-action/commit/22cc5cdb057c153ecc8e9369a18e217e9e15b897))
- use reusable dependabot-auto-merge workflow ([a7ba28d](https://github.com/somaz94/go-changelog-action/commit/a7ba28d1fdc6e376eb05b84e97af0093988e80fa))
- use reusable contributors workflow ([9457767](https://github.com/somaz94/go-changelog-action/commit/94577670b80df5a0a5fa8bf26de0210507ff495b))
- add ok-to-test workflow stub ([06cc11e](https://github.com/somaz94/go-changelog-action/commit/06cc11e1a9ceb4c0b90c1e457030a6bd3b2e3ad3))
- add PR welcome workflow stub ([e05c1b1](https://github.com/somaz94/go-changelog-action/commit/e05c1b1b1ccce86332ce7f5db5ca3c8de49c0ea0))
- add DCO check via shared reusable workflow ([fd273dc](https://github.com/somaz94/go-changelog-action/commit/fd273dc5b49083ae79a4ef26885e67cdfcc975ed))

### Contributors

- somaz

<br/>

## [v1.0.8](https://github.com/somaz94/go-changelog-action/compare/v1.0.7...v1.0.8) (2026-06-02)

### Code Refactoring

- harden path containment, output close, and git stderr ([3040d5b](https://github.com/somaz94/go-changelog-action/commit/3040d5b6df8acc7a4c1a03e5fc96eeff8d30db8b))

### Documentation

- remove duplicate rules covered by global CLAUDE.md ([8a97720](https://github.com/somaz94/go-changelog-action/commit/8a977206ef2c49f304484ee324a781170c467582))

### Builds

- **deps:** bump dependabot/fetch-metadata from 2 to 3 ([362575f](https://github.com/somaz94/go-changelog-action/commit/362575f8a61d3beadcc9fce7feb3820a297c375d))
- **deps:** bump actions/github-script from 8 to 9 ([71b0085](https://github.com/somaz94/go-changelog-action/commit/71b0085cbce894d5f0c8ba81be6cc104d4387ee3))
- **deps:** bump softprops/action-gh-release from 2 to 3 ([85864ea](https://github.com/somaz94/go-changelog-action/commit/85864ea0baf0ef083846743b7749d29dc15af26e))

### Continuous Integration

- add concurrency guards to recurring workflows ([02addc0](https://github.com/somaz94/go-changelog-action/commit/02addc051c03e83fc751e8b92fa78e4a27e48e6f))
- use go-docker-action-ci-action@v1 (replace inline prelude) ([30b6f39](https://github.com/somaz94/go-changelog-action/commit/30b6f390c3aaa019d1e7670d548f52f6aab7ecf9))

### Chores

- remove duplicate rules from CLAUDE.md (moved to global) ([29e3654](https://github.com/somaz94/go-changelog-action/commit/29e36546ea92b423febf9ded3d2e3bc3dd88760d))
- add git config protection to CLAUDE.md ([8fd3b8b](https://github.com/somaz94/go-changelog-action/commit/8fd3b8b3b5ba3ded6b897a3f7132aa71869fba4c))

### Contributors

- somaz

<br/>

## [v1.0.7](https://github.com/somaz94/go-changelog-action/compare/v1.0.6...v1.0.7) (2026-03-25)

### Bug Fixes

- add path traversal protection, git timeout, and isExcludedAuthor tests ([dbc5ddf](https://github.com/somaz94/go-changelog-action/commit/dbc5ddf872154c22536713a0cad0491f0c8bb218))

### Tests

- improve test coverage and fix short hash panic ([b6354a7](https://github.com/somaz94/go-changelog-action/commit/b6354a7ba7ebf748c8fe74f7e3d96bef618150b9))

### Continuous Integration

- skip auto-generated changelog and contributors commits in release notes ([a8cf70b](https://github.com/somaz94/go-changelog-action/commit/a8cf70b52cf82bbbcf8c36ccd91724f649439aad))
- revert to body_path RELEASE.md in release workflow ([def9762](https://github.com/somaz94/go-changelog-action/commit/def9762cb61f5acdc1ff3581ff5cf6e8a0c13e51))
- use generate_release_notes instead of body_path in release workflow ([a038a43](https://github.com/somaz94/go-changelog-action/commit/a038a438fb15aeccccc5d5d79757f8af7e12aa0b))

### Contributors

- somaz

<br/>

## [v1.0.6](https://github.com/somaz94/go-changelog-action/compare/v1.0.5...v1.0.6) (2026-03-20)

### Features

- skip changelog and contributors auto-commits by default ([a508ce6](https://github.com/somaz94/go-changelog-action/commit/a508ce69081498e906e8c6483f940ba62b3a008b))

### Documentation

- add no-push rule to CLAUDE.md ([be8bdf7](https://github.com/somaz94/go-changelog-action/commit/be8bdf75f142554ffeb7dcb2b9e8fab78e8ef822))
- add CLAUDE.md project guide ([9eb5d79](https://github.com/somaz94/go-changelog-action/commit/9eb5d7901a132e479f5ad7c02d48aa4bb10897a9))

### Tests

- add exclude_authors smoke test and CI docker env ([f15a8e1](https://github.com/somaz94/go-changelog-action/commit/f15a8e14767f7418bf84fbe9618a5cf7e72d8f01))

### Continuous Integration

- migrate gitlab-mirror workflow to multi-git-mirror action ([e8664e6](https://github.com/somaz94/go-changelog-action/commit/e8664e687c2e25bf1a3c78772b94f566f3f0f818))

### Contributors

- somaz

<br/>

## [v1.0.5](https://github.com/somaz94/go-changelog-action/compare/v1.0.4...v1.0.5) (2026-03-17)

### Features

- add exclude_authors option to filter bots from contributors ([3171a70](https://github.com/somaz94/go-changelog-action/commit/3171a70292ee2d904939998a7c273400a73ad2c5))

### Documentation

- add exclude_authors usage and default bot filtering to README ([d899e07](https://github.com/somaz94/go-changelog-action/commit/d899e071d95b9336da83b6aa080cd727b3d76b42))

### Continuous Integration

- use somaz94/contributors-action@v1 for contributors generation ([f59cadc](https://github.com/somaz94/go-changelog-action/commit/f59cadcd31e2110dd9dd26e3e8f67e029053ab7e))

### Contributors

- somaz

<br/>

## [v1.0.4](https://github.com/somaz94/go-changelog-action/compare/v1.0.3...v1.0.4) (2026-03-16)

### Features

- add br separator between version entries in changelog ([191c649](https://github.com/somaz94/go-changelog-action/commit/191c649f57e7a2f156cd66bbdf071c34747ee081))

### Contributors

- somaz

<br/>

## [v1.0.3](https://github.com/somaz94/go-changelog-action/compare/v1.0.2...v1.0.3) (2026-03-16)

### Code Refactoring

- improve code quality and consistency ([2015dab](https://github.com/somaz94/go-changelog-action/commit/2015dab65335a4095752334766b83489e0fe2afa))

### Continuous Integration

- use major-tag-action for version tag updates ([69f4ffa](https://github.com/somaz94/go-changelog-action/commit/69f4ffae0f378a2aacb8d9e0184029f3555f8306))

### Contributors

- somaz

<br/>

## [v1.0.2](https://github.com/somaz94/go-changelog-action/compare/v1.0.1...v1.0.2) (2026-03-16)

### Bug Fixes

- improve error handling, deduplicate issues, and raise test coverage to 95% ([4af6c43](https://github.com/somaz94/go-changelog-action/commit/4af6c43f5ed3bae9a3904c9edded563b57493be8))
- changelog-generator.yml ([2fbc07a](https://github.com/somaz94/go-changelog-action/commit/2fbc07a20696c42d041b69c45d1717de076c5586))

### Continuous Integration

- expand test coverage in CI and smoke tests ([7a5865d](https://github.com/somaz94/go-changelog-action/commit/7a5865d67111d0bd093b17f99e50b14252140532))
- add release config and contributors workflow ([a01850b](https://github.com/somaz94/go-changelog-action/commit/a01850b60626a81d4e81f0760b85d47c439511b4))

### Contributors

- somaz

<br/>

## [v1.0.1](https://github.com/somaz94/go-changelog-action/compare/v1.0.0...v1.0.1) (2026-03-16)

### Features

- delete CHANGELOG.md ([3df0a60](https://github.com/somaz94/go-changelog-action/commit/3df0a60d2ff530471df94e45c9f21f34b3771d06))
- delete CHANGELOG.md ([ec7e852](https://github.com/somaz94/go-changelog-action/commit/ec7e852306276f92556a4259f4961011786564c7))

### Bug Fixes

- resolve garbage "Other Changes" entries in changelog ([8c6a9be](https://github.com/somaz94/go-changelog-action/commit/8c6a9be8334f3cb62909b2969d3ac167a43ed8cd))
- action.yml, docs: README.md ([0c3c77e](https://github.com/somaz94/go-changelog-action/commit/0c3c77e7ff4e614e59cdbcf8545169500a1d4466))

### Documentation

- README.md ([4ed9b19](https://github.com/somaz94/go-changelog-action/commit/4ed9b190a83446472270cb5a2f5a5a2d0e637ed0))

### Contributors

- somaz

<br/>

## [v1.0.0](https://github.com/somaz94/go-changelog-action/releases/tag/v1.0.0) (2026-03-16)

### Features

- delete linters ([f6207a0](https://github.com/somaz94/go-changelog-action/commit/f6207a0f72a4a12cdb5fbafb34653f1e5fcf856f))
- add PR/issue links, compare links, contributors, non-conventional commits, tag range, and custom type mapping, closes [#456](https://github.com/somaz94/go-changelog-action/issues/456) ([1fe5642](https://github.com/somaz94/go-changelog-action/commit/1fe5642b33731427a3f2abe58aa9cf4a0a5e254b))
- implement Go-based changelog action with Conventional Commits support ([810e9f6](https://github.com/somaz94/go-changelog-action/commit/810e9f6334a4f441a96f4697dbaf814567384fc4))

### Bug Fixes

- changelog-generator ([edb20ab](https://github.com/somaz94/go-changelog-action/commit/edb20abe0f149cdbfc7569c1263421c9cdd5dc0b))
- update Go to 1.26 and add git safe.directory for container support ([3d32e08](https://github.com/somaz94/go-changelog-action/commit/3d32e08fd31152d23bfa02ffa6450fc750207755))

### Documentation

- add badges, spacing, and contributing section to README ([f8c8f7f](https://github.com/somaz94/go-changelog-action/commit/f8c8f7f8e547f149f3796b251c8352dc3c84e22a))

### Tests

- raise test coverage from 47% to 92% ([7e3c489](https://github.com/somaz94/go-changelog-action/commit/7e3c4891fc324779cee67d2cd95892330f44ae8a))

### Builds

- **deps:** bump the docker-minor group with 2 updates (#1) ([#1](https://github.com/somaz94/go-changelog-action/pull/1)) ([e5f5751](https://github.com/somaz94/go-changelog-action/commit/e5f57511419d540394557b045300aa3308dd425e))

### Contributors

- somaz

<br/>

