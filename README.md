[pac](https://github.com/getlantern/pac) is a simple Go library to toggle on and off pac(proxy auto configuration) for Windows, Mac OSX and Linux.

```go
pac.On(pacUrl string)
pac.Off()
```

### Windows

Install [MinGW-W64](http://sourceforge.net/projects/mingw-w64) as it has up to date SDK headers we require.
