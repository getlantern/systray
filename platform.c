#include "systray.h"

#ifdef WIN32
#include "platform/windows.c"
#endif

#ifdef LINUX
#include "platform/linux.c"
#endif

#ifdef DARWIN
#include "platform/darwin.m"
#endif
