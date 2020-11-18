---
name: Owncast release checklist
about: Use when planning a new public release of Owncast.
title: Owncast release 0.0.x
labels: documentation, Testing
assignees: gabek

---

## Owncast release 0.0.x

The following are intended to address the test scenarios and documentation that should be accomplished for a public release.

- [ ] This release was tested as a fresh install, with no pre-existing data or configuration.
- [ ] This release was tested as an upgrade over a previous install, keeping the data and configuration.
- [ ] This release was tested using a Docker environment.
- [ ] This release was tested with local storage for file distribution.
- [ ] this release was tested with remote S3 storage for file distribution.
- [ ] Release notes and a changelog was written to call out the new features, changes to existing features, and user-impacting behind the scenes updates.
- [ ] Upgrade instructions were written for migrating from the previous release.
- [ ] The API spec `openapi.yaml` is up to date, and the current version of the release is specified in the file.
- [ ] The API documentation was copied to the public documentation site under `/api/0.0.x` and symlinked to `/api/latest`.
