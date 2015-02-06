#import <Foundation/NSArray.h>
#import <Foundation/Foundation.h>
#import <SystemConfiguration/SCPreferences.h>
#import <SystemConfiguration/SCNetworkConfiguration.h>
#include <sys/syslimits.h>
#include <sys/stat.h>
#include <mach-o/dyld.h>
#include "pacon.h"

#pragma GCC diagnostic ignored "-Wdeprecated-declarations"

void runAuthorized(const char *path) {
  OSStatus status;

  AuthorizationItem authItems[1];
  authItems[0].name = kAuthorizationRightExecute;
  authItems[0].valueLength = 0;
  authItems[0].value = NULL;
  authItems[0].flags = 0;

  AuthorizationRights authRights;
  authRights.count = sizeof(authItems) / sizeof(authItems[0]);
  authRights.items = authItems;

  AuthorizationFlags authFlags;
  authFlags = kAuthorizationFlagDefaults | kAuthorizationFlagInteractionAllowed | kAuthorizationFlagExtendRights;

  AuthorizationRef authRef;
  status = AuthorizationCreate(&authRights, kAuthorizationEmptyEnvironment, authFlags, &authRef);

  if(status == errAuthorizationSuccess) {
    FILE *pipe = NULL;
    char readBuffer[256];
    char* argv[] = { "setuid", NULL };
    status = AuthorizationExecuteWithPrivileges(authRef, path, kAuthorizationFlagDefaults, argv, &pipe);
    if(status == errAuthorizationSuccess) {
      read(fileno(pipe), readBuffer, sizeof(readBuffer));
      fclose(pipe);
    }

    status = AuthorizationFree(authRef, kAuthorizationFlagDestroyRights);
  }
}

void togglePacWithHelper(int onOff, const char* cPacUrl, const char* path)
{
  int pid = [[NSProcessInfo processInfo] processIdentifier];
  NSPipe *pipe = [NSPipe pipe];
  NSFileHandle *file = pipe.fileHandleForReading;

  NSTask *task = [[NSTask alloc] init];
  task.launchPath = [[NSString alloc] initWithUTF8String: path];
  NSString* pacUrl = [[NSString alloc] initWithUTF8String: cPacUrl];
  if (onOff == PAC_ON) {
    task.arguments = @[@"on", pacUrl];
  } else {
    task.arguments = @[@"off", pacUrl];
  }
  task.standardOutput = pipe;

  [task launch];

  NSData *data = [file readDataToEndOfFile];
  [file closeFile];

  NSString *grepOutput = [[NSString alloc] initWithData: data encoding: NSUTF8StringEncoding];
  NSLog (@"grep returned:\n%@", grepOutput);
  return;
}
