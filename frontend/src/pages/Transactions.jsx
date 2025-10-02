import { useState, useEffect } from 'preact/hooks';
import { transactionsAPI } from '../lib/api';
import { formatCurrency, formatDate, getTransactionColor, getSourceBadge } from '../lib/utils';
import { Filter } from 'lucide-preact';

export default function Transactions() {
  const [transactions, setTransactions] = useState([]);
  const [loading, setLoading] = useState(true);
  const [filters, setFilters] = useState({
    type: '',
    source: '',
    from: '',
    to: '',
  });

  useEffect(() => {
    loadTransactions();
  }, [filters]);

  const loadTransactions = async () => {
    try {
      const response = await transactionsAPI.list(filters);
      setTransactions(response.data || []);
    } catch (error) {
      console.error('Failed to load transactions:', error);
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

  return (
    <div class="p-6 space-y-6">
      <div>
        <h1 class="text-3xl font-bold">İşlemler</h1>
        <p class="text-base-content/70 mt-1">Tüm bakiye işlemleri</p>
      </div>

      <div class="card bg-base-100 shadow-xl">
        <div class="card-body">
          <div class="flex gap-4 mb-4 flex-wrap">
            <select 
              class="select select-bordered"
              value={filters.type}
              onChange={(e) => setFilters({ ...filters, type: e.target.value })}
            >
              <option value="">Tüm İşlem Tipleri</option>
              <option value="credit">Credit (Yükleme)</option>
              <option value="debit">Debit (Düşüm)</option>
              <option value="refund">Refund (İade)</option>
            </select>

            <select 
              class="select select-bordered"
              value={filters.source}
              onChange={(e) => setFilters({ ...filters, source: e.target.value })}
            >
              <option value="">Tüm Kaynaklar</option>
              <option value="admin">Admin</option>
              <option value="automation">Otomasyon</option>
              <option value="system">Sistem</option>
            </select>
          </div>

          <div class="overflow-x-auto">
            <table class="table">
              <thead>
                <tr>
                  <th>Tarih</th>
                  <th>Kullanıcı</th>
                  <th>İşlem Tipi</th>
                  <th>Tutar</th>
                  <th>Bakiye (Önce → Sonra)</th>
                  <th>Kaynak</th>
                  <th>Admin</th>
                </tr>
              </thead>
              <tbody>
                {transactions.map((tx) => (
                  <tr key={tx.id}>
                    <td class="text-sm">{formatDate(tx.created_at)}</td>
                    <td>
                      <div class="font-semibold">{tx.user?.name}</div>
                      <div class="text-sm text-base-content/70">{tx.user?.rfid_card_id}</div>
                    </td>
                    <td>
                      <span class={`badge ${getTransactionColor(tx.type)}`}>
                        {tx.type}
                      </span>
                    </td>
                    <td class="font-mono font-bold">{formatCurrency(tx.amount)}</td>
                    <td class="font-mono text-sm">
                      {formatCurrency(tx.balance_before)} → {formatCurrency(tx.balance_after)}
                    </td>
                    <td>
                      <span class={`badge ${getSourceBadge(tx.source).class}`}>
                        {getSourceBadge(tx.source).text}
                      </span>
                    </td>
                    <td class="text-sm">{tx.admin?.username || '-'}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          {transactions.length === 0 && (
            <div class="text-center py-8 text-base-content/70">
              İşlem bulunamadı
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
