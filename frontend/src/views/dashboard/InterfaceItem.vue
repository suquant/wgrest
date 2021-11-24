<template>
  <div>
  <router-link :to="{path: `/devices/${item.name}`}">
   <el-card class="box-card" shadow="hover">
    <div slot="header" class="card__header">
      <span>name: {{ item.name }}</span>
     <div>
       <span style="margin-right: 20px">port: {{ item.listen_port }}</span>
       <i class="el-icon-setting" @click.prevent="openDrawer"></i>
     </div>
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
    <el-drawer
      :title="item.name"
      :visible.sync="drawer"
      direction="rtl"
    >
      <div class="interface__device-options">
        <h3>Options</h3>
        <div class="interface__device-options-block">
          <h4>
            Allowed IPS
            <span
              @click="addNewIps"
              class="interface__device-options-add"
            >
              +
            </span>
          </h4>
          <p
            v-for="(item, i) in options.allowed_ips"
            :key="item"
            class="interface__device-options-item"
          >
            {{ item }}
            <span
              @click="deleteIPS(i)"
              class="interface__device-options-delete"
            >
              x
            </span>
          </p>
        </div>
        <div class="interface__device-options-block">
          <h4>
            DNS Servers
            <span
              @click="addNewDNS"
              class="interface__device-options-add"
            >
              +
            </span>
          </h4>
          <p
            v-for="(item, i) in options.dns_servers"
            :key="item"
            class="interface__device-options-item"
          >
            {{ item }}
            <span
              @click="deleteDNS(i)"
              class="interface__device-options-delete"
            >
              x
            </span>
          </p>
        </div>
        <h4>host</h4>
        <el-input style="margin-left: 10px" v-model="options.host"></el-input>
        <el-button
          @click="updateOptionsInterface"
          class="interface__device-options-button"
          type="primary"
        >
          Save
        </el-button>
      </div>
    </el-drawer>
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator'
import { Device, DeviceOptions } from 'wgrest/dist/models'
import { deviceApi } from '@/api/interface'

@Component({
  name: 'interfaceItem'
})
export default class interfaceItem extends Vue {
  @Prop({ default: {} }) item!: Device

  private drawer = false

  private options: DeviceOptions = {
    allowed_ips: [],
    dns_servers: [],
    host: ''
  }

  public async openDrawer(e: Event): Promise<void> {
    e.stopPropagation()

    const { data } = await deviceApi.getDeviceOptions(this.item.name)

    this.options = data

    this.drawer = true
  }

  private deleteIPS(id: number): void {
    this.options.allowed_ips.splice(id, 1)
  }

  private addNewIps(): void {
    this.$prompt('Please input your ips', 'Allowed IPS', {
      confirmButtonText: 'OK',
      cancelButtonText: 'Cancel'
      // eslint-disable-next-line @typescript-eslint/ban-ts-comment
      // @ts-ignore
    }).then(({ value }) => {
      this.options.allowed_ips.push(value)
    })
  }

  private deleteDNS(id: number): void {
    this.options.dns_servers.splice(id, 1)
  }

  private addNewDNS(): void {
    this.$prompt('Please input your dns', 'DNS Server', {
      confirmButtonText: 'OK',
      cancelButtonText: 'Cancel'
      // eslint-disable-next-line @typescript-eslint/ban-ts-comment
      // @ts-ignore
    }).then(({ value }) => {
      this.options.dns_servers.push(value)
    })
  }

  private async updateOptionsInterface(): Promise<void> {
    await deviceApi.updateDeviceOptions(this.options, this.item.name)

    this.$message({
      type: 'success',
      message: 'Options is updated'
    })

    this.drawer = false
  }
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

.interface__device-options {
  padding: 20px;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.interface__device-options-block {
  display: flex;
  flex-direction: column;
}

.interface__device-options-item {
  padding-left: 10px;
  margin: 0 0 10px 0;

  display: flex;
  align-items: center;
}

.interface__device-options-delete {
  cursor: pointer;
  margin-left: 20px;
}

.interface__device-options-add {
  cursor: pointer;
  margin-left: 10px;
}

.interface__device-options-button {
  margin: 40px auto 0;
  width: 150px;
}
</style>
