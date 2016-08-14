/*Package ws implements Websocket transport

Websocket transport implements a HTTP(S) compliable, surveillance proof transport method with plausible deniability.

To configure such a listener, set streamSettings to be ws. A http(s) listener will be listening at the port you have configured.

There is additional configure can be made at transport configure.

"wsSettings":{
      "Path":"ws", // the path our ws handler bind
      "Pto": "wss/ws", // the transport protocol we are using ws or wss(listen ws with tls)
      "Cert":"cert.pem", // if you have configured to use wss, configure your cert here
      "PrivKey":"priv.pem" //if you have configured to use wss, configure your privatekey here
    }


To configure such a Dialer, set streamSettings to be ws.

There is additional configure can be made at transport configure.

"wsSettings":{
      "Path":"ws", // the path our ws handler bind
      "Pto": "wss/ws", // the transport protocol we are using ws or wss(listen ws with tls)
    }

It is worth noting that accepting a non-valid cert is not supported as a self-signed or invalid cert can be a sign of a website that is not correctly configured and lead to additional investigation.


This transport was disscussed at
https://github.com/v2ray/v2ray-core/issues/224
https://trello.com/c/3uCCeBkC/8-add-websocket-transport

*/
package ws
