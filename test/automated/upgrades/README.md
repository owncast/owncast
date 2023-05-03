# Automated Owncast releases and upgrade test

## Upgrades in succession tests

This test will automatically download each release, one after another, to verify they all run in order. It concludes by creating a build from the `develop` branch and running the upgrade to that in order to verify the upcoming release.

## First to last test

This test will automatically download the first release and then upgrade to the development branch to verify that a user can successfully skip versions when upgrading.
