# Neuralink Takehome Retrospective

## How to Build
This repository supports a number of build paths for different use cases. They will be described using
the following language and are ordered by most pure to least pure:
range from most pure to least pure using the following terms:
1. `nix-flake` (or the [`nix` "flake"](https://nixos.wiki/wiki/Flakes) build mode)
1. `nix-legacy` (or the `nix` legacy build mode)
1. `docker`(`podman` is used to avoid the need to use `root` during the build process)
1. `host`

The `./setup.sh` and `./build.sh` scripts support all 4 modes of building based on a fresh Ubuntu 18.04 instance (tested on AWS AMI `ami-090717c950a5c34d3`).

### Quickstart
I encourage you to read the description of all `setup` and `build` paths that can be taken, below. But, if you're in a rush you can simply perform the following steps:

1. Spin up an `Ubuntu 18.04` instance (I used an AWS `ami-090717c950a5c34d3` image with 40GiB of disk space, the default 8GiB is probably too small for the `docker` build *and* a couple of the other builds and disk space is relatively cheap)
2. Run the following lines in the shell

        git clone https://github.com/johnrichardrinehart/Neuralink-Takehome nlth
        cd nlth
        sudo sh -c "./setup.sh && ./build.sh" # this executes the `host` setup and build, described below

from the root of the repository and `nl-client` (the client) and `nl-server` (the server) should be visible in the working directory.

You can use the `--help` or `-h` flag for CLI documentation.

### `nix flake`
### Why?
`nix` builds rely on the [`nix` package manager](https://nixos.org/explore.html). The `nix build` command relies on a `flake.nix`
in the root of the repository which, along with a `flake.lock` file (if present) fetches the proper versions (hashed) of all
dependencies and builds them in an environment with minimal path and no I/O (no file copying or network access).

All dependencies are resolved using the `nix` store  (`/nix/store`) so, depending on the software being built, it can easily 
reference multiple versions of a single dependency, since all dependencies are hashed based on their contents.

In the future, `nix` plans to have a content-addressable store hosted on `IPFS` so no central repository will supply the necessary 
dependencies.

The `nix build` command implicitly uses the `flake.nix` file in the root of the repository, similar to how a plain `Dockerfile` would be 
implicitly picked up by a run of `docker build`/`podman build` from the working directory.

#### How?
    ./setup.sh nix-flake # ensure the user has `sudo` permissions (`wheel` group)
    ./build.sh nix-flake

#### Where?
The binaries should be soft-linked to `./result/bin/{server,client}`

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
The binaries should be soft-linked to `./result/bin/{server,client}`

### `docker` (warning: slow build time)
The `docker` setup and build path actually uses `podman` instead of `docker`, since it doesn't require root permissions to create and manage images or run and control containers.

#### How?
    sudo ./setup.sh docker # sudo is needed to install `podman` into the system directories
    ./build.sh docker # IMPORTANT: no sudo

#### Where?
You can check if the images were built successfully by `podman image ls | grep -e "nl-"`. You should see
two images: `nl-client` and `nl-server`.

To run them together they need to be in the same network. This is challenging on Darwin, but easily accomplished on `linux` by binding the containers to the host network, as follows:

    podman run -d --network=host nl-server --debug # running in the background
    podman run --network=host nl-client --help

### `host`
`host` uses the standard `apt` package manager and build tools for the `Go` programming language to build the `server` and `client` binaries. Dependencies are stored in `go.mod` and `go.sum`.

Builds are somewhat reproducible since the `go` compiler and the project dependencies are referenced explicitly in `setup.sh`, `go.mod`, and `go.sum`. But, nothing prevents someone from changing the compiler version in `setup.sh`. This could break the build.

#### How?
    sudo ./setup.sh [host] # `host` is optional since this is the default setup type
    ./build.sh [host]

#### Where?
The binaries should be located at `./nl-{client,server}` (soft-linked from `/tmp`)

## Conclusion
If I had more time I would have:
1. Written this in `rust` (still learning) - there are some great [image codec libraries for `rust`](https://github.com/image-rs/image) and generally this language produces fast code with no GC overhead like `go`.
1. Used the C FFI for `go` to link into a highly-performance and stress-tested image codec library (like [libjpeg](https://github.com/libjpeg-turbo/libjpeg-turbo))
1. Processed each request from clients in separate goroutines so multiple images could be processed in parallel
1. Set up a (in-memory or external) queue so the `server` could process requests asynchronously
1. Added a `ID` field to the responses so `client`'s could request the status of their message processing
1. Added `prometheus` metrics so we could hook into `grafana` if we wanted
1. Done some `pprof` analysis to see where my code was spending a lot of time or consuming excessive memory
1. Done more extensive tests with extremely large images (hundreds of MiB) to demonstrate a solution that doesn't require the entire image to be in memory during processing
1. Cleaned up some of the loops so we could take advantage of a lot of repeated operation (the mean filter "stencil" added and divided a lot of the same numbers)
1. Implemented structured logging so that this could be more production ready
1. Developed an ASIC (prototyping on an FPGA) that would perform these operations in hardware to save on energy and time (if the scale of the problem demanded it)
1. Flushed runtime statistics and logs to a database so we would have more visibility into the application's behavior over long periods of time
1. Added a health endpoint (and embed some various system statistics in the health endpoint JSON response) so we could trigger autoscaling logic or Ops alerts
1. Added more unit tests to cover expected failure cases and more complex scenarios (95.7% of the `server` package is tested, currently)
1. Added an authentication layer so that only authorized people could perform services within their scope
1. Added retry and customizable timeout behavior for the client
1. Added color detection (if most images are grayscale then this feature is a big network I/O saver)
1. ... and more, but that's a good start to a list of unordered things I would have done if this project was mission-critical

## Design Considerations
### Logging
I considered adding more detailed log information, but it ultimately wasn't helpful in diagnosing 
any problems. I would have added more logging if I needed more visibility into the behavior of the application at runtime.

### Packaging
I implemented 4 different packaging systems (the above `setup` and `build` instructions). I chose to
do this because I think in a corporate/production environment reproducibility and environmental consistency in building is paramount. This is one reason you've seen monolith declarative projects like `k8s` and `bazel` take center stage. However, `nix` is an ops-deployment system (`NixOps`) a day-to-day operating system (`NixOS`) a package manager (`nix`) and a functional programming language (`nix-lang`) that deserves more attention from the broader software community because of the problems it solves, particularly in its ability to declaratively manage dependencies, offer fully-reproducible builds, flexibly construct an entire OS using its semantics (so flexible enough for any software project), and for its ability to manage multiple version of the same dependency (which can be useful for build environments of large software systems whose components parts may develop/mature at different rates).

## GitHub Actions
I also implemented a build system that would make release binaries available to the Neuralink team in case there were any problem building my software. It triggers builds on each push (done with `nix` flakes since it's the most reproducible) and pushes binary releases on each tag event. This is a simple, but nice thing to have and can easily be plugged into AWS ECS/similar to automate deployments. The `docker` build of this repo can similarly be pushed to ECR/similar if the deployment environment is container-based. Neat :).

The component of the GitHub Action that runs the `nix` flake executes all tests and compiles the code. So, any issue with my tests or in being able to compile will be visible as a failure to build in the GitHub Actions history.

## Conclusion
Thanks for taking the time to consider my application and review my code. I hope it was enjoyable for you and I hope to hear from you soon!