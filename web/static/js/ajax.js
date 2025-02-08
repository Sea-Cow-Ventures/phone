async function ajaxLoad(url, options = {}) {
	const defaultOptions = {
		method: 'GET',
		headers: {
			'Content-Type': 'application/json'
		}
	};

	// Merge default options with provided options
	const finalOptions = {
		...defaultOptions,
		...options,
		headers: {
			...defaultOptions.headers,
			...options.headers
		}
	};

	// Convert body to JSON string if it's an object
	if (finalOptions.body && typeof finalOptions.body === 'object') {
		finalOptions.body = JSON.stringify(finalOptions.body);
	}

	try {
		const response = await fetch(url, finalOptions);
		const data = await response.json();
		console.log(data);
		if (!response.ok || !data.success) {
			throw new Error(`<strong>Status:</strong> ${response.status}\n<strong>Error:</strong> ${data.error || 'Operation failed'}`);
		}

		return data.data;
	} catch (error) {
		const errorMessage = error.message || 'An unexpected error occurred';
		console.error('Ajax Error:', {
			url,
			error: errorMessage
		});
		showErrorModal(errorMessage, url);
		return null;
	}
}

function showErrorModal(message, endpoint) {
	if (!document.getElementById('errorModal')) {
		const modalHtml = `
			<div id="errorModal" class="error-modal">
				<div class="error-modal-content relative">
					<button onclick="closeErrorModal()" class="absolute top-2 right-2 text-gray-500 hover:text-gray-700">
						<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
						</svg>
					</button>
					<div class="text-xl font-semibold mb-4">Error</div>
					<div id="errorMessage" class="mb-2 whitespace-pre-wrap"></div>
					<div id="errorEndpoint" class="text-sm text-gray-500"></div>
				</div>
			</div>
		`;
		document.body.insertAdjacentHTML('beforeend', modalHtml);
	}

	const modal = document.getElementById('errorModal');
	const messageElement = document.getElementById('errorMessage');
	const endpointElement = document.getElementById('errorEndpoint');
	
	messageElement.innerHTML = message;
	endpointElement.textContent = `Endpoint: ${endpoint}`;
	modal.style.display = 'flex';
}

function closeErrorModal() {
	const modal = document.getElementById('errorModal');
	
	if (modal) {
		modal.style.display = 'none';
	}
}