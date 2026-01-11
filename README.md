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
go build -o nordlystra .
```

## Getting Started

On first run, Nordlystra will guide you through the setup process:

1. Run the binary:
   ```bash
   ./nordlystra
   ```

2. Enter your Hue Bridge IP address when prompted

3. **Press the link button on your Hue Bridge** when asked, then press Enter

4. The configuration is saved to `./config/hue_config.json` for future use

## Interactive Mode

Run without any flags to enter interactive mode:

```bash
./nordlystra
```

This displays a table of all your devices and provides a command prompt:

```
  Lights in Your Environment
  ------------------------------------------------------------
  ID   | Name             | On/Off | Brightness | Hue   | Saturation
  ------------------------------------------------------------
  1    | Living Room      | On     | 254        | 8417  | 140
  2    | Bedroom          | Off    | 100        | 0     | 0
  ------------------------------------------------------------

nordlystra>
```

### Interactive Commands

| Command | Description |
|---------|-------------|
| `id {id} on` | Turn device on |
| `id {id} off` | Turn device off |
| `id {id} reset` | Reset to brightest white (bri=254, sat=0) |
| `id {id} bri {0-255}` | Set brightness (also turns on) |
| `id {id} hue {0-65535}` | Set hue color (also turns on and sets sat=254) |
| `id {id} sat {0-255}` | Set saturation (also turns on) |
| `help`, `man` | Show help |
| `quit`, `q` | Exit |

### Interactive Examples

```
nordlystra> id 1 on
nordlystra> id 1 bri 200
nordlystra> id 1 hue 5000
nordlystra> id 2 reset
nordlystra> id 1 off
nordlystra> q
```

## CLI Flags

For scripting or one-off commands, you can use flags:

| Flag | Description |
|------|-------------|
| `-list` | List all devices |
| `-id {id}` | Specify device ID |
| `-on` | Turn device on |
| `-off` | Turn device off |
| `-reset` | Reset to brightest white |
| `-bri {0-255}` | Set brightness |
| `-hue {0-65535}` | Set hue |
| `-sat {0-255}` | Set saturation |

### CLI Examples

```bash
# List all devices(press q to exit)
./nordlystra -list 

# Turn on device 1 with full brightness
./nordlystra -id 1 -on -bri 255

# Set device 2 to a warm orange color
./nordlystra -id 2 -hue 5000 -sat 200 -bri 180

# Turn off device 3
./nordlystra -id 3 -off

# Reset device 1 to brightest white
./nordlystra -id 1 -reset
```