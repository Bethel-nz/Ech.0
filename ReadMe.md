# Go Port Forwarder

A simple and efficient TCP port forwarder written in Go. It redirects network traffic from a local port to a specified destination address, making it easy to expose services running on one machine to others on the same network or even locally.

## Features

- **TCP Port Forwarding**: Forwards traffic from any local TCP port to a target address.
- **Dynamic Port Allocation**: If a listen port is not specified or is already in use, the tool automatically scans for and selects an available port within a configurable range.
- **Destination Health Check**: Before starting, it tests the connection to the destination address to ensure it's reachable.
- **Network Interface Display**: Lists available network interfaces and their IP addresses, making it easy to access the forwarded port over the LAN.
- **Graceful Shutdown**: Shuts down cleanly on interruption (Ctrl+C), ensuring all connections are properly closed.
- **Concurrent Connection Handling**: Uses Go's concurrency model to handle multiple client connections simultaneously without blocking.

## Installation

1.  **Clone the repository:**

    ```sh
    git clone https://github.com/Bethel-nz/Ech.0.git
    cd Ech.0
    ```

2.  **Build the executable:**
    ```sh
    go build -o forwarder main.go
    ```
    This will create an executable file named `forwarder` (or `forwarder.exe` on Windows).

## Usage

The application is configured via command-line flags.

```
./forwarder [flags]
```

### Configuration Flags

| Flag           | Description                                                                   | Default | Required |
| -------------- | ----------------------------------------------------------------------------- | ------- | -------- |
| `-destination` | The target address to forward traffic to (e.g., `localhost:3000`).            | ""      | Yes      |
| `-listen`      | The TCP port to listen on (e.g., `:8080`). If empty, finds an available port. | ""      | No       |
| `-start-port`  | The beginning of the port range for automatic port discovery.                 | `8000`  | No       |
| `-end-port`    | The end of the port range for automatic port discovery.                       | `9000`  | No       |

### Examples

#### 1. Forward a specific port

Forward traffic from local port `8080` to a service running on `localhost:3000`.

```sh
./forwarder -listen :8080 -destination localhost:3000
```

#### 2. Let the tool find an available port

Forward traffic to `localhost:3000` by automatically finding a free port between `8000` and `9000` to listen on.

```sh
./forwarder -destination localhost:3000
```

The tool will print the listening port it has selected.

#### 3. Use a custom port range for discovery

If the default port range (`8000`-`9000`) is not suitable, you can specify your own.

```sh
./forwarder -destination localhost:3000 -start-port 9100 -end-port 9200
```

## Development

To run the application in development mode without building the executable, you can use `go run`.

```sh
go run main.go -listen :8080 -destination localhost:3000
```

## License

This project is open-source under the MIT License.
