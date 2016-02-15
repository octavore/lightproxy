default:
	@echo 'example: `make v0.1.0`'

%:
	@GO15VENDOREXPERIMENT=1 go build .
	@tar czf lightproxy-$@.tar.gz lightproxy
	@shasum -a 256 lightproxy-$@.tar.gz

clean:
	@rm *.tar.gz
