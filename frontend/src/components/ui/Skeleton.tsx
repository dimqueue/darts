interface SkeletonProps {
    className?: string;
    variant?: 'text' | 'title' | 'avatar' | 'button' | 'card' | 'icon';
}

export function Skeleton({ className = '', variant = 'text' }: SkeletonProps) {
    const baseClasses = 'animate-pulse rounded bg-gray-200 dark:bg-gray-700';

    const variants = {
        text: 'h-4 w-full',
        title: 'h-6 w-3/4',
        avatar: 'h-12 w-12 rounded-full',
        button: 'h-10 w-24 rounded-xl',
        card: 'h-32 w-full rounded-2xl',
        icon: 'h-8 w-8 rounded-lg',
    };

    return <div className={`${baseClasses} ${variants[variant] || ''} ${className}`} />;
}

interface SkeletonCardProps {
    children?: React.ReactNode;
    className?: string;
}

export function SkeletonCard({ children, className = '' }: SkeletonCardProps) {
    return (
        <div
            className={`rounded-2xl p-6 shadow-sm bg-white dark:bg-gray-800 ${className}`}
        >
            {children}
        </div>
    );
}

export function GamePageSkeleton() {
    return (
        <div className="space-y-6">
            <SkeletonCard>
                <div className="flex items-center justify-between">
                    <div className="flex items-center gap-3">
                        <Skeleton variant="icon" />
                        <div className="space-y-2">
                            <Skeleton className="h-6 w-40" />
                            <Skeleton className="h-4 w-32" />
                        </div>
                    </div>
                    <Skeleton className="h-10 w-28 rounded-xl" />
                </div>
            </SkeletonCard>

            <SkeletonCard>
                <div className="flex gap-3">
                    <Skeleton className="h-12 flex-1 rounded-xl" />
                    <Skeleton className="h-12 w-20 rounded-xl" />
                </div>
            </SkeletonCard>

            <SkeletonCard>
                <div className="flex items-center justify-between mb-4">
                    <Skeleton className="h-5 w-32" />
                    <Skeleton className="h-8 w-24 rounded-xl" />
                </div>
                <div className="space-y-2">
                    {[1, 2, 3].map((i) => (
                        <Skeleton key={i} className="h-12 w-full rounded-xl" />
                    ))}
                </div>
            </SkeletonCard>
        </div>
    );
}

export function LeaderboardSkeleton() {
    return (
        <div className="space-y-6">
            <div className="flex items-center justify-between">
                <Skeleton className="h-8 w-48" />
                <Skeleton className="h-10 w-32 rounded-xl" />
            </div>

            <div className="flex gap-2">
                {[1, 2, 3, 4].map((i) => (
                    <Skeleton key={i} className="h-10 w-20 rounded-xl" />
                ))}
            </div>

            <SkeletonCard className="space-y-3">
                {[1, 2, 3, 4, 5].map((i) => (
                    <div key={i} className="flex items-center gap-4">
                        <Skeleton className="h-8 w-8 rounded-full" />
                        <Skeleton className="h-5 flex-1" />
                        <Skeleton className="h-5 w-16" />
                    </div>
                ))}
            </SkeletonCard>
        </div>
    );
}

export function ProfileSkeleton() {
    return (
        <div className="space-y-6">
            <SkeletonCard>
                <div className="flex items-center gap-4">
                    <Skeleton variant="avatar" className="h-20 w-20" />
                    <div className="space-y-2">
                        <Skeleton className="h-7 w-48" />
                        <Skeleton className="h-4 w-32" />
                    </div>
                </div>
            </SkeletonCard>

            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                {[1, 2, 3, 4].map((i) => (
                    <SkeletonCard key={i} className="p-4">
                        <div className="flex items-center gap-3">
                            <Skeleton className="h-10 w-10 rounded-lg" />
                            <div className="space-y-1">
                                <Skeleton className="h-4 w-16" />
                                <Skeleton className="h-6 w-12" />
                            </div>
                        </div>
                    </SkeletonCard>
                ))}
            </div>

            <SkeletonCard>
                <Skeleton className="h-6 w-40 mb-4" />
                <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
                    {[1, 2, 3, 4, 5, 6].map((i) => (
                        <div key={i} className="space-y-1">
                            <Skeleton className="h-4 w-24" />
                            <Skeleton className="h-6 w-16" />
                        </div>
                    ))}
                </div>
            </SkeletonCard>
        </div>
    );
}

export function SettingsSkeleton() {
    return (
        <div className="max-w-2xl mx-auto space-y-6">
            <SkeletonCard>
                <div className="flex items-center gap-3">
                    <Skeleton className="h-8 w-8 rounded-lg" />
                    <div className="space-y-1">
                        <Skeleton className="h-7 w-32" />
                        <Skeleton className="h-4 w-48" />
                    </div>
                </div>
            </SkeletonCard>

            <SkeletonCard>
                <div className="flex items-center justify-between">
                    <div className="flex items-center gap-3">
                        <Skeleton className="h-5 w-5 rounded" />
                        <div className="space-y-1">
                            <Skeleton className="h-5 w-24" />
                            <Skeleton className="h-4 w-40" />
                        </div>
                    </div>
                    <Skeleton className="h-8 w-14 rounded-full" />
                </div>
            </SkeletonCard>

            <SkeletonCard>
                <div className="flex items-center gap-3 mb-4">
                    <Skeleton className="h-5 w-5 rounded" />
                    <Skeleton className="h-5 w-28" />
                </div>
                <div className="grid grid-cols-3 gap-4">
                    {[1, 2, 3].map((i) => (
                        <div key={i} className="p-4 rounded-xl border-2 border-transparent">
                            <Skeleton className="h-12 w-full rounded-lg mb-2" />
                            <Skeleton className="h-4 w-16 mx-auto" />
                        </div>
                    ))}
                </div>
            </SkeletonCard>

            {[1, 2, 3].map((i) => (
                <SkeletonCard key={i}>
                    <div className="flex items-center gap-3 mb-4">
                        <Skeleton className="h-5 w-5 rounded" />
                        <Skeleton className="h-5 w-32" />
                    </div>
                    <Skeleton className="h-12 w-full rounded-xl" />
                </SkeletonCard>
            ))}

            <Skeleton className="h-12 w-full rounded-xl" />
        </div>
    );
}
