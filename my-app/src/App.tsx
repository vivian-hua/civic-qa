import './App.css';
import React, { useEffect, useState} from "react";
import { Route, NavLink, Redirect, useLocation } from "react-router-dom";
import { ProfileHeader } from "./Profile/ProfileHeader";
import { Dashboard } from './Views/Dashboard';
import { General } from './Views/General';
import { Casework } from './Views/Casework';
import { Responses } from './Views/Responses';
import { EngagementReports } from './Views/EngagementReports';
import { Templates } from './Views/Templates';
import { Settings } from './Views/Settings';
import { Login } from './Views/Login';
import * as Constants from './Constants/Constants';

export default function App() {
  const authToken = localStorage.getItem("Authorization");
  const [auth, setAuth] = useState((authToken != "") && (authToken != null));
  const [path, setPath] = useState("/dashboard");
  const location = useLocation();

  useEffect(() => {
    setPath(location.pathname);
  }, [location]);

  const userLogout = async(e: any) => {
    e.preventDefault();
    var authToken = localStorage.getItem("Authorization") || "";
    const response = await fetch("http://localhost/v0/logout", {
      method: "POST",
      headers: new Headers({
          "Authorization": authToken
      })
    });
    if (response.status >= 300) {
      console.log("Failed to logout");
    }
    localStorage.removeItem("Authorization");
    setAuth(false);
    return(<Redirect to="/login"></Redirect>);
  }

  function userLogin(authToken: string) {
    setAuth(true);
    localStorage.setItem("Authorization", authToken);
    return(<Redirect to="/dashboard"></Redirect>);
  }

  return (
      <div className="App">
        <Route path="/login" component={() => <Login userLogin={userLogin}/>}></Route>
        {auth ? <Redirect to="/dashboard"/> : <Redirect to="/login"/>}
        {auth ?
          <div>
            <div>
              <div className="profile-header">
                <ProfileHeader></ProfileHeader>
              </div>
              <nav className="nav-bar">
                <h1 className="title">{Constants.Title}</h1>
                <ul>
                  <li>{path == "/dashboard" ? <img src="./assets/icons/active-pie.png"/> : <img src="./assets/icons/pie.png"/>}<NavLink className="nav-link" activeClassName="active-link" to="/dashboard">{Constants.Dashboard}</NavLink></li>
                  <li className="dashboard-sub-li"><NavLink className="nav-link" activeClassName="active-link" to="/general">{Constants.GeneralInquiries}</NavLink></li>
                  <li className="dashboard-sub-li"><NavLink className="nav-link" activeClassName="active-link" to="/casework">{Constants.Casework}</NavLink></li>
                  <li>{path =="/responses" ? <img src="./assets/icons/active-inbox.png"/> : <img src="./assets/icons/inbox.png"/>}<NavLink className="nav-link" activeClassName="active-link" to="/responses">{Constants.Responses}</NavLink></li>
                  <li>{path =="/engagement-reports" ? <img src="./assets/icons/active-stats.png"/> :<img src="./assets/icons/stats.png"/>}<NavLink className="nav-link" activeClassName="active-link" to="/engagement-reports">{Constants.EngagementReports}</NavLink></li>
                </ul>
                <div className="compose-email-btn-container">
                  <hr className="solid"/>
                </div>
                <ul>
                <li><img src="./assets/icons/settings.png"/><NavLink className="nav-link" activeClassName="active-link" to="/settings">{Constants.Settings}</NavLink></li>
                  <li><img src="./assets/icons/logout.png"/><button className="logout-btn" onClick={userLogout}>{Constants.Logout}</button></li>
                </ul>
              </nav>
            </div>
            <Route path="/dashboard" component={Dashboard}></Route>
            <Route path="/general" component={General}></Route>
            <Route path="/casework" component={Casework}></Route>
            <Route path="/responses" component={Responses}/>
            <Route path="/engagement-reports" component={EngagementReports}/>
            <Route path="/settings" component={Settings}/>
          </div> : null}
        </div>
  );
}
