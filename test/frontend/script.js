document.addEventListener("DOMContentLoaded", () => {
  const body = document.body;
  const toggleBtn = document.getElementById("themeToggle");
  const form = document.getElementById("taskForm");
  const taskList = document.getElementById("taskList");
  const completedList = document.getElementById("completedList");

  // script.js
  if (!window.taskBackend) {
    window.taskBackend = {
      GetAllTasks: async () => [],
      AddTask: async () => {},
      UpdateStatus: async () => {},
      DeleteTask: async () => {},
    };
  } 

  // --- Переключение темы ---
  body.classList.add("light");
  toggleBtn.addEventListener("click", () => {
    if (body.classList.contains("light")) {
      body.classList.replace("light", "dark");
      toggleBtn.textContent = "☀️";
    } else {
      body.classList.replace("dark", "light");
      toggleBtn.textContent = "🌙";
    }
  });

  // --- Ждем бинды Wails ---
  function waitForWailsBind(callback, retries = 50) {
    if (window.taskBackend) {
      callback();
    } else if (retries > 0) {
      setTimeout(() => waitForWailsBind(callback, retries - 1), 100);
    }
  }

  waitForWailsBind(async () => {
    await fetchTasks(); // сразу загружаем задачи после биндов
  });

  async function fetchTasks() {
    if (!window.taskBackend) return;

    try {
      const tasks = await window.taskBackend.GetAllTasks(); // получаем все задачи
      const filteredTasks = applyFilters(tasks); // фильтруем
      renderTasks(filteredTasks); // отображаем
    } catch (err) {
      taskList.innerHTML = "";
      completedList.innerHTML = "";
    }
  }

  // Навешиваем фильтры на селекты
  [
    "statusFilter",
    "priorityFilter",
    "dateFilter",
    "sortSelect",
    "sortOrder",
  ].forEach((id) => {
    document.getElementById(id).addEventListener("change", fetchTasks);
  });
  
  
  function renderTasks(tasks) {
    taskList.innerHTML = "";
    completedList.innerHTML = "";

    if (!tasks || tasks.length === 0) return;

    tasks.forEach((t) => {
      const statusBool = normalizeStatus(t.status);
      const div = document.createElement("div");
      div.className = "task";
      div.dataset.id = t.id;
      if (statusBool) div.classList.add("completed");

      const priorityText =
        (t.priority || "low").charAt(0).toUpperCase() +
        (t.priority || "low").slice(1);

      let dateHTML = "";
      if (t.due_date && t.due_date !== "0001-01-01T00:00:00Z") {
        const d = new Date(t.due_date);
        if (!isNaN(d.getTime())) {
          dateHTML = `<small>${d.toLocaleDateString()} ${d.toLocaleTimeString(
            [],
            { hour: "2-digit", minute: "2-digit" }
          )}</small>`;
        }
      }

      div.innerHTML = `
        <strong>${priorityText}</strong>
        <span class="task-text">${escapeHtml(t.task || "")}</span>
        ${dateHTML}
        <div class="actions">
          <button class="btn complete-btn">${statusBool ? "↩️" : "✅"}</button>
          <button class="btn delete-btn">🗑️</button>
        </div>
      `;

      if (statusBool) completedList.appendChild(div);
      else taskList.appendChild(div);
    });
  }

  function normalizeStatus(s) {
    if (typeof s === "boolean") return s;
    if (typeof s === "number") return s !== 0;
    if (typeof s === "string") {
      const v = s.toLowerCase();
      return v === "true" || v === "1" || v === "completed" || v === "done";
    }
    return false;
  }

  function escapeHtml(text) {
    return String(text).replace(
      /[&<>"']/g,
      (m) =>
        ({
          "&": "&amp;",
          "<": "&lt;",
          ">": "&gt;",
          '"': "&quot;",
          "'": "&#39;",
        }[m])
    );
  }


  // --- Добавление задачи ---
  form.addEventListener("submit", async (e) => {
    e.preventDefault();
    if (!window.taskBackend) return;

    const text = document.getElementById("taskInput").value.trim();
    if (!text) return;

    const dateVal = document.getElementById("taskDate").value;
    const priority = document.getElementById("prioritySelect").value || "low";

    const taskObj = { task: text, priority, status: false };

    if (dateVal) {
      taskObj.due_date = dateVal + ":00Z";
    }

    try {
      await window.taskBackend.AddTask(taskObj);
      form.reset();
      await fetchTasks();
    } catch (err) {
      console.error("AddTask failed:", err);
    }
  });


  // --- Делегирование кликов (complete / delete) ---
  document.addEventListener("click", async (e) => {
    if (!window.taskBackend) return;
    const taskDiv = e.target.closest(".task");
    if (!taskDiv) return;
    const id = Number(taskDiv.dataset.id);
    if (!id) return;

    if (e.target.classList.contains("complete-btn")) {
      const currentlyCompleted = taskDiv.classList.contains("completed");
      try {
        await window.taskBackend.UpdateStatus(id, !currentlyCompleted);
        await fetchTasks();
      } catch (err) {
        console.error("UpdateStatus failed:", err);
      }
    }

    if (e.target.classList.contains("delete-btn")) {
      if (!confirm("Delete this task?")) return;
      try {
        await window.taskBackend.DeleteTask(id);
        await fetchTasks();
      } catch (err) {
        console.error("DeleteTask failed:", err);
      }
    }
  });


  function applyFilters(tasks) {
    const statusVal = document.getElementById("statusFilter").value;
    const priorityVal = document.getElementById("priorityFilter").value;
    const dateVal = document.getElementById("dateFilter").value;
    const sortVal = document.getElementById("sortSelect").value;
    const orderVal = document.getElementById("sortOrder").value;

    let filtered = tasks.filter((task) => {
      let statusOk = true;
      if (statusVal === "active") statusOk = !normalizeStatus(task.status);
      if (statusVal === "completed") statusOk = normalizeStatus(task.status);

      let priorityOk = true;
      if (priorityVal !== "all") priorityOk = task.priority === priorityVal;

      let dateOk = true;
      const now = new Date();
      const taskDate = task.due_date ? new Date(task.due_date) : null;
      if (taskDate) {
        if (dateVal === "today") {
          dateOk =
            taskDate.getDate() === now.getDate() &&
            taskDate.getMonth() === now.getMonth() &&
            taskDate.getFullYear() === now.getFullYear();
        } else if (dateVal === "week") {
          const weekAhead = new Date();
          weekAhead.setDate(now.getDate() + 7);
          dateOk = taskDate >= now && taskDate <= weekAhead;
        } else if (dateVal === "overdue") {
          dateOk = taskDate < now && !normalizeStatus(task.status);
        }
      } else if (dateVal !== "all") {
        dateOk = false;
      }

      return statusOk && priorityOk && dateOk;
    });

    // --- Сортировка ---
    filtered.sort((a, b) => {
      let aVal, bVal;
      if (sortVal === "date") {
        aVal = a.due_date
          ? new Date(a.due_date)
          : new Date("0001-01-01T00:00:00Z");
        bVal = b.due_date
          ? new Date(b.due_date)
          : new Date("0001-01-01T00:00:00Z");
      } else if (sortVal === "priority") {
        const order = { low: 1, medium: 2, high: 3 };
        aVal = order[a.priority] || 0;
        bVal = order[b.priority] || 0;
      }
      if (orderVal === "asc") return aVal > bVal ? 1 : -1;
      return aVal < bVal ? 1 : -1;
    });

    return filtered;
  }
  


});

