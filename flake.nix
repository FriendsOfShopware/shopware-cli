{
  description = "Shopware CLI";
  inputs.nixpkgs.url = "nixpkgs/nixos-unstable";

  outputs = { self, nixpkgs }:
    let
      version = "0.1.51";
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
          dart-sass-embedded = pkgs.stdenv.mkDerivation rec {
            pname = "dart-sass-embedded";
            version = "1.58.3";

            dontConfigure = true;
            dontBuild = true;

            nativeBuildInputs = pkgs.lib.optionals pkgs.stdenv.hostPlatform.isLinux pkgs.autoPatchelfHook;

            src = pkgs.fetchurl {
              url = {
                "x86_64-linux" = "https://github.com/sass/dart-sass-embedded/releases/download/${version}/sass_embedded-${version}-linux-x64.tar.gz";
                "aarch64-linux" = "https://github.com/sass/dart-sass-embedded/releases/download/${version}/sass_embedded-${version}-linux-arm64.tar.gz";
                "aarch64-darwin" = "https://github.com/sass/dart-sass-embedded/releases/download/${version}/sass_embedded-${version}-macos-arm64.tar.gz";
              }."${pkgs.system}";
              hash = {
                "x86_64-linux" = "sha256-hFhg6FzfJ2ti41YwqvtiDkJ12khWUL5fVKAn/cGlLo8=";
                "aarch64-linux" = "sha256-bYjpOvhjJPXneHc87ZPcsxZpQsOgvZqrknJFyFc67jg=";
                "aarch64-darwin" = "sha256-AihqDuPmDGrjXZV4hYZh//TjWh4L6m5Xqs/18bVgaQw=";
              }."${pkgs.system}";
            };

            installPhase = ''
              mkdir -p $out/bin
              cp -r * $out
              ln -s $out/dart-sass-embedded $out/bin/dart-sass-embedded
            '';
          };

          shopware-cli = pkgs.buildGoModule {
            pname = "shopware-cli";
            inherit version;
            src = ./.;

            nativeBuildInputs = [ pkgs.installShellFiles pkgs.makeWrapper ];

            vendorSha256 = "sha256-i/XZLffThS+/82nBVCzVt4EeFZm552IXQ4sH5FVVTkI=";

            postInstall = ''
              export HOME="$(mktemp -d)"
              installShellCompletion --cmd shopware-cli \
                --bash <($out/bin/shopware-cli completion bash) \
                --zsh <($out/bin/shopware-cli completion zsh) \
                --fish <($out/bin/shopware-cli completion fish)
            '';

            postFixup = ''
              wrapProgram $out/bin/shopware-cli \
                --prefix PATH : ${dart-sass-embedded}/bin
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
