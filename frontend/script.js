document.addEventListener("DOMContentLoaded", () => {
  const toggleBtn = document.getElementById("themeToggle");
  const body = document.body;

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

  const form = document.getElementById("taskForm");
  const taskList = document.getElementById("taskList");
  const completedList = document.getElementById("completedList");

  // Загрузить задачи при старте
  fetchTasks();

  // --- submit формы -> POST /tasks/add (JSON)
  form.addEventListener("submit", async (e) => {
    e.preventDefault();
    const text = document.getElementById("taskInput").value.trim();
    const date = document.getElementById("taskDate").value
      ? document.getElementById("taskDate").value + ":00Z"
      : null;
    const priority = document.getElementById("prioritySelect").value || "low";

    if (!text) return;

    const taskObj = {
      task: text,
      priority,
      status: false, // по умолчанию не выполнена
      due_date: date, // либо null
    };

    try {
      await fetch("/tasks/add", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(taskObj),
      });
      form.reset();
      await fetchTasks();
    } catch (err) {
      console.error("Add failed", err);
    }
  });

  // Делегирование кликов (complete / delete)
  document.addEventListener("click", async (e) => {
    const taskDiv = e.target.closest(".task");
    if (!taskDiv) return;
    const id = taskDiv.dataset.id;
    if (!id) return;

    // complete/restore
    if (e.target.classList.contains("complete-btn")) {
      // определяем текущий статус DOM (и отправляем противоположный)
      const currentlyCompleted = taskDiv.classList.contains("completed");
      const newStatus = !currentlyCompleted; // true -> mark as completed

      try {
        await fetch("/tasks/update", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ id: Number(id), status: newStatus }),
        });
        await fetchTasks();
      } catch (err) {
        console.error("Update status failed", err);
      }
    }

    // delete
    if (e.target.classList.contains("delete-btn")) {
      if (!confirm("Delete this task?")) return;
      try {
        await fetch("/tasks/delete", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ id: Number(id) }),
        });
        await fetchTasks();
      } catch (err) {
        console.error("Delete failed", err);
      }
    }
  });

  // --- API functions ---
  async function fetchTasks() {
    try {
      const res = await fetch("/tasks");
      if (!res.ok) throw new Error("Fetch tasks failed");
      const tasks = await res.json();
      renderTasks(tasks || []);
    } catch (err) {
      console.error("Failed to load tasks:", err);
      taskList.innerHTML = "";
      completedList.innerHTML = "";
    }
  }

  // --- render tasks into two lists
  function renderTasks(tasks) {
    taskList.innerHTML = "";
    completedList.innerHTML = "";

    if (!tasks || tasks.length === 0) {
      // проверка на пустой массив
      console.log("Нет задач для отображения");
      return;
    }

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
      if (
        t.due_date &&
        t.due_date !== "0001-01-01T00:00:00Z" &&
        t.due_date !== null
      ) {
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

  // helper: convert different server status reps to bool
  function normalizeStatus(s) {
    if (typeof s === "boolean") return s;
    if (typeof s === "number") return s !== 0;
    if (typeof s === "string") {
      const v = s.toLowerCase();
      return v === "true" || v === "1" || v === "completed" || v === "done";
    }
    return false;
  }

  // escape minimal
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

  // filters
  // --- APPLY FILTERS ---
  function applyFilters() {
    const priority = document.getElementById("priorityFilter").value;
    const date = document.getElementById("dateFilter").value;
    const sort = document.getElementById("sortSelect").value;
    const order = document.getElementById("sortOrder").value;
    const status = document.getElementById("statusFilter").value;

    fetch(
      `/tasks/filter?priority=${priority}&date=${date}&sort=${sort}&order=${order}&status=${status}`
    )
      .then((res) => {
        if (!res.ok) throw new Error("Failed to fetch filtered tasks");
        return res.json();
      })
      .then((tasks) => renderTasks(tasks))
      .catch((err) => console.error("Filter error:", err));
  }

  // --- ПРИВЯЗКА К SELECT ---
  document
    .getElementById("priorityFilter")
    .addEventListener("change", applyFilters);
  document
    .getElementById("dateFilter")
    .addEventListener("change", applyFilters);
  document
    .getElementById("sortSelect")
    .addEventListener("change", applyFilters);
  document.getElementById("sortOrder").addEventListener("change", applyFilters);
  document
    .getElementById("statusFilter")
    .addEventListener("change", applyFilters); 
});


