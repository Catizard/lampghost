{
  description = "NixOS wails development environment";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
  };

  outputs = { self , nixpkgs ,... }: let
    system = "x86_64-linux";
  in {
    devShells."${system}".default = let
      pkgs = import nixpkgs {
        inherit system;
      };
    in pkgs.mkShell {
    	packages = with pkgs; [
        wails
	gtk3
	webkitgtk_4_0
      ];

      shellHook = with pkgs; ''
        export XDG_DATA_DIRS=${gsettings-desktop-schemas}/share/gsettings-schemas/${gsettings-desktop-schemas.name}:${gtk3}/share/gsettings-schemas/${gtk3.name}:$XDG_DATA_DIRS;
        export GIO_MODULE_DIR="${glib-networking}/lib/gio/modules/";
        exec fish
      '';
    };
  };
}
