import React from 'react'
import { useQuery } from '@tanstack/react-query'
import { api } from '../utils/api'
import { Link } from 'react-router-dom'

interface UserInfoResponse {
    email: string
    firstName: string
    lastName: string
}

const UserInfo: React.FC = () => {
    const { data, isLoading, error } = useQuery<UserInfoResponse, Error>({
        queryKey: ['userInfo'],
        queryFn: () => api.auth(`Bearer ${localStorage.getItem('token')}`).url('/api/user/profile').get().json(),
    })

    if (isLoading) return <div>Loading user info...</div>
    if (error) return <div>Error: {error.message}</div>

    return (
        <div>
            <h1>User Info</h1>
            <p>Email: {data?.email}</p>
            <p>First Name: {data?.firstName}</p>
            <p>Last Name: {data?.lastName}</p>
            <p><Link to="/logout">Logout</Link></p>
        </div>
    )
}

export default UserInfo
