import { Router, Route } from 'preact-router';
import { useEffect } from 'preact/hooks';
import { useAuthStore, useUIStore } from './lib/store';

// Pages
import Login from './pages/Login';
import Dashboard from './pages/Dashboard';
import Users from './pages/Users';
import Transactions from './pages/Transactions';
import Automation from './pages/Automation';

// Components
import Navbar from './components/Navbar';
import Sidebar from './components/Sidebar';

function PrivateRoute({ component: Component, ...rest }) {
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated);
  
  if (!isAuthenticated) {
    if (typeof window !== 'undefined') {
      window.location.href = '/login';
    }
    return null;
  }

  return <Component {...rest} />;
}

export function App() {
  const { checkAuth, isAuthenticated } = useAuthStore();
  const { theme, setTheme } = useUIStore();

  useEffect(() => {
    // Set initial theme
    document.documentElement.setAttribute('data-theme', theme);
    
    // Check authentication
    if (isAuthenticated) {
      checkAuth();
    }
  }, []);

  return (
    <div class="min-h-screen">
      {isAuthenticated && <Navbar />}
      
      <div class="flex">
        {isAuthenticated && <Sidebar />}
        
        <main class="flex-1">
          <Router>
            <Route path="/login" component={Login} />
            <PrivateRoute path="/" component={Dashboard} />
            <PrivateRoute path="/users" component={Users} />
            <PrivateRoute path="/transactions" component={Transactions} />
            <PrivateRoute path="/automation" component={Automation} />
          </Router>
        </main>
      </div>
    </div>
  );
}

