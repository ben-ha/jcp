# jcp

Just cp - A simple cp implementation for the average Joe working on Linux or Mac, who occasionally needs to copy files.


## Features

jcp is a simple copy utility that:
1. Has a copy progress bar
2. Can resume interrupted transfers
3. Can copy directories
4. Supports concurrency
5. Keeps my laptop awake during the transfer!

## Usage

jcp provides a minimal command-line interface, designed to keep it simple. To invoke:

```
jcp <source> <dest>
```

If `<source>` is a directory, it is copied to `<dest>/<source>` recursively.

## Installation

### Linux
Download the jcp binary from the latest release and run.

### MacOS
Download the jcp binary from the latest release and run.

## Limitations
1. jcp is not a sync tool. It does not verify the transfer, detect errors in partial transfers or runs in the background and keeps the directories synced.
2. jcp currently **does not** process symlinks.

## Contributing

If you find this project useful and want to contribute and make it better, you are most welcome!