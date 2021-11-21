import { DeviceApi } from 'wgrest'
import { openModal } from '@/utils/modal'

class API extends DeviceApi {
  constructor() {
    super({
      basePath: process.env.VUE_APP_API_DOMAIN,
      baseOptions: {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('accessToken') || ''}`
        }
      }
    })

    this.axios.interceptors.response.use(
      response => response,
      (error) => {
        const originalConfig = error.config

        if (
          error.response.status === 401 ||
          error.response.data.message.includes('invalid key')
        ) {
          return openModal(this, originalConfig)
        }
        return Promise.reject(error)
      })
  }
}

export const deviceApi = new API()
