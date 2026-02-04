import { useTheme } from '../../contexts/ThemeContext';

interface Tab {
    id: string;
    label: string;
}

interface TabsProps {
    tabs: Tab[];
    activeTab: string;
    onTabChange: (id: string) => void;
}

export default function Tabs({ tabs, activeTab, onTabChange }: TabsProps) {
    const { theme, darkMode } = useTheme();

    const containerClass = darkMode ? 'bg-gray-800' : 'bg-gray-100';
    const inactiveClass = darkMode
        ? 'text-gray-400 hover:text-gray-200 hover:bg-gray-700/50'
        : 'text-gray-600 hover:text-gray-900 hover:bg-white/50';

    return (
        <div className={`flex gap-1 p-1 rounded-xl ${containerClass}`} role="tablist">
            {tabs.map((tab) => (
                <button
                    key={tab.id}
                    onClick={() => onTabChange(tab.id)}
                    className={`px-4 py-2 rounded-lg font-medium transition-all duration-200 ${
                        activeTab === tab.id
                            ? `${theme.gradient} text-white shadow-sm`
                            : inactiveClass
                    }`}
                    role="tab"
                    aria-selected={activeTab === tab.id}
                >
                    {tab.label}
                </button>
            ))}
        </div>
    );
}
