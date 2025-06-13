# 🌐 Miral Ecosystem

**Miral Ecosystem** — это модульная экосистема сервисов, объединённых общей системой аккаунтов (**Miral Account**). Аналогично Google, Miral предоставляет пользователям единый доступ к почте, облачному хранилищу, уведомлениям и другим сервисам.

---

## 📦 Сервисы

### 🔐 Miral Account
> Единая система авторизации и аутентификации пользователей с уникальной системой идентификации Miral ID

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
- Golang
- Модульная архитектура (Django apps)

---

## 🗂️ Структура проекта

chat_backend/ <br/>
├── emails/ # Почтовый сервис <br/>
├── miral_cloud/ # Облачное хранилище (<b>Модуль написан на языке программирования Go</b>) <br/>
├── miralmail/ # Основная конфигурация проекта <br/>
├── pfps/ # Аватары пользователей <br/>
├── pusher/ # Сервис уведомлений <br/>
├── users/ # Miral Account (аккаунты) <br/>
├── manage.py # Django-менеджер <br/>

---

## 🚀 Быстрый старт

# Клонируем проект

``` git clone https://github.com/AbdullaShamuradow28/miral-ecosystem.git ``` <br/> <br/>
``` cd miral-ecosystem ```

# Установка зависимостей

``` pip install -r requirements.txt ```

# Применение миграций

``` python manage.py migrate ```

# Запуск сервера

``` python manage.py runserver ```
