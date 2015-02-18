# pac

[pac](https://github.com/getlantern/pac) is a simple Go library to toggle on and off pac(proxy auto configuration) for Windows, Mac OS and Linux. It will extract a helper tool and use it to actually chage pac.

```go
pac.EnsureHelperToolPresent(fullPath, prompt, iconFullPath)
pac.On(pacUrl string)
pac.Off()
```

See 'example/main.go' for detailed usage.

### Windows

Install [MinGW-W64](http://sourceforge.net/projects/mingw-w64) as it has up to date SDK headers we require.
