# from https://github.com/tweag/gomod2nix
{ pkgs ? import <nixpkgs> { }, lib ? import <nixpkgs/lib> }:

pkgs.buildGoModule rec {
  vendorSha256 = "cCcYG4M07oWkmeJnBy1bWbaoRLrTfoFME8ZUTI5zpXQ=";
  pname = "Neuralink-Takehome";
  version = "0.3.4";

  src = ./.;


  meta = with lib; {
    description = "Neuralink Takehome project";
    homepage = https://github.com/johnrichardrinehart/Neuralink-Takehome;
    license = licenses.mit;
    maintainers = with maintainers; [ johnrichardrinehart ];
    platforms = platforms.linux ++ platforms.darwin;
  };
}
