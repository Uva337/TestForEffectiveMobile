# REST‑сервис управления подписками  
Проект: **Test for Effective Mobile**  
Реализован на Go.  

## 🧩 Описание  
Данный сервис предоставляет REST‑API для управления онлайн‑подписками пользователей:  
- Полный CRUD (создание, чтение, обновление, удаление) подписок.  
- Эндпоинт агрегации: подсчёт суммарной стоимости подписок с гибкой фильтрацией.  
- Валидация входящих данных — обеспечение целостности.  
- Применена «чистая архитектура»: слои Handler → Service → Repository.  
- Контейнеризация: приложение и база данных запускаются через Docker / docker‑compose.  
- Юнит‑тесты бизнес‑логики с использованием моков.  
- Поддержка graceful shutdown, проброс контекста (context), структурированное логирование.  
- Документация API с помощью Swagger (OpenAPI).  

## ⚙️ Технологический стек  
| Направление        | Технология                                  |
|--------------------|---------------------------------------------|
| Язык               | Go                                          |
| База данных        | PostgreSQL                                  |
| Веб‑фреймворк      | chi                                         |
| Миграции           | golang‑migrate/migrate                       |
| Конфигурация       | spf13/viper                                 |
| Валидация          | go‑playground/validator                      |
| Логирование        | slog (стандартная библиотека)               |
| Тестирование       | testify/mock, testify/assert                 |
| Документация API   | swaggo/swag                                 |

## 🚀 Быстрый старт  

### Предварительные условия  
- git  
- Docker  
- docker‑compose  
- Go (локально, если будете запускать тесты без контейнеров)  

### Запуск проекта  
```bash
git clone https://github.com/Uva337/TestForEffectiveMobile.git
cd TestForEffectiveMobile
cp .env.example .env
docker-compose up -d db
migrate -path migrations -database 'postgres://user:password@localhost:5432/subscriptions_db?sslmode=disable' up
docker-compose up --build
```
Приложение станет доступно по адресу: [http://localhost:8080](http://localhost:8080)

Swagger‑документация:  
[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

## 🧪 Запуск тестов  
```bash
go test ./...
```

## 🗂️ Структура проекта  
```
cmd/           — точка входа приложения  
docs/          — документация (например, спецификация OpenAPI)  
internal/      — внутренние пакеты приложения  
migrations/    — SQL‑/скрипты миграций базы данных  
pkg/logger     — модуль логирования  
Dockerfile     — образ приложения  
docker-compose.yml — конфигурация контейнеров  
go.mod, go.sum — модули Go  
```


