# brewls

`brewls` is a command-line interface (CLI) tool that extends Homebrew's `brew ls` functionality by providing a more detailed and human-readable overview of installed Homebrew packages (formulae and casks). It displays the package name, its installed version, and importantly, which other installed packages depend on it (the "Installed By" information). Root packages, i.e., those installed on request and not as dependencies of other installed packages, are marked with an asterisk (`*`).

## 🚀 Features

*   **Detailed Package Information:** View installed version for both formulae and casks.
*   **Reverse Dependency Tracking:** See which other packages rely on a specific installed package.
*   **Root Package Identification:** Easily identify top-level packages that were installed directly by you.
*   **Tabular Output:** Presents information in a clean, easy-to-read table format.

## 🛠️ Installation

### Prerequisites

*   Go (version 1.22 or higher recommended)
*   Homebrew (https://brew.sh/)

### Steps

1.  **Clone the repository:**
    ```bash
    # You would clone from GitHub once the project is hosted.
    # For local development, ensure you are in the project root.
    git clone https://github.com/yourusername/brewls.git # Placeholder
    cd brewls
    ```

2.  **Install the application:**
    ```bash
    go install ./cmd/brewls
    ```
    This will install the `brewls` executable to your Go bin directory, making it available in your PATH.

## 🏃‍♀️ Usage

Simply run `brewls` in your terminal:

```bash
brewls
```

### Example Output Comparison

Compare the default `brew ls` output with the enhanced output from `brewls`.

**`brew ls` Output:**

```
==> Formulae
ada-url                 fmt                     icu4c@78                libunistring            node                    starship
awscli                  fnm                     jpeg-turbo              libuv                   openssl@3               tree
bottom                  fzf                     lazygit                 little-cms2             pcre2                   tree-sitter@0.25
brotli                  gemini-cli              libgit2                 llhttp                  python@3.13             unibilium
bun                     gettext                 libiconv                lpeg                    python@3.14             utf8proc
c-ares                  gh                      libnghttp2              luajit                  readline                uvwasi
ca-certificates         git                     libnghttp3              luv                     ripgrep                 xz
certifi                 go                      libngtcp2               lz4                     simdjson                zstd
deno                    hdrhistogram_c          libssh2                 mpdecimal               specify
eza                     httpie                  libtiff                 neovim                  sqlite

==> Casks
bitwarden               ghostty                 oracle-jdk              rectangle               vlc                     zoom
claude-code             microsoft-edge          pycharm-ce              surfshark               webstorm
flameshot               notion                  raycast                 visual-studio-code      yaak
```

**`brewls` Output:**

```
--- Homebrew Formulae ---
+-----------------+----------+--------------------+
| Name            | Version  | Installed By       |
+-----------------+----------+--------------------+
| awscli          | 2.15.22  |                    |
| go *            | 1.22.0   |                    |
| node            | 20.11.1  |                    |
| openssl@3       | 3.2.1    | node               |
| python@3.13 *   | 3.13.0   | awscli             |
| tree *          | 2.1.1    |                    |
+-----------------+----------+--------------------+

--- Homebrew Casks ---
+--------------------+------------+--------------+
| Name               | Version    | Installed By |
+--------------------+------------+--------------+
| Bitwarden *        | 2023.12.0  |              |
| Visual Studio Code | 1.86.0     |              |
+--------------------+------------+--------------+
```

## 🤝 Contributing

Contributions to `brewls` are welcome! Please follow these guidelines:

1.  **Branching Strategy:** We utilize a Trunk-Based Development approach. All development occurs directly on the `main` branch.
2.  **Pull Request Process:**
    *   Features should be developed behind feature flags to allow safe merging into `main`.
    *   Pull requests should be descriptive and clearly outline the changes.
    *   Code will be reviewed before merging.
3.  **Coding Style Guidelines:**
    *   All Go code must be formatted using `gofmt`.
    *   We recommend running `golangci-lint` to catch potential issues.
4.  **Issue Tracking:** We currently do not use a formal issue tracker. Please submit bug reports or feature requests via GitHub discussions or pull requests.

## 📄 License

This project is licensed under the [MIT License](LICENSE). (Note: A LICENSE file needs to be added.)
```