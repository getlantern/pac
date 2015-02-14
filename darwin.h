#ifndef PACON_H
#define PACON_H

extern const char* NULL_STRING;

#ifndef DARWIN
int togglePac(int onOff, const char* autoProxyConfigFileUrl);
#endif //ndef DARWIN

#ifdef DARWIN
int runAuthorized(char *path, char *prompt, char *iconPath);
int togglePacWithHelper(int onOff, const char* autoProxyConfigFileUrl, const char* path);
#endif //def DARWIN

#endif //ndef PACON_H
