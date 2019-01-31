APP=dirtyci

INSTALL_CONF_DIR=/usr/local/etc/$(APP)/
INSTALL_BIN_DIR=/usr/local/bin/
INSTALL_PLUGIN_DIR=/usr/local/lib/$(APP)/

GO_FLAGS=

SOURCES=$(wildcard *.go */*.go)
UTILS_SRC=$(wildcard utils/*.go)
VALID_PLUGIN_DIRS=$(patsubst %/handler.go,%,$(wildcard plugins/*/handler.go))
PLUGINS=$(addsuffix .so,$(VALID_PLUGIN_DIRS))
PLUGINS_SRCS = $(foreach plugin_dir,$(VALID_PLUGIN_DIRS),$(wildcard $(plugin_dir)/*.go))

all: build

.PHONY: build plugins app run clean heroku

build: plugins app

plugins: $(PLUGINS)

app: $(APP)

%.so: PLUGIN_SRCS=$(wildcard $(patsubst %.so,%,$@)/*.go)
%.so: $(PLUGINS_SRCS) $(UTILS_SRC)
	@echo $<
	go build $(GO_FLAGS) -buildmode=plugin -o $@ $(PLUGIN_SRCS)

$(APP): $(SOURCES)
	go build $(GO_FLAGS) -o $@

run:
	./$(APP)

clean:
	rm -f $(APP) $(PLUGINS)
	rm -rf $(DOCKER_BUILD)

install: installdirs
	mkdir -p $(prefix)/etc/$(APP) $(exec_prefix)/lib/$(APP)
	cp $(APP) $(exec_prefix)/bin/
	cp $(PLUGINS) $(exec_prefix)/lib/$(APP)/

installdirs:
	mkdir -p $(prefix)/etc/$(APP) $(exec_prefix)/lib/$(APP)

uninstall:
	rm $(exec_prefix)/bin/$(APP)
	rm -R $(exec_prefix)/lib/$(APP)/*.so
