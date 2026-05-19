# go-filament-sync

A single binary that syncs your custom filament profiles from your slicer directly to your Creality printer over the network — no printer-side installation required.

Run it as a post-processing script in your slicer and it will sync automatically every time you slice and export G-code.

> **Found a bug or have a suggestion?** Please [open an issue](https://github.com/zaggash/go-filament-sync/issues) — all feedback welcome!

## Table of Contents

- [Download](#download)
- [Quick Start](#quick-start)
- [Run as post-processing script in your slicer](#run-as-post-processing-script-in-your-slicer)
- [Creating custom filament presets (Creality Print)](#creating-custom-filament-presets-creality-print)
- [RFID to CFS Android App](#rfid-to-cfs-android-app)
- [How to build locally (Docker)](#how-to-build-locally-docker)

## Download

Current stable release:
<https://github.com/zaggash/go-filament-sync/releases/latest>

Latest dev pre-releases are also available on the same releases page.

## Quick Start

1. **Download** the binary for your platform from the [Releases page](https://github.com/zaggash/go-filament-sync/releases/latest).
2. **Find your profile path** — use the table in the [Usage & flags](#usage--flags) section below to find the correct path for your slicer and operating system.
3. **Configure** — Windows: edit the variables in `filament-sync-tool.bat`; Linux/macOS: copy the one-liner from the [Linux](#linux) or [macOS](#macos) section below.
4. **Add to your slicer** — in your slicer's settings, go to the **Others** tab and find **Post-processing Scripts**. Paste the bat file path (Windows) or the full one-liner (Linux/macOS) into that field.
5. **Slice and export** — the sync runs automatically each time you export G-code.

See the [Windows](#windows), [Linux](#linux), and [macOS](#macos) sections for detailed setup instructions.

## Run as post-processing script in your slicer

The tool runs automatically each time you slice and export G-code. To set it up, find **Post-processing Scripts** at the bottom of the **Others** tab in your slicer's print settings.

### Usage & flags

```
Usage of filament-sync-tool:

  -password string
        Password for SSH connection to printer (default "creality_2024")
  -printer-ip string
        IP address of the Creality printer (required)
  -profile-path string
        Path to slicer filament profile directory (required)
  -user string
        Username for SSH connection to printer (default "root")
```

`--profile-path` is required. Use the table below to find the correct path for your slicer and operating system:

| Slicer | OS | Path |
|--------|----|------|
| OrcaSlicer | Linux | `~/.config/OrcaSlicer/user/default/filament/base` |
| OrcaSlicer | macOS | `~/Library/Application Support/OrcaSlicer/user/default/filament/base` |
| OrcaSlicer | Windows | `%APPDATA%\OrcaSlicer\user\default\filament\base` |
| Creality Print | Linux | `~/.config/Creality/Creality Print/6.0/user/default/filament/base` |
| Creality Print | macOS | `~/Library/Application Support/Creality/Creality Print/6.0/user/default/filament/base` |
| Creality Print | Windows | `%APPDATA%\Creality\Creality Print\6.0\user\default\filament\base` |

Replace `6.0` with your installed Creality Print version. Replace `default` with your user ID if you are logged into the slicer.

### Windows

1. Download both **filament-sync-tool.bat** and **filament-sync-tool.exe** and place them in the same folder (e.g. `C:\Users\YourName\Downloads\`).
2. Open `filament-sync-tool.bat` in a text editor and set the four variables at the top:
   - `PROFILE_PATH` — full path to your slicer's filament profile directory (see the table above; default is the standard OrcaSlicer path)
   - `PRINTER_IP` — your printer's IP address (find it in your printer's network settings)
   - `SSH_USER` — SSH username (default: `root`, usually no change needed)
   - `SSH_PASSWORD` — SSH password (default: `creality_2024`, usually no change needed)
3. Add the following line to the **Post-processing Scripts** field in your slicer (adjust the path to where you saved the file):

```
"C:\Users\YourName\Downloads\filament-sync-tool.bat"
```

### Linux

Paste the following one-liner into the **Post-processing Scripts** field in your slicer. Replace the binary path, IP address, and profile path with your own values.

```
/home/yourusername/Downloads/filament-sync-tool --printer-ip 192.168.1.100 --profile-path ~/.config/OrcaSlicer/user/default/filament/base --user root --password creality_2024
```

Make the binary executable first if needed: `chmod +x /home/yourusername/Downloads/filament-sync-tool`

### macOS

Paste the following one-liner into the **Post-processing Scripts** field in your slicer. Replace the binary path, IP address, and profile path with your own values.

```
/Users/yourusername/Downloads/filament-sync-tool --printer-ip 192.168.1.100 --profile-path ~/Library/Application\ Support/OrcaSlicer/user/default/filament/base --user root --password creality_2024
```

> **Note:** The `Application Support` path contains a space. If your slicer's post-processing runner does not pass arguments through a shell, you may need to enclose the path in double quotes instead: `--profile-path "~/Library/Application Support/OrcaSlicer/user/default/filament/base"`. Test with your specific slicer if the sync does not run.

Make the binary executable first if needed: `chmod +x /Users/yourusername/Downloads/filament-sync-tool`

## Creating custom filament presets (Creality Print)

To sync a custom filament profile, you first need to create it in Creality Print with a special Notes field that the tool reads. Follow these steps:

1. In Creality Print, open the **Filament** section and click **"Set filaments to use"** on the right side.
2. Click **"Custom Filament"** at the top, then click **"Create New"**.
   > If you don't see this option, go to **Preferences** and set **User Role** to **Professional**.
3. Fill in your filament settings (temperature, flow rate, etc.).
4. In the **Notes** field, paste the following JSON and fill in each value:

```json
{"id":"02345","vendor":"Elegoo","type":"PLA","name":"Fast PLA"}
```

| Field | Description |
|-------|-------------|
| `id` | A unique value (conventionally 5 digits, e.g., `02345`). Used as the primary key — the tool silently skips profiles with a missing or empty `id`. Required even if you are not using RFID tags. |
| `vendor` | Filament brand (e.g., `Elegoo`, `Bambu`) |
| `type` | Filament material type (e.g., `PLA`, `PETG`, `ABS`) |
| `name` | Descriptive profile name (e.g., `Fast PLA`) |

The `id` field is the primary key the tool uses to identify and update filament profiles. Pick any unique value — conventionally 5 digits. The tool accepts any non-empty string.

## RFID to CFS Android App

If you use RFID tags with the CFS app, set the same `id` value in the **Material Code** field in the CFS app. The tool uses the `id` field to identify and update filament profiles — if the CFS Material Code and the profile's `id` do not match, the CFS app will not link the RFID tag to the correct profile.

## How to build locally (Docker)

Cross-platform binaries are built using Docker. No local Go installation is required.

```bash
# 1. Build the Docker image
docker build -f ./docker/Dockerfile.build -t filament-sync-tool-builder:latest .

# 2. Create a temporary container to extract the binaries
docker create --name temp_tool_container filament-sync-tool-builder:latest

# 3. Copy the binaries to your current directory
docker cp temp_tool_container:/app/filament-sync-tool_linux_amd64 ./filament-sync-tool_linux_amd64
docker cp temp_tool_container:/app/filament-sync-tool_macos_amd64 ./filament-sync-tool_macos_amd64
docker cp temp_tool_container:/app/filament-sync-tool_macos_arm64 ./filament-sync-tool_macos_arm64
docker cp temp_tool_container:/app/filament-sync-tool_windows_amd64.exe ./filament-sync-tool_windows_amd64.exe

# 4. Remove the temporary container
docker rm -f temp_tool_container
```

---

*Heavily inspired by [HurricanePrint/Filament-Sync](https://github.com/HurricanePrint/Filament-Sync). This implementation works as a single binary with no printer-side dependencies.*
