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
    console.log(this.newPeerForm)
    const newPeer = {
      ...this.newPeerForm,
      allowed_ips: this.newPeerForm.allowed_ips.split(',')
    }

    console.log(newPeer)
    await deviceApi.createDevicePeer(this.$route.params.id, newPeer)
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
