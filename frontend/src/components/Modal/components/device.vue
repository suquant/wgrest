<template>
  <div>
    <p class="device_form modal_title">Add new Device</p>
    <el-input
      class="device_form"
      placeholder="Name"
      v-model="newDeviceForm.name"
    ></el-input>
    <el-input-number
      class="device_form"
      placeholder="Listen Port"
      v-model="newDeviceForm.listen_port"
    ></el-input-number>
    <el-input
      class="device_form"
      placeholder="Private Key"
      v-model="newDeviceForm.private_key"
    ></el-input>
    <el-input-number
      class="device_form"
      placeholder="Firewall Mark"
      v-model="newDeviceForm.firewall_mark"
    ></el-input-number>
    <el-input
      class="device_form"
      placeholder="List Networks"
      v-model="newDeviceForm.networks"
    ></el-input>
    <div class="modal__buttons">
      <el-button @click="$emit('close')">Cancel</el-button>
      <el-button type="primary" @click="createDevice">Create</el-button>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator'
import { deviceApi } from '@/api/interface'

@Component({
  name: 'DeviceModal'
})
export default class DeviceModal extends Vue {
  private newDeviceForm = {
    name: '',
    listen_port: 0,
    private_key: '',
    firewall_mark: 0,
    networks: ''
  }

  public async createDevice(): Promise<void> {
    console.log(this.newDeviceForm)
    const newDevice = {
      ...this.newDeviceForm,
      networks: this.newDeviceForm.networks.split(',')
    }

    console.log(newDevice)
    await deviceApi.createDevice(newDevice)
  }
}
</script>

<style scoped>
.device_form {
  margin-bottom: 20px;
}

.modal_title {
  text-align: center;
}

.modal__buttons {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
</style>
