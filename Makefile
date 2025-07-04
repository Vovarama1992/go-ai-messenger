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
	@for dir in */; do \
		find "$$dir/internal" -type d -name ports | while read portsdir; do \
			echo "📦 $$dir -> найдена ports: $$portsdir"; \
			mocksdir=$$(dirname $$portsdir)/mocks; \
			mkdir -p "$$mocksdir"; \
			for src in $$portsdir/*.go; do \
				filename=$$(basename $$src .go); \
				mockgen -source=$$src -destination=$$mocksdir/$${filename}_mock.go -package=mocks; \
				echo "  ✅ мок для $$src → $${filename}_mock.go"; \
			done; \
		done; \
	done

swagger-init:
	@for dir in */; do \
		path=$$(find $$dir/internal/**/http -type f -name routes.go 2>/dev/null | head -n 1); \
		if [ ! -z "$$path" ]; then \
			echo "📄 Swagger генерируется в $$dir (из $$path)"; \
			pkg_path=$$(echo $$path | awk -F'/internal/' '{print $$2}' | awk -F'/routes.go' '{print $$1"/http/routes.go"}'); \
			cd $$dir && swag init --parseDependency --parseInternal -g internal/$$pkg_path && cd ..; \
		else \
			echo "⚠️  Не найден routes.go в $$dir//internal/**/http/"; \
		fi; \
	done
	@echo "🧹 Чистим LeftDelim / RightDelim из всех docs.go..."
	@find . -type f -name docs.go | while read file; do \
		sed -i '/LeftDelim/d' $$file; \
		sed -i '/RightDelim/d' $$file; \
		echo "✅ Пофикшен $$file"; \
	done