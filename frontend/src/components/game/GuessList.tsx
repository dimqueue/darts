import { memo } from 'react';
import { RotateCcw } from '../ui/BoxIcon';
import Card from '../ui/Card';
import Button from '../ui/Button';
import { getDistanceColor } from '../../utils/colors';
import type { Guess } from '../../types/game';

interface GuessListProps {
    guesses: Guess[];
    onNewGame: () => void;
}

export default memo(function GuessList({ guesses, onNewGame }: GuessListProps) {
    const validGuesses = guesses.filter((g) => g.distance >= 0);

    const sortedGuesses = validGuesses.slice().sort((a, b) => {
        if (a.distance === 0 && b.distance !== 0) return 1;
        if (b.distance === 0 && a.distance !== 0) return -1;
        return a.distance - b.distance;
    });

    return (
        <Card>
            <div className="flex items-center justify-between mb-4">
                <h3 className="font-semibold text-gray-700 dark:text-gray-200">
                    Your Guesses ({validGuesses.filter((g) => g.distance > 0).length})
                </h3>
                <Button
                    onClick={onNewGame}
                    variant="outline"
                    icon={RotateCcw}
                    className="text-sm px-3 py-1"
                >
                    New Game
                </Button>
            </div>

            {validGuesses.length === 0 ? (
                <div className="text-center py-12 text-gray-400">
                    <p>Make your first guess!</p>
                    <p className="text-sm mt-1">The closer to 0, the closer you are</p>
                </div>
            ) : (
                <ul className="space-y-2" role="list" aria-label="Guesses">
                    {sortedGuesses.map((guess, index) => (
                        <li
                            key={guess.id || index}
                            className={`flex items-center justify-between p-3 rounded-xl ${getDistanceColor(guess.distance)}`}
                        >
                            <span className="font-medium">{guess.guess_word}</span>
                            <span className="font-bold">
                                {guess.distance === 1
                                    ? 'FOUND!'
                                    : guess.distance === 0
                                      ? '\u221E'
                                      : guess.distance}
                            </span>
                        </li>
                    ))}
                </ul>
            )}
        </Card>
    );
});
