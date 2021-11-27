<template>
<div class="peer__list">
  <peer-item
    v-for="peer in peersList"
    :key="peer.public_key"
    :item="peer"
    @delete="deletePeer"
  ></peer-item>
</div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator'
import { deviceApi } from '@/api/interface'
import { Peer } from 'wgrest/dist/models'

import peerItem from '@/views/dashboard/peerItem.vue'
import { emitter } from '@/utils/emmiter'

@Component({
  name: 'Peer',
  components: {
    peerItem
  }
})
export default class extends Vue {
  private peersList: Peer[] = []

  private updatePeersTimer: any

  public async getPeerList(): Promise<void> {
    const { data } = await deviceApi.listDevicePeers(this.$route.params.id)

    this.peersList = data
  }

  private deletePeer(key: string): void {
    const index = this.peersList.findIndex(item => item.public_key === key)

    this.peersList.splice(index, 1)
  }

  created() {
    this.getPeerList()
    emitter.on('updatePeer', this.getPeerList)

    this.updatePeersTimer = setInterval(this.getPeerList, 5000)
  }

  beforeDestroy() {
    clearInterval(this.updatePeersTimer)
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
