import { ChevronLeft, ChevronRight } from 'lucide-react';
import { useTheme } from '../../contexts/ThemeContext';

export default function Pagination({
    currentPage,
    totalPages,
    onPageChange,
}) {
    const { theme } = useTheme();

    if (totalPages <= 1) return null;

    const getPageNumbers = () => {
        const pages = [];
        const showPages = 5;
        let start = Math.max(1, currentPage - Math.floor(showPages / 2));
        let end = Math.min(totalPages, start + showPages - 1);

        if (end - start + 1 < showPages) {
            start = Math.max(1, end - showPages + 1);
        }

        for (let i = start; i <= end; i++) {
            pages.push(i);
        }
        return pages;
    };

    const buttonBase = 'px-3 py-2 rounded-lg font-medium transition-colors disabled:opacity-50 disabled:cursor-not-allowed';

    return (
        <div className="flex items-center justify-center gap-2">
            <button
                onClick={() => onPageChange(currentPage - 1)}
                disabled={currentPage === 1}
                className={`${buttonBase} hover:bg-gray-100`}
            >
                <ChevronLeft className="w-5 h-5" />
            </button>

            {getPageNumbers().map((page) => (
                <button
                    key={page}
                    onClick={() => onPageChange(page)}
                    className={`${buttonBase} ${
                        page === currentPage
                            ? `${theme.gradient} text-white`
                            : 'hover:bg-gray-100 text-gray-700'
                    }`}
                >
                    {page}
                </button>
            ))}

            <button
                onClick={() => onPageChange(currentPage + 1)}
                disabled={currentPage === totalPages}
                className={`${buttonBase} hover:bg-gray-100`}
            >
                <ChevronRight className="w-5 h-5" />
            </button>
        </div>
    );
}