<template>
  <div class="context-menu" v-show="show" :style="style" ref="context" tabindex="0" @blur="close">
    <slot></slot>
  </div>
</template>

<script>
export default {
  name: "ContextMenu",
  props: {
    display: Boolean, // prop detect if we should show context menu
  },
  data() {
    return {
      left: 0, // left position
      top: 0, // top position
      show: false, // affect display of context menu
    };
  },
  computed: {
    // get position of context menu
    style() {
      return {
        top: this.top + "px",
        left: this.left + "px",
      };
    },
  },
  methods: {
    // closes context menu
    close() {
      this.show = false;
      this.left = 0;
      this.top = 0;
    },
    open(evt) {
      // updates position of context menu
      this.left = evt.pageX || evt.clientX;
      this.top = evt.pageY || evt.clientY;
      this.show = true;
    },
    hideOnLeftClick(evt) { // new method to hide the menu on left click outside
      if (evt.button === 0 && !this.$el.contains(evt.target)) { // check if it is a left click and not on the menu element
        this.close(); // close the menu
      }
    }
  },
  mounted() {
    document.addEventListener("mousedown", this.hideOnLeftClick); // add an event listener for mousedown events on document
  },
  beforeDestroy() {
    document.removeEventListener("mousedown", this.hideOnLeftClick); // remove the event listener when component is destroyed
  },
};
</script>

<style scoped>
.context-menu {
  padding: 10px 20px;
  border-radius: 5px;
  background-color: #222222;
  position: fixed;
  z-index: 999;
  outline: none;
  cursor: pointer;
}
</style>
