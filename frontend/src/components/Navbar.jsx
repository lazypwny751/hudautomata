import { useUIStore, useAuthStore } from '../lib/store';
import { Menu, X, Moon, Sun, LogOut, User } from 'lucide-preact';
import { route } from 'preact-router';

export default function Navbar() {
  const { sidebarOpen, toggleSidebar, theme, setTheme } = useUIStore();
  const { user, logout } = useAuthStore();

  const handleLogout = async () => {
    await logout();
    route('/login');
  };

  const toggleTheme = () => {
    setTheme(theme === 'light' ? 'dark' : 'light');
  };

  return (
    <div class="navbar bg-base-100 shadow-md">
      <div class="flex-none">
        <button class="btn btn-square btn-ghost" onClick={toggleSidebar}>
          {sidebarOpen ? <X size={24} /> : <Menu size={24} />}
        </button>
      </div>
      
      <div class="flex-1">
        <a class="btn btn-ghost normal-case text-xl">
          <span class="text-primary">Hud</span>Automata
        </a>
      </div>
      
      <div class="flex-none gap-2">
        <button class="btn btn-ghost btn-circle" onClick={toggleTheme}>
          {theme === 'light' ? <Moon size={20} /> : <Sun size={20} />}
        </button>

        <div class="dropdown dropdown-end">
          <label tabIndex={0} class="btn btn-ghost btn-circle avatar placeholder">
            <div class="bg-neutral text-neutral-content rounded-full w-10">
              <User size={24} />
            </div>
          </label>
          <ul tabIndex={0} class="mt-3 z-[1] p-2 shadow menu menu-sm dropdown-content bg-base-100 rounded-box w-52">
            <li class="menu-title">
              <span>{user?.username}</span>
              <span class="badge badge-sm badge-primary">{user?.role}</span>
            </li>
            <li><a onClick={handleLogout}><LogOut size={16} /> Çıkış Yap</a></li>
          </ul>
        </div>
      </div>
    </div>
  );
}
