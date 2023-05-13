{ pkgs, ... }:

{
  packages = [
    pkgs.php82
    pkgs.php82Packages.composer
    pkgs.nodejs-18_x
    pkgs.golangci-lint
    pkgs.gofumpt
    pkgs.gcc
  ];

  languages.php = {
    enable = true;
    package = pkgs.php82;
  };

  languages.go = {
    enable = true;
  };
}
