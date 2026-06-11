import service from '@/utils/request.js'
// @Summary 用户登录
// @Produce  application/json
// @Param data body {email:"string",password:"string"}
// @Router /base/login [post]
export const login = (data) => {
    return service({
        url: '/user/login',
        method: 'post',
        data: data
    })
}

// @Summary 用户注册
// @Produce  application/json
// @Param data body {nickname:"string",email:"string",password:"string"}
// @Router /user/register [post]
export const register = (data) => {
    return service({
        url: '/user/register',
        method: 'post',
        data: data
    })
}

// @Summary 获取个人中心资料（信息 + 限制 + 用量）
// @Security ApiKeyAuth
// @Router /user/profile [get]
export const getUserProfile = () => {
    return service({
        url: '/user/profile',
        method: 'get'
    })
}

// @Summary 修改个人资料（昵称）
// @Security ApiKeyAuth
// @Router /user/profile [put]
export const updateUserProfile = (data) => {
    return service({
        url: '/user/profile',
        method: 'put',
        data
    })
}
