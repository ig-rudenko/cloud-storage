<h1 align="center">Black Hole Cloud Storage</h1>

<p align="center">
<img alt="img_1.png" src="img/img_1.png"/>
</p>

---

### Используемые языки и фреймворки:

<div>
<img src="https://www.vectorlogo.zone/logos/golang/golang-official.svg" alt="GOLANG" width="80" height="40"/>

<img src="https://raw.githubusercontent.com/devicons/devicon/master/icons/javascript/javascript-original.svg" alt="JS" width="40" height="40"/>

<img src="https://raw.githubusercontent.com/gin-gonic/logo/master/color.png" alt="GIN" width="40" height="60"/>

<img src="https://www.vectorlogo.zone/logos/vuejs/vuejs-icon.svg" alt="VUE.JS" width="40" height="40"/>

<img src="https://www.vectorlogo.zone/logos/axios/axios-icon.svg" alt="AXIOS" width="40" height="40"/>
</div>

### Используемые технологии:

<div>
<img src="https://www.vectorlogo.zone/logos/docker/docker-tile.svg" alt="vue.js" width="40" height="40"/>

<img src="https://www.vectorlogo.zone/logos/nginx/nginx-icon.svg" alt="vue.js" width="40" height="40"/>

<img src="https://www.vectorlogo.zone/logos/mysql/mysql-official.svg" alt="mysql" width="40" height="40"/>
</div>

---

Облачное файловое хранилище с поддержкой drag&drop

* Frontend - Vue.js
* Backend - Go
* Database - MySQL

![img.png](img/img.png)


## Конфигурация

Для настроки подключения к базе данных MySQL используются следующие переменные окружения

```
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=black_hole
DB_USER=user
DB_PASS=password
```

Изменение секретного ключа для генерации JWT

```
SECRET_KEY=my-secret-key
```

В данный момент все файлы пользователей хранятся в локальной директории, 
указанной в переменной окружения

```
STORAGE_DIR=storage
```

Иерархия файлового хранилища:

* storage/
  * <user-id>/
    * <files...>
  * <user-id>/
    * <files...>

## Запуск

Для запуска через docker используем команду

```shell
docker compose up -d
```

Так как пользовательские данные хранятся локально, то при работе backend приложения в
контейнере необходим **bind mount** директории в которой будут храниться файлы в контейнер.

Если переменная окружения `STORAGE_DIR=storage`, то в контейнере это папка
`/app/storage` - это и будет корневая директория хранилища.


### Создание исполняемого файла backend приложения:

```shell
go mod download
go build -v web/backend/cmd/app
```
или
```shell
go mod download
make
```

## Документация

Swagger документация доступна по URL `/api/swagger/index.html`

![img.png](img/img_2.png)