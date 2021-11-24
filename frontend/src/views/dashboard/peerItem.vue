<template>
    <el-card class="box-card" shadow="hover">
      <div slot="header" class="card__header">
        <span>public key: {{ item.url_safe_public_key }}</span>
        <i class="el-icon-setting" @click="drawer = true"></i>
      </div>
      <div class="text item">
        <div class="info__peer">
          <span>receive: {{ item.receive_bytes }}</span>
          <span>transmit: {{ item.transmit_bytes }}</span>
        </div>
      </div>
      <div class="peer__buttons">
        <el-button
          type="primary"
          @click="getQrCode"
        >
          qr code
        </el-button>
        <el-button
          icon="el-icon-download"
          type="info"
          @click="getQuickConf"
        >
          quick.conf
        </el-button>
      </div>
      <el-drawer
        :title="`${item.url_safe_public_key}`"
        :visible.sync="drawer"
        direction="rtl"
        >
        <div class="detail__info">
          <span>url_safe_public_key: {{ item.url_safe_public_key }}</span>
          <span>public_key: {{ item.public_key }}</span>
          <span>last_handshake_time: {{ item.last_handshake_time }}</span>
          <span>persistent_keepalive_interval: {{ item.persistent_keepalive_interval }}</span>
          <span>receive_bytes: {{ item.receive_bytes }}</span>
          <span>transmit_bytes: {{ item.transmit_bytes }}</span>
          <span>endpoint: {{ item.endpoint }}</span>
          <span>allowed_ips: {{ item.allowed_ips.join(',') }}</span>
        </div>
      </el-drawer>
      <el-dialog
        title="QR Code"
        :visible.sync="dialogVisible"
        width="30%"
      >
        <div class="qr__dialog">
          <img :src="qrCode" alt="">
        </div>
      </el-dialog>
    </el-card>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator'
import { Peer } from 'wgrest/dist/models'
import { deviceApi } from '@/api/interface'

@Component({
  name: 'peerItem'
})
export default class peerItem extends Vue {
  @Prop({ default: {} }) item!: Peer

  private drawer = false
  private dialogVisible = false

  private qrCode: string | unknown = ''

  public async getQrCode(): Promise<void> {
    const { data } = await deviceApi.getDevicePeerQuickConfigQRCodePNG(this.$route.params.id, this.item.url_safe_public_key, '', { responseType: 'arraybuffer' })
    const blob = new Blob([data])
    this.qrCode = URL.createObjectURL(blob)
    this.dialogVisible = true
  }

  public async getQuickConf(): Promise<void> {
    const { data } = await deviceApi.getDevicePeerQuickConfig(this.$route.params.id, this.item.url_safe_public_key)
    const blob = new Blob([data])
    const link = document.createElement('a')
    link.href = window.URL.createObjectURL(blob)
    link.download = `${this.item.url_safe_public_key}.conf`
    link.click()
  }
}
</script>

<style lang="scss" scoped>
.text {
  font-size: 14px;
}

.item {
  margin-bottom: 18px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
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
  font-size: 12px;
}

.qr__code {
  width: 200px;
  height: 200px;
  display: flex;
  align-items: center;
  justify-content: center;

  img {
    width: 100%;
    height: 100%;
  }
}

.info__peer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
}

.detail__info {
  display: flex;
  flex-direction: column;

  padding: 20px;

  span {
    margin-bottom: 20px;
  }
}

.peer__buttons {
  margin-top: 30px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.qr__dialog {
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
