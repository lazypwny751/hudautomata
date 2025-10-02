import { useUIStore } from '../lib/store';
import { Link } from 'preact-router/match';
import { 
  LayoutDashboard, 
  Users, 
  CreditCard, 
  Shield, 
  ScrollText,
  Settings,
  Activity
} from 'lucide-preact';

const menuItems = [
  { path: '/', label: 'Dashboard', icon: LayoutDashboard },
  { path: '/users', label: 'Kullanıcılar', icon: Users },
  { path: '/transactions', label: 'İşlemler', icon: CreditCard },
  { path: '/automation', label: 'Otomasyon', icon: Activity },
  { path: '/admins', label: 'Adminler', icon: Shield, adminOnly: true },
  { path: '/logs', label: 'Sistem Logları', icon: ScrollText },
  { path: '/settings', label: 'Ayarlar', icon: Settings },
];

export default function Sidebar() {
  const sidebarOpen = useUIStore((state) => state.sidebarOpen);

  if (!sidebarOpen) return null;

  return (
    <div class="w-64 bg-base-200 min-h-screen p-4">
      <ul class="menu gap-2">
        {menuItems.map((item) => (
          <li key={item.path}>
            <Link 
              href={item.path}
              activeClassName="active bg-primary text-primary-content"
              class="flex items-center gap-3"
            >
              <item.icon size={20} />
              <span>{item.label}</span>
            </Link>
          </li>
        ))}
      </ul>
    </div>
  );
}
