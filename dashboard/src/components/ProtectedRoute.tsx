import React from 'react'
import { Route, Redirect, RouteProps } from 'react-router-dom'

interface ProtectedRouteProps extends RouteProps {
    component: React.ComponentType<any>
}

const ProtectedRoute: React.FC<ProtectedRouteProps> = ({ component: Component, ...rest }) => {
    const isAuthenticated = !!localStorage.getItem('token')

    return (
        <Route
            {...rest}
            render={(props) =>
                isAuthenticated ? (
                    <Component {...props} />
                ) : (
                    <Redirect to="/login" />
                )
            }
        />
    )
}

export default ProtectedRoute
