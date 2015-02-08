// +build ignore

#import <Foundation/NSArray.h>
#import <Foundation/Foundation.h>
#import <SystemConfiguration/SCPreferences.h>
#import <SystemConfiguration/SCNetworkConfiguration.h>
#include <sys/syslimits.h>
#include <sys/stat.h>
#include <mach-o/dyld.h>
#include "pacon.h"

void usage()
{
  puts("Usage: helper [setuid | on | off] <pac url>");
}

int setUid()
{
  AuthorizationRef authRef;
  OSStatus result;
#pragma GCC diagnostic ignored "-Wdeprecated-declarations"
  result = AuthorizationCopyPrivilegedReference(&authRef, kAuthorizationFlagDefaults);
#pragma GCC diagnostic warning "-Wdeprecated-declarations"
  if (result != errAuthorizationSuccess) {
    NSLog(@"Not running as root");
    return EPERM;
  }
  char exeFullPath [PATH_MAX];
  uint32_t size = PATH_MAX;
  if (_NSGetExecutablePath(exeFullPath, &size) != 0)
  {
    NSLog(@"Path longer than %d, should not occur!!!!!", size);
    return E2BIG;
  }
  if (chown(exeFullPath, 0, 0) != 0) // root:wheel
  {
    NSLog(@"Error chown");
    return EPERM;
  }
  if (chmod(exeFullPath, S_IRWXU | S_IRGRP | S_IXGRP | S_IROTH | S_IXOTH | S_ISUID) != 0)
  {
    NSLog(@"Error chmod");
    return EPERM;
  }
  return 0;
}

int togglePac(int onOff, const char* cPacUrl)
{
  NSString* pacUrl = [[NSString alloc] initWithCString: cPacUrl encoding:NSUTF8StringEncoding];
  BOOL success = FALSE;

  SCNetworkSetRef networkSetRef;
  CFArrayRef networkServicesArrayRef;
  SCNetworkServiceRef networkServiceRef;
  SCNetworkProtocolRef proxyProtocolRef;
  NSDictionary *oldPreferences;
  NSMutableDictionary *newPreferences;
  NSString *wantedHost;


  // Get System Preferences Lock
  SCPreferencesRef prefsRef = SCPreferencesCreate(NULL, CFSTR("org.getlantern.lantern"), NULL);

  if(prefsRef==NULL) {
    NSLog(@"Fail to obtain Preferences Ref!!");
    goto freePrefsRef;
  }

  success = SCPreferencesLock(prefsRef, TRUE);
  if (!success) {
    NSLog(@"Fail to obtain PreferencesLock");
    goto freePrefsRef;
  }

  // Get available network services
  networkSetRef = SCNetworkSetCopyCurrent(prefsRef);
  if(networkSetRef == NULL) {
    NSLog(@"Fail to get available network services");
    goto freeNetworkSetRef;
  }

  //Look up interface entry
  networkServicesArrayRef = SCNetworkSetCopyServices(networkSetRef);
  networkServiceRef = NULL;
  for (long i = 0; i < CFArrayGetCount(networkServicesArrayRef); i++) {
    networkServiceRef = CFArrayGetValueAtIndex(networkServicesArrayRef, i);

    // Get proxy protocol
    proxyProtocolRef = SCNetworkServiceCopyProtocol(networkServiceRef, kSCNetworkProtocolTypeProxies);
    if(proxyProtocolRef == NULL) {
      NSLog(@"Couldn't acquire copy of proxyProtocol");
      goto freeProxyProtocolRef;
    }

    oldPreferences = (__bridge NSDictionary*)SCNetworkProtocolGetConfiguration(proxyProtocolRef);
    newPreferences = [NSMutableDictionary dictionaryWithDictionary: oldPreferences];
    wantedHost = @"localhost";

    if(onOff == TRUE) {//Turn proxy configuration ON
      [newPreferences setValue: wantedHost forKey:(NSString*)kSCPropNetProxiesHTTPProxy];
      [newPreferences setValue:[NSNumber numberWithInt:1] forKey:(NSString*)kSCPropNetProxiesProxyAutoConfigEnable];
      [newPreferences setValue:pacUrl forKey:(NSString*)kSCPropNetProxiesProxyAutoConfigURLString];
      NSLog(@"Setting pac ON for device %@ with: %@",
          SCNetworkServiceGetName(networkServiceRef), newPreferences);
    } else {//Turn proxy configuration OFF
      [newPreferences setValue:[NSNumber numberWithInt:0] forKey:(NSString*)kSCPropNetProxiesProxyAutoConfigEnable];
      NSLog(@"Setting pac OFF for device %@", SCNetworkServiceGetName(networkServiceRef));
    }

    success = SCNetworkProtocolSetConfiguration(proxyProtocolRef, (__bridge CFDictionaryRef)newPreferences);
    if(!success) {
      NSLog(@"Failed to set Protocol Configuration");
      goto freeProxyProtocolRef;
    }

freeProxyProtocolRef:
    CFRelease(proxyProtocolRef);
  }

  success = SCPreferencesCommitChanges(prefsRef);
  if(!success) {
    NSLog(@"Failed to Commit Changes");
    goto freeNetworkServicesArrayRef;
  }

  success = SCPreferencesApplyChanges(prefsRef);
  if(!success) {
    NSLog(@"Failed to Apply Changes");
    goto freeNetworkServicesArrayRef;
  }
success = TRUE;
  //Free Resources
freeNetworkServicesArrayRef:
  CFRelease(networkServicesArrayRef);
freeNetworkSetRef:
  CFRelease(networkSetRef);
freePrefsRef:
  SCPreferencesUnlock(prefsRef);
  CFRelease(prefsRef);

  return success == TRUE ? 0 : EPERM;
}


int main() {
  NSArray *args = [[NSProcessInfo processInfo] arguments];
  if (args.count < 2) {
    usage();
    return EINVAL;
  }
  if ([[args objectAtIndex:1] isEqual: @"setuid"]) {
    return setUid();
  } else if ([[args objectAtIndex:1] isEqual: @"on"]) {
    if (args.count < 3) {
      usage();
      return EINVAL;
    }
    return togglePac(PAC_ON, [[args objectAtIndex:2] UTF8String] );
  } else if ([[args objectAtIndex:1] isEqual: @"off"]) {
    return togglePac(PAC_OFF, "");
  } else {
    usage();
    return EINVAL;
  }

  return 0;
}
