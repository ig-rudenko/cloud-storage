<template>
  <div class="login">
    <h1>Пожалуйста, войдите</h1>
    <form @submit.prevent="handleLogin">
      <div class="form-group">
        <label for="username">Username</label>
        <input type="text" id="username" v-model="user.username" required />
      </div>
      <div class="form-group">
        <label for="password">Password</label>
        <input type="password" id="password" v-model="user.password" required />
      </div>
      <button type="submit">Войти</button>
      <br><br>
      <a style="cursor: pointer" type="button" @click="$router.push('/register')">Нет аккаунта</a>
    </form>
  </div>
</template>

<script>

export default {
  data() {
    return {
      user: {  // #
        username: "",
        password: "",
      }
    };
  },
  computed: {  // #
    loggedIn() {
      return this.$store.state.auth.status.loggedIn;
    },
  },
  created() {  // #
    if (this.loggedIn) {
      this.$router.push("/");
    }
  },
  methods: {

    handleLogin() {  // #
      this.$store.dispatch("auth/login", this.user).then(
          () => {
            this.$router.push("/");
          },
          (error) => {
            let message =
                (error.response &&
                    error.response.data &&
                    error.response.data.message) ||
                error.message ||
                error.toString();
            console.log(message)
          }
      );
    },

    submitForm() {
      // send the username and password to the server
      fetch("http://127.0.0.1:8080/token", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          username: this.username,
          password: this.password,
        }),
      })
          .then((response) => response.json())
          .then((data) => {
            // handle the response
            console.log(data);
            // save the JWT tokens in the storage
            localStorage.setItem("access_token", data.access_token);
            localStorage.setItem("refresh_token", data.refresh_token);
          })
          .catch((error) => {
            // handle the error
            console.error(error);
          });
    },
  },
};
</script>

<style scoped>

.login {
  width: 400px;
  margin: auto;
  font-family: Arial, sans-serif;
}

h1 {
  text-align: center;
}

.form-group {
  margin-bottom: 20px;
}

label {
  display: block;
}

input {
  width: 100%;
  padding: 10px;
  border-radius: 5px;
}

button {
  width: 100%;
  padding: 10px;
  border-radius: 5px;
}

button:hover{
  background-color:#00a0ff;
  color:white;
  cursor:pointer;
}
</style>