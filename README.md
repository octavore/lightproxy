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

### Configure your system proxy (MacOS)

1. Go to System Preferences > Network > Proxies.
2. Select and check **Automatic Proxy Configuration** in the sidebar.
3. Set the URL to http://localhost:7999/proxy.pac

![screen shot 2018-06-26 at 8 46 00 pm](https://user-images.githubusercontent.com/1707744/41951981-87e8f856-7982-11e8-8e95-c06cca186eb3.png)

## todo

- daemon mode
- path based routing
- prettier UI


