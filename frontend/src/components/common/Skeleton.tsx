import { cn } from '../../utils/cn'

interface SkeletonProps {
  className?: string
}

export function Skeleton({ className }: SkeletonProps) {
  return (
    <div
      className={cn(
        'animate-pulse bg-gray-700 rounded',
        className
      )}
    />
  )
}

export function SeatGridSkeleton() {
  return (
    <div className="grid grid-cols-3 sm:grid-cols-6 gap-3">
      {Array.from({ length: 18 }).map((_, i) => (
        <div key={i} className="flex flex-col items-center gap-1">
          <Skeleton className="w-12 h-12 rounded-full" />
          <Skeleton className="w-10 h-3" />
        </div>
      ))}
    </div>
  )
}

export function MemberListSkeleton() {
  return (
    <div className="space-y-2">
      {Array.from({ length: 4 }).map((_, i) => (
        <div key={i} className="flex items-center gap-3 p-2">
          <Skeleton className="w-8 h-8 rounded-full" />
          <Skeleton className="flex-1 h-4" />
        </div>
      ))}
    </div>
  )
}

export function RoomCardSkeleton() {
  return (
    <div className="bg-gray-800 rounded-xl p-4 space-y-2">
      <Skeleton className="h-5 w-3/4" />
      <Skeleton className="h-3 w-1/2" />
      <Skeleton className="h-8 w-full mt-3" />
    </div>
  )
}
