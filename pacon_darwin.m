#import <Foundation/NSArray.h>
#import <Foundation/Foundation.h>
#import <SystemConfiguration/SCPreferences.h>
#import <SystemConfiguration/SCNetworkConfiguration.h>
#include <sys/syslimits.h>
#include <sys/stat.h>
#include <mach-o/dyld.h>
#include "pacon.h"

AuthorizationRef setUid()
{
  // Get Authorization
  AuthorizationFlags rootFlags = kAuthorizationFlagDefaults
    |  kAuthorizationFlagExtendRights
    |  kAuthorizationFlagInteractionAllowed
    |  kAuthorizationFlagPreAuthorize;
  AuthorizationRef auth;
  OSStatus authErr = AuthorizationCreate(NULL, kAuthorizationEmptyEnvironment, rootFlags, &auth);
  if (authErr != errAuthorizationSuccess) {
    NSLog(@"No Authorization!!!!!");
    auth = NULL;
  }
  char exeFullPath [PATH_MAX];
  uint32_t size = PATH_MAX;
  if (_NSGetExecutablePath(exeFullPath, &size) != 0)
  {
    NSLog(@"Path longer than %d, should not occur!!!!!", size);
    exit(-1);
  }
  if (chown(exeFullPath, 0, 0) != 0) // root:wheel
  {
    NSLog(@"Error chown");
    exit(-1);
  }
  if (chmod(exeFullPath, 0x4755) != 0)
  {
    NSLog(@"Error chmod");
    exit(-1);
  }
  return auth;
}

void togglePac(int onOff, const char* cPacUrl)
{
  AuthorizationRef auth = setUid();
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
  SCPreferencesRef prefsRef = SCPreferencesCreateWithAuthorization(NULL, CFSTR("org.getlantern.lantern"), NULL, auth);

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

    if(onOff == PAC_ON) {//Turn proxy configuration ON
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
  //Free Resources
freeNetworkServicesArrayRef:
  CFRelease(networkServicesArrayRef);
freeNetworkSetRef:
  CFRelease(networkSetRef);
freePrefsRef:
  SCPreferencesUnlock(prefsRef);
  CFRelease(prefsRef);

  return;
}
