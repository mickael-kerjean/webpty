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

# Documentation

Installation: http://www.webpty.io/

# Fleet server

WebPty works with an associated web fleet server to centralise all your devices in one place even though they don't sit in the same network:

<img src="https://raw.githubusercontent.com/mickael-kerjean/webpty/main/.assets/webfleet-banner.png" alt="WebFleet server" />
