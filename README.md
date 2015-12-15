# dirty-ci

A tiny and dirty continuous integration script for web projects deployed on a web server

Receives [GitHub](https://github.com) and [GitLab](https://gitlab.com) webhooks for selected projects and `git pull` their content into specified directory.

## Requirements

* Web server with PHP
* Git

## Install

1. Create a Deploy Key<sup>[*](#deploy-key)</sup> for the web server account (i.e. www-data or http)
2. Put the script and the config file on your web server (I personnaly use a specific virtual host with a dedicated subdomain)
  - avoid exposing the config file and the log file to the Internet
  - filenames (and paths) to the config file and the log file can be customized in the first lines of the script
3. Change the `git` parameters in config.json
  - `root` is intended to be the web server root, i.e. `/var/www` or `/srv/http` (see _Usage_ for more explanations)
  - `branches.remote` and `branches.local` are default values for `git pull` commands


## Usage

1. Add the Deploy Key to your GitHub or GitLab project
2. `git clone` your project on your web server
3. Create a new webhook pointing the script
  - on GitHub, choose `application/json` for content type and check `Just the push event`
  - on GitLab, check `Push events`
3. Add the project to the config file
  - to the corresponding `github` or `gitlab` arrays
  - the project __must__ have an `remote.url` property (the SSH URL<sup>[*](#ssh-key)</sup> to your repository, e.g. `git@github.com:username/project.git`)
  - the project __must__ have a `local.url` property, pointing to the project directory relatively to the `root` set up before 
  - the project _might_ have a `remote.branch` and/or a `local.branch` properties if your don't want to use the default values set before.

## Notes

### Deploy Key

Considering we'll only pull from the repository, a read-only access is sufficent.

### SSH URL

The SSH URL are used as unique references to map repositories using informations contained in web hooks.

Even if you do not use SSH to pull from the repository, set the SSH URL in the configuration file.

## License

MIT License

