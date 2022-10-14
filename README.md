<p align="center">
    <a href="http://www.webpty.io" target="_blank" rel="noopener">
        <img src="https://raw.githubusercontent.com/mickael-kerjean/webpty/main/.assets/github_readme.png" alt="WebPty - open source web shell from 1 binary" />
    </a>
</p>

[WebPty](http://www.webpty) is an open source web shell made in Golang which features:

- full fledge terminal from your browser that works well with emacs/vim, ...
- safe by default: creates its own self signed certificates if you haven't supply your own, rely on SSH to authorise users
- good looking
- zero config fat binary that "just work" from anything anywhere
- (in the roadmap) make it work even from some closed network so you can access your raspberry pi at home from your office sitting behind your company VPN
- (in the roadmap) plugin mechanism to support additional user management mechanism, custom theme, ...
- (in the roadmap) can be integrated as a library to provide a webshell to your application
- (in the roadmap) emoji and other weird characters

# Install

```
# download binary from github release
curl http://TODO

# launch it
chmod +x webpty.bin
./webpty.bin
```
