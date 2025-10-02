import { useState, useEffect } from 'preact/hooks';
import { dashboardAPI } from '../lib/api';
import { formatCurrency, formatRelativeTime, getTransactionColor, getSourceBadge } from '../lib/utils';
import { Users, Wallet, TrendingUp, Activity } from 'lucide-preact';

export default function Dashboard() {
  const [stats, setStats] = useState(null);
  const [recentActivities, setRecentActivities] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadDashboard();
  }, []);

  const loadDashboard = async () => {
    try {
      const [statsData, activitiesData] = await Promise.all([
        dashboardAPI.getStats(),
        dashboardAPI.getRecent(),
      ]);
      setStats(statsData);
      setRecentActivities(activitiesData);
    } catch (error) {
      console.error('Failed to load dashboard:', error);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div class="flex items-center justify-center min-h-screen">
        <span class="loading loading-spinner loading-lg"></span>
      </div>
    );
  }

  const statCards = [
    {
      title: 'Toplam Kullanıcı',
      value: stats?.total_users || 0,
      icon: Users,
      color: 'text-primary',
      bgColor: 'bg-primary/10',
    },
    {
      title: 'Aktif Kullanıcı',
      value: stats?.active_users || 0,
      icon: Activity,
      color: 'text-success',
      bgColor: 'bg-success/10',
    },
    {
      title: 'Toplam Bakiye',
      value: formatCurrency(stats?.total_balance || 0),
      icon: Wallet,
      color: 'text-secondary',
      bgColor: 'bg-secondary/10',
    },
    {
      title: "Bugünkü Gelir",
      value: formatCurrency(stats?.today_revenue || 0),
      icon: TrendingUp,
      color: 'text-accent',
      bgColor: 'bg-accent/10',
    },
  ];

  return (
    <div class="p-6 space-y-6">
      <div>
        <h1 class="text-3xl font-bold">Dashboard</h1>
        <p class="text-base-content/70 mt-1">Sistem genel bakış ve istatistikler</p>
      </div>

      {/* Stats Grid */}
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        {statCards.map((stat, index) => (
          <div key={index} class="stat-card">
            <div class="flex items-center justify-between">
              <div>
                <p class="text-base-content/70 text-sm">{stat.title}</p>
                <p class="text-2xl font-bold mt-2">{stat.value}</p>
              </div>
              <div class={`p-4 rounded-full ${stat.bgColor}`}>
                <stat.icon size={24} class={stat.color} />
              </div>
            </div>
          </div>
        ))}
      </div>

      {/* Recent Activities */}
      <div class="card bg-base-100 shadow-xl">
        <div class="card-body">
          <h2 class="card-title">Son İşlemler</h2>
          
          <div class="overflow-x-auto">
            <table class="table">
              <thead>
                <tr>
                  <th>Kullanıcı</th>
                  <th>İşlem</th>
                  <th>Tutar</th>
                  <th>Kaynak</th>
                  <th>Zaman</th>
                </tr>
              </thead>
              <tbody>
                {recentActivities.map((activity) => (
                  <tr key={activity.id}>
                    <td>
                      <div class="font-semibold">{activity.user?.name || '-'}</div>
                      <div class="text-sm text-base-content/70">{activity.user?.rfid_card_id}</div>
                    </td>
                    <td>
                      <span class={`badge ${getTransactionColor(activity.type)}`}>
                        {activity.type}
                      </span>
                    </td>
                    <td class="font-mono">{formatCurrency(activity.amount)}</td>
                    <td>
                      <span class={`badge ${getSourceBadge(activity.source).class}`}>
                        {getSourceBadge(activity.source).text}
                      </span>
                    </td>
                    <td class="text-sm text-base-content/70">
                      {formatRelativeTime(activity.created_at)}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          {recentActivities.length === 0 && (
            <div class="text-center py-8 text-base-content/70">
              Henüz işlem bulunmuyor
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
