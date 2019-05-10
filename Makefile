default:
	@echo 'example: `make v0.1.0`'

%:
	@go build .
	@tar czf lightproxy-$@.tar.gz lightproxy
	@shasum -a 256 lightproxy-$@.tar.gz

clean:
	@rm *.tar.gz

release:
	goreleaser

release-nr:
	goreleaser --skip-publish
