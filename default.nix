# from https://github.com/tweag/gomod2nix
{ pkgs ? import <nixpkgs> { }, lib ? import <nixpkgs/lib> }:

pkgs.buildGoModule rec {
  vendorSha256 = "pQpattmS9VmO3ZIQUFn66az8GSmB4IvYhTTCFn6SUmo=";
  pname = "Neuralink-Takehome";
  version = "0.3.4";

  src = ./.;

  #  modSha256 = "Gjw1dRrgM8D3G7v6WIM2+50r4HmTXvx0Xxme2fH9TlQ=";

  meta = with lib; {
    description = "Neuralink Takehome project";
    homepage = https://github.com/johnrichardrinehart/Neuralink-Takehome;
    license = licenses.mit;
    maintainers = with maintainers; [ johnrichardrinehart ];
    platforms = platforms.linux ++ platforms.darwin;
  };
}
