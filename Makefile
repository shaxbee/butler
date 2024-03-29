BUILD := build
OS := $(shell uname -s)
ARCH := $(shell uname -m)

TEMPL_VERSION := v0.2.513
TEMPL_ROOT := $(BUILD)/bin/templ-$(TEMPL_VERSION)
TEMPL := $(TEMPL_ROOT)/templ

$(TEMPL):
	@mkdir -p $(TEMPL_ROOT)
	curl -fsL https://github.com/a-h/templ/releases/download/$(TEMPL_VERSION)/templ_$(OS)_$(ARCH).tar.gz | tar -xf - -C $(TEMPL_ROOT) templ

TAILWINDCSS_VERSION := v3.3.7
TAILWINDCSS_ROOT := $(BUILD)/bin/tailwindcss-$(TAILWINDCSS_VERSION)
TAILWINDCSS := $(TAILWINDCSS_ROOT)/tailwindcss

$(TAILWINDCSS):
	@mkdir -p $(TAILWINDCSS_ROOT)
ifeq ($(OS),Darwin)
	curl -fsL https://github.com/tailwindlabs/tailwindcss/releases/download/$(TAILWINDCSS_VERSION)/tailwindcss-macos-$(ARCH) -o $(TAILWINDCSS)
endif # OS==Darwin
	chmod a+x $(TAILWINDCSS)

WGO_VERSION := v0.5.4
WGO_ROOT := $(BUILD)/bin/wgo-$(WGO_VERSION)
WGO := $(WGO_ROOT)/wgo

$(WGO):
	GOBIN=$(abspath $(WGO_ROOT)) go install github.com/bokwoon95/wgo@$(WGO_VERSION)

.PHONY: generate generate-templ generate-tailwindcss

generate: generate-templ generate-tailwindcss

generate-templ: $(TEMPL)
	$(TEMPL) generate -path templates/

generate-tailwindcss: $(TAILWINDCSS)
	$(TAILWINDCSS) -i assets/tailwind.css -o assets/dist/styles.css

.PHONY: migrate migrate-schema

.PHONY: dev

dev: generate $(WGO)
	$(WGO) -file 'templates/(.*).templ' -file assets/tailwind.css make generate :: \
	wgo -file '.go' go run ./cmd/server --addr localhost:8080