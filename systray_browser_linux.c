#include <gtk/gtk.h>
#include <webkit2/webkit2.h>

static GtkWidget *web_window = NULL;
static WebKitWebView *webView = NULL;

static gboolean closeWebViewCb(WebKitWebView* webView, GtkWidget* window);

void prepareBrowser()
{
    // Create an 800x600 window that will contain the browser instance
    web_window = gtk_window_new(GTK_WINDOW_TOPLEVEL);
    gtk_window_set_default_size(GTK_WINDOW(web_window), 800, 600);

    // Create a browser instance
    webView = WEBKIT_WEB_VIEW(webkit_web_view_new());

    // Put the browser area into the web window
    gtk_container_add(GTK_CONTAINER(web_window), GTK_WIDGET(webView));

    g_signal_connect(webView, "close", G_CALLBACK(closeWebViewCb), web_window);

    // Make sure that when the browser area becomes visible, it will get mouse
    // and keyboard events
    gtk_widget_grab_focus(GTK_WIDGET(webView));
}

gboolean do_open_in_browser(gpointer data)
{
    gtk_widget_show_all(web_window);
    return TRUE;
}

void openInBrowser(char* url)
{
        // Load a web page into the browser instance
    webkit_web_view_load_uri(webView, url);

    gdk_threads_add_idle(do_open_in_browser, url);
}

static gboolean closeWebViewCb(WebKitWebView* webView, GtkWidget* window)
{
    gtk_widget_destroy(window);
    return TRUE;
}