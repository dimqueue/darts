import { memo } from 'react';
import { ChevronUp, ChevronDown } from '../ui/BoxIcon';

interface SortableHeaderProps {
    label: string;
    field: string;
    sortBy: string;
    sortDirection: 'asc' | 'desc';
    onSort: (field: string) => void;
    align?: 'left' | 'center' | 'right';
}

export default memo(function SortableHeader({
    label,
    field,
    sortBy,
    sortDirection,
    onSort,
    align = 'center',
}: SortableHeaderProps) {
    const isActive = sortBy === field;
    const alignMap = {
        left: 'text-left',
        center: 'text-center justify-center',
        right: 'text-right justify-end',
    };

    return (
        <th
            className={`px-6 py-4 text-xs font-semibold uppercase cursor-pointer select-none ${alignMap[align]} text-gray-500 dark:text-gray-400 hover:bg-theme-hover-bg`}
            onClick={() => onSort(field)}
            aria-sort={isActive ? (sortDirection === 'asc' ? 'ascending' : 'descending') : 'none'}
        >
            <div className={`flex items-center gap-1 ${alignMap[align]}`}>
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
});