{
  description = "A flake for totipo project spec";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    llm-agents.url = "github:numtide/llm-agents.nix";
    jailed-agents = {
      url = "github:andersonjoseph/jailed-agents";
      inputs.llm-agents.follows = "llm-agents";
    };
  };

  outputs = { nixpkgs, flake-utils, jailed-agents, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs {
          inherit system;
        };
        # Only qualification inputs; the archived history is read by check_spec.py.
        qualificationSource = pkgs.lib.fileset.toSource {
          root = ./.;
          fileset = pkgs.lib.fileset.unions [
            ./spec
            ./vectors
            ./requirements
            (pkgs.lib.fileset.fileFilter
              (file: file.hasExt "go" || builtins.elem file.name [ "go.mod" "go.sum" ])
              ./conformance)
            (pkgs.lib.fileset.fileFilter (file: file.hasExt "py") ./tools)
            ./Makefile
            ./go.work
            ./review/V1_PRE_R16_REVISION_HISTORY.md
          ];
        };
        # Fixed supply-chain input, using nixpkgs' Go dependency derivation.
        # Workspace vendoring preserves the root-level go run commands in Makefile.
        goDependencies = (pkgs.buildGoModule {
          pname = "totipo-spec-dependencies";
          version = "r19";
          src = qualificationSource;
          vendorHash = "sha256-637XdgBok7hRhMw5s9hdY3Mr5jlPXqc1viSYoYh9Iqw=";
          env.GOTOOLCHAIN = "local";
          modBuildPhase = ''
            runHook preBuild
            go work vendor
            runHook postBuild
          '';
        }).goModules;
        specQualification = pkgs.stdenv.mkDerivation {
          pname = "totipo-spec-qualification";
          version = "r19";
          src = qualificationSource;
          nativeBuildInputs = with pkgs; [ go python3 gnumake ];
          strictDeps = true;
          dontConfigure = true;
          buildPhase = ''
            runHook preBuild
            export HOME="$TMPDIR/home"
            export GOCACHE="$TMPDIR/go-cache"
            export GOTMPDIR="$TMPDIR/go-tmp"
            export GOPATH="$TMPDIR/go"
            export GOMODCACHE="$TMPDIR/go-mod-cache"
            export GOENV=off GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off
            export CGO_ENABLED=1 GOFLAGS=-mod=vendor
            export PYTHONDONTWRITEBYTECODE=1
            mkdir -p "$HOME" "$GOCACHE" "$GOTMPDIR" "$GOPATH" "$GOMODCACHE"
            cp -r ${goDependencies} vendor
            go version
            python3 --version
            test -z "$(gofmt -l conformance)"
            make check
            go -C conformance vet ./...
            make race
            make fuzz
            runHook postBuild
          '';
          installPhase = ''
            runHook preInstall
            mkdir -p "$out"
            printf '%s\n' 'Totipo v1/r19 specification qualification passed' > "$out/qualified"
            runHook postInstall
          '';
        };
      in
      {
        checks = pkgs.lib.optionalAttrs pkgs.stdenv.hostPlatform.isLinux {
          spec = specQualification;
        };
        formatter = pkgs.nixpkgs-fmt;
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            gnumake
            gopls
            golangci-lint
            golangci-lint-langserver
            python3
          ];
          packages = [
            (jailed-agents.lib.${system}.makeJailedCodex {
              fwdEnv = [ "GOPATH" "GOBIN" ];
              extraPkgs = with pkgs; [
                go
                gnumake
                gopls
                golangci-lint
                golangci-lint-langserver
                libgcc
                gcc
                python3
              ];
            })
            (jailed-agents.lib.${system}.makeJailedPi {
              fwdEnv = [ "GOPATH" "GOBIN" ];
              extraPkgs = with pkgs; [
                go
                gnumake
                gopls
                golangci-lint
                golangci-lint-langserver
                libgcc
                gcc
                python3
              ];
            })
          ];
        };
      });

  nixConfig = {
    extra-substituters = [ "https://cache.numtide.com" ];
    extra-trusted-public-keys = [ "niks3.numtide.com-1:DTx8wZduET09hRmMtKdQDxNNthLQETkc/yaX7M4qK0g=" ];
  };
}
