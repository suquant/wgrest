<template>
  <router-link :to="{path: `/devices/${item.name}`}">
   <el-card class="box-card" shadow="hover">
    <div slot="header" class="card__header">
      <span>name: {{ item.name }}</span>
      <span>port: {{ item.listen_port }}</span>
    </div>
   <div class="text item">
     <svg-icon name="user" /> :   {{ item.peers_count }}
   </div>
    <div class="text item">
      firewall mark: {{ item.firewall_mark }}
    </div>
    <div class="text item">
      Receive bytes: {{ item.total_receive_bytes }}
    </div>
    <div class="text item">
      Transmit bytes: {{ item.total_transmit_bytes }}
    </div>
  </el-card>
  </router-link>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator'
import { Device } from 'wgrest/dist/models'

@Component({
  name: 'interfaceItem'
})
export default class interfaceItem extends Vue {
  @Prop({ default: {} }) item!: Device
}
</script>

<style scoped>
.text {
  font-size: 14px;
}

.item {
  margin-bottom: 18px;
}

.box-card {
  width: 480px;

  cursor: pointer;
  margin: 0 50px 50px 0;
}

.card__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
