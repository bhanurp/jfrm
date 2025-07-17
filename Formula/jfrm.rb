class Jfrm < Formula
  desc "JFrog Release Manager - Manage releases and dependencies for JFrog projects"
  homepage "https://github.com/bhanurp/jfrm"
  version "0.0.1"
  
  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/bhanurp/jfrm/releases/download/v0.0.1/jfrm_darwin_arm64.tar.gz"
      sha256 "de1ca1268a73255d117d7b7395e17f0b6b37d0c97936338284d779fbd61d8d43"
    else
      url "https://github.com/bhanurp/jfrm/releases/download/v0.0.1/jfrm_darwin_amd64.tar.gz"
      sha256 "bdd8664f7181525c32a2da1b403806f95ff6b892421bd6e7ccf2208ab2aaa4cd"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/bhanurp/jfrm/releases/download/v0.0.1/jfrm_linux_arm64.tar.gz"
      sha256 "6f27d85b260a0845e8af10f0ddf10839b9df864c3a6c6e5dbc89c1b97cffee4f"
    else
      url "https://github.com/bhanurp/jfrm/releases/download/v0.0.1/jfrm_linux_amd64.tar.gz"
      sha256 "575798095430becd4d4b1a88098eb494dfd74450a1b7bf7f6dd6e066947cd1b8"
    end
  end

  def install
    bin.install "jfrm"
  end

  test do
    system "#{bin}/jfrm", "--help"
  end
end 