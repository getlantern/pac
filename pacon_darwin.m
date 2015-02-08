#import <Foundation/NSArray.h>
#import <Foundation/Foundation.h>
#import <SystemConfiguration/SCPreferences.h>
#import <SystemConfiguration/SCNetworkConfiguration.h>
#include <sys/syslimits.h>
#include <sys/stat.h>
#include <mach-o/dyld.h>
#include "pacon.h"

int runAuthorized(const char *path) {
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
  OSStatus status = AuthorizationCreate(&authRights, kAuthorizationEmptyEnvironment, authFlags, &authRef);
  if(status != errAuthorizationSuccess) {
    return -1;
  }
  FILE *pipe = NULL;
  char readBuffer[256];
  char* argv[] = { "setuid", NULL };
#pragma GCC diagnostic ignored "-Wdeprecated-declarations"
  status = AuthorizationExecuteWithPrivileges(authRef, path, kAuthorizationFlagDefaults, argv, &pipe);
#pragma GCC diagnostic warning "-Wdeprecated-declarations"
  if(status == errAuthorizationSuccess) {
    while (read(fileno(pipe), readBuffer, sizeof(readBuffer)) > 0) {
      ;
    }
    fclose(pipe);
  }
  AuthorizationFree(authRef, kAuthorizationFlagDestroyRights);
  return status == errAuthorizationSuccess ? 0 : -1;
}

int togglePacWithHelper(int onOff, const char* cPacUrl, const char* path)
{
  NSTask *task = [[NSTask alloc] init];
  task.launchPath = [[NSString alloc] initWithUTF8String: path];
  NSString* pacUrl = [[NSString alloc] initWithUTF8String: cPacUrl];
  if (onOff == PAC_ON) {
    task.arguments = @[@"on", pacUrl];
  } else {
    task.arguments = @[@"off", @""];
  }
  [task launch];
  [task waitUntilExit];
  return [task terminationStatus];
}
