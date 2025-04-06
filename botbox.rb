# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class Botbox < Formula
  desc "A discord bot template generator to help create discord bots quickly and easily. Forget about the boilerplate and focus on what really matters, what your bot will do. Bot Box is built using Golang, Cobra, and Huh, offering an intuitive cli tool to quickly build Discord bot projects. It includes a cog-based architecture, `.env` management, and built-in utilities for automating bot configuration and extension development."
  homepage "https://github.com/choice404/botbox"
  version "2.0.2"
  license "MIT"

  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/choice404/botbox/releases/download/v2.0.2/botbox_Darwin_x86_64.tar.gz"
      sha256 "b31be7762019abb6605c9e3684fa384da1cd690dc19470134f1b2576b0df6d6f"

      def install
        bin.install "botbox"
      end
    end
    if Hardware::CPU.arm?
      url "https://github.com/choice404/botbox/releases/download/v2.0.2/botbox_Darwin_arm64.tar.gz"
      sha256 "a824d62b2d1a53ee1eeb9fc1af7e54a7b2356fc427dd90d69bdba567053bd092"

      def install
        bin.install "botbox"
      end
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      if Hardware::CPU.is_64_bit?
        url "https://github.com/choice404/botbox/releases/download/v2.0.2/botbox_Linux_x86_64.tar.gz"
        sha256 "2b70f712ebcd9f9c310433c3605793574ebbe92d19ddcf87e573b407e4ad5615"

        def install
          bin.install "botbox"
        end
      end
    end
    if Hardware::CPU.arm?
      if Hardware::CPU.is_64_bit?
        url "https://github.com/choice404/botbox/releases/download/v2.0.2/botbox_Linux_arm64.tar.gz"
        sha256 "5ae663a107bc73ecf0e5f25dbfbd7b93d5f0acf4562afc7368acf0739dcfd89f"

        def install
          bin.install "botbox"
        end
      end
    end
  end

  test do
    system "#{bin}/botbox --version"
  end
end
