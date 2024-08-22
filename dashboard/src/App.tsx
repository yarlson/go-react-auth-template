import { Switch, Route, Redirect } from 'react-router-dom'
import Login from './components/Login'
import ProtectedRoute from './components/ProtectedRoute'
import UserInfo from './components/UserInfo'
import { Callback, Logout } from "./components/Callback.tsx";

const App = () => {
    return (
        <Switch>
            <Route exact path="/login" component={Login} />
            <Route path="/callback" component={Callback} />
            <ProtectedRoute path="/user-info" component={UserInfo} />
            <ProtectedRoute path="/logout" component={Logout} />
            <Redirect from="/" to="/user-info" />
        </Switch>
    )
}

export default App