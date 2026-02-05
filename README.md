# Comment Tree — древовидные комментарии с навигацией и поиском
Сервис для работы с древовидными комментариями, поддерживающий неограниченную вложенность и дополнительные функции - поиск, сортировка и постраничный вывод.

[Старт](https://github.com/andreyxaxa/URL-Shortener?tab=readme-ov-file#%D0%B7%D0%B0%D0%BF%D1%83%D1%81%D0%BA)

## Обзор

- UI - http://localhost:8080/v1/ui
- Документация API - Swagger - http://localhost:8080/swagger
- Конфиг - [config/config.go](https://github.com/andreyxaxa/Comment-Tree/blob/main/config/config.go). Читается из `.env` файла.
- Логгер - [pkg/logger/logger.go](https://github.com/andreyxaxa/Comment-Tree/blob/main/pkg/logger/logger.go). Интерфейс позволяет подменить логгер.
- Удобная и гибкая конфигурация HTTP сервера - [pkg/httpserver/options.go](https://github.com/andreyxaxa/Comment-Tree/blob/main/pkg/httpserver/options.go).
  Позволяет конфигурировать сервер в конструкторе таким образом:
  ```go
  httpServer := httpserver.New(httpserver.Port(cfg.HTTP.Port))
  ```
- В слое контроллеров применяется версионирование - [internal/controller/restapi/v1](https://github.com/andreyxaxa/Comment-Tree/tree/main/internal/controller/restapi/v1).
  Для версии v2 нужно будет просто добавить папку `restapi/v2` с таким же содержимым, в файле [internal/controller/restapi/router.go](https://github.com/andreyxaxa/Comment-Tree/blob/main/internal/controller/restapi/router.go) добавить строку:
```go
{
		v1.NewCommentRoutes(apiV1Group, c, l)
}

{
		v2.NewCommentRoutes(apiV1Group, c, l)
}
```
- Graceful shutdown - [internal/app/app.go](https://github.com/andreyxaxa/Comment-Tree/blob/main/internal/app/app.go).

## Запуск

1. Клонируйте репозиторий
2. В корне создайте `.env` файл, скопируйте туда содержимое [env.example](https://github.com/andreyxaxa/Comment-Tree/blob/main/.env.example) - `cp .env.example .env`
3. Выполните, дождитесь запуска сервиса
   ```
   make compose-up
   ```
4. Перейдите на http://localhost:8080/v1/ui и пользуйтесь сервисом.
<img width="1429" height="1017" alt="image" src="https://github.com/user-attachments/assets/db5f6a4d-d995-4eb3-8787-c42da564e456" />
- Перейдите на http://localhost:8080/swagger и ознакомьтесь с API, если хотите взаимодействовать с сервисом вручную или из стороннего сервиса.

## Прочие `make` команды
Зависимости:
```
make deps
```
docker compose down -v:
```
make compose-down
```
