# mdd

Markdown HTTP server for local editing.

## Install

```bash
make install
```

## Usage

Run `mdd` from any directory containing `.md` files. It serves them as rendered HTML at `http://127.0.0.1:1980/`.

- `/` and directory paths default to `README.md`
- Only `.md` files within the working directory are served

## Testing

```bash
make test
```
