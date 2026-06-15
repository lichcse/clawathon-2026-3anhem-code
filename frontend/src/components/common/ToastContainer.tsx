import { Toast } from '@/stores/useUIStore';
import clsx from 'clsx';

interface ToastContainerProps {
  toasts: Toast[];
}

export default function ToastContainer({ toasts }: ToastContainerProps) {
  return (
    <div className="fixed bottom-4 right-4 space-y-2 z-50">
      {toasts.map((toast) => (
        <div
          key={toast.id}
          className={clsx(
            'px-4 py-3 rounded-lg shadow-lg text-white animate-pulse',
            {
              'bg-green-600': toast.type === 'success',
              'bg-red-600': toast.type === 'error',
              'bg-blue-600': toast.type === 'info',
              'bg-yellow-600': toast.type === 'warning',
            }
          )}
        >
          {toast.message}
        </div>
      ))}
    </div>
  );
}
