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
    const containerClass = 'bg-gray-100 dark:bg-gray-800';
    const inactiveClass = 'text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-200 hover:bg-white/50 dark:hover:bg-gray-700/50';

    return (
        <div className={`flex gap-1 p-1 rounded-xl ${containerClass}`} role="tablist">
            {tabs.map((tab) => (
                <button
                    key={tab.id}
                    onClick={() => onTabChange(tab.id)}
                    className={`px-4 py-2 rounded-lg font-medium transition-all duration-200 ${
                        activeTab === tab.id
                            ? 'bg-theme-gradient text-white shadow-sm'
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
