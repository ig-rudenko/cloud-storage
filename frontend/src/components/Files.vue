<template>
  <h1>Облачное файловое хранилище</h1>
  <div class="container" @click="$refs.menu.close()">
    <div class="breadcrumbs">
      <span @click="selectFolder(null)">Главная</span>
      <span v-for="(folder, index) in folderPath">
        / <span @click="selectFolder(folderPath.slice(0, index+1))">{{ folder }}</span>
      </span>
    </div>

    <Progress :progress="uploadProgress" :loading="uploading">
      <p>Загрузка...</p>
    </Progress>
    <Progress :progress="downloadProgress" :loading="downloading">
      <p>Скачиваем файл...</p>
    </Progress>

    <div class="tiles">
      <div
          class="tile"
          v-for="(item, index) in items"
          :key="item.id"
          :class="{ folder: item.isDir, file: !item.isDir }"
          @click="openItem(item)"
          @contextmenu.prevent.stop="showMenu($event, item)"
      >

<!--        FOLDER-->
        <div v-if="item.isDir">
          <svg style="vertical-align: middle" xmlns="http://www.w3.org/2000/svg" width="80" height="80" fill="currentColor" viewBox="0 0 16 16">
            <path d="M9.828 3h3.982a2 2 0 0 1 1.992 2.181l-.637 7A2 2 0 0 1 13.174 14H2.825a2 2 0 0 1-1.991-1.819l-.637-7a1.99 1.99 0 0 1 .342-1.31L.5 3a2 2 0 0 1 2-2h3.672a2 2 0 0 1 1.414.586l.828.828A2 2 0 0 0 9.828 3zm-8.322.12C1.72 3.042 1.95 3 2.19 3h5.396l-.707-.707A1 1 0 0 0 6.172 2H2.5a1 1 0 0 0-1 .981l.006.139z"/>
          </svg>
        </div>

<!--        FILE-->
        <div v-else>
          <svg xmlns="http://www.w3.org/2000/svg" width="80" height="80" fill="currentColor" viewBox="0 0 16 16">
            <path d="M14 4.5V14a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V2a2 2 0 0 1 2-2h5.5L14 4.5zm-3 0A1.5 1.5 0 0 1 9.5 3V1H4a1 1 0 0 0-1 1v12a1 1 0 0 0 1 1h8a1 1 0 0 0 1-1V4.5h-2z"/>
          </svg>
        </div>

        <div style="overflow-wrap: break-word;">
          {{ item.name }}
        </div>
      </div>

<!--      + FOLDER-->
      <div @click.stop="createFolder" class="tile new-folder">
        <svg xmlns="http://www.w3.org/2000/svg" width="80" height="80" fill="currentColor" class="bi bi-folder-plus" viewBox="0 0 16 16">
          <path d="m.5 3 .04.87a1.99 1.99 0 0 0-.342 1.311l.637 7A2 2 0 0 0 2.826 14H9v-1H2.826a1 1 0 0 1-.995-.91l-.637-7A1 1 0 0 1 2.19 4h11.62a1 1 0 0 1 .996 1.09L14.54 8h1.005l.256-2.819A2 2 0 0 0 13.81 3H9.828a2 2 0 0 1-1.414-.586l-.828-.828A2 2 0 0 0 6.172 1H2.5a2 2 0 0 0-2 2zm5.672-1a1 1 0 0 1 .707.293L7.586 3H2.19c-.24 0-.47.042-.683.12L1.5 2.98a1 1 0 0 1 1-.98h3.672z"/>
          <path d="M13.5 10a.5.5 0 0 1 .5.5V12h1.5a.5.5 0 1 1 0 1H14v1.5a.5.5 0 1 1-1 0V13h-1.5a.5.5 0 0 1 0-1H13v-1.5a.5.5 0 0 1 .5-.5z"/>
        </svg>
      </div>
    </div>

    <context-menu :display="showContextMenu" ref="menu">
        <p @click="downloadFile">Скачать</p>
        <p @click="renameFile">Переименовать</p>
        <p @click="deleteItem">Удалить</p>

        <p style="border-top: 1px solid; margin: 5px 0;"></p>

        <p v-if="!selectedFile.isDir">{{formatBytes(selectedFile.size)}}</p>
        <p>{{selectedFile.modTime}}</p>
    </context-menu>
  </div>
</template>

<script>

import ContextMenu from "@/components/ContextMenu.vue";
import Progress from "@/components/Progress.vue";
import api from "@/services/api";

export default {

  components: {
    ContextMenu,
    Progress
  },

  data() {
    return {
      items: [],
      folderPath: [],
      showContextMenu: false,
      selectedFile: null,
      uploadProgress: null,
      uploading: false,
      downloadProgress: null,
      downloading: false
    };
  },

  async mounted() {
    if (!this.currentUser) {
      this.$router.push('/login');
    }

    await this.updateFolder();
    this.addDragAndDropListeners();
    setTimeout(this.updateFolder, 5000)
  },

  computed: {
    currentUser() {
      return this.$store.state.auth.user;
    }
  },

  methods: {

    formatBytes(bytes) {
      let marker = 1024; // Измените на 1000 при необходимости
      let decimal = 1;
      let kiloBytes = marker; // Один килобайт - это 1024 байта
      let megaBytes = marker * marker; // Один мегабайт - это 1024 килобайта
      let gigaBytes = marker * marker * marker; // Один гигабайт - это 1024 мегабайта
      let teraBytes = marker * marker * marker * marker; // Один терабайт - это 1024 гигабайта

      if (bytes < kiloBytes) return bytes + " Б";
      else if (bytes < megaBytes) return (bytes / kiloBytes).toFixed(decimal) + " КБ";
      else if (bytes < gigaBytes) return (bytes / megaBytes).toFixed(decimal) + " МБ";
      else if (bytes < teraBytes) return (bytes / gigaBytes).toFixed(decimal) + " ГБ";
    },

    showMenu(event, file) {
      this.selectedFile = file;
      this.$refs.menu.open(event);
    },

    async selectFolder(folder_path) {
      if (folder_path === null) {
        folder_path = []
      }
      this.folderPath = folder_path
      console.log("this.folderPath", this.folderPath)
      let path = "/items/" + this.folderPath.join("/")
      console.log(path)
      let resp = await api.get(path)
      this.items = resp.data
    },

    async updateFolder() {
      await this.selectFolder(this.folderPath)
    },

    openItem(item) {
      if (item.isDir) {
        this.selectFolder([...this.folderPath, item.name]);
      }
    },

    async downloadFile () {
      // Если это файл и скачивание другого файла не происходит
      if (!this.selectedFile.isDir && !this.downloading) {
        // определяем url для скачивания файла
        const url = "/item/" + [...this.folderPath, this.selectedFile.name].join("/");

        this.downloading = true
        // делаем запрос к серверу с типом ответа 'blob'
        api.get(
            url, {
              responseType: "blob",
              onDownloadProgress: (progressEvent) => { //
                // функция для обновления прогресса загрузки
                this.downloadProgress = Math.round( (progressEvent.loaded * 100) / progressEvent.total ); // вычисляем и сохраняем прогресс в процентах
              }
            }
        )
          .then(response => {
            // получаем blob из ответа
            const blob = response.data;
            // создаем ссылку на файл из blob
            const fileURL = URL.createObjectURL(blob);
            // создаем элемент <a> с атрибутом href равным ссылке на файл
            const link = document.createElement("a");
            link.href = fileURL;
            // задаем имя файла для загрузки (можно получить из ответа сервера)
            link.download = this.selectedFile.name;
            // добавляем элемент <a> в документ
            document.body.appendChild(link);
            // имитируем клик по элементу <a>
            link.click();
            // удаляем элемент <a> из документа
            document.body.removeChild(link);
            this.downloading = false
          })
          .catch(error => {
            // обрабатываем ошибку запроса или загрузки файла
            console.error(error.message);
            this.downloading = false
          });
      }
    },

    async deleteItem() {
      await api.delete("/item/" + this.folderPath.join("/") + "/" + this.selectedFile.name)
      await this.updateFolder()
    },

    async renameFile() {
      let newName = prompt("Введите название папки", this.selectedFile.name)
      const url = "/item/rename/" + [...this.folderPath, this.selectedFile.name].join("/")
      await api.post(url, {"newName": newName})
      await this.updateFolder()
    },

    async createFolder(){
      let name = prompt("Введите название папки");
        if(name){
          // Здесь можно сделать запрос к серверу для создания новой папки в выбранной папке
          // Для простоты я просто вывожу сообщение в консоль и добавляю папку в список

          console.log("Создание папки:", name);
          let url = "/items/create-folder/" + [...this.folderPath, name].join("/")
          let res = await api.post(url)
          await this.updateFolder()
        }
      },

    async uploadFiles (e) {
      // Если уже загружаем файл
      if (this.uploading) return;

      e.preventDefault();
      let files = e.dataTransfer.files;
      if (!files) return;

      const formData = new FormData();
      for (let file of files) {
        formData.append("files", file);
      }
      let url = "/items/upload/" + this.folderPath.join('/')

      this.uploading = true
      let resp = await api.post(
          url,
          formData,
          {
            headers: {
              "Content-Type": "multipart/form-data",
            },
            onUploadProgress: (progressEvent) => { //
              // функция для обновления прогресса загрузки
              this.uploadProgress = Math.round( (progressEvent.loaded * 100) / progressEvent.total ); // вычисляем и сохраняем прогресс в процентах
            }
          }
      )
      this.uploading = false
      console.log(resp)
      await this.updateFolder()
    },

    addDragAndDropListeners() {
      let container = document.querySelector("body");
      container.addEventListener("dragover", e => e.preventDefault());
      container.addEventListener("drop", (e) => this.uploadFiles(e));
    }
  }

};
</script>

<style scoped>

.container {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  height: 100%;
}

.breadcrumbs{
  margin: 20px;
  cursor: pointer;
  padding-bottom: 5px;
}

.tiles{
  display: flex;
  flex-wrap: wrap;
}

.tile{
  width: 140px;
  margin: 10px;
  border-radius: 10px;
  align-items:center;
  justify-content:center;
  text-align:center;
  cursor:pointer;
  user-select: none;
}

.folder{
  color: #e3e3e3;
}

.file{
  color: #cec0c0;
}

.new-folder{
  color: #656565;
}
</style>
