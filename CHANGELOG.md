# Changelog

All notable changes to this project will be documented in this file.

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

