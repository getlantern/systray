#include <stdlib.h>
#include <gtk/gtk.h>
#include <webkit2/webkit2.h>
#include "webview.h"

static GtkWindow *web_window = NULL;
static WebKitWebView *web_view = NULL;

static gint x, y;
static bool needsMove = false;

gboolean on_window_deleted(GtkWidget *window, GdkEvent *event, gpointer data)
{
    gtk_window_get_position(GTK_WINDOW(window), &x, &y);
    needsMove = true;
    gtk_widget_hide(window);
    return TRUE;
}

void configureAppWindow(char* title, int width, int height)
{
    // Create an 800x600 window that will contain the browser instance
    web_window = GTK_WINDOW(gtk_window_new(GTK_WINDOW_TOPLEVEL));
    gtk_window_set_title(web_window, title);
    gtk_window_set_default_size(web_window, width, height);
    gtk_window_set_skip_taskbar_hint (web_window, TRUE);
    g_signal_connect(G_OBJECT(web_window), "delete-event", G_CALLBACK(on_window_deleted), NULL);

    // Create a browser instance
    web_view = WEBKIT_WEB_VIEW(webkit_web_view_new());

    // Put the browser area into the web window
    gtk_container_add(GTK_CONTAINER(web_window), GTK_WIDGET(web_view));

    // Make sure that when the browser area becomes visible, it will get mouse
    // and keyboard events
    gtk_widget_grab_focus(GTK_WIDGET(web_view));
    free(title);
    gtk_main();
}

gboolean do_show_app_window(gpointer data)
{
    gtk_widget_show_all(GTK_WIDGET(web_window));
    if (needsMove) {
        gtk_window_move(web_window, x, y);
        needsMove = false;
    }
    gtk_window_present(web_window);
    return FALSE;
}

void showAppWindow(char* url)
{
    // Load a web page into the browser instance
    webkit_web_view_load_uri(web_view, url);

    g_idle_add(do_show_app_window, NULL);
    free(url);
}
