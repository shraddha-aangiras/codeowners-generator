# CODEOWNERS Generator

`codeowners-generator` is a CLI tool that automates the creation of a `CODEOWNERS` file for your GitHub repository. It identifies the top contributors based on commit history and assigns them as code owners, ensuring efficient code review and ownership management.

## Features

- Fetches top contributors from a GitHub repository.
- Generates a `CODEOWNERS` file automatically.
- Configurable lookback duration for commit history.
- Supports custom GitHub server URLs (e.g., for GitHub Enterprise).
- CLI-based configuration for flexibility.

## Installation

1. Clone the repository:
    ```bash
    git clone https://github.com/shraddha-aangiras/codeowners-generator.git
    cd codeowners-generator
    ```

2. Build the binary:
    ```bash
    go build -o codeowners-generator ./cmd/codeowners-generator
    ```

3. (Optional) Add the binary to your `PATH` for easier usage.

## Usage

Run the tool with the required flags:

```bash
./codeowners-generator \
  --organization-name <organization> \
  --repository-name <repository> \
  --github-token <token> \
  --duration <lookback-duration> \
  --code-reviewers-count <number-of-reviewers>
```

### Example

```bash
./codeowners-generator \
  --organization-name my-org \
  --repository-name my-repo \
  --github-token ghp_1234567890abcdef \
  --duration 30d \
  --code-reviewers-count 3
```

### Flags

| Flag                   | Description                                                                 | Default Value       |
|------------------------|-----------------------------------------------------------------------------|---------------------|
| `--github-server-url`  | GitHub API server URL (e.g., `https://api.github.com`).                     | `https://api.github.com` |
| `--organization-name`  | GitHub organization name.                                                   | Required            |
| `--repository-name`    | GitHub repository name.                                                     | Required            |
| `--github-token`       | GitHub personal access token.                                               | Required            |
| `--duration`           | Lookback duration for commits (e.g., `2w3d5h`).                             | `30d`               |
| `--code-reviewers-count` | Number of top contributors to assign as code owners.                      | `3`                 |

## How It Works

1. The tool fetches commit history from the specified repository using the GitHub API.
2. It identifies the top contributors based on the number of commits within the specified duration.
3. A `CODEOWNERS` file is generated, assigning the top contributors as code owners.

## GitHub Actions Integration

You can integrate `codeowners-generator` into your CI/CD pipeline using GitHub Actions. See the provided [`release.yml`](./.github/workflows/release.yml) for an example workflow that builds the binary and creates a release.

## Development

### Prerequisites

- Go 1.21 or later
- GitHub personal access token with `repo` scope

### Running Locally

1. Install dependencies:
    ```bash
    go mod tidy
    ```

2. Run the tool:
    ```bash
    go run ./cmd/codeowners-generator \
      --organization-name my-org \
      --repository-name my-repo \
      --github-token ghp_1234567890abcdef \
      --duration 30d \
      --code-reviewers-count 3
    ```

### Testing

Run the tests using:
```bash
go test ./...
```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request with your changes.

## License

This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for details.

## Acknowledgments

- Built using [go-github](https://github.com/google/go-github) for GitHub API integration.
- CLI powered by [urfave/cli](https://github.com/urfave/cli).
