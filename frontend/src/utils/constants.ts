import { $fetch } from 'ofetch'

export const MAX_PASSWORD_LEN = 30
export const MAX_USERNAME_LEN = 30
export const backendFetch = $fetch.create({
    baseURL: '/api/backend',
})
