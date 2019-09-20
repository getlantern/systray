#include <gtk/gtk.h>
#include <webkit2/webkit2.h>

static GtkWidget *web_window = NULL;
static WebKitWebView *webView = NULL;

static gint x, y;
static bool hasPosition = false;

gboolean on_window_deleted(GtkWidget *window, GdkEvent *event, gpointer data)
{
    gtk_window_get_position(GTK_WINDOW(window), &x, &y);
    hasPosition = true;
    gtk_widget_hide(window);
    return TRUE;
}

void configureAppWindow(char* title, int width, int height)
{
    // Create an 800x600 window that will contain the browser instance
    web_window = gtk_window_new(GTK_WINDOW_TOPLEVEL);
    gtk_window_set_title(GTK_WINDOW(web_window), title);
    gtk_window_set_default_size(GTK_WINDOW(web_window), width, height);
    g_signal_connect(G_OBJECT(web_window), "delete-event", G_CALLBACK(on_window_deleted), NULL);

    // Create a browser instance
    webView = WEBKIT_WEB_VIEW(webkit_web_view_new());

    // Put the browser area into the web window
    gtk_container_add(GTK_CONTAINER(web_window), GTK_WIDGET(webView));

    // Make sure that when the browser area becomes visible, it will get mouse
    // and keyboard events
    gtk_widget_grab_focus(GTK_WIDGET(webView));
    free(title);
}

gboolean do_show_app_window(gpointer data)
{
    gtk_widget_show_all(web_window);
    if (hasPosition) {
        gtk_window_move(GTK_WINDOW(web_window), x, y);
    }
    return FALSE;
}

void showAppWindow(char* url)
{
    // Load a web page into the browser instance
    webkit_web_view_load_uri(webView, url);

    gdk_threads_add_idle(do_show_app_window, NULL);
    free(url);
}