{
  description = "Publishing tool for minimalist photographers";

  # Nixpkgs / NixOS version to use.
  inputs.nixpkgs.url = "nixpkgs/nixos-unstable";

  outputs =
    { self, nixpkgs }:
    let

      # to work with older version of flakes
      lastModifiedDate = self.lastModifiedDate or self.lastModified or "19700101";

      # Generate a user-friendly version number.
      version = builtins.substring 0 8 lastModifiedDate;

      # System types to support.
      supportedSystems = [
        "x86_64-linux"
        "x86_64-darwin"
        "aarch64-linux"
        "aarch64-darwin"
      ];

      # Helper function to generate an attrset '{ x86_64-linux = f "x86_64-linux"; ... }'.
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;

      # Nixpkgs instantiated for supported system types.
      nixpkgsFor = forAllSystems (system: import nixpkgs { inherit system; });

    in
    {

      # Provide some binary packages for selected system types.
      packages = forAllSystems (
        system:
        let
          pkgs = nixpkgsFor.${system};
          lib = pkgs.lib;
        in
        rec {
          foto = pkgs.buildGoModule {
            pname = "foto";
            inherit version;
            src = ./.;
            vendorHash = "sha256-GiCLg/b+ZF5nAXZh/yIH34yyRFPh2LxEKfXiCp929LI=";

            meta = with lib;{
              homepage = "https://github.com/waynezhang/foto";
              description = "Yet another publishing tool for minimalist photographers";
              license = licenses.mit;
              mainProgram = "foto";
              maintainers = with maintainers; [ pinpox ];
            };
          };

          default = foto;
        }
      );
    };
}
