#include <gtk/gtk.h>
#include "pacon.h"

int togglePac(int onOff, const char* pacUrl)
{
	int ret = 0;
	g_type_init(); // deprecated since version 2.36
	GSettings* setting = g_settings_new("org.gnome.system.proxy");
	if (onOff == PAC_ON) {
		gboolean success = g_settings_set_string(setting, "mode", "auto");
		if (!success) {
			printf("error setting mode to auto\n");
			ret = -1;
			goto cleanup;
		}
		success = g_settings_set_string(setting, "autoconfig-url", pacUrl);
		if (!success) {
			printf("error setting autoconfig-url to %s\n", pacUrl);
			ret = -1;
			goto cleanup;
		}
	}
	else {
		gboolean success = g_settings_set_string(setting, "mode", "none");
		if (!success) {
			printf("error setting mode to none\n");
			ret = -1;
			goto cleanup;
		}
		g_settings_reset(setting, "autoconfig-url");
	}
cleanup:
	g_settings_sync();
	g_object_unref(setting);
}
