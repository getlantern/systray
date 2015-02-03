// +build ignore

#include <stdlib.h>
#include <gtk/gtk.h>
#include <libappindicator/app-indicator.h>
#include "systray.h"

static AppIndicator *indicator;
static GtkWidget *indicator_menu;
static void activate_action (GtkAction *action);

static GtkActionEntry entries[] = {
	{ "FileMenu", NULL, "_File" },
	{ "New",      "document-new", "_New", "<control>N",
		"Create a new file", G_CALLBACK (activate_action) },
	{ "Open",     "document-open", "_Open", "<control>O",
		"Open a file", G_CALLBACK (activate_action) },
	{ "Save",     "document-save", "_Save", "<control>S",
		"Save file", G_CALLBACK (activate_action) },
	{ "Quit",     "application-exit", "_Quit", "<control>Q",
		"Exit the application", G_CALLBACK (activate_action) },
};
static guint n_entries = G_N_ELEMENTS (entries);


static const gchar *ui_info =
"<ui>"
"  <menubar name='MenuBar'>"
"    <menu action='FileMenu'>"
"      <menuitem action='New'/>"
"      <menuitem action='Open'/>"
"      <menuitem action='Save'/>"
"      <separator/>"
"      <menuitem action='Quit'/>"
"    </menu>"
"  </menubar>"
"  <popup name='IndicatorPopup'>"
"    <menuitem action='New' />"
"    <menuitem action='Open' />"
"    <menuitem action='Save' />"
"    <menuitem action='Quit' />"
"  </popup>"
"</ui>";

	static void
activate_action (GtkAction *action)
{
	const gchar *name = gtk_action_get_name (action);
	GtkWidget *dialog;

	dialog = gtk_message_dialog_new (NULL,
			GTK_DIALOG_DESTROY_WITH_PARENT,
			GTK_MESSAGE_INFO,
			GTK_BUTTONS_CLOSE,
			"You activated action: \"%s\"",
			name);

	g_signal_connect (dialog, "response",
			G_CALLBACK (gtk_widget_destroy), NULL);

	gtk_widget_show (dialog);
}


int nativeLoop(void) {
	gtk_init(0, NULL);
	indicator = app_indicator_new("git_indicator", "icon.png",
			APP_INDICATOR_CATEGORY_APPLICATION_STATUS);
	app_indicator_set_status (indicator, APP_INDICATOR_STATUS_ACTIVE);
	indicator_menu = gtk_menu_new();
	app_indicator_set_menu (indicator, GTK_MENU (indicator_menu));
	systray_ready();
	gtk_main();
}


void setIcon(const char* iconBytes, int length) {
	app_indicator_set_icon_full (indicator, "icon.png", "");
	app_indicator_set_attention_icon_full (indicator, "icon.png", "");
	app_indicator_set_label (indicator, "", "Git Indicator");
	printf("icon created\n");
	printf("current icon: %s\n", app_indicator_get_icon(indicator));
}

void setTitle(char* ctitle) {
	printf("set title\n");
	app_indicator_set_label (indicator, "", ctitle);
	free(ctitle);
}

void setTooltip(char* ctooltip) {
	printf("set tooltip\n");
	free(ctooltip);
}

void addMenuItem(char* menuId, char* title, char* tooltip) {
	GtkWidget *item = (GtkWidget *) malloc (sizeof (GtkWidget*));
	item = gtk_menu_item_new_with_label (title);
	gtk_menu_shell_append (GTK_MENU_SHELL (indicator_menu), item);

	GtkWidget *sep = gtk_separator_menu_item_new ();
	gtk_menu_shell_append (GTK_MENU_SHELL (indicator_menu), sep);
	gtk_widget_show_all (indicator_menu);
	printf("addMenuItem\n");

	// free(menuId);
	free(title);
	free(tooltip);
}

void quit() {
}
