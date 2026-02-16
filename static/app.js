// In-memory device storage (simulating backend for demo)
let devices = [
    {
        id: '1',
        name: 'Router-BDG-01',
        ip_address: '192.168.100.11',
        location: 'Bandung',
        status: 'online',
        last_seen: new Date().toISOString(),
        created_at: new Date(Date.now() - 86400000 * 3).toISOString(),
        updated_at: new Date().toISOString()
    },
    {
        id: '2',
        name: 'Router-JKT-01',
        ip_address: '192.168.100.12',
        location: 'Jakarta',
        status: 'online',
        last_seen: new Date().toISOString(),
        created_at: new Date(Date.now() - 86400000 * 2).toISOString(),
        updated_at: new Date().toISOString()
    },
    {
        id: '3',
        name: 'Router-SBY-01',
        ip_address: '192.168.100.13',
        location: 'Surabaya',
        status: 'offline',
        last_seen: new Date(Date.now() - 3600000).toISOString(),
        created_at: new Date(Date.now() - 86400000 * 5).toISOString(),
        updated_at: new Date(Date.now() - 3600000).toISOString()
    }
];

// Initialize dashboard
document.addEventListener('DOMContentLoaded', () => {
    loadDevices();
    updateStats();
    
    // Setup form handler
    document.getElementById('add-device-form').addEventListener('submit', handleAddDevice);
    document.getElementById('search-input').addEventListener('input', loadDevices);
    document.getElementById('status-filter').addEventListener('change', loadDevices);
    
    // Auto-refresh every 30 seconds
    setInterval(() => {
        loadDevices();
        updateStats();
    }, 30000);
});

// Load and display devices
function loadDevices() {
    const tbody = document.getElementById('devices-table');
    const searchTerm = (document.getElementById('search-input')?.value || '').toLowerCase().trim();
    const statusFilter = document.getElementById('status-filter')?.value || 'all';
    const filteredDevices = devices.filter((device) => {
        const matchStatus = statusFilter === 'all' ? true : device.status === statusFilter;
        const matchSearch = !searchTerm
            ? true
            : device.name.toLowerCase().includes(searchTerm)
              || device.ip_address.toLowerCase().includes(searchTerm)
              || device.location.toLowerCase().includes(searchTerm);
        return matchStatus && matchSearch;
    });
    
    if (filteredDevices.length === 0) {
        tbody.innerHTML = `
            <tr>
                <td colspan="7" class="px-6 py-8 text-center text-gray-500">
                    <div class="text-4xl mb-2">📡</div>
                    <p>No matching devices found.</p>
                </td>
            </tr>
        `;
        return;
    }
    
    tbody.innerHTML = filteredDevices.map(device => `
        <tr class="hover:bg-gray-50">
            <td class="px-6 py-4">
                <div class="font-medium text-gray-900">${device.name}</div>
            </td>
            <td class="px-6 py-4 text-gray-600">${device.ip_address}</td>
            <td class="px-6 py-4 text-gray-600">${device.location}</td>
            <td class="px-6 py-4">
                <span class="px-3 py-1 rounded-full text-sm font-medium ${
                    device.status === 'online' 
                        ? 'bg-green-100 text-green-800' 
                        : 'bg-red-100 text-red-800'
                }">
                    ${device.status === 'online' ? '🟢' : '🔴'} ${device.status.toUpperCase()}
                </span>
            </td>
            <td class="px-6 py-4 text-gray-600 text-sm">
                ${formatTimeAgo(device.last_seen)}
            </td>
            <td class="px-6 py-4 text-gray-600 text-sm">
                ${formatDateTime(device.updated_at)}
            </td>
            <td class="px-6 py-4 text-gray-600 text-sm">
                ${formatDateTime(device.created_at)}
            </td>
        </tr>
    `).join('');
}

// Update statistics
function updateStats() {
    const total = devices.length;
    const online = devices.filter(d => d.status === 'online').length;
    const offline = total - online;
    
    document.getElementById('total-devices').textContent = total;
    document.getElementById('online-devices').textContent = online;
    document.getElementById('offline-devices').textContent = offline;
    document.getElementById('active-commands').textContent = '0';
}

// Show add device modal
function showAddDeviceModal() {
    document.getElementById('add-device-modal').classList.remove('hidden');
}

// Hide add device modal
function hideAddDeviceModal() {
    document.getElementById('add-device-modal').classList.add('hidden');
    document.getElementById('add-device-form').reset();
}

// Handle add device form submission
function handleAddDevice(e) {
    e.preventDefault();
    
    const newDevice = {
        id: Date.now().toString(),
        name: document.getElementById('device-name').value,
        ip_address: document.getElementById('device-ip').value,
        location: document.getElementById('device-location').value,
        status: 'online',
        last_seen: new Date().toISOString(),
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString()
    };
    
    devices.push(newDevice);
    loadDevices();
    updateStats();
    hideAddDeviceModal();
    
    // Show success message
    showNotification('✅ Device added successfully!');
}

// Bulk action handler
function bulkAction(action) {
    const onlineDevices = devices.filter(d => d.status === 'online');
    
    if (onlineDevices.length === 0) {
        showNotification('⚠️ No online devices to execute command', 'warning');
        return;
    }
    
    // Show execution panel
    const panel = document.getElementById('execution-panel');
    panel.classList.remove('hidden');
    
    // Simulate bulk execution
    let completed = 0;
    const total = onlineDevices.length;
    
    const interval = setInterval(() => {
        completed++;
        const progress = (completed / total) * 100;
        
        document.getElementById('progress-bar').style.width = progress + '%';
        document.getElementById('progress-percent').textContent = Math.round(progress) + '%';
        document.getElementById('progress-text').textContent = `Processing ${completed} of ${total} devices...`;
        document.getElementById('completed-count').textContent = completed;
        
        if (completed >= total) {
            clearInterval(interval);
            document.getElementById('progress-text').textContent = '✅ Operation completed!';
            showNotification(`✅ ${action.replace('_', ' ')} completed on ${total} devices`);
            
            // Hide panel after 5 seconds
            setTimeout(() => {
                panel.classList.add('hidden');
                document.getElementById('progress-bar').style.width = '0%';
            }, 5000);
        }
    }, 500);
}

// Refresh devices
function refreshDevices() {
    // Simulate API call
    showNotification('🔄 Refreshing devices...', 'info');
    
    setTimeout(() => {
        // Update last_seen for online devices
        devices.forEach(device => {
            if (device.status === 'online') {
                device.last_seen = new Date().toISOString();
                device.updated_at = new Date().toISOString();
            }
        });
        
        loadDevices();
        updateStats();
        showNotification('✅ Devices refreshed!');
    }, 1000);
}

// Format time ago
function formatTimeAgo(dateString) {
    const date = new Date(dateString);
    const now = new Date();
    const seconds = Math.floor((now - date) / 1000);
    
    if (seconds < 60) return 'Just now';
    if (seconds < 3600) return Math.floor(seconds / 60) + ' minutes ago';
    if (seconds < 86400) return Math.floor(seconds / 3600) + ' hours ago';
    return Math.floor(seconds / 86400) + ' days ago';
}

function formatDateTime(dateString) {
    if (!dateString) return '-';
    const date = new Date(dateString);
    return date.toLocaleString('en-US', {
        year: 'numeric',
        month: 'short',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit'
    });
}

// Show notification
function showNotification(message, type = 'success') {
    const notification = document.createElement('div');
    notification.className = `fixed right-4 top-4 z-50 rounded-xl border px-4 py-3 text-sm font-medium shadow-sm transition-all ${
        type === 'success' ? 'border-green-200 bg-green-50 text-green-700' :
        type === 'warning' ? 'border-amber-200 bg-amber-50 text-amber-700' :
        type === 'error' ? 'border-red-200 bg-red-50 text-red-700' :
        'border-blue-200 bg-blue-50 text-blue-700'
    }`;
    notification.textContent = message;
    
    document.body.appendChild(notification);
    
    setTimeout(() => {
        notification.style.opacity = '0';
        notification.style.transform = 'translateY(-8px)';
        setTimeout(() => notification.remove(), 180);
    }, 2600);
}

// Update clock
function updateClock() {
    const now = new Date();
    const timeString = now.toLocaleTimeString(undefined, {
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit'
    });
    const dateString = now.toLocaleDateString(undefined, {
        weekday: 'short',
        day: '2-digit',
        month: 'short'
    });
    const clockEl = document.getElementById('current-time');
    if (clockEl) clockEl.textContent = `🕒 ${dateString} • ${timeString}`;
}

// Initialize clock
setInterval(updateClock, 1000);
updateClock();
