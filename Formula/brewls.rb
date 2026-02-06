class Brewls < Formula
  desc "Extended Homebrew ls output with installed versions and reverse deps"
  homepage "https://github.com/rahulballal/brewls"
  url "https://github.com/rahulballal/brewls/archive/refs/tags/v0.1.0.tar.gz"
  sha256 "REPLACE_WITH_TARBALL_SHA256"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", "-o", bin/"brewls", "./cmd/brewls"
  end

  test do
    system "#{bin}/brewls", "-h"
  end
end
