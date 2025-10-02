import { useState } from 'preact/hooks';
import { useAuthStore } from '../lib/store';
import { route } from 'preact-router';

export default function Login() {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const login = useAuthStore((state) => state.login);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      await login({ username, password });
      route('/');
    } catch (err) {
      setError(err.message || 'Giriş başarısız. Lütfen tekrar deneyin.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-primary/20 to-secondary/20">
      <div class="card w-full max-w-md bg-base-100 shadow-2xl">
        <div class="card-body">
          <div class="text-center mb-6">
            <h1 class="text-3xl font-bold text-primary">HudAutomata</h1>
            <p class="text-base-content/70 mt-2">RFID Bakiye Yönetim Sistemi</p>
          </div>

          {error && (
            <div class="alert alert-error mb-4">
              <span>{error}</span>
            </div>
          )}

          <form onSubmit={handleSubmit}>
            <div class="form-control">
              <label class="label">
                <span class="label-text">Kullanıcı Adı</span>
              </label>
              <input
                type="text"
                placeholder="admin"
                class="input input-bordered"
                value={username}
                onInput={(e) => setUsername(e.target.value)}
                required
                disabled={loading}
              />
            </div>

            <div class="form-control mt-4">
              <label class="label">
                <span class="label-text">Şifre</span>
              </label>
              <input
                type="password"
                placeholder="••••••••"
                class="input input-bordered"
                value={password}
                onInput={(e) => setPassword(e.target.value)}
                required
                disabled={loading}
              />
            </div>

            <div class="form-control mt-6">
              <button 
                type="submit" 
                class="btn btn-primary"
                disabled={loading}
              >
                {loading ? (
                  <span class="loading loading-spinner"></span>
                ) : (
                  'Giriş Yap'
                )}
              </button>
            </div>
          </form>

          <div class="divider">veya</div>

          <div class="text-center text-sm text-base-content/60">
            <p>Varsayılan Kullanıcı:</p>
            <p class="font-mono mt-1">admin / admin123</p>
          </div>
        </div>
      </div>
    </div>
  );
}
