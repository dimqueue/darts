export function getDistanceColor(distance: number): string {
    if (distance === 1) return 'bg-green-500 text-white';
    if (distance === 0) {
        return 'bg-gray-100 dark:bg-gray-700/50 text-gray-500 dark:text-gray-400 border border-gray-300 dark:border-gray-600';
    }
    if (distance < 100) {
        return 'bg-green-100 dark:bg-emerald-500/30 text-green-800 dark:text-emerald-200 dark:border dark:border-emerald-500/50';
    }
    if (distance < 500) {
        return 'bg-yellow-100 dark:bg-amber-500/30 text-yellow-800 dark:text-amber-200 dark:border dark:border-amber-500/50';
    }
    if (distance < 1000) {
        return 'bg-orange-100 dark:bg-orange-500/30 text-orange-800 dark:text-orange-200 dark:border dark:border-orange-500/50';
    }
    return 'bg-red-100 dark:bg-red-500/30 text-red-800 dark:text-red-200 dark:border dark:border-red-500/50';
}
