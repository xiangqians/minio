#!/bin/bash

# beta
docker build \
  -f Dockerfile \
  -t minio:2026-08-beta .

# release
