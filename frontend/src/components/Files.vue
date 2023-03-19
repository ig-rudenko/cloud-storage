<template>
  <h1>Облачное файловое хранилище</h1>
  <div class="container" @click="$refs.menu.close()">
    <div class="breadcrumbs" style="cursor:pointer;">
      <span @click="selectFolder(null)">Главная</span>
      <span v-for="(folder, index) in folderPath" :key="folder.id">
        / <span @click="selectFolder(folder)">{{ folder.name }}</span>
      </span>
    </div>

    <div class="tiles">
      <div
          class="tile"
          v-for="(item, index) in items"
          :key="item.id"
          :class="{ folder: item.type === 'folder', file: item.type === 'file' }"
          @click="openItem(item)"
          @contextmenu.prevent.stop="showMenu($event, item)"
      >

<!--        FOLDER-->
        <div v-if="item.type === 'folder'">
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

        <div>
          {{ shortFileName(item.name) }}
        </div>
      </div>

<!--      + FOLDER-->
      <div @click.stop="createFolder" class="tile">
        <svg xmlns="http://www.w3.org/2000/svg" width="80" height="80" fill="#afafaf" class="bi bi-folder-plus" viewBox="0 0 16 16">
          <path d="m.5 3 .04.87a1.99 1.99 0 0 0-.342 1.311l.637 7A2 2 0 0 0 2.826 14H9v-1H2.826a1 1 0 0 1-.995-.91l-.637-7A1 1 0 0 1 2.19 4h11.62a1 1 0 0 1 .996 1.09L14.54 8h1.005l.256-2.819A2 2 0 0 0 13.81 3H9.828a2 2 0 0 1-1.414-.586l-.828-.828A2 2 0 0 0 6.172 1H2.5a2 2 0 0 0-2 2zm5.672-1a1 1 0 0 1 .707.293L7.586 3H2.19c-.24 0-.47.042-.683.12L1.5 2.98a1 1 0 0 1 1-.98h3.672z"/>
          <path d="M13.5 10a.5.5 0 0 1 .5.5V12h1.5a.5.5 0 1 1 0 1H14v1.5a.5.5 0 1 1-1 0V13h-1.5a.5.5 0 0 1 0-1H13v-1.5a.5.5 0 0 1 .5-.5z"/>
        </svg>
      </div>
    </div>

    <context-menu :display="showContextMenu" ref="menu">
        <p>Переименовать</p>
        <p>Удалить</p>
    </context-menu>
    <input type="file" ref="fileInput" hidden @change="uploadFile" />
  </div>
</template>

<script>

import ContextMenu from "@/components/ContextMenu.vue";

export default {
  components: {ContextMenu},
  data() {
    return {
      folders: [
        { id: 1, name: "Документы", parentId: null, type: "folder"},
        { id: 2, name: "Фотографии", parentId: null, type: "folder" },
        { id: 3, name: "Видео", parentId: null, type: "folder" },
        { id: 4, name: "Сканы", parentId: 1, type: "folder" },
      ],
      files: [
        { id: 5,
          name: "Резюме.docx",
          size: "25 KB",
          date: "18.03.2023",
          folderId: 1,
          type: "file"
        },
        { id: 6,
          name: "Отчет.pdf",
          size: "150 KB",
          date: "17.03.2023",
          folderId: 1,
          type: "file"
        },
        { id: 7,
          name: "Презентация.pptx",
          size: "300 KB",
          date: "16.03.2023",
          folderId: 1,
          type: "file"
        },
        { id: 8,
          name: "Паспорт.jpg",
          size: "50 KB",
          date: "15.03.2023",
          folderId: 4,
          type: "file"
        }
        ],
      selectedFolder: null,
      items: [],
      folderPath: [],
      showContextMenu: false,
      selectedFile: null,

    };
  },
  mounted() {
    if (!this.currentUser) {  // #
      this.$router.push('/login');
    }

    this.selectFolder(null);
    this.addDragAndDropListeners();
  },

  computed: {  // #
    currentUser() {
      return this.$store.state.auth.user;
    }
  },

  methods: {

    shortFileName(file_name) {
      if (file_name.length > 12) {
        return file_name.slice(0, 5) + "..." + file_name.slice(-6)
      }
      return file_name
    },

    showMenu(event, file) {
      this.selectedFile = file;
      this.$refs.menu.open(event);
    },

    selectFolder(folder) {
      this.selectedFolder = folder;
      console.log("folder", folder)
      // Здесь можно сделать запрос к серверу для получения папок и файлов в выбранной папке
      // Для простоты я использую фиктивные данные
      this.items = [
        ...this.folders.filter((f) => f.parentId === (folder ? folder.id : null)),
        ...this.files.filter((f) => f.folderId === (folder ? folder.id : null)),
      ];
      console.log(this.items)
      this.folderPath = [];
      let currentFolder = folder;
      while (currentFolder) {
        this.folderPath.unshift(currentFolder);
        currentFolder = this.folders.find((f) => f.id === currentFolder.parentId);
      }
    },
    openItem(item) {
      console.log(item)
      if (item.type === "folder") {
        this.selectFolder(item);
      } else if (item.type === "file") {
        // Здесь можно сделать запрос к серверу для скачивания или открытия файла
        // Для простоты я просто вывожу сообщение в консоль
        console.log("Открытие файла:", item.name);
      }
    },
    createFolder(){
      let name = prompt("Введите название папки");
        if(name){
          // Здесь можно сделать запрос к серверу для создания новой папки в выбранной папке
          // Для простоты я просто вывожу сообщение в консоль и добавляю папку в список

          console.log("Создание папки:", name);
          let newFolder = {
            id: this.folders.length + 1,
            name: name,
            parentId: this.selectedFolder ? this.selectedFolder.id : null,
            type: "folder"
          };
          this.folders.push(newFolder);
          this.items.push(newFolder);
        }
      },
    uploadFile() {
      let file = this.$refs.fileInput.files[0];
      if (file) {
        // Здесь можно сделать запрос к серверу для загрузки файла в выбранную папку
        // Для простоты я просто вывожу сообщение в консоль и добавляю файл в список

        console.log("Загрузка файла:", file.name);
        let newFile = {
          id: this.files.length + 1,
          name: file.name,
          size: file.size,
          date: new Date().toLocaleDateString(),
          folderId: this.selectedFolder ? this.selectedFolder.id : null,
        };
        this.files.push(newFile);
        this.items.push(newFile);
        this.$refs.fileInput.value = "";
      }
      },
      addDragAndDropListeners() {
        let container = document.querySelector(".container");
        container.addEventListener("dragover", e => e.preventDefault());
        container.addEventListener("drop", e => {
          e.preventDefault();
          let file = e.dataTransfer.files[0];
          if(file){
            // Здесь можно сделать запрос к серверу для загрузки файла в выбранную папку
            // Для простоты я просто вывожу сообщение в консоль и добавляю файл в список

            console.log("Загрузка файла:", file.name);
            let newFile = {
              id: this.files.length + 1,
              name: file.name,
              size: file.size,
              date: new Date().toLocaleDateString(),
              folderId: this.selectedFolder ? this.selectedFolder.id : null,
            };
            this.files.push(newFile);
            this.items.push(newFile);}
        }
        );
    }
  }
};
</script>

<style scoped>

.container {
  /*padding: 100px;*/
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  height: 100%;
}

.breadcrumbs{
  margin: 10px;
}

.tiles{
  display: flex;
  flex-wrap: wrap;
  /*justify-content: center;*/
  /*padding: 50px;*/
}

.tile{
  /*max-width: 100%;*/
  /*width: 100px;*/
  padding: 20px;
  /*height: 100px;*/
  margin: 10px;
  border-radius: 10px;
  /*box-shadow: 0 2px 5px rgba(248, 240, 240, 0.3);*/
  /*display: flex;*/
  align-items:center;
  justify-content:center;
  text-align:center;
  cursor:pointer;
}

.folder{
  /*background-color: #2d2d2d;*/
  color: #ffffff;
}

.file{
  /*background-color: #919191;*/
  color: #cec0c0;
}

.new-folder{
  background-color:#4caf50;
  color:#fff;
}
</style>