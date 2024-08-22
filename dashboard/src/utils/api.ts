import wretch, { Wretch } from 'wretch'

const BASE_URL = 'http://localhost:8080'

interface RefreshResponse {
    token: string
    refreshToken: string
}

export const api: Wretch = wretch(BASE_URL)
    .middlewares([
        next => async (url, opts) => {
            try {
                return await next(url, opts)
            } catch (error) {
                if (error instanceof Error && 'status' in error && error.status === 401) {
                    const refreshToken = localStorage.getItem('refreshToken')
                    if (refreshToken) {
                        try {
                            const response = await wretch(BASE_URL)
                                .url('/auth/refresh')
                                .post({refreshToken})
                                .json<RefreshResponse>()

                            localStorage.setItem('token', response.token)
                            localStorage.setItem('refreshToken', response.refreshToken)

                            if (opts.headers) {
                                opts.headers['Authorization'] = `Bearer ${response.token}`
                            }
                            return next(url, opts)
                        } catch (refreshError) {
                            localStorage.removeItem('token')
                            localStorage.removeItem('refreshToken')
                            window.location.href = '/login'
                        }
                    }
                }
                throw error
            }
        }
    ])