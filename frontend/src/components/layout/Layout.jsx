import Navbar from './Navbar';
import { useTheme } from '../../contexts/ThemeContext';

export default function Layout({ children }) {
    const { theme } = useTheme();

    return (
        <div className={`min-h-screen bg-gradient-to-br ${theme.bgGradient}`}>
            <Navbar />
            <main className="max-w-6xl mx-auto px-4 py-8">
                {children}
            </main>
        </div>
    );
}