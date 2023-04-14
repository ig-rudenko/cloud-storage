<template>
  <div class="logo">
    <img src="../../assets/logo.png" alt="logo">
  </div>
</template>

<style>
.logo {
  display: flex;
  justify-content: center;
  height: 350px;
  margin-bottom: 20px;
  animation: spin 60s infinite linear;
}


@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}
</style>
