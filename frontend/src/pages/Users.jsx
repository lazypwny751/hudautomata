import { useState, useEffect } from 'preact/hooks';
import { usersAPI, transactionsAPI } from '../lib/api';
import { formatCurrency, getBalanceColor } from '../lib/utils';
import { Plus, Search, CreditCard, Wallet } from 'lucide-preact';

export default function Users() {
  const [users, setUsers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState('');
  const [showAddModal, setShowAddModal] = useState(false);
  const [showBalanceModal, setShowBalanceModal] = useState(false);
  const [selectedUser, setSelectedUser] = useState(null);

  useEffect(() => {
    loadUsers();
  }, [search]);

  const loadUsers = async () => {
    try {
      const response = await usersAPI.list({ search });
      setUsers(response.data || []);
    } catch (error) {
      console.error('Failed to load users:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleAddBalance = (user) => {
    setSelectedUser(user);
    setShowBalanceModal(true);
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
      <div class="flex justify-between items-center">
        <div>
          <h1 class="text-3xl font-bold">Kullanıcılar</h1>
          <p class="text-base-content/70 mt-1">RFID kullanıcı yönetimi</p>
        </div>
        <button 
          class="btn btn-primary gap-2"
          onClick={() => setShowAddModal(true)}
        >
          <Plus size={20} />
          Yeni Kullanıcı
        </button>
      </div>

      <div class="card bg-base-100 shadow-xl">
        <div class="card-body">
          <div class="flex gap-4 mb-4">
            <div class="form-control flex-1">
              <div class="input-group">
                <span><Search size={20} /></span>
                <input
                  type="text"
                  placeholder="İsim veya RFID kartı ile ara..."
                  class="input input-bordered w-full"
                  value={search}
                  onInput={(e) => setSearch(e.target.value)}
                />
              </div>
            </div>
          </div>

          <div class="overflow-x-auto">
            <table class="table">
              <thead>
                <tr>
                  <th>RFID Kart ID</th>
                  <th>İsim</th>
                  <th>İletişim</th>
                  <th>Bakiye</th>
                  <th>Durum</th>
                  <th>İşlemler</th>
                </tr>
              </thead>
              <tbody>
                {users.map((user) => (
                  <tr key={user.id}>
                    <td>
                      <div class="flex items-center gap-2">
                        <CreditCard size={16} class="text-primary" />
                        <span class="font-mono">{user.rfid_card_id}</span>
                      </div>
                    </td>
                    <td class="font-semibold">{user.name}</td>
                    <td>
                      <div class="text-sm">
                        <div>{user.email || '-'}</div>
                        <div class="text-base-content/70">{user.phone || '-'}</div>
                      </div>
                    </td>
                    <td>
                      <div class={`font-mono font-bold ${getBalanceColor(user.balance)}`}>
                        {formatCurrency(user.balance)}
                      </div>
                    </td>
                    <td>
                      <span class={`badge ${user.is_active ? 'badge-success' : 'badge-error'}`}>
                        {user.is_active ? 'Aktif' : 'Pasif'}
                      </span>
                    </td>
                    <td>
                      <button 
                        class="btn btn-sm btn-primary gap-2"
                        onClick={() => handleAddBalance(user)}
                      >
                        <Wallet size={16} />
                        Bakiye Yükle
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          {users.length === 0 && (
            <div class="text-center py-8 text-base-content/70">
              {search ? 'Kullanıcı bulunamadı' : 'Henüz kullanıcı eklenmemiş'}
            </div>
          )}
        </div>
      </div>

      {showAddModal && <AddUserModal onClose={() => { setShowAddModal(false); loadUsers(); }} />}
      {showBalanceModal && selectedUser && (
        <BalanceModal 
          user={selectedUser} 
          onClose={() => { setShowBalanceModal(false); setSelectedUser(null); loadUsers(); }} 
        />
      )}
    </div>
  );
}

function AddUserModal({ onClose }) {
  const [formData, setFormData] = useState({
    rfid_card_id: '',
    name: '',
    email: '',
    phone: '',
    balance: 0,
  });
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    try {
      await usersAPI.create(formData);
      onClose();
    } catch (error) {
      alert('Hata: ' + (error.message || 'Kullanıcı oluşturulamadı'));
    } finally {
      setLoading(false);
    }
  };

  return (
    <div class="modal modal-open">
      <div class="modal-box">
        <h3 class="font-bold text-lg mb-4">Yeni Kullanıcı Ekle</h3>
        
        <form onSubmit={handleSubmit} class="space-y-4">
          <div class="form-control">
            <label class="label"><span class="label-text">RFID Kart ID *</span></label>
            <input
              type="text"
              class="input input-bordered"
              value={formData.rfid_card_id}
              onInput={(e) => setFormData({ ...formData, rfid_card_id: e.target.value })}
              required
            />
          </div>

          <div class="form-control">
            <label class="label"><span class="label-text">İsim *</span></label>
            <input
              type="text"
              class="input input-bordered"
              value={formData.name}
              onInput={(e) => setFormData({ ...formData, name: e.target.value })}
              required
            />
          </div>

          <div class="form-control">
            <label class="label"><span class="label-text">Email (Opsiyonel)</span></label>
            <input
              type="email"
              class="input input-bordered"
              value={formData.email}
              onInput={(e) => setFormData({ ...formData, email: e.target.value })}
            />
          </div>

          <div class="form-control">
            <label class="label"><span class="label-text">Telefon (Opsiyonel)</span></label>
            <input
              type="tel"
              class="input input-bordered"
              value={formData.phone}
              onInput={(e) => setFormData({ ...formData, phone: e.target.value })}
            />
          </div>

          <div class="form-control">
            <label class="label"><span class="label-text">Başlangıç Bakiyesi</span></label>
            <input
              type="number"
              step="0.01"
              class="input input-bordered"
              value={formData.balance}
              onInput={(e) => setFormData({ ...formData, balance: parseFloat(e.target.value) })}
            />
          </div>

          <div class="modal-action">
            <button type="button" class="btn" onClick={onClose}>İptal</button>
            <button type="submit" class="btn btn-primary" disabled={loading}>
              {loading ? <span class="loading loading-spinner"></span> : 'Ekle'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}

function BalanceModal({ user, onClose }) {
  const [amount, setAmount] = useState('');
  const [description, setDescription] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    try {
      await transactionsAPI.create({
        user_id: user.id,
        type: 'credit',
        amount: parseFloat(amount),
        description: description || 'Bakiye yükleme',
      });
      onClose();
    } catch (error) {
      alert('Hata: ' + (error.message || 'Bakiye yüklenemedi'));
    } finally {
      setLoading(false);
    }
  };

  return (
    <div class="modal modal-open">
      <div class="modal-box">
        <h3 class="font-bold text-lg mb-4">Bakiye Yükle</h3>
        
        <div class="bg-base-200 p-4 rounded-lg mb-4">
          <div class="flex justify-between items-center">
            <div>
              <p class="font-semibold">{user.name}</p>
              <p class="text-sm text-base-content/70">{user.rfid_card_id}</p>
            </div>
            <div class="text-right">
              <p class="text-sm text-base-content/70">Mevcut Bakiye</p>
              <p class={`font-mono font-bold ${getBalanceColor(user.balance)}`}>
                {formatCurrency(user.balance)}
              </p>
            </div>
          </div>
        </div>

        <form onSubmit={handleSubmit} class="space-y-4">
          <div class="form-control">
            <label class="label"><span class="label-text">Yüklenecek Tutar</span></label>
            <input
              type="number"
              step="0.01"
              min="0.01"
              class="input input-bordered"
              value={amount}
              onInput={(e) => setAmount(e.target.value)}
              required
            />
          </div>

          <div class="form-control">
            <label class="label"><span class="label-text">Açıklama (Opsiyonel)</span></label>
            <textarea
              class="textarea textarea-bordered"
              value={description}
              onInput={(e) => setDescription(e.target.value)}
            />
          </div>

          <div class="modal-action">
            <button type="button" class="btn" onClick={onClose}>İptal</button>
            <button type="submit" class="btn btn-primary" disabled={loading}>
              {loading ? <span class="loading loading-spinner"></span> : 'Yükle'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
