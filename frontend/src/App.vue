<template>
  <div id="app">
    <router-view />
    <Modal />
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator'
import Modal from '@/components/Modal/index.vue'

@Component({
  name: 'App',

  components: {
    Modal
  }
})
export default class extends Vue {}
</script>
