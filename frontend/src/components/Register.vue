<template>
  <Logo/>
  <div class="register">
    <h1>Регистрация</h1>
    <form @submit.prevent="handleRegister">
      <div class="form-group">
        <label for="username">Username</label>
        <input type="text" id="username" v-model="user.username" required />
      </div>
      <div class="form-group">
        <label for="password">Password</label>
        <input type="password" id="password" v-model="user.password" required />
      </div>
      <button type="submit">Создать аккаунт</button>
      <br><br>
      <a style="cursor: pointer" type="button" @click="$router.push('/login')">Уже есть аккаунт</a>
    </form>
  </div>
</template>

<script>

import Logo from "@/components/icons/Logo.vue";

export default {
  components: {Logo},
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
  mounted() {  // #
    if (this.loggedIn) {
      this.$router.push("/");
    }
  },

  methods: {

    handleRegister() {
      console.log(this.user)
      this.$store.dispatch("auth/register", this.user).then(
          (data) => {
            console.log(data)
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

  },
};
</script>

<style scoped>

.register {
  width: 330px;
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
  font-size: 16px;
  padding: 12px;
  border-radius: 5px;
  border: none;
}

button:hover{
  background-color: #333333;
  color: white;
  border: none;
  cursor: pointer;
}
</style>