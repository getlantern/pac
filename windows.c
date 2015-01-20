#include <stdlib.h>
#include <windows.h>
#include <wininet.h>
#include <ras.h>
#include <tchar.h>
#pragma comment( lib, "wininet" )
#pragma comment( lib, "rasapi32" )


// Figure out which Dial-Up or VPN connection is active; in a normal LAN connection, this should
// return NULL. NOTE: For some reason this method fails when compiled in Debug mode but works
// every time in Release mode.
LPTSTR FindActiveConnection() {
  DWORD dwCb = sizeof(RASCONN);
  DWORD dwErr = ERROR_SUCCESS;
  DWORD dwRetries = 5;
  DWORD dwConnections = 0;
  RASCONN* lpRasConn = NULL;
  RASCONNSTATUS rasconnstatus;
  rasconnstatus.dwSize = sizeof(RASCONNSTATUS);

  //
  // Loop through in case the information from RAS changes between calls.
  //
  while (dwRetries--) {
    // If the memory is allocated, free it.
    if (NULL != lpRasConn) {
      HeapFree(GetProcessHeap(), 0, lpRasConn);
      lpRasConn = NULL;
    }

    // Allocate the size needed for the RAS structure.
    lpRasConn = (RASCONN*)HeapAlloc(GetProcessHeap(), 0, dwCb);
    if (NULL == lpRasConn) {
      dwErr = ERROR_NOT_ENOUGH_MEMORY;
      break;
    }

    // Set the structure size for version checking purposes.
    lpRasConn->dwSize = sizeof(RASCONN);

    // Call the RAS API then exit the loop if we are successful or an unknown
    // error occurs.
    dwErr = RasEnumConnections(lpRasConn, &dwCb, &dwConnections);
    if (ERROR_INSUFFICIENT_BUFFER != dwErr) {
      break;
    }
  }
  //
  // In the success case, print the names of the connections.
  //
  if (ERROR_SUCCESS == dwErr) {
    DWORD i;
    for (i = 0; i < dwConnections; i++) {
      RasGetConnectStatus(lpRasConn[i].hrasconn, &rasconnstatus);
      if (rasconnstatus.rasconnstate == RASCS_Connected){
        return lpRasConn[i].szEntryName;
      }

    }
  }
  return NULL; // Couldn't find an active dial-up/VPN connection; return NULL
}

void toggleAutoProxyConfigFile(const char* cOnOff, const char* cAutoProxyConfigFileUrl)
{
  INTERNET_PER_CONN_OPTION_LIST options;
  BOOL    bReturn;
  DWORD   dwBufferSize = sizeof(options);
  options.dwSize = dwBufferSize;
  options.pszConnection = FindActiveConnection();

  options.dwOptionCount = 2;
  options.pOptions = calloc(2, sizeof(INTERNET_PER_CONN_OPTION));

  if(!options.pOptions) {
    return FALSE;
  }

  options.pOptions[0].dwOption = INTERNET_PER_CONN_FLAGS;
  if strcmp(cOnOff, "on") == 0 {
    options.pOptions[0].Value.dwValue = PROXY_TYPE_AUTO_PROXY_URL;
  }
  else {
    options.pOptions[0].Value.dwValue = PROXY_TYPE_AUTO_DETECT
  }

  options.pOptions[1].dwOption = INTERNET_PER_CONN_AUTOCONFIG_URL;

  options.pOptions[1].Value.pszValue = proxyAddressStr;

  bReturn = InternetSetOption(NULL,INTERNET_OPTION_PER_CONNECTION_OPTION, &options, dwBufferSize);

  free(options.pOptions);

  InternetSetOption(NULL, INTERNET_OPTION_SETTINGS_CHANGED, NULL, 0);
  InternetSetOption(NULL, INTERNET_OPTION_REFRESH , NULL, 0);

  return bReturn;

