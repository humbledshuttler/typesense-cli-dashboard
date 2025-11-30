# Typesense CLI

A terminal UI for browsing Typesense collections, viewing metrics, and managing documents. Inspired by [lazydocker](https://github.com/jesseduffield/lazydocker).

## Features

- **Collections Browser**: View all collections, their schemas, and document counts
- **Metrics Dashboard**: Real-time server health and performance metrics
- **Document Viewer**: Browse and inspect documents within collections
- **Keyboard Navigation**: Full keyboard control with intuitive shortcuts
- **Multi-page Interface**: Switch between different views seamlessly

## Installation

```bash
go build -o typesense-cli .
```

Or install directly:

```bash
go install github.com/humbledshuttler/typesense-cli@latest
```

## Usage

```bash
typesense-cli [flags]
```

### Flags

- `--host`: Typesense server host (default: localhost)
- `-p, --port`: Typesense server port (default: 8108)
- `--protocol`: Typesense server protocol (default: http)
- `-k, --api-key`: Typesense API key (default: xyz)
- `-d, --debug`: Enable debug mode
- `-h, --help`: Show help information
- `--version`: Show version information

### Examples

```bash
# Connect to local Typesense instance
typesense-cli

# Connect to remote Typesense server
typesense-cli --host typesense.example.com --port 443 --protocol https --api-key your-api-key

# Enable debug mode
typesense-cli --debug
```

## Keyboard Shortcuts

### Global Navigation
- `F1`: Switch to Collections view
- `F2`: Switch to Metrics view
- `F3`: Switch to Documents view
- `Ctrl+C`: Quit application

### Collections View
- `r`: Refresh collections list
- `Enter`: View documents in selected collection
- Arrow keys: Navigate collections

### Metrics View
- `r`: Refresh metrics
- Auto-refreshes every 5 seconds

### Documents View
- `r`: Refresh documents
- `n`: Next page
- `p`: Previous page
- Arrow keys: Navigate documents

## Views

### Collections View
Displays all collections in your Typesense server. Select a collection to view its schema, field definitions, and document count. Press Enter to browse documents in that collection.

### Metrics View
Shows server health status and performance metrics including:
- Server health status
- Collection count
- Total document count
- Auto-refreshes every 5 seconds

### Documents View
Browse documents within a selected collection. View full document details in the side panel. Navigate through pages using `n` and `p` keys.

## Requirements

- Go 1.21 or later
- Typesense server (local or remote)

## Libraries Used

- [tview](https://github.com/rivo/tview): Terminal UI framework
- [tcell](https://github.com/gdamore/tcell): Terminal cell library
- [typesense-go](https://github.com/typesense/typesense-go): Typesense Go client
- [flaggy](https://github.com/integrii/flaggy): Command-line flag parsing
- [go-errors](https://github.com/go-errors/errors): Enhanced error handling

## License

MIT

