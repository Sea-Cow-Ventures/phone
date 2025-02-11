// When the service worker is installed
self.addEventListener('install', (event) => {
    console.log('Service Worker installed');
});

// When the service worker is activated
self.addEventListener('activate', (event) => {
    console.log('Service Worker activated');
});

// Handle incoming push messages
self.addEventListener('push', (event) => {
    if (!event.data) {
        console.log('Push event but no data');
        return;
    }

    try {
        const data = event.data.json();
        
        const options = {
            body: data.body || 'New notification',
            icon: '/static/images/icon.png',  // Your notification icon
            badge: '/static/images/badge.png', // Small icon for notification tray
            vibrate: [100, 50, 100], // Vibration pattern
            data: {
                url: data.url || '/',  // URL to open when notification is clicked
                timestamp: Date.now()
            },
            actions: [
                {
                    action: 'open',
                    title: 'Open'
                },
                {
                    action: 'close',
                    title: 'Close'
                }
            ],
            requireInteraction: true  // Notification stays until user interacts
        };

        event.waitUntil(
            self.registration.showNotification(data.title || 'New Message', options)
        );
    } catch (error) {
        console.error('Error showing notification:', error);
    }
});

// Handle notification clicks
self.addEventListener('notificationclick', (event) => {
    event.notification.close();

    // Handle action buttons
    if (event.action === 'open') {
        // Custom open action
        const url = event.notification.data.url;
        event.waitUntil(
            clients.openWindow(url)
        );
    }
    // Default click behavior
    else if (!event.action) {
        // Get notification data
        const data = event.notification.data;
        
        event.waitUntil(
            clients.matchAll({
                type: 'window'
            })
            .then((clientList) => {
                // If a window exists, focus it
                for (const client of clientList) {
                    if (client.url === data.url && 'focus' in client) {
                        return client.focus();
                    }
                }
                // If no window exists, open new one
                if (clients.openWindow) {
                    return clients.openWindow(data.url);
                }
            })
        );
    }
});

// Handle notification close
self.addEventListener('notificationclose', (event) => {
    console.log('Notification was closed', event.notification);
});