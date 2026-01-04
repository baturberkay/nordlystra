# Nordlystra

Nordlystra is a lightweight CLI tool written in Go for controlling Philips Hue smart lights.

## Prerequisites

- [Go](https://go.dev/dl/) 1.18 or later
- A Philips Hue Bridge on your local network
- Physical access to the Hue Bridge (to press the link button during setup)

## Installation

```bash
git clone https://github.com/baturberkay/nordlystra.git
cd nordlystra
go build -o nordlystra main.go
```

## Getting Started

On first run, Nordlystra will guide you through the setup process:

1. Run the binary:
   ```bash
   ./nordlystra -list
   ```

2. Enter your Hue Bridge IP address when prompted

3. **Press the link button on your Hue Bridge** when asked, then press Enter

4. The configuration is saved to `./config/hue_config.json` for future use

## Commands

### `-list`

List all available lights and their current states.

```bash
./nordlystra -list
```

### `-light-id`

Specify the ID of the light to control. Required for all state changes.

```bash
./nordlystra -light-id 1 -bri 200 -o true
```

### `-o`

Turn the light on or off. Values: `true` or `false`

```bash
./nordlystra -light-id 1 -o true
```

### `-bri`

Set the brightness value of the light. Range: `0-255`

```bash
./nordlystra -light-id 1 -bri 128
```

### `-hue`

Set the hue value for the light color. Range: `0-65535`

```bash
./nordlystra -light-id 1 -hue 30000
```

### `-sat`

Set the saturation value of the light. Range: `0-255`

```bash
./nordlystra -light-id 1 -sat 150
```

## Examples

Turn on light 1 with full brightness:
```bash
./nordlystra -light-id 1 -o true -bri 255
```

Set light 2 to a warm orange color:
```bash
./nordlystra -light-id 2 -hue 5000 -sat 200 -bri 180
```

Turn off light 3:
```bash
./nordlystra -light-id 3 -o false
```