import { format, formatDistance, parseISO } from 'date-fns';
import { tr } from 'date-fns/locale';

export const formatDate = (date, formatStr = 'dd MMM yyyy HH:mm') => {
  if (!date) return '-';
  const d = typeof date === 'string' ? parseISO(date) : date;
  return format(d, formatStr, { locale: tr });
};

export const formatRelativeTime = (date) => {
  if (!date) return '-';
  const d = typeof date === 'string' ? parseISO(date) : date;
  return formatDistance(d, new Date(), { addSuffix: true, locale: tr });
};

export const formatCurrency = (amount) => {
  return new Intl.NumberFormat('tr-TR', {
    style: 'currency',
    currency: 'TRY',
  }).format(amount || 0);
};

export const getTransactionColor = (type) => {
  const colors = {
    credit: 'badge-success',
    debit: 'badge-error',
    refund: 'badge-info',
  };
  return colors[type] || 'badge-ghost';
};

export const getSourceBadge = (source) => {
  const badges = {
    admin: { text: 'Admin', class: 'badge-primary' },
    automation: { text: 'Otomasyon', class: 'badge-secondary' },
    system: { text: 'Sistem', class: 'badge-accent' },
  };
  return badges[source] || { text: source, class: 'badge-ghost' };
};

export const getBalanceColor = (balance) => {
  if (balance >= 50) return 'text-success';
  if (balance >= 20) return 'text-warning';
  return 'text-error';
};

export const cn = (...classes) => {
  return classes.filter(Boolean).join(' ');
};
