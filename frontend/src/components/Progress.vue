<template>
<div v-if="loading">
  <slot>
    <p>Загрузка…</p>
  </slot>
  <progress :value="progress" max="100"></progress>
</div>
</template>

<script>
export default {
  name: "Progress",
  props: {
    progress: Number,
    loading: Boolean
  }
}
</script>
