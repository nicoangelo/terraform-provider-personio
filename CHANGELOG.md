# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.2.0] - 2023-04-26

### Added

- Add `format` block to allow phone number formatting of dynamic attributes

### BREAKING CHANGES

- `personio_employee` data source provides the employee record in `employee` attribute instead of root level

## [0.1.0] - 2023-03-14

### Added

- Reading all employees with `personio_employees` data source
- Reading a single employee by ID with `personio_employee` data source

[unreleased]: https://github.com/nicoangelo/terraform-provider-personio/compare/v0.2.0...HEAD
[0.2.0]: https://github.com/nicoangelo/terraform-provider-personio/releases/tag/v0.2.0
[0.1.0]: https://github.com/nicoangelo/terraform-provider-personio/releases/tag/v0.1.0
