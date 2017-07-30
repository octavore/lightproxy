# lightproxy

Lightweight proxy, useful for giving your local services memorable local domain names like `myproject.wip` instead of `localhost:8000`.

## Usage

### Installation (brew)

```
brew install octavore/tools/lightproxy
```

### Registering URL mappings

To map `foo.wip` to `localhost:3000`:

```bash
lightproxy set-dest foo.wip 3000
```

To map `foo.wip` to a folder:

```bash
lightproxy set-dir foo.wip ~/Code/foo/
```

### Starting lightproxy

Open up a new terminal shell and run `lightproxy`.

## todo

- daemon mode
- path based routing
- prettier UI
