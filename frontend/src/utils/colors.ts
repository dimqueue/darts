export function getDistanceColor(distance: number, darkMode: boolean): string {
    if (distance === 1) return 'bg-green-500 text-white';
    if (distance === 0) {
        return darkMode
            ? 'bg-gray-700/50 text-gray-400 border border-gray-600'
            : 'bg-gray-100 text-gray-500 border border-gray-300';
    }
    if (distance < 100) {
        return darkMode
            ? 'bg-emerald-500/30 text-emerald-200 border border-emerald-500/50'
            : 'bg-green-100 text-green-800';
    }
    if (distance < 500) {
        return darkMode
            ? 'bg-amber-500/30 text-amber-200 border border-amber-500/50'
            : 'bg-yellow-100 text-yellow-800';
    }
    if (distance < 1000) {
        return darkMode
            ? 'bg-orange-500/30 text-orange-200 border border-orange-500/50'
            : 'bg-orange-100 text-orange-800';
    }
    return darkMode
        ? 'bg-red-500/30 text-red-200 border border-red-500/50'
        : 'bg-red-100 text-red-800';
}
