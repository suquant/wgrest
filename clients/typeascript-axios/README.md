## @

Manage WireGuard VPN tunnels by RESTful manner.

Supported features:

  * Manage device: create, update, and delete wireguard interface
  * Manage device's ip addresses: attache or detach ip addresses to the netowrk interface
  * Manage device's peers: create, update, and delete peers
  * Peer's QR code, for use in WireGuard & ForestVPN client

ForestVPN client may be used as alternative client with enabled P2P technology over WireGuard tunnelling.
Read more on https://forestvpn.com/

Environment
* Node.js
* Webpack
* Browserify

Language level
* ES5 - you must have a Promises/A+ library installed
* ES6

Module system
* CommonJS
* ES6 module system

It can be used in both TypeScript and JavaScript. In TypeScript, the definition should be automatically resolved via `package.json`. ([Reference](http://www.typescriptlang.org/docs/handbook/typings-for-npm-packages.html))

### Building

To build and compile the typescript sources to javascript use:
```
npm install
npm run build
```

### Publishing

First build the package then run ```npm publish```

### Consuming

navigate to the folder of your consuming project and run one of the following commands.

_published:_

```
npm install @ --save
```

_unPublished (not recommended):_

```
npm install PATH_TO_GENERATED_PACKAGE --save
