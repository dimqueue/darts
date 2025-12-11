import { useTheme } from '../../contexts/ThemeContext';

export default function Tabs({ tabs, activeTab, onTabChange }) {
    const { theme } = useTheme();

    return (
        <div className="flex gap-1 bg-gray-100 p-1 rounded-xl">
            {tabs.map((tab) => (
                <button
                    key={tab.id}
                    onClick={() => onTabChange(tab.id)}
                    className={`px-4 py-2 rounded-lg font-medium transition-all duration-200 ${
                        activeTab === tab.id
                            ? `${theme.gradient} text-white shadow-sm`
                            : 'text-gray-600 hover:text-gray-900 hover:bg-white/50'
                    }`}
                >
                    {tab.label}
                </button>
            ))}
        </div>
    );
}
