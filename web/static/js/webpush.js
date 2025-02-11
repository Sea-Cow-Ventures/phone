async function setupWebPush() {
	try {
		// Check if service workers and push messaging are supported
		if (!('serviceWorker' in navigator) || !('PushManager' in window)) {
			console.log('Push messaging is not supported');
			return;
		}

		console.log('Starting service worker registration...');
		const registration = await navigator.serviceWorker.register('/static/sw.js');
		console.log(registration);
		console.log('Service Worker registered');

		// Check existing subscription
		const subscription = await registration.pushManager.getSubscription();
		if (subscription) {
			console.log('Already subscribed:', subscription);
			return subscription;
		}

		console.log('No existing subscription, creating new...');
		
		// Get the server's public key
		const publicKey = await ajaxLoad('/webpush/vapidkey');
		if (!publicKey) {
			throw new Error('No VAPID public key available');
		}

		// Create new subscription
		const newSubscription = await registration.pushManager.subscribe({
			userVisibleOnly: true,
			applicationServerKey: publicKey
		});

		console.log('Created new subscription:', newSubscription);

		// Send subscription to server
		await ajaxLoad('/webpush/subscribe', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
			},
			body: JSON.stringify({
				endpoint: newSubscription.endpoint,
				keys: newSubscription.keys,
				userAgent: navigator.userAgent
			})
		});

		console.log('Push notification subscription complete');
		return newSubscription;

	} catch (error) {
		console.error('Failed to setup push notification:', error);
		throw error;
	}
}