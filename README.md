[pacon](https://github.com/getlantern/pacon) is a simple Go library to toggle on and off pac(proxy auto configuration) for Windows and Mac OSX.

```go
pacon.PacOn(pacUrl string)
pacon.PacOff(pacUrl string)
```

### Windows

Install [MinGW-W64](http://sourceforge.net/projects/mingw-w64) as it has up to date SDK headers we require.
