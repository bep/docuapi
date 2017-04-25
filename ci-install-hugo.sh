HUGO_VERSION=0.20.5

set -x
set -e

# Install Hugo if not already cached or upgrade an old version.
if [ ! -e $CIRCLE_BUILD_DIR/bin/hugo ] || ! [[ `hugo version` =~ v${HUGO_VERSION} ]]; then
  wget https://github.com/spf13/hugo/releases/download/v${HUGO_VERSION}/hugo_${HUGO_VERSION}_Linux-64bit.tar.gz
  tar xvzf hugo_${HUGO_VERSION}_Linux-64bit.tar.gz -C $CIRCLE_BUILD_DIR/bin
fi
