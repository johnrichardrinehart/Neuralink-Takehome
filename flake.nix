{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";
    gomod2nix.url = "github:tweag/gomod2nix";
  };

  outputs = { self, nixpkgs, gomod2nix, ... }:
    #    let pkgs = import inputs.nixpkgs;
    let pkgs = import nixpkgs { system = "x86_64-linux"; };
    in
    {
      defaultPackage.x86_64-linux = with nixpkgs.legacyPackages.x86_64-linux;
        (import ./.) { inherit pkgs; lib = nixpkgs.lib; };

      devShell.x86_64-linux = (import ./shell.nix) { inherit pkgs; };

    };
}
