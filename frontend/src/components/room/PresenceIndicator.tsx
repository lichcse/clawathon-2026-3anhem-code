interface PresenceIndicatorProps {
  isOnline: boolean
  size?: 'sm' | 'md'
}

export function PresenceIndicator({ isOnline, size = 'sm' }: PresenceIndicatorProps) {
  const dim = size === 'sm' ? 'w-2.5 h-2.5' : 'w-3.5 h-3.5'
  return (
    <span
      className={`inline-block rounded-full ${dim} ${
        isOnline ? 'bg-green-400 shadow-[0_0_6px_#4ade80]' : 'bg-gray-500'
      }`}
      title={isOnline ? 'Online' : 'Offline'}
    />
  )
}
