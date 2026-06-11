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

export const updateHtmlPublishMode = (id, publishMode) => {
  return service({
    url: `/html/${id}/publish-mode`,
    method: 'put',
    data: { publishMode }
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

export const getPublicHtmlList = (params) => {
  return service({
    url: '/html/public',
    method: 'get',
    params,
    donNotShowLoading: true
  })
}

export const likeHtmlRecord = (id) => {
  return service({
    url: `/html/${id}/like`,
    method: 'post'
  })
}

export const unlikeHtmlRecord = (id) => {
  return service({
    url: `/html/${id}/like`,
    method: 'delete'
  })
}

export const getMyLikedList = () => {
  return service({
    url: '/html/liked',
    method: 'get'
  })
}

export const getMySyncDataList = () => {
  return service({
    url: '/html/data/my',
    method: 'get'
  })
}

export const exportMySyncData = (id) => {
  return service({
    url: `/html/data/my/${id}/export`,
    method: 'get'
  })
}

export const deleteMySyncData = (id) => {
  return service({
    url: `/html/data/my/${id}`,
    method: 'delete'
  })
}

export const clearMySyncData = () => {
  return service({
    url: '/html/data/my',
    method: 'delete'
  })
}
