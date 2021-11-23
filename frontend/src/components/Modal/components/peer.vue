<template>
  <div>
    <p class="device_form modal_title">Add new Peer</p>
    <el-input
      class="device_form"
      placeholder="Private Key"
      v-model="newPeerForm.private_key"
    ></el-input>
    <el-input
      class="device_form"
      placeholder="Public Key"
      v-model="newPeerForm.public_key"
    ></el-input>
    <el-input
      class="device_form"
      placeholder="Preshared Key"
      v-model="newPeerForm.preshared_key"
    ></el-input>
    <el-input
      class="device_form"
      placeholder="Allowed Ips"
      v-model="newPeerForm.allowed_ips"
    ></el-input>
    <span class="input__info">127.0.0.1, 127.0.0.1</span>
    <el-input
      class="device_form"
      placeholder="Persistent KeepAlive Interval"
      v-model="newPeerForm.persistent_keepalive_interval"
    ></el-input>
    <el-input
      class="device_form"
      placeholder="Endpoint"
      v-model="newPeerForm.endpoint"
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

import { emitter } from '@/utils/emmiter'

@Component({
  name: 'DeviceModal'
})
export default class DeviceModal extends Vue {
  private newPeerForm = {
    private_key: '',
    public_key: '',
    preshared_key: '',
    allowed_ips: '',
    persistent_keepalive_interval: '',
    endpoint: ''
  }

  public async createDevice(): Promise<void> {
    const newPeerData = JSON.parse(JSON.stringify(this.newPeerForm))

    Object.keys(newPeerData).forEach(key => {
      if (!newPeerData[key]) {
        delete newPeerData[key]
      }
    })

    const newPeer = {
      ...newPeerData,
      allowed_ips: Object.keys(newPeerData).includes('allowed_ips') ? newPeerData.allowed_ips.split(',') : null
    }

    await deviceApi.createDevicePeer(this.$route.params.id, newPeer)
    this.$emit('close')
    emitter.emit('updatePeer')
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

.input__info {
  color: gray;
  font-size: 10px;
  display: block;
  position: relative;
  top: -20px;
}
</style>
