#!/bin/bash

IMAGES=""

add_image() {
  local image="${1}"

  #podman pull "${image}"
  if [[ -n "${IMAGES}" ]]; then
    IMAGES+=",${image}"
  else
    IMAGES="${image}"
  fi
}

add_image "quay.io/openshift-community-operators/prometheus:v0.14.0"
add_image "quay.io/openshift-community-operators/prometheus:v0.15.0"
add_image "quay.io/openshift-community-operators/prometheus:v0.22.2"
add_image "quay.io/openshift-community-operators/prometheus:v0.27.0"
add_image "quay.io/openshift-community-operators/prometheus:v0.32.0"
add_image "quay.io/openshift-community-operators/prometheus:v0.37.0"
add_image "quay.io/openshift-community-operators/prometheus:v0.47.0"
add_image "quay.io/openshift-community-operators/prometheus:v0.56.3"

add_image "quay.io/openshift-community-operators/redis-operator:v0.0.1"
add_image "quay.io/openshift-community-operators/redis-operator:v0.2.0"
add_image "quay.io/openshift-community-operators/redis-operator:v0.3.0"
add_image "quay.io/openshift-community-operators/redis-operator:v0.4.0"
add_image "quay.io/openshift-community-operators/redis-operator:v0.5.0"
add_image "quay.io/openshift-community-operators/redis-operator:v0.6.0"
add_image "quay.io/openshift-community-operators/redis-operator:v0.8.0"
add_image "quay.io/openshift-community-operators/redis-operator:v0.13.0"

opm index add --bundles="${IMAGES}" --tag=quay.io/apodhrad/iib-test:v0.0.1 -c podman
