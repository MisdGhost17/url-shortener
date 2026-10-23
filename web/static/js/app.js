'use strict';

const API_BASE = '';

const api = {
  async shorten(url) {
    const res = await fetch(`${API_BASE}/api/shorten`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ url }),
    });
    const data = await res.json();
    if (!res.ok) throw new Error(data.error || 'Ошибка сервера');
    return data;
  },

  async delete(shortCode) {
    const res = await fetch(`${API_BASE}/api/delete/${shortCode}`, {
      method: 'DELETE',
    });
    const data = await res.json();
    if (!res.ok) throw new Error(data.error || 'Ошибка удаления');
    return data;
  },
};

const ui = {
  urlInput:    document.getElementById('urlInput'),
  shortenBtn:  document.getElementById('shortenBtn'),
  resultCard:  document.getElementById('resultCard'),
  resultLink:  document.getElementById('resultLink'),
  resultOrig:  document.getElementById('resultOriginal'),
  copyBtn:     document.getElementById('copyBtn'),
  copyLabel:   document.getElementById('copyLabel'),
  closeResult: document.getElementById('closeResult'),
  errorToast:  document.getElementById('errorToast'),
  errorMsg:    document.getElementById('errorMsg'),
  linkCount:   document.getElementById('linkCount'),
  linksTable:  document.getElementById('linksTable'),
  linksBody:   document.getElementById('linksBody'),
  emptyState:  document.getElementById('emptyState'),
  refreshBtn:  document.getElementById('refreshBtn'),
  copiedTip:   document.getElementById('copiedTooltip'),

  showResult(data) {
    ui.resultLink.href      = data.short_url;
    ui.resultLink.textContent = data.short_url;
    ui.resultOrig.textContent = data.original;
    ui.resultCard.classList.remove('hidden');
    ui.hideError();
  },

  hideResult() {
    ui.resultCard.classList.add('hidden');
  },

  showError(msg) {
    ui.errorMsg.textContent = msg;
    ui.errorToast.classList.remove('hidden');
    ui.hideResult();
    setTimeout(() => ui.hideError(), 4000);
  },

  hideError() {
    ui.errorToast.classList.add('hidden');
  },

  setLoading(isLoading) {
    ui.shortenBtn.disabled = isLoading;
    ui.shortenBtn.classList.toggle('loading', isLoading);
  },

  updateCount(n) {
    ui.linkCount.textContent = `${n} ССЫЛОК`;
  },

  renderTable(links) {
    ui.updateCount(links.length);

    if (links.length === 0) {
      ui.linksTable.classList.add('hidden');
      ui.emptyState.style.display = 'block';
      return;
    }

    ui.emptyState.style.display = 'none';
    ui.linksTable.classList.remove('hidden');

    const sorted = [...links].sort(
      (a, b) => new Date(b.created_at) - new Date(a.created_at)
    );

    ui.linksBody.innerHTML = sorted.map((link, i) => `
      <tr class="row-enter" style="animation-delay: ${i * 30}ms">
        <td class="code-cell">
          <a class="code-link"
             href="/${link.short_code}"
             target="_blank"
             title="Перейти по ссылке">
            ${escapeHtml(link.short_code)}
          </a>
        </td>
        <td class="url-cell">
          <span class="url-text" title="${escapeHtml(link.original_url)}">
            ${escapeHtml(link.original_url)}
          </span>
        </td>
        <td class="date-cell">${formatDate(link.created_at)}</td>
        <td class="clicks-cell">
          <span class="clicks-badge">
            <span class="n">${link.clicks}</span>
            <span class="label">CLICKS</span>
          </span>
        </td>
        <td class="action-cell">
          <button class="btn-delete"
                  data-code="${escapeHtml(link.short_code)}"
                  title="Удалить">
            DELETE
          </button>
        </td>
      </tr>
    `).join('');
  },

  showCopied() {
    ui.copiedTip.classList.remove('hidden');
    clearTimeout(ui._copiedTimer);
    ui._copiedTimer = setTimeout(() => {
      ui.copiedTip.classList.add('hidden');
    }, 1800);
  },
};

function escapeHtml(str) {
  return String(str)
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;');
}

function formatDate(iso) {
  const d = new Date(iso);
  return d.toLocaleString('en-GB', {
    day: '2-digit',
    month: 'short',
    hour: '2-digit',
    minute: '2-digit',
  });
}

async function loadLinks() {
  const saved = JSON.parse(localStorage.getItem('myCodes') || '[]');

  if (saved.length === 0) {
    ui.renderTable([]);
    return;
  }

  try {
    const results = await Promise.all(
      saved.map(code =>
        fetch(`/api/stats/${code}`).then(r => r.json())
      )
    );
    const valid = results.filter(r => !r.error);
    const validCodes = valid.map(r => r.short_code);
    localStorage.setItem('myCodes', JSON.stringify(validCodes));

    ui.renderTable(valid);
  } catch (err) {
    console.error('Ошибка загрузки ссылок:', err);
  }
}

async function handleShorten() {
  const raw = ui.urlInput.value.trim();
  if (!raw) {
    ui.showError('Введите URL');
    ui.urlInput.focus();
    return;
  }

  const url = /^https?:\/\//i.test(raw) ? raw : `https://${raw}`;

  ui.setLoading(true);
  ui.hideError();

  try {
    const data = await api.shorten(url);

    const saved = JSON.parse(localStorage.getItem('myCodes') || '[]');
    saved.push(data.short_code);
    localStorage.setItem('myCodes', JSON.stringify(saved));

    ui.showResult(data);
    ui.urlInput.value = '';
    await loadLinks();
  } catch (err) {
    ui.showError(err.message);
  } finally {
    ui.setLoading(false);
  }
}

async function handleCopy() {
  const url = ui.resultLink.href;
  try {
    await navigator.clipboard.writeText(url);
    ui.showCopied();
    ui.copyLabel.textContent = '✓ СКОПИРОВАНО';
    setTimeout(() => { ui.copyLabel.textContent = 'COPY'; }, 2000);
  } catch {
    const sel = window.getSelection();
    const range = document.createRange();
    range.selectNode(ui.resultLink);
    sel.removeAllRanges();
    sel.addRange(range);
  }
}

async function handleTableClick(e) {
  const btn = e.target.closest('.btn-delete');
  if (!btn) return;

  const code = btn.dataset.code;
  if (!code) return;

  if (!confirm(`Удалить /${code}?`)) return;

  const row = btn.closest('tr');
  if (row) {
    row.style.opacity = '0.3';
    row.style.pointerEvents = 'none';
  }

  try {
    await api.delete(code);

    const saved = JSON.parse(localStorage.getItem('myCodes') || '[]');
    const updated = saved.filter(c => c !== code);
    localStorage.setItem('myCodes', JSON.stringify(updated));

    await loadLinks();
  } catch (err) {
    if (row) {
      row.style.opacity = '';
      row.style.pointerEvents = '';
    }
    ui.showError(err.message);
  }
}
function init() {
  ui.shortenBtn.addEventListener('click', handleShorten);
  ui.urlInput.addEventListener('keydown', (e) => {
    if (e.key === 'Enter') handleShorten();
  });
  ui.closeResult.addEventListener('click', ui.hideResult);
  ui.copyBtn.addEventListener('click', handleCopy);
  ui.refreshBtn.addEventListener('click', loadLinks);
  ui.linksBody.addEventListener('click', handleTableClick);
  loadLinks();
}

if (document.readyState === 'loading') {
  document.addEventListener('DOMContentLoaded', init);
} else {
  init();
}