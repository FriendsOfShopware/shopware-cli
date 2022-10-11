let eventSource = new EventSource('/__internal-admin-proxy/events');

eventSource.onmessage = function (message) {
    if (message.data === 'reloadCss') {
        document.querySelectorAll('link').forEach(link => {
            if (link.href.indexOf('extension.css') !== -1) {
                const newURL = new URL(link.href)
        
                newURL.searchParams.set('t', new Date().getTime())
        
                link.href = newURL.toString()
            }
        });
    } else {
        window.location.reload();
    }
}