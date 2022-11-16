# Simple dumb fuzzer

A simple dumb fuzzer than ever. 🤣

## Useage

```shell
$ STATSD_SERVER="eff.g3un.com:8125" gor . ../test $(which lldb)
```

## Supported mutate methods

- Insert
- Delete

## Supported debuggers

- GDB (for Linux)
- LLDB (for macOS)
- WinDbg (for Windows)
