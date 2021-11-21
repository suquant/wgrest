<template>
  <transition name="modal">
    <div v-if="getModal.length" class="modal">
      <div class="modal__container">
        <div class="modal__overlay"  @click.self="closeModal"></div>
        <div class="modal__content">
          <component @close="closeModal" :is="getModal" />
        </div>
      </div>
    </div>
  </transition>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator'
import { ModalModule } from '@/store/modules/modal'

function loadModalsComponents() {
  const locales = require.context(
    './components',
    true,
    /[A-Za-z0-9-_,\s]+\.vue$/i
  )
  const components = {}
  locales.keys().forEach(key => {
    const matched = key.match(/([A-Za-z0-9-_]+)\./i)
    if (matched && matched.length > 1) {
      const component = matched[1]
      // eslint-disable-next-line @typescript-eslint/ban-ts-comment
      // @ts-ignore
      components[component] = locales(key).default
    }
  })
  console.log(components)
  return components
}

@Component({
  name: 'Modal',
  components: loadModalsComponents()
})
export default class Modal extends Vue {
  get getModal(): string {
    return ModalModule.modal
  }

  public closeModal(): void {
    ModalModule.closeModal()
  }
}
</script>
<style lang="scss" scoped>
.modal {
  position: fixed;
  z-index: 1100;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  overflow-y: auto;
  -webkit-overflow-scrolling: touch;
  background: rgba(0, 0, 0, 0.8);
  backdrop-filter: blur(60px);
}

.modal__container {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  box-sizing: border-box;
  min-height: 100%;
  width: 100%;
  padding: 40px 0;
}

.modal__content {
  z-index: 5;
  display: flex;
  justify-content: center;
  max-width: 460px;
  width: 100%;
  padding: 40px 60px;
  margin: 0 auto;
  overflow: hidden;
  background: #fff;
  border-radius: 16px;
}

.modal__overlay {
  position: absolute;
  z-index: 1;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  cursor: pointer;
}

//TRANSITION
.modal-enter,
.modal-leave-to {
  background: transparent;

  .modal__content {
    transform: scale(0);
  }
}

.modal-enter-active,
.modal-leave-active {
  transition: all 0.3s;

  .modal__content {
    transition: all 0.3s ease-out;
  }
}

</style>
