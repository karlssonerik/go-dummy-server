DOCKER ?= docker
CURL ?= curl
TAR ?= tar
INSTALL ?= install
FIND ?= find
MV ?= mv
SED ?= sed
PWD ?= pwd
TOUCH ?= touch

SWAGGER_UI_VERSION ?= 3.38.0

.PHONY: help

all: models/web openapi/docs

openapi/docs/service.yaml:
	$(INSTALL) -d openapi/docs
	cp openapi/oas.yaml $@

models/web:
	$(DOCKER) run --rm \
		--volume "`pwd`/openapi:/local" \
		--user "$(shell id -u):$(shell id -g)" \
		openapitools/openapi-generator-cli \
		generate \
		--input-spec /local/oas.yaml \
		--global-property models,modelDocs=false \
		--generator-name go \
		--additional-properties packageName=models \
		--output /local/model
	$(INSTALL) -d $@
	$(TOUCH) $@
	$(RM) -r models/web/model_*
	mv openapi/model/* models/web/

swagger-ui-%.tar.gz:
	$(RM) -f "$@"
	$(CURL) -sSLJO "https://github.com/swagger-api/swagger-ui/archive/v$*.tar.gz"

openapi/docs: swagger-ui-${SWAGGER_UI_VERSION}.tar.gz openapi/docs/service.yaml
# Clean any previous built files and ensure the folder exists for tar to extract into.
	$(INSTALL) -d "openapi/docs"
# Extract compiled Swagger UI from source package.
	$(TAR) --extract \
        --file="$<" \
        --gzip \
        --directory=$@ \
        --strip-components=2 \
        "swagger-ui-${SWAGGER_UI_VERSION}/dist"
# Remove unnecessary files (removes ~90% of extracted size).
	cd openapi/docs && $(RM) $@*.map $@*-es-bundle*.js $@swagger-ui.js
# Insert location of OpenAPI spec to be rendered.
	$(SED) -i 's|url: "[^"]*",|url: "{{`{{.EndpointBaseURL}}`}}/docs/swagger/service.yaml",|g' "openapi/docs/index.html"

clean:
	$(RM) -r models/web models/web/model_* openapi/docs openapi/model swagger-ui-*.tar.gz openapi/bundle.yaml

help:
	docker run --rm -v "${PWD}:/local" openapitools/openapi-generator-cli config-help -g go
