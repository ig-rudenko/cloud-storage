# Используем официальный образ node:alpine как родительский
FROM node:alpine as builder

# Устанавливаем рабочую директорию в контейнере
WORKDIR /app

# Копируем файлы package.json и package-lock.json в рабочую директорию
COPY package*.json ./

# Устанавливаем зависимости для приложения
RUN npm install

# Копируем остальные файлы проекта в рабочую директорию
COPY . .

# Собираем приложение для продакшена
RUN npm run build


# Используем официальный образ nginx:stable-alpine как базовый для запуска приложения
FROM nginx:stable-alpine

# Копируем собранные файлы приложения из предыдущего этапа в папку /usr/share/nginx/html
COPY --from=builder /app/dist /usr/share/nginx/html

# Открываем порт 80 для доступа к приложению
EXPOSE 80

# Запускаем nginx с конфигурацией по умолчанию
CMD ["nginx", "-g", "daemon off;"]