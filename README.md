# lightproxy

Lightweight local proxy, useful for giving your local services memorable local domain names like `myproject.wip` instead of `localhost:8000`.

**Key features!**

- Super simple proxy-based setup
- Choose your own local tld, e.g. `.localhost` or `.wip`
- Wildcard subdomains, e.g. `*.devserver.wip`
- Zero config HTTPS support for all mapped URLs
- Proxy mode: route requests to a local server
- File server mode: serve files directly from a folder

## Usage

### Installation (brew)

```
brew install octavore/tools/lightproxy
```

To have homebrew start lightproxy when you start your computer

```
brew services start octavore/tools/lightproxy
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

### Configure your system proxy (Windows)

[See here](https://pypac.readthedocs.io/en/latest/about_pac.html#windows)

## TLS

All proxied URLS are also automatically available over https, i.e. http://foo.wip is also served with https at https://foo.wip.

Internally, lightproxy does this by listening for TLS connections on a separate port, 7998. You can change this port if it conflicts with another app you have running. You should never need to connect to it directly.

Please note that you will see browser warnings because lightproxy generates a self-signed certificate. To suppress the warnings, you can tell lightproxy to sign certificates using a locally trusted root CA. There are two ways to do this: using [`mkcert`](https://github.com/FiloSottile/mkcert), or generating your own local CA using openssl.

### Option 1: Generate a local root CA with openssl

1. Install [`mkcert`](https://github.com/FiloSottile/mkcert) with `brew install mkcert`
2. Run `mkcert -install` to generate and install its root CA
3. Set `mkcert` to `true` in `config.json`.

### Option 2: Generate a local root CA with openssl

1. Generate a local CA key file

   ```
   openssl genrsa -out lightproxyCA.key 4096
   openssl req -x509 -new -nodes -key lightproxyCA.key -sha256 -days 1024 -out lightproxyCA.crt
   ```

2. Trust the key file

   If you are on a Mac, open **Keychain Access** in System Preference, select the **System** keychain, select the **Certificates** category, and drag your key file into the list of certificates.

   Double click the certificate, open up the **Trust** section, and make sure _When using this certificate:_ is set to _Always Trust_. Restart your browsers for this to take effect.

3. Set `ca_key_file` in lightproxy's config

   Tell lightproxy where to find your files by setting the `ca_key_file` parameter in `config.json` (see below). Set it to the keyfile path you generated in step 1. Restart lightproxy.

## config.json

URL mappings can be added easily using the commands above, but you can also edit the config file directly.

To view the current config and where it is saved, run this:

```
lightproxy config
```

**Example config**

```jsonc
{
  # only hosts ending in this tld are handled by the proxy
  "tld": "wip",

  # the local addr the proxy listens on
  "addr": "localhost:7999",

  # the local addr the proxy listens on internally to handle
  # TLSrequests
  "tls_addr": "localhost:7998",


  # if you are using mkcert for your local CA, this should be set to true
  "mkcert": false,

  # if you are using a local CA, this should be set to the value
  # of the *.key file. The corresponding *.crt should be in the same
  # folder.
  "ca_key_file": "...",

  # if you are using a local CA, this should be set to the value
  # of the certificate file, if it cannot be inferred from the ca_key_file value.
  "ca_cert_file": "...",

  # all hosts to map. Entries must have `host` set, and
  # either `dest` (a local addr), or `dest_folder` (a folder)
  "entries": [
    {
      "host": "ketchup.wip",
      "dest": "localhost:8000"
    },
    {
      "host": "goatcodes.wip",
      "dest": "localhost:8100"
    },
    {
      "host": "go.goatcodes.wip",
      "dest": "localhost:8100"
    },
    {
      "host": "blog.goatcodes.wip",
      "dest_folder": "/sites/goatblog"
    }
  ]
}
```

## Under the hood

lightproxy uses a [proxy auto-config file](<https://developer.mozilla.org/en-US/docs/Web/HTTP/Proxy_servers_and_tunneling/Proxy_Auto-Configuration_(PAC)_file>) to tell your system what URLs to send to the proxy.

By default, this PAC file is served at `http://localhost:7999/proxy.pac`. It contains a simple javascript function:

```js
function FindProxyForURL(url, host) {
  if (shExpMatch(host, "*.wip")) {
    return "PROXY 127.0.0.1:7999";
  }
  return "DIRECT";
}
```

Resources

- https://findproxyforurl.com/example-pac-file/
- https://pypac.readthedocs.io/en/latest/about_pac.html

## Known issues

- Secure websockets do not proxy correctly in Safari.

## Changelog

### 2022-08-20 / 1.3.0

- [Fixed support for secure websockets](https://github.com/octavore/lightproxy/pull/4) (credit: @jonian)
- Add support for `mkcert` in config.
- Add support for explicitly defining CA cert file in config: `ca_cert_file`.

### 2019-12-11 / 1.2.1

- Search both `$XDG_CONFIG_HOME` and `$HOME/.config/lightproxy for`config.json`(preferring the former if set). This change improves the`brew installation` experience while maintaining backward compatability.

### 2019-05-10 / 1.2

- Support wildcards in source url

### 2019-05-10 / 1.1

- Added TLS support
- Added `tls_addr` and `ca_key_file` config option
