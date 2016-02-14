# lightproxy

Lightweight proxy, useful for developing local services.

## Example Usage

### Installation (brew)

```
brew install octavore/tools/lightproxy
```

### Registering URL mapping

To map `foo.dev` to `localhost:3000`:

```bash
# add mapping to hosts file
echo "127.0.0.1 foo.dev" | sudo tee --append /etc/hosts
lightproxy set-dest foo.dev localhost
```

### Starting server

```
# sudo is used to allow lightproxy to listen on port 80
sudo lightproxy start
```

## todo

- daemon mode
- path based routing
- automatically add to host file
- drop privileges
- prettier UI
