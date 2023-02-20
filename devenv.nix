{ pkgs, ... }:

{
  packages = [
    pkgs.php81Packages.composer
    pkgs.nodejs-18_x
  ];

  languages.go = {
    enable = true;
    package = pkgs.go_1_20;
  };
}
