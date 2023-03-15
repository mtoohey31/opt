{
  description = "opt";

  inputs = {
    nixpkgs.url = "nixpkgs/nixpkgs-unstable";
    utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, utils }:
    utils.lib.eachDefaultSystem (system: with import nixpkgs
      {
        overlays = [
          (final: _: {
            gow = final.buildGoModule rec {
              pname = "gow";
              version = builtins.substring 0 7 src.rev;
              src = final.fetchFromGitHub {
                owner = "mitranim";
                repo = pname;
                rev = "36c8536a96b851631e800bb00f73383fc506f210";
                sha256 = "q56s97j+Npurb942TeQhJPqq1vl/XFe7a2Dj5fw7EtQ=";
              };
              vendorSha256 = "o6KltbjmAN2w9LMeS9oozB0qz9tSMYmdDW3CwUNChzA=";
            };
          })
        ]; inherit system;
      }; {
      devShells.default = mkShell {
        packages = [ go gopls gow revive ];
      };
    });
}
