<p align="center">
    <a href="http://www.webpty.io" target="_blank" rel="noopener">
        <img src="https://raw.githubusercontent.com/mickael-kerjean/webpty/main/.assets/hero-banner.png" alt="WebPty - open source web shell from 1 binary" />
    </a>
</p>

WebPty makes it possible to run your server terminal from your browser. It is a zero config fat binary that "just work" from any linux server with any program

I've build WebPty so you can't missuse it in a way that could do something bad meaning:
- It won't work on HTTP even if you try really hard
- it will create its own self signed certificates if you don't supply your own
- rely on SSH to authorise users meaning only people who already have an account on the box can connect to webpty
- I have it exposed on the internet to manage my servers which run in prod at [Filestash](https://github.com/mickael-kerjean/filestash)

# Install

```
# arm build
curl -L -o webpty.bin "https://github.com/mickael-kerjean/webpty/releases/download/stable/webpty_linux_arm.bin"
# intel/amd
curl -L -o webpty.bin "https://github.com/mickael-kerjean/webpty/releases/download/stable/webpty_linux_amd64.bin"


# launch it
chmod +x webpty.bin
./webpty.bin
```

# The vision / roadmap

In the long term vision, WebPty is a small agent who can establish a connection to its accompanying WebFleet server to have a 1 stop shop to manage everything, making it possible to do cool stuff like:
- give restricted access to a range of authorised people
- make it possible for someone to request access to a server
- enable screensharing session
- plugin mechanism to support various user management system and not just basic authentication
- run webpty as a library to expose a webshell to your application
- use other protocols / platform like RDP, VNC and more

So far we have an alpha release of the webfleet server that is capable of nat traversal and firewall punching to get things to work accross etherogeneous network:

<img src="https://raw.githubusercontent.com/mickael-kerjean/webpty/main/.assets/webfleet-banner.png" alt="WebFleet server" />

Expect dragons and unfinished things, the tunnel is functional but at this very early stage, webfleet is just a cool demo of what the whole thing could look like if I could put enough time to build it all.
