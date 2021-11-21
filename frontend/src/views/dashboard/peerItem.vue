<template>
    <el-card class="box-card" shadow="hover">
      <div slot="header" class="card__header">
        <span>public key: {{ item.url_safe_public_key }}</span>
        <i class="el-icon-setting" @click="drawer = true"></i>
      </div>
      <div class="text item">
        <div class="qr__code">
          qr
        </div>
        <div class="info__peer">
          <span>receive: {{ item.receive_bytes }}</span>
          <span>transmit: {{ item.transmit_bytes }}</span>
        </div>
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

  private qrCode = ''

  public async getQrCode(): Promise<void> {
    const { data } = await deviceApi.getDevicePeerQuickConfigQRCodePNG(this.$route.params.id, this.item.url_safe_public_key)
    this.qrCode = data
  }

  created() {
    this.getQrCode()
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
</style>
