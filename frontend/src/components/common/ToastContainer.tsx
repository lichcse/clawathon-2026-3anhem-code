import { Toast } from '@/stores/useUIStore';
import { useUIStore } from '@/stores/useUIStore';
import clsx from 'clsx';

interface ToastContainerProps {
  toasts: Toast[];
}

export default function ToastContainer({ toasts }: ToastContainerProps) {
  const removeToast = useUIStore((state) => state.removeToast);

  return (
    <div className="fixed bottom-20 sm:bottom-4 right-4 space-y-2 z-50 max-w-xs w-full pointer-events-none">
      {toasts.map((toast) => (
        <div
          key={toast.id}
          onClick={() => removeToast(toast.id)}
          className={clsx(
            'px-4 py-3 rounded-lg shadow-xl text-white text-sm font-medium',
            'transition-all duration-300 ease-in-out pointer-events-auto cursor-pointer',
            'animate-[fadeInUp_0.3s_ease-out]',
            {
              'bg-green-600 hover:bg-green-700': toast.type === 'success',
              'bg-red-600 hover:bg-red-700': toast.type === 'error',
              'bg-blue-600 hover:bg-blue-700': toast.type === 'info',
              'bg-yellow-600 hover:bg-yellow-700': toast.type === 'warning',
            }
          )}
        >
          {toast.message}
        </div>
      ))}
    </div>
  );
}
