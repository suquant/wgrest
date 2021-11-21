import { MessageBox, Message } from 'element-ui'
export const openModal = (instance: any, config: any) => {
  MessageBox.prompt('Please input your access token', 'Access Token', {
    confirmButtonText: 'OK',
    cancelButtonText: 'Cancel'
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore
  }).then(({ value }) => {
    localStorage.setItem('accessToken', value)
    Message({
      type: 'success',
      message: 'Your token is saved'
    })

    config.headers.Authorization = `Bearer ${value}`

    instance.configuration.baseOptions.headers.Authorization = `Bearer ${value}`

    location.reload()

    return instance.axios(config)
  }).catch(() => {
    Message({
      type: 'info',
      message: 'Input canceled'
    })
  })
}
