import { ChevronUp, ChevronDown } from 'lucide-react';

interface SortableHeaderProps {
    label: string;
    field: string;
    sortBy: string;
    sortDirection: 'asc' | 'desc';
    onSort: (field: string) => void;
    align?: 'left' | 'right';
    darkMode: boolean;
}

export default function SortableHeader({
    label,
    field,
    sortBy,
    sortDirection,
    onSort,
    align = 'left',
    darkMode,
}: SortableHeaderProps) {
    const isActive = sortBy === field;
    const alignClass = align === 'right' ? 'text-right justify-end' : 'text-left';

    return (
        <th
            className={`px-6 py-4 text-xs font-semibold uppercase cursor-pointer transition-colors select-none ${alignClass} ${
                darkMode ? 'text-gray-400 hover:bg-gray-700' : 'text-gray-500 hover:bg-gray-100'
            }`}
            onClick={() => onSort(field)}
            aria-sort={isActive ? (sortDirection === 'asc' ? 'ascending' : 'descending') : 'none'}
        >
            <div className={`flex items-center gap-1 ${align === 'right' ? 'justify-end' : ''}`}>
                {label}
                {isActive &&
                    (sortDirection === 'asc' ? (
                        <ChevronUp className="w-4 h-4" />
                    ) : (
                        <ChevronDown className="w-4 h-4" />
                    ))}
            </div>
        </th>
    );
}
