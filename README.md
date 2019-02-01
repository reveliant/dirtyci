# dirtyci

A tiny and dirty continuous integration script to deploy web projects

Receives webhooks from [GitHub](https://github.com) or [GitLab](https://gitlab.com)
(or any other solution via a new plugin) for selected projects and pull their
content into specified directory.

## Requirements

* Go environment
* `libgit2`

## Build

First, `go get github.com/reveliant/dirtyci` then :
- either build directly with Go :
  ```sh
  go build github.com/reveliant/dirtyci
  go build -buildmode=plugin github.com/reveliant/dirtyci/plugins/github
  go build -buildmode=plugin github.com/reveliant/dirtyci/plugins/gitlab
  ```
- or build with `make` (plugins are built in the `plugins` directory) :
  ```sh
  cd $GOPATH/src/github.com/reveliant/dirtyci; make
  ```

## Usage

1. Create a Deploy Key for the dirtyci process owner (not necessary if you
   indend to only fetch via HTTPS access).
2. Create a configuration file (see Configuration section below)
Configure dirtyci (example configuration files are available in the repository)
3. Start dirtyci. The following options are supported:
   * `-c`: configuration file to use. Default: `config.toml`
   * `-d`: enable Gin debug mode
   * `-host IFACE`: interface to listen to (address or hostname).
     Default: `127.0.0.1`
   * `-port PORT`: port to listen to. Default: `26979` (`ci` string in decimal)

## Add a new project

1. Add the Deploy Key to your project
    - A read-only access is sufficent
    - Not necessary if you only fetch repository via HTTPS
2. `git clone` your project on your web server
3. Create a new webhook pointing the script
    - on GitHub, choose `application/json` for content type and check `Just the
     push event`
    - on GitLab, check `Push events`
3. Add the project to the `repositories` section of the config file
   (see Configuration section below)

## Configuration

The configuration file, either at [JSON](https://json.org/),
[YAML](https://yaml.org/) or [TOML](https://github.com/toml-lang/toml) format,
has 3 sections:

### Plugins

```toml
pluginsDir = "./plugins"

[plugins]
github = "/github"
gitlab = "/gitlab"
```
`pluginsDir` define path to the directory containing plugins (`.so` files).

If the path is relative, it will be to the current working dir. Therefore, it
is recommanded to rather use absolute path.

Each plugin in the `[plugins]` section will be loaded and set as a
new handler for GET and POST requests on the specified URI. Here :
* `github` plugin (`github.so`) will handle `GET /github` and `POST /github`
  requests
* `gitlab` plugin (`github.so`) will handle `GET /gitlab` and `POST /gitlab`
  requests


### Defaults settings for repositories

```toml
[defaults]
publicKeyPath = "$HOME/.ssh/id_rsa.pub"
privateKeyPath = "$HOME/.ssh/id_rsa"
remoteName = "origin"
remoteBranch = "master"
localBranch = "master"
```
(Harcoded default settings are indicated above, but note that **variables are
not supported**)

* `publicKeyPath` and `privateKeyPath` : paths to SSH public and private key
* `remoteName`: name of Git remote (`git clone` set it to `origin` by default)
* `remoteBranch`: name of Git remote branch (`git clone` will clone `master` by
  default)
* `localBranch`: name of Git local branch to merge into
  (**currently not supported**)

## Repositories section

```toml
[[repositories]]
remoteUrl = "https://github.com/reveliant/dirtyci"
localUrl = "/var/www/dirtyci"
```

Additionnal to settings already defined in the `[defaults]` section (which would
then be overwritten), the following settings are required for each repository:
* `remoteUrl`: unique identifier for external repository (read by the plugin)
* `localUrl`: local path to the cloned repository to merge into

Several repositories entries can have the same `localUrl` (for example to sync
both from GitHub and Gitlab, just set the `remoteName` accordingly)

## Notes

### Deploy Key

Considering we'll only pull from the repository, a read-only access is
sufficent.

This step is not necessary if you use HTTPS access to repository.

## License

MIT License
