# Gentoo System Check (GSC)

![Status](https://img.shields.io/badge/status-stable-green) ![Go Version](https://img.shields.io/badge/go-1.16-blue)

GSC stands for Gentoo System Check. It's designed to help you maintain the integrity and cleanliness of your Gentoo system by using your local package database.

## Overview

Over time, a Gentoo system might accumulate files that deviate from the local package database.

While it's expected for files to be modified or added, there may be situations where you encounter orphan files, dead symlinks, or even manually installed packages that have been forgotten over time. GSC gives you a snapshot of these anomalies, helping you maintain a clean system.

## Installation

```bash
go install github.com/jsurloppe/gsc/cmd/gsc@latest
```

## Why GSC?

- **Efficient System Management**: With GSC, quickly identify and manage anomalies in your system.
- **Flexibility**: Choose between standard or JSON outputs, and even specify paths to ignore.
- **Stay Updated**: Regularly comparing against the package database ensures your system remains streamlined.

## Usage

```bash
gsc [flags] [path]
```

usually
```bash
gsc /etc
```

and

```
gsc /usr
```

## Available Flags

- `--json`: Output the logs in a JSON format. Ideal for parsing.
- `-i`: Utilize an ignore file to bypass specific files or paths.
- `-V`: Display the GSC version.


## Feedback and Contributions

We welcome feedback, issues, and PRs! Please use the GitHub issue tracker for any bugs or feature requests.

## License

This project is licensed under the BSD 3-Clause License - see the LICENSE file for details.