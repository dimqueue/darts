interface SettingsToggleProps {
    label: string;
    description?: string;
    checked: boolean;
    onChange: (checked: boolean) => void;
}

export default function SettingsToggle({ label, description, checked, onChange }: SettingsToggleProps) {
    return (
        <label className="flex items-center justify-between cursor-pointer">
            <div>
                <span className="text-gray-700 dark:text-gray-200">{label}</span>
                {description && (
                    <p className="text-sm text-gray-500 dark:text-gray-400">{description}</p>
                )}
            </div>
            <input
                type="checkbox"
                checked={checked}
                onChange={(e) => onChange(e.target.checked)}
                className="w-5 h-5 rounded accent-theme"
            />
        </label>
    );
}
