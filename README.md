# Neuralink Takehome Retrospective

## How to Build
This repository supports a number of build paths for different use cases. They will be described using
the following language and are ordered by most pure to least pure:
range from most pure to least pure using the following terms:
1. `nix-flake` (or the [`nix` "flake"](https://nixos.wiki/wiki/Flakes) build mode)
1. `nix-legacy` (or the `nix` legacy build mode)
1. `docker`(`podman` is used to avoid the need to use `root` during the build process)
1. `host`

The `./setup.sh` and `./build.sh` scripts support all 4 modes of building based on a fresh Ubuntu 18.04 instance.

### `nix flake`
### Why?
`nix` builds rely on the [`nix` package manager](https://nixos.org/explore.html). The `nix build` command relies on a `flake.nix`
in the root of the repository which, along with a `flake.lock` file (if present) fetches the proper versions (hashed) of all
dependencies and builds them in an environment with minimal path and no I/O (no file copying or network access).

All dependencies are resolved using the `nix` store  (`/nix/store`) so, depending on the software being built, it can easily 
reference multiple versions of a single dependency, since all depenendencies are hashed based on their contents.

In the future, `nix` plans to have a content-addressable store hosted on `IPFS` so no central repository will supply the necessary 
dependencies.

The `nix build` command implicitly uses the `flake.nix` file in the root of the repository, similar to how a plain `Dockerfile` would be 
implicitly picked up by a run of `docker build`/`podman build` from the working directory.

#### How?
    ./setup.sh nix-flake # ensure the user has `sudo` permissions (`wheel` group)
    ./build.sh nix-flake

#### Where?
The binaries should be soft-linked to `./result/bin/{server,client}

### `nix-legacy`
#### Why?
`nix-legacy` uses the `default.nix` in the root of the repository. Builds are not cached in `legacy` `nix` and builds are not fully 
reproducible since the dependency hashes are not referenced at build-time (`flake.lock`).

However, the build process is more reproducible than `Docker` (and definitely `host`) since every dependency version is pinned to 
a particular commit of the `github.com/NixOS/nixpkgs` repository.

#### How?
    ./setup.sh nix-legacy # ensure the user has `sudo` permissions (`wheel` group)
    ./build.sh nix-legacy

#### Where?
The binaries should be soft-linked to `./result/bin/{server,client}

### `docker`
The `docker` setup and build path actually uses `podman` instead of `docker`, since it doesn't require root permissions to create and manage images or run and control containers.

#### How?
    sudo ./setup.sh docker # sudo is needed to install `podman` into the system directories
    ./build.sh docker # note no sudo

#### Where?
The binaries should be located at `/etc/nl-{client,server}`

### `host`
`host` uses the standard `apt` package manager and build tools for the `Go` programming language to build the `server` and `client` binaries. Dependencies are stored in `go.mod` and `go.sum`.

Builds are somewhat reproducible since the `go` compiler and the project dependencies are referenced explicitly in `setup.sh`, `go.mod`, and `go.sum`. But, nothing prevents someone from changing the compiler version in `setup.sh`. This could break the build.

#### How?
    sudo ./setup.sh [host] # `host` is optional since this is the default setup type
    ./build.sh [host]

#### Where?
The binaries should be located at `/etc/nl-{client,server}`