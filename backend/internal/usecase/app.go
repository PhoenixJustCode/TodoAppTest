package usecase

type App struct {
    // можно хранить ссылки на сервисы/репозитории
}

func NewApp() *App {
    return &App{}
}

// пример метода для фронта
func (a *App) Hello() string {
    return "Hello from Go!"
}
