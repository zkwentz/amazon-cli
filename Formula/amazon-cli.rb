# typed: false
# frozen_string_literal: true

# Homebrew formula for amazon-cli
class AmazonCli < Formula
  desc "CLI for Amazon shopping - orders, returns, purchases, subscriptions"
  homepage "https://github.com/zkwentz/amazon-cli"
  license "MIT"
  version "0.1.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/zkwentz/amazon-cli/releases/download/v#{version}/amazon-cli_#{version}_darwin_arm64.tar.gz"
      sha256 "6b7cbd4c1f3395ef4c6537a83d7b4bcc563dc2ff21b26e510fd80ab6b1b12213"
    else
      url "https://github.com/zkwentz/amazon-cli/releases/download/v#{version}/amazon-cli_#{version}_darwin_amd64.tar.gz"
      sha256 "bdae49ffac71d50f60d6d9647031f38c3d1b669637650fcc400fe6086637ec0c"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/zkwentz/amazon-cli/releases/download/v#{version}/amazon-cli_#{version}_linux_arm64.tar.gz"
      sha256 "57066a6e26c537c5c8e2df6f52973998e234ffc6e7c6d6cb38ff24df742ec50a"
    else
      url "https://github.com/zkwentz/amazon-cli/releases/download/v#{version}/amazon-cli_#{version}_linux_amd64.tar.gz"
      sha256 "4c722b9092f4a3a6d13ea3a5391da33aeee6a3e447f5952f3361304ad0705e2e"
    end
  end

  def install
    bin.install "amazon-cli"
  end

  test do
    system "#{bin}/amazon-cli", "--help"
  end
end
