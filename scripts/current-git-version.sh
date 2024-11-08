#!/usr/bin/env bash

set -e

git describe --tags --always --dirty
