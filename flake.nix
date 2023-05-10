{
  description = "Shopware CLI";
  inputs.nixpkgs.url = "nixpkgs/nixos-unstable";

  outputs = { self, nixpkgs }:
    let
      version = "0.1.61";
      supportedSystems = [ "x86_64-linux" "x86_64-darwin" "aarch64-linux" "aarch64-darwin" ];
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;
      nixpkgsFor = forAllSystems (system: import nixpkgs { inherit system; });
    in
    {
      packages = forAllSystems (system:
        let
          pkgs = nixpkgsFor.${system};
        in
        rec {
          shopware-cli = pkgs.buildGoModule {
            pname = "shopware-cli";
            inherit version;
            src = ./.;

            nativeBuildInputs = [ pkgs.installShellFiles pkgs.makeWrapper ];

            vendorSha256 = "sha256-abeKokkYV6yrjIJmknd2umwO4sVHds3P0oZqJhifikg=";

            postInstall = ''
              export HOME="$(mktemp -d)"
              installShellCompletion --cmd shopware-cli \
                --bash <($out/bin/shopware-cli completion bash) \
                --zsh <($out/bin/shopware-cli completion zsh) \
                --fish <($out/bin/shopware-cli completion fish)
            '';

            postFixup = ''
              wrapProgram $out/bin/shopware-cli \
                --prefix PATH : ${pkgs.dart-sass-embedded}/bin
            '';

            CGO_ENABLED = 0;

            ldflags = [
              "-s"
              "-w"
              "-X 'github.com/FriendsOfShopware/shopware-cli/cmd.version=${version}'"
            ];
          };
          default = shopware-cli;
        });

      apps = forAllSystems (system: rec {
        shopware-cli = {
          type = "app";
          program = "${self.packages.${system}.shopware-cli}/bin/shopware-cli";
        };
        default = shopware-cli;
      });

      formatter = forAllSystems (
        system:
        nixpkgs.legacyPackages.${system}.nixpkgs-fmt
      );
    };
}
