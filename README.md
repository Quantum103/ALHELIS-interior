[![Go CI](https://github.com/Quantum103/ALHELIS-interior/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/Quantum103/ALHELIS-interior/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/github/go-mod/v/Quantum103/ALHELIS-interior?filename=auth-service%2Fgo.mod)](https://github.com/Quantum103/ALHELIS-interior)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![GitHub repo size](https://img.shields.io/github/repo-size/Quantum103/ALHELIS-interior)](https://github.com/Quantum103/ALHELIS-interior)
# ALHELIS Interior

Платформа для заказа дизайн-проектов интерьеров. Проект задуман как микросервисное веб-приложение, где клиенты могут смотреть портфолио, регистрироваться, оставлять заявки и связываться со специалистами. 

Разрабатываю проект с нуля, постепенно выстраивая и усложняя архитектуру. Сейчас реализован базовый каркас, количество микросервисов будет расти по мере расширения функционала.

## Стек технологий

* **Язык бэкенда:** Go (Golang)
* **Взаимодействие:** gRPC, REST API
* **База данных:** PostgreSQL 15 (драйвер `pgx/v5`)
* **Инфраструктура:** Docker, Docker Compose
* **API Gateway / Веб-сервер:** Nginx
* **Фронтенд:** HTML5, CSS3, JavaScript (Vanilla)

## Архитектура

Проект разбит на изолированные сервисы для удобства масштабирования и разработки:

* **`gateway` (Nginx)** — единая точка входа. Отдает статические файлы фронтенда, проксирует HTTP-запросы от веб-форм на бэкенд и маршрутизирует gRPC-трафик.
* **`auth-service`** — микросервис аутентификации. Обрабатывает регистрацию, логин и выдачу токенов. Общается как по HTTP (принимает данные с фронта), так и по gRPC (внутри сети).
* **`user-service`** — микросервис управления профилями пользователей и их данными (gRPC).
* **`postgres_db`** — общая база данных для надежного хранения учетных записей и связанных сущностей.

## Как запустить локально

Для запуска вам потребуются установленные **Docker** и **Docker Compose**.

1. Склонируйте репозиторий и перейдите в папку проекта:
   ```bash
   git clone https://github.com/Quantum103/ALHELIS-interior.git
   cd ALHELIS-interior
   

2. Создайте файл `.env` в корне проекта:
   ```env
   DB_HOST=postgres_db
   DB_PORT=5432
   DB_USER=alhelis
   DB_PASSWORD=password
   DB_NAME=alhelis
   JWT_SECRET=super-secret-key-for-my-app

3. Запустите все сервисы в фоновом режиме
   ```bash
   docker compose up --build -d