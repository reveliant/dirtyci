# dirty-ci

A tiny and dirty continuous integration script for web projects deployed on a
web server

Receives [GitHub](https://github.com) and [GitLab](https://gitlab.com) webhooks
for selected projects and `git pull` their content into specified directory.

## Requirements

* Web server with PHP
* Git

## Install

1. Create a Deploy Key<sup>[1](#deploy-key)</sup> for the web server account
   (i.e. www-data or http)
2. Put the script and the config file on your web server (I personnaly use a
   specific virtual host with a dedicated subdomain)
  - avoid exposing the config file and the log file to the Internet
  - filenames (and paths) to the config file and the log file can be customized
    in the first lines of the script
3. Change the `git` parameters in config.json
  - `root` is intended to be the web server root, e.g. `/var/www` or
    `/srv/http`, containing your deployed projects
  - `branches.remote` and `branches.local` are default values for `git pull`
    commands


## Adding a new project

1. Add the Deploy Key to your GitHub or GitLab project
2. `git clone` your project on your web server
3. Create a new webhook pointing the script
  - on GitHub, choose `application/json` for content type and check `Just the
    push event`
  - on GitLab, check `Push events`
3. Add the project to the `repositories` section of the config file:
  - the project __must__ have an `remote.url` property, set with the SSH
    URL<sup>[2](#ssh-url)</sup> of remote repository
  - the project __must__ have a `local.url` property, pointing to the project
    directory relatively to the `root` set up before 
  - the project _might_ have `remote.branch` and `local.branch` properties if
    your don't want to use the default values set before.

## Notes

### Deploy Key

Considering we'll only pull from the repository, a read-only access is
sufficent.

This step is not necessary if you use HTTPS access to repository.

### SSH URL

SSH URLs are used as unique, simple-to-find references to map repositories
using informations contained in web hooks, even if you do not use SSH to pull
from the repository.

## License

MIT License

