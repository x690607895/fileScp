class FlashCat < Formula
    desc "a simple file scp"
    homepage "https://github.com/x690607895/fileScp"
  
    if Hardware::CPU.intel?
      url "https://github.com/x690607895/fileScp/releases/download/release-v1.0.0/fileScp-cli-macos-amd64.tar.gz"
      sha256 "57005806a8645bfbb7c4e745c15c0089840c5293edd73b7a8934e7eb0243b244"
    elsif Hardware::CPU.arm?
      url "https://github.com/x690607895/fileScp/releases/download/release-v1.0.0/fileScp-cli-macos-arm64.tar.gz"
      sha256 "6a63ec2e18d0d291ff260144eea6b4a53f95dc516315614bdc2843dd696a16fc"
    end
  
    version "1.0.0"
  
    def install
      if Hardware::CPU.intel?
        bin.install "flash-cat"
      elsif Hardware::CPU.arm?
        bin.install "flash-cat"
      end
    end
  
  end
  