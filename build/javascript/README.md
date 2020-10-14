## Third party web dependencies

Owncast's web frontend utilizes a few third party Javascript and CSS dependencies that we ship with the application.

To add, remove, or update one of these components:

1. Perform your `npm install/uninstall/etc`, or edit the `package.json` file to reflect the change you want to make.
2. Edit the `snowpack` `install` block of `package.json` to specify what files you want to add to the Owncast project.  This can be an entire library (such as `preact`) or a single file (such as `video.js/dist/video.min.js`).  These paths point to files that live in `node_modules`.
3. Run `npm run build`.  This will download the requested module from NPM, package up the assets you specified, and then copy them to the Owncast web app in the `webroot/js/web_modules` directory.
4. Your new web dependency is now available for use in your web code.

## VideoJS versions

Currently Videojs version 7.8.3 and http-streaming version 2.2.0 are hardcoded because these are versions that have been found to work properly with our HLS stream.  Other versions have had issues with things like discontinuities causing a loading spinner.

So if you update videojs or vhs make sure you do an end-to-end test of a stream and make sure the "this stream is offline" ending video displays properly.
