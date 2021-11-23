import { VuexModule, Module, Action, Mutation, getModule } from 'vuex-module-decorators'
import store from '@/store'

export interface ModalState {
  modal: string
  payload?: unknown
}

@Module({ dynamic: true, store, name: 'modal' })
class Modal extends VuexModule implements ModalState {
  public modal = ''

  @Mutation
  private SET_MODAL(modal: string) {
    this.modal = modal
  }

  @Mutation
  private CLEAR_MODAL() {
    this.modal = ''
  }

  @Action
  public openModal(modal: string): void {
    this.SET_MODAL(modal)
  }

  @Action
  public closeModal(): void {
    this.CLEAR_MODAL()
  }
}

export const ModalModule = getModule(Modal)
