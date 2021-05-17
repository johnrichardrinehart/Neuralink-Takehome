{ pkgs ? import <nixpkgs> { } }:
let
  protoc-gen-go = import ./protoc-gen-go.nix { inherit pkgs; lib = pkgs.lib; };

  # To use this shell.nix on NixOS your user needs to be configured as such:
  # users.extraUsers.adisbladis = {
  #   subUidRanges = [{ startUid = 100000; count = 65536; }];
  #   subGidRanges = [{ startGid = 100000; count = 65536; }];
  # };

  # Provides a script that copies required files to ~/
  podmanSetupScript = let
    registriesConf = pkgs.writeText "registries.conf" ''
      [registries.search]
      registries = ['docker.io']
      [registries.block]
      registries = []
    '';
  in pkgs.writeScript "podman-setup" ''
    #!${pkgs.runtimeShell}
    # Dont overwrite customised configuration
    if ! test -f ~/.config/containers/policy.json; then
      install -Dm555 ${pkgs.skopeo.src}/default-policy.json ~/.config/containers/policy.json
    fi
    if ! test -f ~/.config/containers/registries.conf; then
      install -Dm555 ${registriesConf} ~/.config/containers/registries.conf
    fi
  '';

  # Provides a fake "docker" binary mapping to podman
  dockerCompat = pkgs.runCommandNoCC "docker-podman-compat" {} ''
    mkdir -p $out/bin
    ln -s ${pkgs.podman}/bin/podman $out/bin/docker
  '';

protoPackageSetupScript = pkgs.writeScript "proto-package-setup" ''
  ${pkgs.protobuf}/bin/protoc --proto_path=./proto \
	--go_out=./proto --go_opt=Mimage.proto="/;image" \
	--go-grpc_out=./proto --go-grpc_opt=paths=source_relative --go-grpc_opt=Mimage.proto=/ \
	./proto/image.proto
'';

in
pkgs.mkShell {
  # nativeBuildInputs is usually what you want -- tools you need to run
  nativeBuildInputs = with pkgs.buildPackages; [ go protobuf protoc-gen-go protoc-gen-go-grpc podman ];

  shellHook = ''
    # Install required configuration
    ${podmanSetupScript}
    ${protoPackageSetupScript}
  '';
}
