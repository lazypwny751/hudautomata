import { useState, useEffect } from 'preact/hooks';
import { automationAPI } from '../lib/api';
import { formatCurrency, formatDate, getSourceBadge } from '../lib/utils';
import { Scan, Activity } from 'lucide-preact';

export default function Automation() {
  const [history, setHistory] = useState([]);
  const [loading, setLoading] = useState(true);
  const [rfidInput, setRfidInput] = useState('');
  const [scanResult, setScanResult] = useState(null);

  useEffect(() => {
    loadHistory();
  }, []);

  const loadHistory = async () => {
    try {
      const data = await automationAPI.getHistory();
      setHistory(data);
    } catch (error) {
      console.error('Failed to load history:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleScan = async (e) => {
    e.preventDefault();
    try {
      const result = await automationAPI.checkBalance(rfidInput);
      setScanResult(result);
    } catch (error) {
      setScanResult({ error: error.message });
    }
  };

  return (
    <div class="p-6 space-y-6">
      <div>
        <h1 class="text-3xl font-bold">Otomasyon</h1>
        <p class="text-base-content/70 mt-1">RFID okutma ve otomasyon geçmişi</p>
      </div>

      {/* RFID Scan Simulator */}
      <div class="card bg-base-100 shadow-xl">
        <div class="card-body">
          <h2 class="card-title">
            <Scan size={24} />
            RFID Kart Okuyucu Simülatörü
          </h2>
          
          <form onSubmit={handleScan} class="flex gap-4">
            <input
              type="text"
              placeholder="RFID Kart ID'sini girin..."
              class="input input-bordered flex-1"
              value={rfidInput}
              onInput={(e) => setRfidInput(e.target.value)}
            />
            <button type="submit" class="btn btn-primary">
              Oku
            </button>
          </form>

          {scanResult && (
            <div class={`alert ${scanResult.error ? 'alert-error' : 'alert-success'} mt-4`}>
              {scanResult.error ? (
                <span>{scanResult.error}</span>
              ) : (
                <div class="w-full">
                  <div class="flex justify-between items-center">
                    <div>
                      <p class="font-semibold">{scanResult.user_name}</p>
                      <p class="text-sm">Bakiye: {formatCurrency(scanResult.balance)}</p>
                    </div>
                    <span class={`badge ${scanResult.is_active ? 'badge-success' : 'badge-error'}`}>
                      {scanResult.is_active ? 'Aktif' : 'Pasif'}
                    </span>
                  </div>
                </div>
              )}
            </div>
          )}
        </div>
      </div>

      {/* Automation History */}
      <div class="card bg-base-100 shadow-xl">
        <div class="card-body">
          <h2 class="card-title">
            <Activity size={24} />
            Otomasyon Geçmişi
          </h2>

          {loading ? (
            <div class="flex justify-center py-8">
              <span class="loading loading-spinner"></span>
            </div>
          ) : (
            <div class="overflow-x-auto">
              <table class="table">
                <thead>
                  <tr>
                    <th>Tarih</th>
                    <th>Kullanıcı</th>
                    <th>RFID</th>
                    <th>Tutar</th>
                    <th>Açıklama</th>
                    <th>Durum</th>
                  </tr>
                </thead>
                <tbody>
                  {history.map((item) => (
                    <tr key={item.id}>
                      <td class="text-sm">{formatDate(item.created_at)}</td>
                      <td class="font-semibold">{item.user?.name}</td>
                      <td class="font-mono text-sm">{item.user?.rfid_card_id}</td>
                      <td class="font-mono">{formatCurrency(item.amount)}</td>
                      <td class="text-sm">{item.description || '-'}</td>
                      <td>
                        <span class="badge badge-secondary">Otomatik</span>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>

              {history.length === 0 && (
                <div class="text-center py-8 text-base-content/70">
                  Henüz otomasyon işlemi yok
                </div>
              )}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
