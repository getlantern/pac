#import <Foundation/NSArray.h>
#import <Foundation/Foundation.h>
#import <SystemConfiguration/SCPreferences.h>
#import <SystemConfiguration/SCNetworkConfiguration.h>
#include <sys/syslimits.h>
#include <sys/stat.h>
#include <mach-o/dyld.h>
#include "darwin.h"

int runAuthorized(char *path, char *prompt, char *iconPath)
{
  AuthorizationEnvironment authEnv;
  AuthorizationItem kAuthEnv[2];
  authEnv.items = kAuthEnv;
  authEnv.count = 0;

  if (prompt != NULL_STRING)
  {
    kAuthEnv[authEnv.count].name = kAuthorizationEnvironmentPrompt;
    kAuthEnv[authEnv.count].valueLength = strlen(prompt);
    kAuthEnv[authEnv.count].value = prompt;
    kAuthEnv[authEnv.count].flags = 0;
    authEnv.count++;
  }
  if (iconPath != NULL_STRING)
  {
    kAuthEnv[authEnv.count].name = kAuthorizationEnvironmentIcon;
    kAuthEnv[authEnv.count].valueLength = strlen(iconPath);
    kAuthEnv[authEnv.count].value = iconPath;
    kAuthEnv[authEnv.count].flags = 0;
    authEnv.count++;
  }

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
  OSStatus status = AuthorizationCreate(&authRights, &authEnv, authFlags, &authRef);
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
