import { useEffect, useRef } from 'react'
import { useHistory, useLocation } from 'react-router-dom'
import { useMutation } from '@tanstack/react-query'
import { api } from '../utils/api'

interface CallbackResponse {
    token: string
    refreshToken: string
}

export const Callback = () => {
    const history = useHistory()
    const location = useLocation()
    const hasMutated = useRef(false)

    const mutation = useMutation<CallbackResponse, Error, string>({
        mutationFn: (callbackUrl: string) =>
            api.url(`/auth/google/callback${callbackUrl}&provider=google`)
                .get()
                .json(),
        onSuccess: (data) => {
            localStorage.setItem('token', data.token)
            localStorage.setItem('refreshToken', data.refreshToken)
            history.push('/user-info')
        },
        onError: (error) => {
            console.error('Callback error:', error)
        }
    })

    useEffect(() => {
        if (!hasMutated.current && location.search && !mutation.isSuccess) {
            hasMutated.current = true
            mutation.mutate(location.search)
        }
    }, [location.search, mutation])

    return (
        <div>
            <h1>Login</h1>
            {mutation.isPending && <div>Processing login...</div>}
            {mutation.isError && <div>Error: {mutation.error.message}</div>}
            {mutation.isSuccess && <div>Login successful, redirecting...</div>}
        </div>
    )
}

export const Logout = () => {
    const history = useHistory()

    localStorage.removeItem('token')
    localStorage.removeItem('refreshToken')
    history.push('/')

    return null
}
