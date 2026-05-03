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

export const deleteHtmlRecord = (id) => {
  return service({
    url: `/html/${id}`,
    method: 'delete'
  })
}

export const updateHtmlVisibility = (id, visibility) => {
  return service({
    url: `/html/${id}/visibility`,
    method: 'put',
    data: { visibility }
  })
}

export const updateHtmlDescription = (id, description) => {
  return service({
    url: `/html/${id}/description`,
    method: 'put',
    data: { description }
  })
}

export const updateHtmlContent = (id, htmlContent) => {
  return service({
    url: `/html/${id}/content`,
    method: 'put',
    data: { htmlContent }
  })
}
