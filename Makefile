PROTO_DIR := proto
PROTO_FILES := $(shell find $(PROTO_DIR) -name "*.proto")
GO_OUT := .
PROTOC := protoc

.PHONY: proto generate-mocks

proto:
	$(PROTOC) \
		--go_out=$(GO_OUT) --go_opt=paths=source_relative \
		--go-grpc_out=$(GO_OUT) --go-grpc_opt=paths=source_relative \
		$(PROTO_FILES)

generate-mocks:
	@find . -type d -name ports | while read portdir; do \
		svc=$$(dirname $$portdir); \
		for file in $$portdir/*.go; do \
			grep -E '^type [A-Za-z0-9_]+ interface' $$file | while read line; do \
				intf=$$(echo $$line | awk '{print $$2}'); \
				outdir=$$svc/mocks; \
				outfile=$$outdir/mock_$${intf}.go; \
				if [ -f $$outfile ]; then \
					echo "✅ Mock for $$intf already exists: $$outfile"; \
				else \
					mkdir -p $$outdir; \
					echo "🚧 Generating mock for $$intf into $$outfile"; \
					mockgen "$$(go list -m)/$$portdir" $$intf > $$outfile; \
				fi \
			done; \
		done; \
	done