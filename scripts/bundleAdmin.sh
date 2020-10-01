PROJECT_SOURCE_DIR=$(pwd)
mkdir $TMPDIR/admin 2> /dev/null
cd $TMPDIR/admin
git clone --depth 1 https://github.com/owncast/owncast-admin 2> /dev/null
cd owncast-admin
npm --silent install 2> /dev/null
(node_modules/.bin/next build && node_modules/.bin/next export) | grep info
ADMIN_BUILD_DIR=$(pwd)
cd $PROJECT_SOURCE_DIR
mkdir webroot/admin 2> /dev/null
cd webroot/admin
cp -R ${ADMIN_BUILD_DIR}/out/* .
rm -rf $TMPDIR/admin
