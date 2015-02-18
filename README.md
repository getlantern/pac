# pac

[pac](https://github.com/getlantern/pac) is a simple Go library to toggle on and off pac(proxy auto configuration) for Windows, Mac OS and Linux. It will extract a helper tool to actually chage pac.

```go
pac.On(pacUrl string)
pac.Off()
```
Optionally, you can `SetHelperPath` to designate the full path to where helper tool saved.

Keep in mind, don't call it from different coroutine.

See 'example/main.go' for detailed usage.

### Windows

Install [MinGW-W64](http://sourceforge.net/projects/mingw-w64) as it has up to date SDK headers we require.

### Mac OS
Changing network configuration is privileged operation on Mac OS. So after extracting help tool, we need to grant root privilege to it. Operation system will show a dialog requesting user to input password to gain .
Optionally, you can `SetIconPathOnMacOS` to designate the full path to where helper tool saved.

