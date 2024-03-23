release:
	rm -rf dist; \
	sh -c '( \
		export PNAME=`./scripts/provider_name.sh` ; \
		export GPG_FINGERPRINT=`./scripts/key_fingerprint.sh` ; \
		go run github.com/goreleaser/goreleaser release --clean --snapshot; \
	)'
