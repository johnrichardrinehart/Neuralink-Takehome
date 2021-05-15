{ pkgs ? import <nixpkgs> { }, lib ? import <nixpkgs/lib> }:
pkgs.buildGoModule rec {
  vendorSha256 = "yb8l4ooZwqfvenlxDRg95rqiL+hmsn0weS/dPv/oD2Y=";
  pname = "protoc-gen-go";
  version = "1.26.0";

  src = pkgs.fetchFromGitHub {
    owner = "protocolbuffers";
    repo = "protobuf-go";
    rev = "f2d1f6cbe10b90d22296ea09a7217081c2798009"; # v1.26.0
    sha256 = "n2LHI8DXQFFWhTPOFCegBgwi/0tFvRE226AZfRW8Bnc=";
  };

  subPackages = [ "cmd/protoc-gen-go" ];
  meta = with lib; {
    description = "Go support for Google's protocol buffers";
    homepage = https://github.com/protocolbuffers/protobuf-go;
    maintainers = with maintainers; [ johnrichardrinehart ];
    platforms = platforms.linux ++ platforms.darwin;
  };
}
