# asciiify

A library and cli for converting gif to animated ascii-art and playing the result in terminal.

Forked from https://github.com/theMomax/asciiify



* Removed the embedded binary "Happy B Day 716.gif" gif
* CLI now accepts commandline arguments:

```
asciiify [GIF URL] [LOOP COUNT]
```

## Debian/Ubuntu Build Instructions

Install goncurses dependencies
```
sudo apt install pkg-config libncurses5-dev
```

Build the dependencies
```
go get
```

Build executable
```
go build
```