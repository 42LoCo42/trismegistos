{
  inputs.obscura.url = "github:42loco42/obscura";

  outputs = { obscura, ... }:
    obscura.inputs.flake-utils.lib.eachDefaultSystem (system:
      let pkgs = import obscura.inputs.nixpkgs { inherit system; }; in rec {
        defaultPackage = pkgs.buildGoModule {
          pname = "trismegistos";
          version = "0";
          src = ./.;
          vendorHash = "sha256-rFMqxlMXyAibQvCxqJafb8DozLp9McGWLkzB1qkA6KM=";
        };

        devShell = pkgs.mkShell {
          inputsFrom = [ defaultPackage ];
          packages = with pkgs; [
            gopls
            nodePackages.prettier
            redis

            (pkgs.writeShellScriptBin "ppgt-setup" ''
              mkdir -p node_modules
              cd node_modules
              ln -sf "${obscura.packages.${system}.prettier-plugin-go-template}/lib/node_modules/"*
            '')
          ];
        };
      });
}
