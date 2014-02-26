# srvdir - Expose any directory a secure public file server

## Homepage
[https://www.srvdir.net](https://www.srvdir.net)

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
