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

### How to use:
### Windows 10 and Docker Desktop:
PowerShell
```powershell
docker build --file .\cloud-storage\Dockerfile .\cloud-storage\
```
where path .\cloud-storage is directory with all elements of system, root.

or

```powershell
docker build -t cloud-storage .
```
where cloud-storage is name of container. Execute by root directory

If you wanna use MySQL on Host you must use `host.docker.internal` as `localhost`, then in `backend/config/config.go` replace `DSN` with:
```golang
// DSN Get database connection settings from environment variables
var DSN = fmt.Sprintf("%s:%s@tcp(host.docker.internal:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local",
	getEnv("DB_USER", "root"),
	getEnv("DB_PASS", "root"),
	getEnv("DB_NAME", "black_hole"))
```
