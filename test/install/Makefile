NAME=xeol

IMAGE_NAME=$(NAME)-install.sh-env
UBUNTU_IMAGE=$(IMAGE_NAME):ubuntu-20.04
ALPINE_IMAGE=$(IMAGE_NAME):alpine-3.6
BUSYBOX_IMAGE=busybox:1.36.1-musl

ENVS=./environments
DOCKER_RUN=docker run --rm -t -w /project/test/install -v $(shell pwd)/../../:/project
UNIT=make unit-local

# acceptance testing is running the current install.sh against the latest release. Note: this could be a problem down
# the line if there are breaking changes made that don't align with the latest release (but will be OK with the next
# release)
ACCEPTANCE_CMD=sh -c '../../install.sh -b /usr/local/bin  && xeol version'
# we also want to test against a previous release to ensure that install.sh defers execution to a former install.sh
PREVIOUS_RELEASE=v0.2.3
ACCEPTANCE_PREVIOUS_RELEASE_CMD=sh -c "../../install.sh -b /usr/local/bin $(PREVIOUS_RELEASE) && xeol version"

# CI cache busting values; change these if you want CI to not use previous stored cache
INSTALL_TEST_CACHE_BUSTER=894d8ca

define title
    @printf '\n≡≡≡[ $(1) ]≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡≡\n'
endef

.PHONY: test
test: unit acceptance

.PHONY: ci-test-mac
ci-test-mac: unit-local acceptance-local

# note: do not add acceptance-local to this list
.PHONY: acceptance
acceptance: acceptance-ubuntu-20.04 acceptance-alpine-3.6 acceptance-busybox

.PHONY: unit
unit: unit-ubuntu-20.04

.PHONY: unit-local
unit-local:
	$(call title,unit tests)
	@for f in $(shell ls *_test.sh); do echo "Running unit test suite '$${f}'"; bash $${f} || exit 1; done

.PHONY: acceptance-local
acceptance-local: acceptance-current-release-local acceptance-previous-release-local

.PHONY: acceptance-current-release-local
acceptance-current-release-local:
	$(ACCEPTANCE_CMD)

.PHONY: acceptance-previous-release-local
acceptance-previous-release-local:
	$(ACCEPTANCE_PREVIOUS_RELEASE_CMD)
	xeol version | grep $(shell echo $(PREVIOUS_RELEASE)| tr -d "v")

.PHONY: save
save: ubuntu-20.04 alpine-3.6 pull-busybox
	@mkdir cache || true
	docker image save -o cache/ubuntu-env.tar $(UBUNTU_IMAGE)
	docker image save -o cache/alpine-env.tar $(ALPINE_IMAGE)

.PHONY: load
load:
	docker image load -i cache/ubuntu-env.tar
	docker image load -i cache/alpine-env.tar

## UBUNTU #######################################################

.PHONY: acceptance-ubuntu-20.04
acceptance-ubuntu-20.04: ubuntu-20.04
	$(call title,ubuntu:20.04 - acceptance)
	$(DOCKER_RUN) $(UBUNTU_IMAGE) \
		$(ACCEPTANCE_CMD)

.PHONY: unit-ubuntu-20.04
unit-ubuntu-20.04: ubuntu-20.04
	$(call title,ubuntu:20.04 - unit)
	$(DOCKER_RUN) $(UBUNTU_IMAGE) \
		$(UNIT)

.PHONY: ubuntu-20.04
ubuntu-20.04:
	$(call title,ubuntu:20.04 - build environment)
	docker build -t $(UBUNTU_IMAGE) -f $(ENVS)/Dockerfile-ubuntu-20.04 .

## ALPINE #######################################################

# note: unit tests cannot be run with sh (alpine dosn't have bash by default)

.PHONY: acceptance-alpine-3.6
acceptance-alpine-3.6: alpine-3.6
	$(call title,alpine:3.6 - acceptance)
	$(DOCKER_RUN) $(ALPINE_IMAGE) \
		$(ACCEPTANCE_CMD)

.PHONY: alpine-3.6
alpine-3.6:
	$(call title,alpine:3.6 - build environment)
	docker build -t $(ALPINE_IMAGE) -f $(ENVS)/Dockerfile-alpine-3.6 .

## BUSYBOX #######################################################

# note: unit tests cannot be run with sh (busybox dosn't have bash by default)

# note: busybox by default will not have cacerts, so you will get TLS warnings (we want to test under these conditions)

.PHONY: acceptance-busybox
acceptance-busybox: pull-busybox
	$(call title,busybox - acceptance)
	$(DOCKER_RUN) $(BUSYBOX_IMAGE) \
		$(ACCEPTANCE_CMD)
	@echo "\n*** test note: you should see xeol spit out a 'x509: certificate signed by unknown authority' error --this is expected ***"

.PHONY: pull-busybox
pull-busybox:
	$(call title,busybox - build environment)
	docker pull $(BUSYBOX_IMAGE)

## For CI ########################################################

.PHONY: cache.fingerprint
cache.fingerprint:
	$(call title,Install test fixture fingerprint)
	@find ./environments/* -type f -exec md5sum {} + | awk '{print $1}' | sort | tee /dev/stderr | md5sum | tee cache.fingerprint && echo "$(INSTALL_TEST_CACHE_BUSTER)" >> cache.fingerprint

