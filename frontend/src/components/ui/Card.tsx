import { useTheme } from '../../contexts/ThemeContext';

interface CardProps {
    children: React.ReactNode;
    className?: string;
    padding?: string;
}

export default function Card({ children, className = '', padding = 'p-6' }: CardProps) {
    const { darkMode } = useTheme();

    return (
        <div
            className={`rounded-card shadow-card transition-colors duration-300 ${padding} ${darkMode ? 'bg-gray-800 text-gray-100' : 'bg-white text-gray-800'} ${className}`}
        >
            {children}
        </div>
    );
}
