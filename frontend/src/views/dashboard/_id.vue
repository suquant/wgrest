<template>
  <div class="peer__container">
    <div class="peer__list">
      <peer-item
        v-for="peer in peersList"
        :key="peer.public_key"
        :item="peer"
        @delete="deletePeer"
      ></peer-item>
    </div>
<!--    <el-button-->
<!--      type="primary"-->
<!--      @click="loadMorePeers"-->
<!--    >-->
<!--      Load More-->
<!--    </el-button>-->
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

  private queries = {
    name: this.$route.params.id,
    per_page: 100,
    page: 0,
    q: '',
    sort: ''
  }

  public async getPeerList(): Promise<void> {
    // eslint-disable-next-line camelcase
    const { name, per_page, page, q, sort } = this.queries
    const { data } = await deviceApi.listDevicePeers(name, per_page, page, q, sort)

    const flatten = (a: any) => [].concat.apply([], a)
    const noDuplicateProps = (a: any, b: any) => a.url_safe_public_key === b.url_safe_public_key

    const combineAndDeDup = (...arrs: any[]) => {
      return flatten(arrs).reduce((acc, item) => {
        const uniqueItem = acc.findIndex(i => noDuplicateProps(i, item)) === -1

        if (uniqueItem) return acc.concat([item])

        return acc
      }, [])
    }

    this.peersList = combineAndDeDup(this.peersList, data)
  }

  private deletePeer(key: string): void {
    const index = this.peersList.findIndex(item => item.public_key === key)

    this.peersList.splice(index, 1)
  }

  private loadMorePeers(): void {
    this.queries.page++
    this.getPeerList()
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
.peer__container {
  display: flex;
  flex-direction: column;
  align-items: center;
}
.peer__list {
  padding-top: 30px;
  display: flex;
  flex-direction: row;
  align-items: center;
  flex-wrap: wrap;
}
</style>
