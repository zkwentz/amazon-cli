# typed: false
# frozen_string_literal: true

# Homebrew formula for amazon-cli
# This formula is automatically updated by GoReleaser on new releases
class AmazonCli < Formula
  desc "CLI for Amazon shopping - orders, returns, purchases, subscriptions"
  homepage "https://github.com/zkwentz/amazon-cli"
  license "MIT"
  version "0.1.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/zkwentz/amazon-cli/releases/download/v#{version}/amazon-cli_#{version}_darwin_arm64.tar.gz"
      sha256 "PLACEHOLDER_SHA256_DARWIN_ARM64"
    else
      url "https://github.com/zkwentz/amazon-cli/releases/download/v#{version}/amazon-cli_#{version}_darwin_amd64.tar.gz"
      sha256 "PLACEHOLDER_SHA256_DARWIN_AMD64"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/zkwentz/amazon-cli/releases/download/v#{version}/amazon-cli_#{version}_linux_arm64.tar.gz"
      sha256 "PLACEHOLDER_SHA256_LINUX_ARM64"
    else
      url "https://github.com/zkwentz/amazon-cli/releases/download/v#{version}/amazon-cli_#{version}_linux_amd64.tar.gz"
      sha256 "PLACEHOLDER_SHA256_LINUX_AMD64"
    end
  end

  def install
    bin.install "amazon-cli"
  end

  test do
    system "#{bin}/amazon-cli", "--help"
  end
end
