import { memo } from 'react';
import { Trophy } from '../ui/BoxIcon';
import Card from '../ui/Card';
import LanguageSelect from '../ui/LanguageSelect';
import { LANGUAGES } from '../../config/constants';

interface GameHeaderProps {
    language: string;
    onLanguageChange: (language: string) => void;
}

export default memo(function GameHeader({ language, onLanguageChange }: GameHeaderProps) {
    return (
        <Card>
            <div className="flex items-center justify-between">
                <div className="flex items-center gap-3">
                    <Trophy
                        className="w-8 h-8 text-theme-text"
                    />
                    <div>
                        <h1 className="text-2xl font-bold text-gray-800 dark:text-white">
                            Daily Challenge
                        </h1>
                        <p className="text-sm text-gray-500 dark:text-gray-400">
                            Guess the secret word
                        </p>
                    </div>
                </div>
                <LanguageSelect
                    id="language-select"
                    value={language}
                    onChange={onLanguageChange}
                    languages={LANGUAGES}
                />
            </div>
        </Card>
    );
});
