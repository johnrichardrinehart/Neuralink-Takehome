# from https://github.com/tweag/gomod2nix
{ pkgs ? import <nixpkgs> { }, lib ? import <nixpkgs/lib> }:
pkgs.buildGoModule rec {
  vendorSha256 = "zXRytszXa0NFIYnN3PDZq2tSiLJbfB+LL+4OAYtcQJ4=";
  pname = "Neuralink-Takehome";
  version = "1.0.7";

  nativeBuildInputs = with pkgs.buildPackages; [ go protobuf protoc-gen-go protoc-gen-go-grpc ];

  preBuild = ''
  ${pkgs.protobuf}/bin/protoc --proto_path=./proto \
	--go_out=./proto --go_opt=Mimage.proto="/;image" \
	--go-grpc_out=./proto --go-grpc_opt=paths=source_relative --go-grpc_opt=Mimage.proto="/;image" \
	./proto/image.proto
  '';

  src = ./.;

  meta = with lib; {
    description = "Neuralink Takehome project";
    homepage = https://github.com/johnrichardrinehart/Neuralink-Takehome;
    license = licenses.mit;
    maintainers = with maintainers; [ johnrichardrinehart ];
    platforms = platforms.linux ++ platforms.darwin;
  };
}
