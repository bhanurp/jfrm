class Jfrm < Formula
  desc "JFrog Release Manager - Manage releases and dependencies for JFrog projects"
  homepage "https://github.com/bhanurp/jfrm"
  version "0.0.1"
  
  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/bhanurp/jfrm/releases/download/v0.0.1/jfrm_darwin_arm64.tar.gz"
      sha256 "YOUR_SHA256_HERE"
    else
      url "https://github.com/bhanurp/jfrm/releases/download/v0.0.1/jfrm_darwin_amd64.tar.gz"
      sha256 "YOUR_SHA256_HERE"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/bhanurp/jfrm/releases/download/v0.0.1/jfrm_linux_arm64.tar.gz"
      sha256 "YOUR_SHA256_HERE"
    else
      url "https://github.com/bhanurp/jfrm/releases/download/v0.0.1/jfrm_linux_amd64.tar.gz"
      sha256 "YOUR_SHA256_HERE"
    end
  end

  def install
    bin.install "jfrm"
  end

  test do
    system "#{bin}/jfrm", "--help"
  end
end 