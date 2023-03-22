let pkgs = import <nixpkgs> {};

in pkgs.mkShell rec {
  name = "owncast-dev";
  
  buildInputs = with pkgs; [
    nodejs
    yarn
    go
    gcc
    ffmpeg
  ];
}    
