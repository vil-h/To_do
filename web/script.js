const API_BASE = '';

const state = {
    filter: 'all',
    tasks: [],
    loading: false
};

const elements = {
    form: null,
    input: null,
    list: null,
    submit: null,
    status: null,
    filters: null,
    deadline: null
};

function $(selector, root = document) {
    return root.querySelector(selector);
}

function createStatus() {
    const node = document.createElement('div');
    node.className = 'js-status';
    node.setAttribute('aria-live', 'polite');
    return node;
}

function createDeadlineInput() {
    const input = document.createElement('input');
    input.type = 'datetime-local';
    input.className = 'main__input js-deadline';
    input.id = 'deadline';
    input.name = 'deadline';
    return input;
}

function createFilters() {
    const wrap = document.createElement('div');
    wrap.className = 'js-filters';
    wrap.innerHTML = `
    <button type="button" data-filter="all" class="js-filter-btn is-active">Все</button>
    <button type="button" data-filter="active" class="js-filter-btn">Активные</button>
    <button type="button" data-filter="completed" class="js-filter-btn">Готово</button>
    <button type="button" data-filter="archived" class="js-filter-btn">Архив</button>
  `;
    return wrap;
}

function injectStyles() {
    const style = document.createElement('style');
    style.textContent = `
    .main__form { max-width: 900px; margin: 0 auto; }
    .main__nav { display: flex; flex-direction: column; gap: 14px; }
    .main__input { width: 100%; }
    .js-status { min-height: 24px; font-family: 'Comfortaa', sans-serif; }
    .js-status.is-error { color: #c62828; }
    .js-status.is-success { color: #1b7d34; }
    .js-status.is-loading { color: #555; }
    .js-filters { display: flex; flex-wrap: wrap; gap: 10px; margin-top: 8px; }
    .js-filter-btn {
      border: none; border-radius: 10px; padding: 12px 18px; cursor: pointer;
      font-family: 'Comfortaa', sans-serif; font-weight: 700; background: #efe6d6;
    }
    .js-filter-btn.is-active { background: #00a627; color: #fff; }
    .todoList { list-style: none; display: grid; gap: 14px; padding: 0; margin-top: 10px; }
    .todo-item {
      display: flex; justify-content: space-between; align-items: flex-start; gap: 16px;
      padding: 18px; border-radius: 14px; background: #f8f3ea; box-shadow: 0 6px 18px rgba(0,0,0,.08);
    }
    .todo-item.is-completed { opacity: .78; }
    .todo-item.is-completed .todo-item__name { text-decoration: line-through; }
    .todo-item__content { display: grid; gap: 8px; }
    .todo-item__name { font-family: 'Comfortaa', sans-serif; font-weight: 700; }
    .todo-item__meta { color: #666; font-size: 14px; font-family: 'Comfortaa', sans-serif; }
    .todo-item__actions { display: flex; flex-wrap: wrap; gap: 8px; }
    .todo-action {
      border: none; border-radius: 10px; padding: 10px 14px; cursor: pointer;
      font-family: 'Comfortaa', sans-serif; font-weight: 700;
    }
    .todo-action--toggle { background: #00a627; color: #fff; }
    .todo-action--delete { background: #d45252; color: #fff; }
    .todo-empty {
      padding: 18px; border-radius: 14px; background: #f8f3ea; font-family: 'Comfortaa', sans-serif;
    }
    @media (max-width: 700px) {
      .todo-item { flex-direction: column; }
      .todo-item__actions { width: 100%; }
      .todo-action { width: 100%; }
    }
  `;
    document.head.appendChild(style);
}

function setStatus(message, type = '') {
    if (!elements.status) return;
    elements.status.textContent = message;
    elements.status.className = `js-status ${type}`.trim();
}

function formatDate(value) {
    if (!value) return 'Без дедлайна';
    const date = new Date(value);
    if (Number.isNaN(date.getTime())) return value;
    return new Intl.DateTimeFormat('ru-RU', {
        dateStyle: 'medium',
        timeStyle: 'short'
    }).format(date);
}

function toBackendDate(value) {
    if (!value) return '';
    return `${value.replace('T', ' ')}:00`;
}

async function request(path, options = {}) {
    const response = await fetch(`${API_BASE}${path}`, {
        headers: {
            'Content-Type': 'application/json',
            ...(options.headers || {})
        },
        ...options
    });

    const text = await response.text();
    let data = null;

    try {
        data = text ? JSON.parse(text) : null;
    } catch (_) {
        data = text;
    }

    if (!response.ok) {
        throw new Error(typeof data === 'string' ? data : `Ошибка ${response.status}`);
    }

    return data;
}

function renderTasks() {
    if (!elements.list) return;

    if (!state.tasks.length) {
        elements.list.innerHTML = '<li class="todo-empty">Задач пока нет.</li>';
        return;
    }

    elements.list.innerHTML = state.tasks.map(task => `
    <li class="todo-item ${task.is_active ? 'is-completed' : ''}" data-id="${task.id}">
      <div class="todo-item__content">
        <p class="todo-item__name">${escapeHtml(task.name)}</p>
        <p class="todo-item__meta">Создано: ${formatDate(task.created_time)}</p>
        <p class="todo-item__meta">Дедлайн: ${formatDate(task.deadline_time)}</p>
        <p class="todo-item__meta">Статус: ${task.is_active ? 'Выполнено' : 'В процессе'}</p>
      </div>
      <div class="todo-item__actions">
        ${task.is_archived ? '' : `<button type="button" class="todo-action todo-action--toggle" data-action="toggle" data-id="${task.id}">${task.is_active ? 'Вернуть' : 'Готово'}</button>`}
        ${task.is_archived ? '' : `<button type="button" class="todo-action todo-action--delete" data-action="delete" data-id="${task.id}">В архив</button>`}
      </div>
    </li>
  `).join('');
}

function escapeHtml(value) {
    return String(value)
        .replaceAll('&', '&amp;')
        .replaceAll('<', '&lt;')
        .replaceAll('>', '&gt;')
        .replaceAll('"', '&quot;')
        .replaceAll("'", '&#39;');
}

async function loadTasks() {
    state.loading = true;
    setStatus('Загружаем задачи...', 'is-loading');
    try {
        const data = await request(`/tasks?filter=${encodeURIComponent(state.filter)}`, {
            method: 'GET'
        });
        state.tasks = Array.isArray(data) ? data : [];
        renderTasks();
        setStatus(`Загружено задач: ${state.tasks.length}`, 'is-success');
    } catch (error) {
        setStatus(error.message || 'Не удалось загрузить задачи.', 'is-error');
    } finally {
        state.loading = false;
    }
}

async function createTask(event) {
    event.preventDefault();
    const name = elements.input?.value.trim() || '';
    const deadline = elements.deadline?.value || '';

    if (!name) {
        setStatus('Введите название задачи.', 'is-error');
        elements.input?.focus();
        return;
    }

    elements.submit.disabled = true;
    setStatus('Создаём задачу...', 'is-loading');

    try {
        await request('/create', {
            method: 'POST',
            body: JSON.stringify({
                name,
                deadline_time: toBackendDate(deadline)
            })
        });
        elements.form.reset();
        setStatus('Задача создана.', 'is-success');
        await loadTasks();
    } catch (error) {
        setStatus(error.message || 'Не удалось создать задачу.', 'is-error');
    } finally {
        elements.submit.disabled = false;
    }
}

async function toggleTask(id) {
    setStatus('Меняем статус...', 'is-loading');
    try {
        await request('/tasks/change', {
            method: 'PATCH',
            body: JSON.stringify({ id: Number(id) })
        });
        await loadTasks();
    } catch (error) {
        setStatus(error.message || 'Не удалось изменить статус.', 'is-error');
    }
}

async function archiveTask(id) {
    setStatus('Архивируем задачу...', 'is-loading');
    try {
        await request('/tasks/delete', {
            method: 'DELETE',
            body: JSON.stringify({ id: Number(id) })
        });
        await loadTasks();
    } catch (error) {
        setStatus(error.message || 'Не удалось архивировать задачу.', 'is-error');
    }
}

function bindListActions() {
    elements.list?.addEventListener('click', async event => {
        const button = event.target.closest('[data-action]');
        if (!button) return;

        const { action, id } = button.dataset;
        if (!id) return;

        if (action === 'toggle') {
            await toggleTask(id);
        }

        if (action === 'delete') {
            await archiveTask(id);
        }
    });
}

function bindFilters() {
    elements.filters?.addEventListener('click', async event => {
        const button = event.target.closest('[data-filter]');
        if (!button) return;

        state.filter = button.dataset.filter;
        elements.filters.querySelectorAll('[data-filter]').forEach(item => item.classList.remove('is-active'));
        button.classList.add('is-active');
        await loadTasks();
    });
}

function setup() {
    injectStyles();

    elements.form = $('.main__form');
    elements.input = $('.main__input');
    elements.list = $('.todoList');
    elements.submit = $('.main__btn');

    if (!elements.form || !elements.input || !elements.list || !elements.submit) {
        return;
    }

    const deadline = createDeadlineInput();
    deadline.placeholder = 'Дедлайн';
    elements.input.insertAdjacentElement('afterend', deadline);
    elements.deadline = deadline;

    const filters = createFilters();
    elements.submit.insertAdjacentElement('afterend', filters);
    elements.filters = filters;

    const status = createStatus();
    elements.list.insertAdjacentElement('beforebegin', status);
    elements.status = status;

    elements.form.addEventListener('submit', createTask);
    bindListActions();
    bindFilters();
    loadTasks();
}

document.addEventListener('DOMContentLoaded', setup);