#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
REPO_ROOT=$(git rev-parse --show-toplevel)

set -e

# If npm is not installed, error out
if ! [ -x "$(command -v npm)" ]; then
  echo "### Error: npm is not installed. Please install node & npm and try again."
  exit 1
fi

# If node_modules is not here, run npm install
if [ ! -d "$DIR/node_modules" ]; then
  echo "### Installing dependencies"
  npm install
fi

echo "### Generating OpenAPI spec from TypeScript"
node $DIR/node_modules/.bin/tsp compile --output-dir $DIR $DIR/main.tsp

echo "### Copying generated spec to swagger-ui"
cp $DIR/@typespec/openapi3/openapi.json $REPO_ROOT/cmd/swagger-ui/openapi.json

echo "### Removing output"
rm -rf $DIR/@typespec