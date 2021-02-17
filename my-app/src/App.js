import './App.css';
import { BrowserRouter as Router, Route, Link } from "react-router-dom";
import Dashboard from './Views/Dashboard';
import Inbox from './Views/Inbox';
import EngagementReports from './Views/EngagementReports';
import Templates from './Views/Templates';

function App() {
  return (
    <Router>
      <nav class="nav-bar">
        <ul>
          <li><Link class="nav-link" to="/dashboard">Dashboard</Link></li>
          <li><Link class="nav-link" to="/inbox">Inbox</Link></li>
          <li><Link class="nav-link" to="/engagement-reports">Engagement Reports</Link></li>
          <li><Link class="nav-link" to="/templates">Your Templates</Link></li>
        </ul>
      </nav>
      <Route path="/" exact component={Dashboard}></Route>
      <Route path="/dashboard" exact component={Dashboard}></Route>
      <Route path="/inbox" exact component={Inbox}/>
      <Route path="/engagement-reports" exact component={EngagementReports}/>
      <Route path="/templates" exact component={Templates}/>
      <div className="App">
      </div>
    </Router>
  );
}

export default App;
