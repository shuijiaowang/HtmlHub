import service from '@/utils/request.js'

export const uploadHtml = (data) => {
  return service({
    url: '/html/upload',
    method: 'post',
    data
  })
}

export const getMyHtmlList = () => {
  return service({
    url: '/html/my',
    method: 'get'
  })
}
