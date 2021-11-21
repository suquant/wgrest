<template>
  <div class="dashboard-container">
    <interface-item
      v-for="device in devices"
      :key="device.public_key"
      :item="device"
    />
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator'
import { DevicesModule } from '@/store/modules/devices'
import InterfaceItem from '@/views/dashboard/InterfaceItem.vue'
import { Device } from 'wgrest/dist/models'

@Component({
  name: 'Dashboard',
  components: { InterfaceItem },

  mounted() {
    DevicesModule.getDevicesList()
  }
})
export default class extends Vue {
  get devices(): Device[] {
    return DevicesModule.devices
  }
}
</script>

<style lang="scss" scoped>
.dashboard {
  &-container {
    margin: 30px;
    display: flex;
    align-items: center;
    flex-direction: row;
    flex-wrap: wrap;
  }

  &-text {
    font-size: 30px;
    line-height: 46px;
  }
}
</style>
