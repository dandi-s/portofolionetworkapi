let devices = [
    {
        id: '1',
        name: 'Router-BDG-01',
        ip_address: '192.168.100.11',
        location: 'Bandung',
        status: 'online',
        last_seen: new Date().toISOString(),
        last_downtime: null,
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
        last_downtime: null,
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
        last_downtime: new Date(Date.now() - 3600000).toISOString(),
        created_at: new Date(Date.now() - 86400000 * 5).toISOString(),
        updated_at: new Date(Date.now() - 3600000).toISOString()
    }
];

const selectedDeviceIds = new Set();
let currentPage = 1;
const PAGE_SIZE = 50;
const DEVICE_LIMIT = 25;
const DEVICE_LIMIT_WARNING_RATIO = 0.8;
const LIMIT_RESET_MS = 60 * 60 * 1000;
const LIMIT_RESET_KEY = 'netops_device_limit_reset_at';
let deviceLimitWarningShown = false;
let deviceLimitReachedShown = false;
let isAutoResetting = false;

document.addEventListener('DOMContentLoaded', () => {
    loadDevices();
    updateStats();

    const addDeviceForm = document.getElementById('add-device-form');
    const editDeviceForm = document.getElementById('edit-device-form');
    const bulkEditForm = document.getElementById('bulk-edit-form');
    const searchInput = document.getElementById('search-input');
    const statusFilter = document.getElementById('status-filter');
    const openAddDeviceBtn = document.getElementById('open-add-device-btn');
    const selectAllCheckbox = document.getElementById('select-all-devices');
    const tableBody = document.getElementById('devices-table');
    const selectionMenuToggle = document.getElementById('selection-menu-toggle');

    if (addDeviceForm) addDeviceForm.addEventListener('submit', handleAddDevice);
    if (editDeviceForm) editDeviceForm.addEventListener('submit', handleEditDevice);
    if (bulkEditForm) bulkEditForm.addEventListener('submit', handleBulkEdit);

    if (searchInput) {
        searchInput.addEventListener('input', () => {
            currentPage = 1;
            loadDevices();
        });
    }

    if (statusFilter) {
        statusFilter.addEventListener('change', () => {
            currentPage = 1;
            loadDevices();
        });
    }

    if (openAddDeviceBtn) openAddDeviceBtn.addEventListener('click', showAddDeviceModal);

    if (selectAllCheckbox) {
        selectAllCheckbox.addEventListener('change', (event) => {
            const isChecked = event.target.checked;
            const pageDevices = getPaginatedDevices().pageDevices;
            pageDevices.forEach((device) => {
                if (isChecked) {
                    selectedDeviceIds.add(device.id);
                } else {
                    selectedDeviceIds.delete(device.id);
                }
            });
            loadDevices();
        });
    }

    if (tableBody) {
        tableBody.addEventListener('change', (event) => {
            const target = event.target;
            if (!target.classList.contains('device-row-checkbox')) return;
            const deviceId = target.getAttribute('data-device-id');
            if (!deviceId) return;

            if (target.checked) {
                selectedDeviceIds.add(deviceId);
            } else {
                selectedDeviceIds.delete(deviceId);
            }

            updateSelectionControls();
            updateSelectAllState();
        });
    }

    if (selectionMenuToggle) {
        selectionMenuToggle.addEventListener('click', (event) => {
            event.stopPropagation();
            const menu = document.getElementById('selection-menu');
            if (menu) menu.classList.toggle('hidden');
        });
    }

    document.addEventListener('click', (event) => {
        const actions = document.getElementById('selection-actions');
        const menu = document.getElementById('selection-menu');
        if (!actions || !menu) return;
        if (!actions.contains(event.target)) {
            menu.classList.add('hidden');
        }
    });

    setInterval(() => {
        loadDevices();
        updateStats();
    }, 30000);

    setInterval(() => {
        checkDeviceLimit({ total: devices.length, limit: DEVICE_LIMIT, warningRatio: DEVICE_LIMIT_WARNING_RATIO });
    }, 1000);
});

function getFilteredDevices() {
    const searchTerm = (document.getElementById('search-input')?.value || '').toLowerCase().trim();
    const statusFilter = document.getElementById('status-filter')?.value || 'all';

    return devices.filter((device) => {
        const matchStatus = statusFilter === 'all' ? true : device.status === statusFilter;
        const matchSearch = !searchTerm
            ? true
            : device.name.toLowerCase().includes(searchTerm)
                || device.ip_address.toLowerCase().includes(searchTerm)
                || device.location.toLowerCase().includes(searchTerm);
        return matchStatus && matchSearch;
    });
}

function getPaginatedDevices() {
    const filtered = getFilteredDevices();
    const totalPages = Math.max(1, Math.ceil(filtered.length / PAGE_SIZE));

    if (currentPage > totalPages) currentPage = totalPages;
    if (currentPage < 1) currentPage = 1;

    const startIndex = (currentPage - 1) * PAGE_SIZE;
    const endIndex = Math.min(startIndex + PAGE_SIZE, filtered.length);
    const pageDevices = filtered.slice(startIndex, endIndex);

    return {
        filtered,
        pageDevices,
        startIndex,
        endIndex,
        totalPages
    };
}

function loadDevices() {
    const tbody = document.getElementById('devices-table');
    if (!tbody) return;

    const { filtered, pageDevices, startIndex, endIndex, totalPages } = getPaginatedDevices();

    if (filtered.length === 0) {
        tbody.innerHTML = `
            <tr>
                <td colspan="9" class="px-6 py-8 text-center text-gray-500">
                    <div class="text-4xl mb-2">📡</div>
                    <p>No matching devices found.</p>
                </td>
            </tr>
        `;
        updatePaginationControls(0, 0, 0, 1);
        updateSelectionControls();
        updateSelectAllState();
        return;
    }

    tbody.innerHTML = pageDevices.map(device => `
        <tr class="hover:bg-gray-50">
            <td class="px-4 py-4">
                <input
                    type="checkbox"
                    class="device-row-checkbox h-4 w-4 rounded border-gray-300"
                    data-device-id="${device.id}"
                    ${selectedDeviceIds.has(device.id) ? 'checked' : ''}
                >
            </td>
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
            <td class="px-6 py-4 text-gray-600 text-sm">${formatTimeAgo(device.last_seen)}</td>
            <td class="px-6 py-4 text-gray-600 text-sm">${formatDateTime(device.last_downtime)}</td>
            <td class="px-6 py-4 text-gray-600 text-sm">${formatDateTime(device.updated_at)}</td>
            <td class="px-6 py-4 text-gray-600 text-sm">${formatDateTime(device.created_at)}</td>
        </tr>
    `).join('');

    updatePaginationControls(startIndex + 1, endIndex, filtered.length, totalPages);
    updateSelectionControls();
    updateSelectAllState();
}

function updatePaginationControls(start, end, total, totalPages) {
    const info = document.getElementById('pagination-info');
    const indicator = document.getElementById('page-indicator');
    const prevBtn = document.getElementById('prev-page-btn');
    const nextBtn = document.getElementById('next-page-btn');

    if (info) info.textContent = `Showing ${start} to ${end} of ${total} devices`;
    if (indicator) indicator.textContent = `Page ${currentPage} of ${totalPages}`;
    if (prevBtn) prevBtn.disabled = currentPage <= 1;
    if (nextBtn) nextBtn.disabled = currentPage >= totalPages;
}

function goToPreviousPage() {
    currentPage = Math.max(1, currentPage - 1);
    loadDevices();
}

function goToNextPage() {
    const { totalPages } = getPaginatedDevices();
    currentPage = Math.min(totalPages, currentPage + 1);
    loadDevices();
}

function updateSelectionControls() {
    const selectedCountEl = document.getElementById('selected-count');
    const selectionActions = document.getElementById('selection-actions');
    const selectionMenu = document.getElementById('selection-menu');
    const bulkEditBtn = document.getElementById('open-bulk-edit-btn');
    const bulkDeleteBtn = document.getElementById('bulk-delete-btn');

    const selectedCount = selectedDeviceIds.size;
    if (selectedCountEl) selectedCountEl.textContent = selectedCount;

    const hasSelection = selectedCount > 0;
    if (selectionActions) selectionActions.classList.toggle('hidden', !hasSelection);
    if (!hasSelection && selectionMenu) selectionMenu.classList.add('hidden');

    if (bulkEditBtn) bulkEditBtn.disabled = !hasSelection;
    if (bulkDeleteBtn) bulkDeleteBtn.disabled = !hasSelection;
}

function updateSelectAllState() {
    const selectAll = document.getElementById('select-all-devices');
    if (!selectAll) return;

    const { pageDevices } = getPaginatedDevices();
    if (pageDevices.length === 0) {
        selectAll.checked = false;
        selectAll.indeterminate = false;
        return;
    }

    const selectedOnPage = pageDevices.filter((device) => selectedDeviceIds.has(device.id)).length;
    selectAll.checked = selectedOnPage === pageDevices.length;
    selectAll.indeterminate = selectedOnPage > 0 && selectedOnPage < pageDevices.length;
}

function updateStats() {
    const total = devices.length;
    const online = devices.filter(d => d.status === 'online').length;
    const offline = total - online;

    document.getElementById('total-devices').textContent = total;
    document.getElementById('online-devices').textContent = online;
    document.getElementById('offline-devices').textContent = offline;
    document.getElementById('active-commands').textContent = '0';

    checkDeviceLimit({ total, limit: DEVICE_LIMIT, warningRatio: DEVICE_LIMIT_WARNING_RATIO });
}

function toLocalInputValue(date) {
    const d = new Date(date);
    return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}T${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`;
}

function showAddDeviceModal() {
    const modal = document.getElementById('add-device-modal');
    if (!modal) return;
    modal.classList.remove('hidden');

    const now = new Date();
    const localValue = toLocalInputValue(now);

    const lastSeen = document.getElementById('device-last-seen');
    const updatedAt = document.getElementById('device-updated-at');
    const createdAt = document.getElementById('device-created-at');

    if (lastSeen) lastSeen.value = localValue;
    if (updatedAt) updatedAt.value = localValue;
    if (createdAt) createdAt.value = localValue;
}

function hideAddDeviceModal() {
    const modal = document.getElementById('add-device-modal');
    const form = document.getElementById('add-device-form');
    if (modal) modal.classList.add('hidden');
    if (form) form.reset();
}

function handleAddDevice(e) {
    e.preventDefault();

    const limitState = checkDeviceLimit({ total: devices.length, limit: DEVICE_LIMIT });
    if (limitState.reached) {
        showNotification(`⚠️ Device limit reached (${limitState.limit}). Delete old devices before adding new ones.`, 'warning');
        return;
    }

    const nowIso = new Date().toISOString();
    const parseDateTimeLocal = (value) => (value ? new Date(value).toISOString() : null);

    const newDevice = {
        id: Date.now().toString(),
        name: document.getElementById('device-name').value,
        ip_address: document.getElementById('device-ip').value,
        location: document.getElementById('device-location').value,
        status: document.getElementById('device-status').value || 'online',
        last_seen: parseDateTimeLocal(document.getElementById('device-last-seen').value) || nowIso,
        last_downtime: null,
        created_at: parseDateTimeLocal(document.getElementById('device-created-at').value) || nowIso,
        updated_at: parseDateTimeLocal(document.getElementById('device-updated-at').value) || nowIso
    };

    if (newDevice.status === 'offline') {
        newDevice.last_downtime = newDevice.updated_at || nowIso;
    }

    devices.push(newDevice);
    currentPage = Math.ceil(devices.length / PAGE_SIZE);
    loadDevices();
    updateStats();
    hideAddDeviceModal();
    showNotification('✅ Device added successfully!');
}

// Show limit warning if needed
function checkDeviceLimit(meta = {}) {
    const total = Number.isFinite(meta.total) ? meta.total : devices.length;
    const limit = Number.isFinite(meta.limit) ? meta.limit : DEVICE_LIMIT;
    const warningRatio = Number.isFinite(meta.warningRatio) ? meta.warningRatio : DEVICE_LIMIT_WARNING_RATIO;
    const warningAt = Math.max(1, Math.floor(limit * warningRatio));
    const reached = total >= limit;

    if (reached) {
        ensureLimitResetTimer();
    } else {
        clearLimitResetTimer();
    }

    const resetAt = getLimitResetAt();
    const remainingMs = resetAt ? Math.max(0, resetAt - Date.now()) : 0;

    if (reached && remainingMs === 0 && resetAt && !isAutoResetting) {
        runDeviceLimitReset();
        return { total, limit, warningAt, reached: false, resetAt: null };
    }

    updateLimitUI({ total, limit, warningAt, reached, remainingMs, resetAt });

    if (reached) {
        if (!deviceLimitReachedShown) {
            showNotification(`⚠️ Device limit reached (${total}/${limit}). Auto reset in 1 hour.`, 'warning');
            deviceLimitReachedShown = true;
        }
        return { total, limit, warningAt, reached: true, resetAt };
    }

    deviceLimitReachedShown = false;

    if (total >= warningAt) {
        if (!deviceLimitWarningShown) {
            showNotification(`⚠️ Near device limit (${total}/${limit}).`, 'warning');
            deviceLimitWarningShown = true;
        }
    } else {
        deviceLimitWarningShown = false;
    }

    return { total, limit, warningAt, reached: false, resetAt: null };
}

function updateLimitUI(state) {
    const addBtn = document.getElementById('open-add-device-btn');
    const banner = document.getElementById('device-limit-banner');
    const title = document.getElementById('device-limit-title');
    const message = document.getElementById('device-limit-message');
    const countdown = document.getElementById('device-limit-countdown');

    if (addBtn) {
        addBtn.disabled = state.reached;
        addBtn.classList.toggle('opacity-50', state.reached);
        addBtn.classList.toggle('cursor-not-allowed', state.reached);
        addBtn.title = state.reached
            ? `Limit reached (${state.total}/${state.limit})`
            : `Add device (${state.total}/${state.limit})`;
    }

    if (!banner || !title || !message || !countdown) return;

    if (state.reached) {
        banner.classList.remove('hidden', 'border-amber-200', 'bg-amber-50');
        banner.classList.add('border-red-200', 'bg-red-50');
        title.className = 'text-sm font-semibold text-red-800';
        message.className = 'text-sm text-red-700';
        countdown.className = 'rounded-md border border-red-300 bg-white px-2.5 py-1 text-xs font-semibold text-red-800';
        title.textContent = 'Device limit reached';
        message.textContent = `Maximum ${state.limit} devices. System will auto-reset and restore 3 dummy devices.`;
        countdown.textContent = formatCountdown(state.remainingMs);
        return;
    }

    if (state.total >= state.warningAt) {
        banner.classList.remove('hidden', 'border-red-200', 'bg-red-50');
        banner.classList.add('border-amber-200', 'bg-amber-50');
        title.className = 'text-sm font-semibold text-amber-800';
        message.className = 'text-sm text-amber-700';
        countdown.className = 'rounded-md border border-amber-300 bg-white px-2.5 py-1 text-xs font-semibold text-amber-800';
        title.textContent = 'Approaching device limit';
        message.textContent = `You are using ${state.total} of ${state.limit} devices.`;
        countdown.textContent = `${state.limit - state.total} slots left`;
        return;
    }

    banner.classList.add('hidden');
}

function ensureLimitResetTimer() {
    if (getLimitResetAt()) return;
    setLimitResetAt(Date.now() + LIMIT_RESET_MS);
}

function clearLimitResetTimer() {
    localStorage.removeItem(LIMIT_RESET_KEY);
}

function getLimitResetAt() {
    const raw = localStorage.getItem(LIMIT_RESET_KEY);
    if (!raw) return null;
    const timestamp = Number(raw);
    return Number.isFinite(timestamp) ? timestamp : null;
}

function setLimitResetAt(timestamp) {
    localStorage.setItem(LIMIT_RESET_KEY, String(timestamp));
}

function runDeviceLimitReset() {
    isAutoResetting = true;
    clearLimitResetTimer();
    selectedDeviceIds.clear();
    currentPage = 1;
    deviceLimitReachedShown = false;
    deviceLimitWarningShown = false;
    devices = buildResetDummyDevices();
    loadDevices();
    updateStats();
    showNotification('✅ Protection reset complete. 3 dummy devices restored.', 'success');
    isAutoResetting = false;
}

function buildResetDummyDevices() {
    const now = Date.now();
    return [
        {
            id: `reset-${now}-1`,
            name: 'Router-RESET-01',
            ip_address: '10.10.10.11',
            location: 'Jakarta',
            status: 'online',
            last_seen: new Date(now).toISOString(),
            last_downtime: null,
            created_at: new Date(now - 86400000).toISOString(),
            updated_at: new Date(now).toISOString()
        },
        {
            id: `reset-${now}-2`,
            name: 'Switch-RESET-01',
            ip_address: '10.10.10.12',
            location: 'Bandung',
            status: 'online',
            last_seen: new Date(now).toISOString(),
            last_downtime: null,
            created_at: new Date(now - 86400000 * 2).toISOString(),
            updated_at: new Date(now).toISOString()
        },
        {
            id: `reset-${now}-3`,
            name: 'Firewall-RESET-01',
            ip_address: '10.10.10.13',
            location: 'Surabaya',
            status: 'offline',
            last_seen: new Date(now - 3600000).toISOString(),
            last_downtime: new Date(now - 3600000).toISOString(),
            created_at: new Date(now - 86400000 * 3).toISOString(),
            updated_at: new Date(now - 3600000).toISOString()
        }
    ];
}

function formatCountdown(ms) {
    const totalSeconds = Math.max(0, Math.floor(ms / 1000));
    const hours = String(Math.floor(totalSeconds / 3600)).padStart(2, '0');
    const minutes = String(Math.floor((totalSeconds % 3600) / 60)).padStart(2, '0');
    const seconds = String(totalSeconds % 60).padStart(2, '0');
    return `${hours}:${minutes}:${seconds}`;
}

function closeSelectionMenu() {
    const menu = document.getElementById('selection-menu');
    if (menu) menu.classList.add('hidden');
}

function showEditDeviceModal() {
    closeSelectionMenu();
    if (selectedDeviceIds.size === 0) {
        showNotification('⚠️ Select one device to edit', 'warning');
        return;
    }

    if (selectedDeviceIds.size > 1) {
        showBulkEditModal();
        return;
    }

    const selectedId = Array.from(selectedDeviceIds)[0];
    const device = devices.find((item) => item.id === selectedId);
    if (!device) {
        showNotification('⚠️ Selected device was not found', 'warning');
        return;
    }

    const modal = document.getElementById('edit-device-modal');
    if (modal) modal.classList.remove('hidden');

    document.getElementById('edit-device-name').value = device.name;
    document.getElementById('edit-device-ip').value = device.ip_address;
    document.getElementById('edit-device-location').value = device.location;
    document.getElementById('edit-device-status').value = device.status;
    document.getElementById('edit-device-last-seen').value = device.last_seen ? toLocalInputValue(device.last_seen) : '';
    document.getElementById('edit-device-last-downtime').value = device.last_downtime ? toLocalInputValue(device.last_downtime) : '';
    document.getElementById('edit-device-updated-at').value = device.updated_at ? toLocalInputValue(device.updated_at) : toLocalInputValue(new Date());
    document.getElementById('edit-device-created-at').value = device.created_at ? toLocalInputValue(device.created_at) : toLocalInputValue(new Date());
}

function hideEditDeviceModal() {
    const modal = document.getElementById('edit-device-modal');
    const form = document.getElementById('edit-device-form');
    if (modal) modal.classList.add('hidden');
    if (form) form.reset();
}

function handleEditDevice(event) {
    event.preventDefault();

    if (selectedDeviceIds.size !== 1) {
        showNotification('⚠️ Please select exactly one device to edit', 'warning');
        return;
    }

    const selectedId = Array.from(selectedDeviceIds)[0];
    const parseDateTimeLocal = (value) => (value ? new Date(value).toISOString() : null);
    const nowIso = new Date().toISOString();

    const payload = {
        name: document.getElementById('edit-device-name').value.trim(),
        ip_address: document.getElementById('edit-device-ip').value.trim(),
        location: document.getElementById('edit-device-location').value.trim(),
        status: document.getElementById('edit-device-status').value || 'online',
        last_seen: parseDateTimeLocal(document.getElementById('edit-device-last-seen').value),
        last_downtime: parseDateTimeLocal(document.getElementById('edit-device-last-downtime').value),
        updated_at: parseDateTimeLocal(document.getElementById('edit-device-updated-at').value) || nowIso,
        created_at: parseDateTimeLocal(document.getElementById('edit-device-created-at').value)
    };

    devices = devices.map((device) => {
        if (device.id !== selectedId) return device;

        const next = {
            ...device,
            name: payload.name || device.name,
            ip_address: payload.ip_address || device.ip_address,
            location: payload.location || device.location,
            status: payload.status,
            last_seen: payload.last_seen || device.last_seen,
            updated_at: payload.updated_at,
            created_at: payload.created_at || device.created_at
        };

        if (payload.last_downtime) {
            next.last_downtime = payload.last_downtime;
        } else if (device.status !== 'offline' && payload.status === 'offline') {
            next.last_downtime = payload.updated_at || nowIso;
        }

        return next;
    });

    hideEditDeviceModal();
    loadDevices();
    updateStats();
    showNotification('✅ Device updated successfully');
}

function showBulkEditModal() {
    const modal = document.getElementById('bulk-edit-modal');
    if (modal) modal.classList.remove('hidden');

    const statusEl = document.getElementById('bulk-edit-status');
    const locationEl = document.getElementById('bulk-edit-location');
    if (statusEl) statusEl.value = 'keep';
    if (locationEl) locationEl.value = '';
}

function hideBulkEditModal() {
    const modal = document.getElementById('bulk-edit-modal');
    const form = document.getElementById('bulk-edit-form');
    if (modal) modal.classList.add('hidden');
    if (form) form.reset();
}

function handleBulkEdit(event) {
    event.preventDefault();

    if (selectedDeviceIds.size < 2) {
        hideBulkEditModal();
        showNotification('⚠️ Bulk edit requires at least 2 selected devices', 'warning');
        return;
    }

    const statusValue = document.getElementById('bulk-edit-status')?.value || 'keep';
    const locationValue = document.getElementById('bulk-edit-location')?.value.trim() || '';
    const nowIso = new Date().toISOString();

    let changedCount = 0;

    devices = devices.map((device) => {
        if (!selectedDeviceIds.has(device.id)) return device;

        const next = { ...device };
        let changed = false;

        if (statusValue !== 'keep' && next.status !== statusValue) {
            next.status = statusValue;
            if (statusValue === 'offline') {
                next.last_downtime = nowIso;
            }
            changed = true;
        }

        if (locationValue && next.location !== locationValue) {
            next.location = locationValue;
            changed = true;
        }

        if (changed) {
            next.updated_at = nowIso;
            changedCount += 1;
        }

        return next;
    });

    hideBulkEditModal();
    loadDevices();
    updateStats();
    showNotification(`✅ Updated ${changedCount} device(s)`);
}

function deleteSelectedDevices() {
    closeSelectionMenu();
    if (selectedDeviceIds.size === 0) {
        showNotification('⚠️ Select at least one device first', 'warning');
        return;
    }

    const count = selectedDeviceIds.size;
    const confirmed = window.confirm(`Delete ${count} selected device(s)?`);
    if (!confirmed) return;

    devices = devices.filter((device) => !selectedDeviceIds.has(device.id));
    selectedDeviceIds.clear();

    loadDevices();
    updateStats();
    showNotification(`✅ Deleted ${count} device(s)`);
}

function bulkAction(action) {
    const onlineDevices = devices.filter(d => d.status === 'online');

    if (onlineDevices.length === 0) {
        showNotification('⚠️ No online devices to execute command', 'warning');
        return;
    }

    const panel = document.getElementById('execution-panel');
    panel.classList.remove('hidden');

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

            setTimeout(() => {
                panel.classList.add('hidden');
                document.getElementById('progress-bar').style.width = '0%';
            }, 5000);
        }
    }, 500);
}

function refreshDevices() {
    showNotification('🔄 Refreshing devices...', 'info');

    setTimeout(() => {
        const nowIso = new Date().toISOString();
        devices.forEach(device => {
            if (device.status === 'online') {
                device.last_seen = nowIso;
                device.updated_at = nowIso;
            }
        });

        loadDevices();
        updateStats();
        showNotification('✅ Devices refreshed!');
    }, 1000);
}

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
    return date.toLocaleString(undefined, {
        year: 'numeric',
        month: 'short',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit'
    });
}

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

setInterval(updateClock, 1000);
updateClock();
