{ pkgs, ... }:

{
  packages = [
    pkgs.php81Packages.composer
    pkgs.nodejs-18_x
    pkgs.golangci-lint
    pkgs.gofumpt
  ];

  languages.go = {
    enable = true;
  };
}
