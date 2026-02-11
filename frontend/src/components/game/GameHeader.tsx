import { memo } from 'react';
import { Trophy } from '../ui/BoxIcon';
import Card from '../ui/Card';
import LanguageSelect from '../ui/LanguageSelect';
import { LANGUAGES } from '../../config/constants';
import { getGameMode } from '../../config/gameModes';
import type { GameModeId } from '../../types/game';

interface GameHeaderProps {
    language: string;
    onLanguageChange: (language: string) => void;
    mode?: GameModeId;
}

export default memo(function GameHeader({ language, onLanguageChange, mode }: GameHeaderProps) {
    const modeConfig = mode ? getGameMode(mode) : null;
    const Icon = modeConfig?.icon ?? Trophy;
    const title = modeConfig?.name ?? 'Daily Challenge';

    return (
        <Card>
            <div className="flex items-center justify-between">
                <div className="flex items-center gap-3">
                    <Icon
                        className="w-8 h-8 text-theme-text"
                    />
                    <div>
                        <h1 className="text-2xl font-bold text-gray-800 dark:text-white">
                            {title}
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
