#ifndef PACON_H
#define PACON_H

extern const char* NULL_STRING;

#ifdef DARWIN
int runAuthorized(char *path, char *prompt, char *iconPath);
#endif //def DARWIN

#endif //ndef PACON_H
