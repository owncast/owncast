# Build + Distribute Official Owncast Releases

Owncast is released both as standalone archives that can be downloaded and installed themselves, as well as Docker images that can be pulled from Docker Hub.

The original Docker Hub image was [gabekangas/owncast](https://hub.docker.com/repository/docker/gabekangas/owncast) but it has been deprecated in favor of [owncast/owncast](https://hub.docker.com/repository/docker/owncast/owncast). In the short term both images will need to be updated with new releases and in the future we can deprecate the old one.

## Dependencies

1. Install [Earthly](https://earthly.dev/get-earthly), a build automation tool. It uses our [Earthfile](https://github.com/owncast/owncast/blob/develop/Earthfile) to reproducably build the release files and Docker images.
2. Be [logged into Docker Hub](https://docs.docker.com/engine/reference/commandline/login/) with an account that has access to `gabekangas/owncast` and `owncast/owncast` so the images can be pushed to Docker Hub.

## Build release files

1. Create the release archive files for all the different architectures. Specify the human readable version number in the `version` flag such as `v0.1.0`, `nightly`, `develop`, etc. It will be used to identify this binary when running Owncast. You'll find the archives for this release in the `dist` directory when it's complete.

**Run**:

```bash
earthly +package-all --version="v0.1.0"
```

2. Create a release on GitHub with release notes and Changelog for the version.

3. Upload the release archive files to the release on GitHub via the web interface.

### Tip: For releasing only a single architecture

If you require building only a single architecture you can save some time by specifying the architecture you want to build. For example, if you only want to build the 64bit amd64 Linux version you can run:

**Run**:

```bash
earthly +package --platform="linux/amd64"
```

## Build and upload Docker images

Specify the human readable version number in the `version` flag such as `v0.1.0`, `nightly`, `develop`, etc. It will be used to identify this binary when running Owncast.

Create and push the image to Docker Hub with a list of tags. You'll want to tag the image with both the new version number and `latest`.

**Run**: `earthly --push +docker-all --images="owncast/owncast:0.1.0 owncast/owncast:latest gabekangas/owncast:0.1.0 gabekangas/owncast:latest" --version="v0.1.0"`

Omit `--push` if you don't want to push the image to Docker Hub and want to just build and test the image locally first.

## Update installer script

Once you have uploaded the release archive files and made the new files public and are confident the release is working and available you can update the installer script to point to the new release.

Edit the `OWNCAST_VERSION` in [`install.sh`](https://github.com/owncast/owncast.github.io/blob/master/static/install.sh).

## Final

Once the installer is pointing to the new release number and Docker Hub has new images tagged as `latest` the new version is released to the public.
