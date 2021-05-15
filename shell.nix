{ pkgs ? import <nixpkgs> { } }:
let
  protoc-gen-go = import ./protoc-gen-go.nix { inherit pkgs; lib = pkgs.lib; };
in
pkgs.mkShell {
  # nativeBuildInputs is usually what you want -- tools you need to run
  nativeBuildInputs = with pkgs.buildPackages; [ go protobuf protoc-gen-go protoc-gen-go-grpc ];
}
