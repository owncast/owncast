---
name: Owncast release checklist
about: A template with steps required for releasing a build of Owncast
title: Owncast release 0.0.x
labels: documentation, Testing
assignees: gabek

---

## Owncast release 0.0.x

The following are intended to address the test scenarios and documentation that should be accomplished for a public release.

- [ ] This release was tested as a fresh install, with no pre-existing data or configuration.
- [ ] This release was tested as an upgrade over a previous install, keeping the data and configuration.
- [ ] This release was tested using a Docker environment.
- [ ] Release notes and a changelog was written to call out the new features, changes to existing features, and user-impacting behind the scenes updates.
- [ ] Upgrade instructions were written for migrating from the previous release.
- [ ] The API spec `openapi.yaml` is up to date, and the current version of the release is specified in the file.
- [ ] The API documentation was copied to the public documentation site under `/api`.
