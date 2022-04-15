# go-docker-vscode

Template for working with Go in Docker via [VSCode remote containers](https://code.visualstudio.com/docs/remote/containers).

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/majodev/go-docker-vscode/blob/master/LICENSE)
[![Build and Test](https://github.com/majodev/go-docker-vscode/actions/workflows/build-test.yml/badge.svg)](https://github.com/majodev/go-docker-vscode/actions)

![go starter overview](https://public.allaboutapps.at/go-starter-wiki/go-starter-main-overview.png)

Extracted from [allaboutapps/go-starter](https://github.com/allaboutapps/go-starter).

**ToC**:

- [go-docker-vscode](#go-docker-vscode)
    - [Requirements](#requirements)
    - [Quickstart](#quickstart)
    - [Visual Studio Code](#visual-studio-code)
    - [Building and testing](#building-and-testing)
    - [Uninstall](#uninstall)
  - [Maintainers](#maintainers)
  - [License](#license)

### Requirements

Requires the following local setup for development:

- [Docker CE](https://docs.docker.com/install/) (19.03 or above)
- [Docker Compose](https://docs.docker.com/compose/install/) (1.25 or above)
- [VSCode Extension: Remote - Containers](https://code.visualstudio.com/docs/remote/containers) (`ms-vscode-remote.remote-containers`)

This project makes use of the [Remote - Containers extension](https://code.visualstudio.com/docs/remote/containers) provided by [Visual Studio Code](https://code.visualstudio.com/). A local installation of the Go tool-chain is **no longer required** when using this setup.

Please refer to the [official installation guide](https://code.visualstudio.com/docs/remote/containers) how this works for your host OS and head to our [FAQ: How does our VSCode setup work?](https://github.com/allaboutapps/go-starter/wiki/FAQ#how-does-our-vscode-setup-work) if you encounter issues.

### Quickstart

Create a new git repository through the GitHub template repository feature ([use this template](https://github.com/majodev/go-docker-vscode/generate)). You will then start with a **single initial commit** in your own repository. 

```bash
# Clone your new repository, cd into it, then easily start the docker-compose dev environment through our helper
./docker-helper.sh --up
```

You should be inside the 'service' docker container with a bash shell.

```bash
development@94242c61cf2b:/app$ # inside your container...

# Shortcut for make init, make build, make info and make test
make all

# Print all available make targets
make help
```

### Visual Studio Code

> If you are new to VSCode Remote - Containers feature, see our [FAQ: How does our VSCode setup work?](https://github.com/allaboutapps/go-starter/wiki/FAQ#how-does-our-vscode-setup-work).

Run `CMD+SHIFT+P` `Go: Install/Update Tools` **after** attaching to the container with VSCode to auto-install all golang related vscode extensions.

### Building and testing

Other useful commands while developing your service:

```bash
development@94242c61cf2b:/app$ # inside your container...

# Print all available make targets
make help

# Shortcut for make init, make build, make info and make test
make all

# Init install/cache dependencies and install tools to bin
make init

# Rebuild only after changes to files
make

# Execute all tests
make test
```

Full docker build:

```bash
docker build . -t go-docker-vscode

docker run go-docker-vscode
# Hello World
```

### Uninstall

Simply run `./docker-helper --destroy` in your working directory (on your host machine) to wipe all docker related traces of this project (and its volumes!).

## Maintainers

- [Mario Ranftl - @majodev](https://github.com/majodev)


## License

[MIT](LICENSE) © Mario Ranftl