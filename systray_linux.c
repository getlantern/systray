#include <stdlib.h>
#include <string.h>
#include <errno.h>
#include <libappindicator/app-indicator.h>
#include "systray.h"

static AppIndicator *global_app_indicator;
static GtkWidget *global_tray_menu = NULL;
static GList *global_menu_items = NULL;
typedef struct {
	GtkWidget *menu_item;
	char *menu_id;
} MenuItemNode;

int nativeLoop(void) {
	gdk_threads_init();
	gtk_init(0, NULL);
	global_app_indicator = app_indicator_new("git_indicator", "",
			APP_INDICATOR_CATEGORY_APPLICATION_STATUS);
	app_indicator_set_status(global_app_indicator, APP_INDICATOR_STATUS_ACTIVE);
	global_tray_menu = gtk_menu_new();
	app_indicator_set_menu(global_app_indicator, GTK_MENU(global_tray_menu));
	systray_ready();
	gtk_main();
}


void setIcon(const char* iconBytes, int length) {
	char temp_file_name[PATH_MAX];
	strcpy(temp_file_name, "/tmp/systray_XXXXXX");
	int fd = mkstemp(temp_file_name);
	ssize_t written = write(fd, iconBytes, length);
	close(fd);
	if(written != length) {
		printf("failed to write temp icon file: %s\n", strerror(errno));
		return;
	}
	app_indicator_set_icon_full(global_app_indicator, temp_file_name, "");
	app_indicator_set_attention_icon_full(global_app_indicator, temp_file_name, "");
}

void setTitle(char* ctitle) {
	app_indicator_set_label(global_app_indicator, ctitle, "");
	free(ctitle);
}

void setTooltip(char* ctooltip) {
	free(ctooltip);
}

void add_or_update_menu_item(char* menuId, char* title, char* tooltip) {
	gdk_threads_enter();
	GList* it;
	for(it = global_menu_items; it != NULL; it = it->next) {
		MenuItemNode* item = (MenuItemNode*)(it->data);
		if(strcmp(item->menu_id, menuId) == 0){
			gtk_menu_item_set_label(GTK_MENU_ITEM(item->menu_item), title);
			free(menuId);
			break;
		}
	}

	// menu id doesn't exist, add new item
	if(it == NULL) {
		GtkWidget *titleMenuItem = gtk_menu_item_new_with_label(title);
		g_signal_connect_swapped(G_OBJECT(titleMenuItem), "activate", G_CALLBACK(systray_menu_item_selected), menuId);
		gtk_menu_shell_append(GTK_MENU_SHELL(global_tray_menu), titleMenuItem);

		MenuItemNode* new_item = malloc(sizeof(MenuItemNode));
		new_item->menu_id = menuId;
		new_item->menu_item = titleMenuItem;
		GList* new_node = malloc(sizeof(GList));
		new_node->data = new_item;
		new_node->next = global_menu_items;
		if(global_menu_items != NULL) {
			global_menu_items->prev = new_node;
		}
		global_menu_items = new_node;
	}
	gtk_widget_show_all(global_tray_menu);
	gdk_threads_leave();

	free(title);
	free(tooltip);
}

void quit() {
	gdk_threads_enter();
	gtk_main_quit();
	gdk_threads_leave();
}
