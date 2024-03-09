# Node.js Dependency Updater

Node.js Dependency Updater is a command-line tool written in Go for updating dependencies in Node.js projects based on the package.json file.

## Features

- Update dependencies listed in the `dependencies` and `devDependencies` sections of the package.json file.
- Option to ignore updating dependencies or devDependencies.
- Automatic detection of package.json file in the project directory.
- Color-coded output for easy readability.

## Usage

### Installation

To use the Node.js Dependency Updater, you need to have Go installed on your system.

Clone the repository:

```bash
git clone https://github.com/yourusername/dependency-updater.git
```

Navigate to the project directory:

```bash
cd dependency-updater
```

Build the executable:

```bash
go build
```

### Command-line options

```
Usage: ./dependency-updater [options]

Options:
  --ignore-dependencies       Ignore updating dependencies
  --ignore-devDependencies    Ignore updating devDependencies
  --path                      Path of the Node.js project (default: current directory)
  --help                      Display this help message
```

### Examples

Update dependencies in the current directory:

```bash
./dependency-updater
```

Update dependencies in a specific directory:

```bash
./dependency-updater --path /path/to/project
```

Update dependencies while ignoring devDependencies:

```bash
./dependency-updater --ignore-devDependencies
```

## License

This project is licensed under the [Apache 2.0 License](LICENSE).

## Contributing

Contributions are welcome! Feel free to submit bug reports, feature requests, or pull requests.