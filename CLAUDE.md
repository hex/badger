# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Badger is a cross-platform command-line tool that adds customizable labels/badges to app icons. It's written in C# (.NET 7.0) and uses the SixLabors.ImageSharp library for image manipulation.

## Build and Development Commands

### Prerequisites
- .NET 7.0 SDK
- For macOS release builds: Ruby with Bundler (for Fastlane notarization)

### Common Commands

```bash
# Restore dependencies
dotnet restore

# Build the project
dotnet build

# Build in release mode
dotnet build -c Release

# Run the application
dotnet run -- --text "ALPHA" --icon path/to/icon.png

# Publish for specific platforms (creates self-contained single-file executables)
dotnet publish -c Release -r osx-x64 --self-contained true
dotnet publish -c Release -r osx-arm64 --self-contained true
dotnet publish -c Release -r win-x64 --self-contained true
dotnet publish -c Release -r linux-x64 --self-contained true

# Create NuGet package
dotnet pack -c Release -o .
```

## Code Architecture

### Single-File Design
The entire application logic resides in `Program.cs` - there are no separate classes or modules. This is a deliberate design choice for a simple, focused CLI tool.

### Core Image Processing Flow

1. **Badge Creation** (`CreateBadge` function):
   - Loads the input icon image
   - Creates a blank canvas matching the icon dimensions
   - Calculates badge dimensions as percentages of the icon size
   - Determines badge placement using pivot points (top, bottom, left, right, corners, center)
   - Measures text and scales font to fit within badge dimensions with padding
   - Renders the badge background and text
   - Applies rotation if specified
   - Saves temporary badge overlay to `badgerOutput/` directory

2. **Badge Application** (`AddBadge` function):
   - Loads both the original icon and the generated badge overlay
   - Composites the badge onto the icon at the specified offset with opacity
   - Saves the final result (either overwriting the original or creating a new file in `badgerOutput/`)
   - Cleans up the temporary badge overlay file

3. **Main Entry Point** (`Badger` function):
   - Decorated with ConsoleAppFramework attributes for automatic CLI argument parsing
   - Handles both single-file and directory processing
   - For directories: recursively processes all .png, .jpg, and .jpeg files

### Key Helper Functions

- `GetPivotPoint`: Converts string pivot names to actual coordinates on the image
- `GetTextAlignment`: Maps alignment strings to SixLabors TextAlignment enum
- `SanitizeFileName`: Cleans filenames for output files
- `EmptyDirectory`: Clears the output directory before processing

### Output Behavior

- By default, creates `badgerOutput/` directory and saves badged images there
- With `--overwrite` flag: modifies input files in place and removes `badgerOutput/` directory
- Temporary badge overlays are always cleaned up after compositing

## Project Configuration

### Build Settings (badger.csproj)
- **PublishSingleFile**: Produces a single executable
- **SelfContained**: Bundles .NET runtime with the application
- **PublishTrimmed**: Reduces binary size by removing unused code
- **IncludeNativeLibrariesForSelfExtract**: Embeds native ImageSharp dependencies
- **TrimMode**: Set to "partial" for compatibility with ImageSharp

### Version Management
- Version is defined in `badger.csproj` (currently 2022.12.6)
- GitHub Actions automatically updates version from git tags during release builds

## Release Process

Releases are fully automated via `.github/workflows/main.yml`:

1. Triggered by pushing a git tag matching `v*` pattern (e.g., `v2022.12.6`)
2. Builds for all four platforms (osx-x64, osx-arm64, win-x64, linux-x64)
3. For macOS binaries:
   - Code signs with Developer ID certificate
   - Notarizes via Fastlane with App Store Connect API
4. Creates GitHub release with platform-specific ZIP files
5. Publishes NuGet package to both nuget.org and GitHub Packages
6. Updates `badger.json` for Scoop package manager

### Platform-Specific Notes
- **macOS**: Requires code signing and notarization (configured in workflow)
- **entitlements.plist**: Grants necessary permissions for macOS binaries (JIT, unsigned memory, etc.)
- **Scoop (Windows)**: Package manifest in `badger.json`
- **Homebrew (macOS/Linux)**: Formula in separate repository (hex/formulae)

## Testing

This project currently has no automated tests. When adding tests:
- Create a separate test project (e.g., `badger.Tests.csproj`)
- Use `dotnet test` to run tests
- Consider testing: pivot point calculations, text alignment, color parsing, file sanitization

## Dependencies

- **ConsoleAppFramework** (4.2.4): CLI framework for automatic argument parsing and help generation
- **SixLabors.ImageSharp** (2.1.3): Core image manipulation library
- **SixLabors.ImageSharp.Drawing** (1.0.0-beta15): Drawing primitives for text and shapes

Note: ImageSharp requires specific trim settings in the .csproj to work correctly with single-file publishing.
