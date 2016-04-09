# srvdir - Expose any directory as a secure public file server

## Homepage
[http://srvdir.net](http://srvdir.net)

## Project status

srvdir is no longer developed, supported or maintained by its author, except to ensure that the project continues to compile. The contribution policy has the following guidelines:

1. All issues against this repository will be closed unless they demonstrate a crash or other complete failure of srvdir's functionality.
2. No new features will be added. Any pull requests with new features will be closed. Please fork the project instead.
3. Pull requests fixing existing bugs or improving documentation are welcomed.

The public srvdir.net service that used to run is no longer available. It was shut down in April, 2016.

## Usage
srvdir allows you to turn any directory into a secure public file server. You can think of it like Python's SimpleHTTPServer, but with a TLS-secured public URL. It is similar to an on-demand public Dropbox folder.

srvdir is extremely simple to use. Serve the files in the current directory:

    srvdir

Serve the files in a different directory:

    srvdir /path/to/directory

Serve the files in two directories:

    srvdir /path/to/directory /path/to/other/directory

Give the public url a custom name:

    srvdir example:/path/to/directory

Require a username/password to access the file server:

    srvdir -auth="user:password" /path/to/directory

Putting it all together, let's serve the directory at /usr/local on ulocal.srvdir.net and also serve the current directory at current.srvdir.net. Let's also make sure you need to enter a username and password as well:

    srvdir -auth="root:12345" ulocal:/usr/local current:.

That's it!

## Code
Most of the code from srvdir is either a modified version of the standard library's http.FileServer or from the go-tunnel library: [https://github.com/inconshreveable/go-tunnel](https://github.com/inconshreveable/go-tunnel)

## Future Plans
I'd love to add native support to serve the files over SFTP/FTP as well as HTTP.

## Licenese
Apache
