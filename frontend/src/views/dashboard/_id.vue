<template>
<div class="peer__list">
  <peer-item
    v-for="peer in peersList"
    :key="peer.public_key"
    :item="peer"
  ></peer-item>
</div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator'
import { deviceApi } from '@/api/interface'
import { Peer } from 'wgrest/dist/models'

import peerItem from '@/views/dashboard/peerItem.vue'

@Component({
  name: 'Peer',
  components: {
    peerItem
  }
})
export default class extends Vue {
  private peersList: Peer[] = []

  public async getPeerList(): Promise<void> {
    const { data } = await deviceApi.listDevicePeers(this.$route.params.id)

    this.peersList = data
  }

  created() {
    this.getPeerList()
  }
}
</script>

<style scoped>
.peer__list {
  padding-top: 30px;
  display: flex;
  flex-direction: row;
  align-items: center;
  flex-wrap: wrap;

}
</style>
