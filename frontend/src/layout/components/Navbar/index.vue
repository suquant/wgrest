<template>
  <div class="navbar">
    <hamburger
      id="hamburger-container"
      :is-active="sidebar.opened"
      class="hamburger-container"
      @toggle-click="toggleSideBar"
    />
    <breadcrumb
      id="breadcrumb-container"
      class="breadcrumb-container"
    />
    <div class="right-menu">
      <el-button
        class="navbar__button"
        type="primary"
        icon="el-icon-refresh"
        @click="updateList"
      >
        Update List
      </el-button>
      <el-button
        class="navbar__button"
        type="primary"
        icon="el-icon-plus"
        v-if="CurrentButton === NewAddButtonEnum.Device"
        @click="openModal('device')"
      >
        New Device
      </el-button>
      <el-button
        class="navbar__button"
        type="primary"
        icon="el-icon-plus"
        @click="openModal('peer')"
        v-if="CurrentButton === NewAddButtonEnum.Peer"
      >
        New Peer
      </el-button>
<!--      <el-dropdown-->
<!--        class="avatar-container right-menu-item hover-effect"-->
<!--        trigger="click"-->
<!--      >-->
<!--        <div class="avatar-wrapper">-->
<!--          <img-->
<!--            :src="avatar+'?imageView2/1/w/80/h/80'"-->
<!--            class="user-avatar"-->
<!--          >-->
<!--          <i class="el-icon-caret-bottom" />-->
<!--        </div>-->
<!--      </el-dropdown>-->
<!--      <el-dropdown-menu slot="dropdown">-->
<!--        <router-link to="/">-->
<!--          <el-dropdown-item>-->
<!--            Home-->
<!--          </el-dropdown-item>-->
<!--        </router-link>-->
<!--        <el-dropdown-item divided>-->
<!--            <span-->
<!--              style="display:block;"-->
<!--              @click="logout"-->
<!--            >LogOut</span>-->
<!--        </el-dropdown-item>-->
<!--      </el-dropdown-menu>-->
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator'
import { AppModule } from '@/store/modules/app'
import { UserModule } from '@/store/modules/user'
import { ModalModule } from '@/store/modules/modal'
import Breadcrumb from '@/components/Breadcrumb/index.vue'
import Hamburger from '@/components/Hamburger/index.vue'

import { emitter } from '@/utils/emmiter'

export enum NewAddButtonEnum {
  Device = 'Device',
  Peer = 'Peer'
}

@Component({
  name: 'Navbar',
  components: {
    Breadcrumb,
    Hamburger
  }
})
export default class extends Vue {
  NewAddButtonEnum = NewAddButtonEnum
  get sidebar() {
    return AppModule.sidebar
  }

  get device() {
    return AppModule.device.toString()
  }

  get avatar() {
    return UserModule.avatar
  }

  private toggleSideBar() {
    AppModule.ToggleSideBar(false)
  }

  public openModal(name: string): void {
    ModalModule.openModal(name)
  }

  get CurrentButton(): String {
    if (this.$route.name === 'Device') {
      return NewAddButtonEnum.Device
    }
    if (this.$route.name === 'Peer') {
      return NewAddButtonEnum.Peer
    }

    return ''
  }

  private async logout() {
    await UserModule.LogOut()
    this.$router.push(`/login?redirect=${this.$route.fullPath}`)
  }

  private updateList(): void {
    if (this.$route.name === 'Device') {
      emitter.emit('updateDevice')

      this.$message({
        type: 'success',
        message: 'Devices updated'
      })
    }
    if (this.$route.name === 'Peer') {
      emitter.emit('updatePeer')

      this.$message({
        type: 'success',
        message: 'Peers updated'
      })
    }
  }
}
</script>

<style lang="scss" scoped>
.navbar {
  height: 50px;
  overflow: hidden;
  position: relative;
  background: #fff;
  box-shadow: 0 1px 4px rgba(0,21,41,.08);

  .hamburger-container {
    line-height: 46px;
    height: 100%;
    float: left;
    padding: 0 15px;
    cursor: pointer;
    transition: background .3s;
    -webkit-tap-highlight-color:transparent;

    &:hover {
      background: rgba(0, 0, 0, .025)
    }
  }

  .breadcrumb-container {
    float: left;
  }

  .right-menu {
    float: right;
    height: 100%;
    line-height: 50px;

    &:focus {
      outline: none;
    }

    .right-menu-item {
      display: inline-block;
      padding: 0 8px;
      height: 100%;
      font-size: 18px;
      color: #5a5e66;
      vertical-align: text-bottom;

      &.hover-effect {
        cursor: pointer;
        transition: background .3s;

        &:hover {
          background: rgba(0, 0, 0, .025)
        }
      }
    }

    .avatar-container {
      margin-right: 30px;

      .avatar-wrapper {
        margin-top: 5px;
        position: relative;

        .user-avatar {
          cursor: pointer;
          width: 40px;
          height: 40px;
          border-radius: 10px;
        }

        .el-icon-caret-bottom {
          cursor: pointer;
          position: absolute;
          right: -20px;
          top: 25px;
          font-size: 12px;
        }
      }
    }
  }
}

.navbar__button {
  margin-right: 30px;
}
</style>
