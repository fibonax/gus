# Git Uncommitted Scanner (gus)

A command-line tool (CLI) to scan and find Git repositories with uncommitted changes in the directory system.

## 🚀 Features

- Recursively scan all subdirectories to find Git repositories
- Check the status of each found repository
- List repositories with uncommitted changes
- Support JSON output for CI/CD integration
- Simple and user-friendly CLI interface

## 📋 System Requirements

- Go 1.20 or later
- Git installed and available in PATH
- Operating System: Linux, macOS, Windows

## 🔧 Installation

### From source code

```bash
# Clone repository
git clone https://github.com/nguyendangminh/gus.git
cd gus/cmd/gus

# Build and install
go build
```

### Using Go

```bash
go install github.com/nguyendangminh/gus@latest
```

## 💻 Usage

### Basic Syntax

```bash
gus [flags] [path]
```

### Options

```bash
Flags:
  -h, --help       Show help
  -j, --json       Output in JSON format
  -p, --path       Directory path to scan (default: current directory)
  -v, --verbose    Show detailed information
```

### Examples

1. Scan current directory:

```bash
gus
```

2. Scan specific directory:

```bash
gus /path/to/directory
```

3. Output in JSON format:

```bash
gus --json
```

4. Show detailed information:

```bash
gus --verbose
```

## 📝 Output

### Text Format

```
Found 2 repositories with uncommitted changes:

1. /home/user/projects/project-a
   - modified: main.go
   - deleted: old_file.go

2. /home/user/projects/utils/helper
   - new file: helper_test.go
```

### JSON Format

```json
{
  "repositories": [
    {
      "path": "/home/user/projects/project-a",
      "changes": [
        "modified: main.go",
        "deleted: old_file.go"
      ]
    },
    {
      "path": "/home/user/projects/utils/helper",
      "changes": [
        "new file: helper_test.go"
      ]
    }
  ],
  "metadata": {
    "scan_time": "2024-03-20T10:30:00Z",
    "total_repositories": 2,
    "version": "1.0.0"
  }
}
```

## 🧪 Running Tests

```bash
go test ./...
```

## 📦 Project Structure

```
.
├── cmd/
│   ├── gus/        # Entry point
│   └── root/       # Root command
├── pkg/
│   ├── core/       # Core functionality
│   ├── formatter/  # Output formatting
│   ├── git/        # Git operations
│   └── scanner/    # Directory scanning
└── README.md
```

## 🤝 Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a new branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Create a Pull Request

## 📄 License

MIT License - see the [LICENSE](LICENSE) file for details.
