#include "pacon.h"

const int PAC_ON = 1;
const int PAC_OFF = 0;

#ifdef WIN32
#include "platform/windows.c"
#endif

#ifdef LINUX
#include "platform/linux.c"
#endif

#ifdef DARWIN
#include "platform/darwin.m"
#endif
