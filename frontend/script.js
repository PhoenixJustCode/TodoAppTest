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
});

const form = document.getElementById("taskForm");
const taskList = document.getElementById("taskList");
const completedList = document.getElementById("completedList");

form.addEventListener("submit", (e) => {
  e.preventDefault();

  const text = document.getElementById("taskInput").value.trim();
  const date = document.getElementById("taskDate").value; // может быть пустым
  const priority = document.getElementById("prioritySelect").value;

  if (!text) return; // текст обязателен

  let dateHTML = "";
  if (date) {
    const [datePart, timePart] = date.split("T");
    const formattedTime = timePart ? timePart.slice(0, 5) : ""; // HH:MM
    dateHTML = `<small>${datePart} - ${formattedTime}</small>`;
  }

  const div = document.createElement("div");
  div.className = "task";
  const priorityText = priority.charAt(0).toUpperCase() + priority.slice(1);

  div.innerHTML = `
    <strong>${priorityText}</strong>
    <span>${text}</span>
    ${dateHTML}
    <div class="actions">
      <button class="btn complete-btn">✅</button>
      <button class="btn delete-btn">🗑️</button>
    </div>
  `;

  taskList.appendChild(div);
  form.reset();
});


// делегирование событий для кнопок
document.addEventListener("click", (e) => {
  const task = e.target.closest(".task");
  if (!task) return;

  //завершение задачи
  if (e.target.classList.contains("complete-btn")) {
    if (task.parentElement.id === "taskList") {
      task.classList.add("completed");
      task.querySelector(".complete-btn").textContent = "↩️";
      completedList.appendChild(task);
    } else {
      task.classList.remove("completed");
      task.querySelector(".complete-btn").textContent = "✅";
      taskList.appendChild(task);
    }
  }

  // Удаление задачи (с confirm)
  if (e.target.classList.contains("delete-btn")) {
    if (confirm("Delete this task?")) {
      task.remove();
    }
  }
});



// подключение к GO
// Получаем все задачи
async function fetchTasks() {
  const res = await fetch("/tasks"); // путь к твоему API
  const tasks = await res.json();
  console.log(tasks);
  renderTasks(tasks);
}

// Добавляем задачу
async function addTask(taskObj) {
  const res = await fetch("/tasks/add", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(taskObj)
  });
  const data = await res.json();
  console.log(data);
}
