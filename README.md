# Nexus CLI

## Why?

For the sake of learning Go.

## Building

```bash
make all_os
```

## Installation

Download the latest [release](https://github.com/bzon/nexus-cli/releases) to any directory in your system. Rename it as `nexus-cli` or for Windows, `nexus-cli.exe`.

Run it as `nexus-cli --help` or `nexus-cli.exe --help`.

## Usage

Environment Variables to avoid using flags `-H`, `-U`, `-P` for host and authentication settings.

```bash
NEXUS_HOST=http://localhost:8081/nexus
NEXUS_USERNAME=admin
NEXUS_PASSWORD=admin123
```

### Downloading an Artifact

Using `download` subcommand.

Help: 

```bash
nexus-cli download --help
```

Usage:

```bash
nexus-cli download -g com.example -a artifactA -p jar -v 1.0.1 -H http://localhost:8081/nexus -U admin -P admin123
```

### Downloading Multiple Artifacts

Using `multi-download` subcommand.

Help: 

```bash
nexus-cli multi-download --help
```

Usage:

Create a file named 'artifacts.txt'.

```bash
com.example:artifactA:1.0.1:jar
com.example:artifactB:1.0-SNAPSHOT:war
com.example:artifactC:LATEST:war
```

```bash
nexus-cli multi-download -f artifacts.txt -h http://localhost:8081/nexus -U admin -P admin123
```

### Uploading an Artifact

WIP
