# 🌐 Miral Ecosystem

**Miral** — это модульная экосистема сервисов, объединённых общей системой аккаунтов (**Miral Account**). Аналогично Google, Miral предоставляет пользователям единый доступ к почте, облачному хранилищу, уведомлениям и другим сервисам.

---

## 📦 Сервисы

### 🔐 Miral Account
> Единая система авторизации и аутентификации пользователей

- Регистрация и вход
- Авторизация через Django-сессии
- Используется всеми приложениями экосистемы

### 📩 Miral Mail
> Встроенный почтовый сервис

- Отправка системных и пользовательских писем
- Шаблоны email-сообщений
- Интеграция с аккаунтами пользователей

### ☁️ Miral Cloud
> Облачное хранилище данных

- Загрузка и хранение файлов
- Привязка файлов к пользователям
- API-доступ к файлам

### 📬 Pusher
> Сервис уведомлений

- Отправка уведомлений по событиям
- REST API
- Возможная интеграция с FCM/WebSocket

### 🗨️ Chat Backend
> Бэкенд для мессенджера

- Django + Django REST Framework
- Поддержка хранения сообщений
- Потенциальная реализация WebSocket

---

## 🧰 Используемые технологии

- Python 3.x
- Django
- Django REST Framework (DRF)
- SQLite (на этапе разработки)
- Django Sessions
- REST API
- Модульная архитектура (Django apps)

---

## 🗂️ Структура проекта

chat_backend/
├── emails/ # Почтовый сервис
├── miral_cloud/ # Облачное хранилище
├── miralmail/ # Основная конфигурация проекта
├── pfps/ # Аватары пользователей
├── pusher/ # Сервис уведомлений
├── users/ # Miral Account (аккаунты)
├── manage.py # Django-менеджер

---

## 🚀 Быстрый старт

# Клонируем проект
git clone https://github.com/yourname/miral-ecosystem.git
cd miral-ecosystem

# Установка зависимостей
pip install -r requirements.txt

# Применение миграций
python manage.py migrate

# Запуск сервера
python manage.py runserver

Telegram (автор): @MiralAbbastada
