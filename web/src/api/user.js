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

// @Summary 获取用户信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} json "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /user/getUserInfo [get]
export const getUserInfo = () => {
    return service({
        url: '/user/getUserInfo',
        method: 'get'
    })
}
