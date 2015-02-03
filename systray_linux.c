#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <gtk/gtk.h>
#include <gio/gio.h>
#include <gdk-pixbuf/gdk-pixbuf.h>
#include "systray.h"

static GtkWidget *global_tray_menu = NULL;
static GtkStatusIcon *global_tray_icon = NULL;
static GList *global_menu_items = NULL;
typedef struct {
	GtkWidget *menu_item;
	char *menu_id;
} MenuItemNode;

static void show_menu(GtkStatusIcon *status_icon, guint button, guint activate_time, gpointer user_data)
{
	gtk_widget_show_all(global_tray_menu);
	gtk_menu_popup(GTK_MENU(global_tray_menu), NULL, NULL, NULL, NULL, 0, gtk_get_current_event_time());
}

int nativeLoop(void) {
	gdk_threads_init();
	gtk_init(0, NULL);
	global_tray_menu = gtk_menu_new();
	systray_ready();
	gtk_main();
}


void setIcon(const char* iconBytes, int length) {
	GInputStream *stream = g_memory_input_stream_new_from_data(iconBytes, length, NULL);
	GError *error = NULL;
	GdkPixbuf *pixbuf = gdk_pixbuf_new_from_stream(stream, NULL, &error);
	if (error)
		fprintf(stderr, "Unable to create PixBuf: %s\n", error->message);
	gdk_threads_enter();
	global_tray_icon = gtk_status_icon_new_from_pixbuf(pixbuf);
	g_signal_connect(G_OBJECT(global_tray_icon), "activate", G_CALLBACK(show_menu), NULL);
	g_signal_connect(G_OBJECT(global_tray_icon), "popup-menu", G_CALLBACK(show_menu), NULL);
	gtk_status_icon_set_visible(global_tray_icon, TRUE);
	gdk_threads_leave();
}

void setTitle(char* ctitle) {
	gdk_threads_enter();
	gtk_status_icon_set_title(global_tray_icon, ctitle);
	gdk_threads_leave();
	free(ctitle);
}

void setTooltip(char* ctooltip) {
	gdk_threads_enter();
	gtk_status_icon_set_tooltip_text(global_tray_icon, ctooltip);
	gdk_threads_leave();
	free(ctooltip);
}

void addMenuItem(char* menuId, char* title, char* tooltip) {

	gdk_threads_enter();
	GList* it;
	for(it = global_menu_items; it != NULL; it = it->next) {
		MenuItemNode* item = (MenuItemNode*)(it->data);
		if (strcmp(item->menu_id, menuId) == 0){
			gtk_menu_item_set_label(GTK_MENU_ITEM(item->menu_item), title);
			break;
		}
	}

	// menu id doesn't exist, add new item
	if (it == NULL) {
		GtkWidget *titleMenuItem = gtk_menu_item_new_with_label(title);
		g_signal_connect_swapped(G_OBJECT(titleMenuItem), "activate", G_CALLBACK(systray_menu_item_selected), menuId);
		gtk_menu_shell_append(GTK_MENU_SHELL(global_tray_menu), titleMenuItem);

		MenuItemNode* new_item = malloc(sizeof(MenuItemNode));
		new_item->menu_id = menuId;
		new_item->menu_item = titleMenuItem;
		GList* new_node = malloc(sizeof(GList));
		new_node->data = new_item;
		new_node->next = global_menu_items;
		if (global_menu_items != NULL) {
			global_menu_items->prev = new_node;
		}
		global_menu_items = new_node;
	}

	gdk_threads_leave();

	free(title);
	free(tooltip);
}

void quit() {
	gdk_threads_enter();
	gtk_main_quit();
	gdk_threads_leave();
}
