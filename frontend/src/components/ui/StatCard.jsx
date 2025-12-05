import { useTheme } from '../../contexts/ThemeContext';

export default function StatCard({ label, value, icon: Icon }) {
    const { theme } = useTheme();

    return (
        <div className="bg-white rounded-xl p-4 shadow-card border border-gray-100">
            <div className="flex items-center gap-3">
                {Icon && (
                    <div className={`p-2 rounded-lg ${theme.gradient}`}>
                        <Icon className="w-5 h-5 text-white" />
                    </div>
                )}
                <div>
                    <p className="text-sm text-gray-500">{label}</p>
                    <p className="text-xl font-bold text-gray-800">{value}</p>
                </div>
            </div>
        </div>
    );
}