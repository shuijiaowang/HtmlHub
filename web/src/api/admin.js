import service from '@/utils/request.js'

export const getAdminUsers = (params) => {
  return service({
    url: '/admin/users',
    method: 'get',
    params
  })
}

export const getAdminHtmlList = (params) => {
  return service({
    url: '/admin/html',
    method: 'get',
    params
  })
}

export const getAdminHtmlDetail = (id) => {
  return service({
    url: `/admin/html/${id}`,
    method: 'get'
  })
}

export const updateAdminHtmlApproval = (id, approvalStatus) => {
  return service({
    url: `/admin/html/${id}/approval`,
    method: 'put',
    data: { approvalStatus }
  })
}

export const updateAdminHtmlSubdomain = (id, subdomain) => {
  return service({
    url: `/admin/html/${id}/subdomain`,
    method: 'put',
    data: { subdomain }
  })
}

export const deleteAdminHtml = (id) => {
  return service({
    url: `/admin/html/${id}`,
    method: 'delete'
  })
}
