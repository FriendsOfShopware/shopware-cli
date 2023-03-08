{ pkgs, ... }:

{
  packages = [
    pkgs.php81Packages.composer
    pkgs.nodejs-18_x
    pkgs.golangci-lint
  ];

  languages.go = {
    enable = true;
  };
}
