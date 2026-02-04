import { ChevronLeft, ChevronRight } from './BoxIcon';
import { useTheme } from '../../contexts/ThemeContext';

interface PaginationProps {
    currentPage: number;
    totalPages: number;
    onPageChange: (page: number) => void;
}

export default function Pagination({ currentPage, totalPages, onPageChange }: PaginationProps) {
    const { theme } = useTheme();

    if (totalPages <= 1) return null;

    const getPageNumbers = (): number[] => {
        const pages: number[] = [];
        const showPages = 5;
        let start = Math.max(1, currentPage - Math.floor(showPages / 2));
        const end = Math.min(totalPages, start + showPages - 1);

        if (end - start + 1 < showPages) {
            start = Math.max(1, end - showPages + 1);
        }

        for (let i = start; i <= end; i++) {
            pages.push(i);
        }
        return pages;
    };

    const buttonBase =
        'px-3 py-2 rounded-lg font-medium transition-colors disabled:opacity-50 disabled:cursor-not-allowed';

    const navButtonClass = 'hover:bg-gray-100 dark:hover:bg-gray-700 text-gray-700 dark:text-gray-300';

    const inactivePageClass = 'hover:bg-gray-100 dark:hover:bg-gray-700 text-gray-700 dark:text-gray-300';

    return (
        <nav className="flex items-center justify-center gap-2" aria-label="Pagination">
            <button
                onClick={() => onPageChange(currentPage - 1)}
                disabled={currentPage === 1}
                className={`${buttonBase} ${navButtonClass}`}
                aria-label="Previous page"
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
                            : inactivePageClass
                    }`}
                    aria-label={`Page ${page}`}
                    aria-current={page === currentPage ? 'page' : undefined}
                >
                    {page}
                </button>
            ))}

            <button
                onClick={() => onPageChange(currentPage + 1)}
                disabled={currentPage === totalPages}
                className={`${buttonBase} ${navButtonClass}`}
                aria-label="Next page"
            >
                <ChevronRight className="w-5 h-5" />
            </button>
        </nav>
    );
}
