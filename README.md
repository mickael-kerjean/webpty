<p align="center">
    <a href="http://www.webpty.io" target="_blank" rel="noopener">
        <img src="https://raw.githubusercontent.com/mickael-kerjean/webpty/main/.assets/hero-banner.png" alt="WebPty - open source web shell from 1 binary" />
    </a>
</p>

I love bash and I spend most of my time in a terminal so I needed a way to bring that environment over to my browser to make it confortable to run my servers.

Webpty is a zero config fat binary that "just work" from any server with any program (obviously including the likes of tmux, vim, emacs, ...). 

It is safe by default meaning:
- you can't use webpty using HTTP at all
- it will create its own self signed certificates if you don't supply your own
- rely on SSH to authorise users meaning only people who already have an account on the box can connect to webpty
- I have it exposed on the internet to manage my servers which run in prod at [Filestash](https://github.com/mickael-kerjean/filestash)

In the roadmap:
- make it possible to expose a server from the internet via a tunnel
- plugin mechanism to support additional user management mechanism, custom theme, logs, ...
- run webpty as a library so you can make a webshell for your own web application

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
