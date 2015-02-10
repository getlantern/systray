#include <stdlib.h>
#include <stdio.h>
#include <dlfcn.h>
#include <string.h>
#include <errno.h>
#include <limits.h>
#include <libappindicator/app-indicator.h>
#include "systray.h"

static AppIndicator *global_app_indicator;
static GtkWidget *global_tray_menu = NULL;
static GList *global_menu_items = NULL;
// Keep track of all generated temp files to remove when app quits
static GArray *global_temp_icon_file_names = NULL;

// Dynamically loaded libraries and functions from libappindicator and libgtk. We do this to support cross-compiling from platforms that don't have these libraries

void *libglib2;
GArray* (*d_g_array_new());

void *libgtk3;
void (*d_gtk_main)();
GtkWidget* (*d_gtk_menu_new())();

void *libappindicator3;
AppIndicator* (*d_app_indicator_new)(const gchar*, const gchar*, AppIndicatorCategory);
void (*d_app_indicator_set_status)(AppIndicator*, AppIndicatorStatus);
void (*d_app_indicator_set_menu)(AppIndicator*, GtkMenu*);
void (*d_app_indicator_set_label)(AppIndicator*, const gchar*, const gchar*);
void (*d_app_indicator_set_icon_full)(AppIndicator*, const gchar*, const gchar*);
void (*d_app_indicator_set_attention_icon_full)(AppIndicator*, const gchar*, const gchar*);

void *library(const char* name) {
    void *handle = dlopen(name, RTLD_LAZY);
    if (!handle) {
        fputs(dlerror(), stderr);
        exit(1);
    }
    return handle;
}

void *symbol(void *library, const char* name) {
    char *error;
    void *sym = dlsym(library, name);
    if ((error = dlerror()) != NULL) {
        fputs(error, stderr);
        exit(1);
    }
    return sym;
}

typedef struct {
	GtkWidget *menu_item;
	char *menu_id;
} MenuItemNode;

typedef struct {
	char* menu_id;
	char* title;
	char* tooltip;
	short disabled;
	short checked;
} MenuItemInfo;

void load_libraries() {
	libglib2 = library("libglib-2.0");
	d_g_array_new = symbol(libglib2, "g_array_new");
	d_g_array_append_val = symbol(libglib2, "g_array_append_val");
    d_ = symbol(libglib2, "replaceme");
    d_ = symbol(libglib2, "replaceme");
    d_ = symbol(libglib2, "replaceme");
    d_ = symbol(libglib2, "replaceme");
    d_ = symbol(libglib2, "replaceme");
    d_ = symbol(libglib2, "replaceme");
    d_ = symbol(libglib2, "replaceme");
    d_ = symbol(libglib2, "replaceme");
    d_ = symbol(libglib2, "replaceme");
    d_ = symbol(libglib2, "replaceme");
        

	libgtk3 = library("libgtk-3");
    d_gtk_init = symbol(libgtk3, "gtk_init");
    d_gtk_menu_new = symbol(libgtk3, "gtk_menu_new");
    d_gtk_main = symbol(libgtk3, "gtk_main");
    d_ = symbol(libgtk3, "replaceme");
    d_ = symbol(libgtk3, "replaceme");
    d_ = symbol(libgtk3, "replaceme");
    d_ = symbol(libgtk3, "replaceme");
    d_ = symbol(libgtk3, "replaceme");
    d_ = symbol(libgtk3, "replaceme");
    d_ = symbol(libgtk3, "replaceme");
    d_ = symbol(libgtk3, "replaceme");
    d_ = symbol(libgtk3, "replaceme");

    libappindicator3 = library("libappindicator3");
    d_app_indicator_new = symbol(libappindicator3, "app_indicator_new");
    d_app_indicator_set_status = symbol(libappindicator3, "app_indicator_set_status");
    d_app_indicator_set_menu = symbol(libappindicator3, "app_indicator_set_menu");
    d_app_indicator_set_label = symbol(libappindicator3, "app_indicator_set_label");
    d_app_indicator_set_icon_full = symbol(libappindicator3, "app_indicator_set_icon_full");
    d_app_indicator_set_attention_icon_full = symbol(libappindicator3, "app_indicator_set_attention_icon_full");
    

    
    
}

int nativeLoop(void) {
    load_libraries();
	d_gtk_init(0, NULL);
	global_app_indicator = d_app_indicator_new("systray", "",
			APP_INDICATOR_CATEGORY_APPLICATION_STATUS);
	d_app_indicator_set_status(global_app_indicator, APP_INDICATOR_STATUS_ACTIVE);
	global_tray_menu = d_gtk_menu_new();
	d_app_indicator_set_menu(global_app_indicator, GTK_MENU(global_tray_menu));
	global_temp_icon_file_names = d_g_array_new(TRUE, FALSE, sizeof(char*));
	systray_ready();
	d_gtk_main();
	return;
}

// runs in main thread, should always return FALSE to prevent gtk to execute it again
gboolean do_set_icon(gpointer data) {
	GBytes* bytes = (GBytes*)data;
	char* temp_file_name = malloc(PATH_MAX);
	strcpy(temp_file_name, "/tmp/systray_XXXXXX");
	int fd = mkstemp(temp_file_name);
	if (fd == -1) {
		printf("failed to create temp icon file %s: %s\n", temp_file_name, strerror(errno));
		return FALSE;
	}
	g_array_append_val(global_temp_icon_file_names, temp_file_name);
	gsize size = 0;
	gconstpointer icon_data = g_bytes_get_data(bytes, &size);
	ssize_t written = write(fd, icon_data, size);
	close(fd);
	if(written != size) {
		printf("failed to write temp icon file %s: %s\n", temp_file_name, strerror(errno));
		return FALSE;
	}
	d_app_indicator_set_icon_full(global_app_indicator, temp_file_name, "");
	d_app_indicator_set_attention_icon_full(global_app_indicator, temp_file_name, "");
	g_bytes_unref(bytes);
	return FALSE;
}

// runs in main thread, should always return FALSE to prevent gtk to execute it again
gboolean do_add_or_update_menu_item(gpointer data) {
	MenuItemInfo *mii = (MenuItemInfo*)data;
	GList* it;
	for(it = global_menu_items; it != NULL; it = it->next) {
		MenuItemNode* item = (MenuItemNode*)(it->data);
		if(strcmp(item->menu_id, mii->menu_id) == 0){
			gtk_menu_item_set_label(GTK_MENU_ITEM(item->menu_item), mii->title);
			free(mii->menu_id);
			break;
		}
	}

	// menu id doesn't exist, add new item
	if(it == NULL) {
		GtkWidget *menu_item = gtk_menu_item_new_with_label(mii->title);
		g_signal_connect_swapped(G_OBJECT(menu_item), "activate", G_CALLBACK(systray_menu_item_selected), mii->menu_id);
		gtk_menu_shell_append(GTK_MENU_SHELL(global_tray_menu), menu_item);

		MenuItemNode* new_item = malloc(sizeof(MenuItemNode));
		new_item->menu_id = mii->menu_id;
		new_item->menu_item = menu_item;
		GList* new_node = malloc(sizeof(GList));
		new_node->data = new_item;
		new_node->next = global_menu_items;
		if(global_menu_items != NULL) {
			global_menu_items->prev = new_node;
		}
		global_menu_items = new_node;
		it = new_node;
	}
	GtkWidget * menu_item = GTK_WIDGET(((MenuItemNode*)(it->data))->menu_item);
	gtk_widget_set_sensitive(menu_item, mii->disabled == 1 ? FALSE : TRUE);
	gtk_widget_show_all(global_tray_menu);

	free(mii->title);
	free(mii->tooltip);
	free(mii);
	return FALSE;
}

// runs in main thread, should always return FALSE to prevent gtk to execute it again
gboolean do_quit(gpointer data) {
	int i;
	for (i = 0; i < INT_MAX; ++i) {
		char * temp_file_name = g_array_index(global_temp_icon_file_names, char*, i);
		if (temp_file_name == NULL) {
			break;
		}
		int ret = unlink(temp_file_name);
		if (ret == -1) {
			printf("failed to remove temp icon file %s: %s\n", temp_file_name, strerror(errno));
		}
	}
	gtk_main_quit();
	return FALSE;
}

void setIcon(const char* iconBytes, int length) {
	GBytes* bytes = g_bytes_new_static(iconBytes, length);
	g_idle_add(do_set_icon, bytes);
}

void setTitle(char* ctitle) {
	d_app_indicator_set_label(global_app_indicator, ctitle, "");
	free(ctitle);
}

void setTooltip(char* ctooltip) {
	free(ctooltip);
}

void add_or_update_menu_item(char* menu_id, char* title, char* tooltip, short disabled, short checked) {
	MenuItemInfo *mii = malloc(sizeof(MenuItemInfo));
	mii->menu_id = menu_id;
	mii->title = title;
	mii->tooltip = tooltip;
	mii->disabled = disabled;
	mii->checked = checked;
	g_idle_add(do_add_or_update_menu_item, mii);
}

void quit() {
	g_idle_add(do_quit, NULL);
}